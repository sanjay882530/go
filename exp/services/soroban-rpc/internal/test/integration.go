package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/hcnet/go/clients/auroraclient"
	"github.com/hcnet/go/clients/hcnetcore"
	"github.com/hcnet/go/exp/services/soroban-rpc/internal"
	"github.com/hcnet/go/exp/services/soroban-rpc/internal/methods"
	"github.com/hcnet/go/support/log"
)

const (
	StandaloneNetworkPassphrase = "Standalone Network ; February 2017"
	hcnetCoreProtocolVersion  = 19
	hcnetCorePort             = 11626
)

type Test struct {
	t *testing.T

	composePath string

	handler       internal.Handler
	server        *httptest.Server
	auroraClient *auroraclient.Client

	coreClient *hcnetcore.Client

	shutdownOnce  sync.Once
	shutdownCalls []func()
}

func NewTest(t *testing.T) *Test {
	if os.Getenv("SOROBAN_RPC_INTEGRATION_TESTS_ENABLED") == "" {
		t.Skip("skipping integration test: SOROBAN_RPC_INTEGRATION_TESTS_ENABLED not set")
	}

	composePath := findDockerComposePath()
	i := &Test{
		t:           t,
		composePath: composePath,
	}

	// Only run Hcnet Core container and its dependencies.
	i.runComposeCommand("up", "--detach", "--quiet-pull", "--no-color")
	i.prepareShutdownHandlers()
	i.coreClient = &hcnetcore.Client{URL: "http://localhost:" + strconv.Itoa(hcnetCorePort)}
	i.auroraClient = &auroraclient.Client{AuroraURL: "http://localhost:8000"}
	i.waitForCore()
	i.waitForAurora()
	i.configureJSONRPCServer()

	return i
}

func (i *Test) configureJSONRPCServer() {
	logger := log.New()

	var err error
	i.handler, err = internal.NewJSONRPCHandler(internal.HandlerParams{
		AccountStore: methods.AccountStore{
			Client: i.auroraClient,
		},
		Logger: logger,
	})
	if err != nil {
		i.t.Fatalf("cannot create handler: %v", err)
	}
	i.server = httptest.NewServer(i.handler)
}

// Runs a docker-compose command applied to the above configs
func (i *Test) runComposeCommand(args ...string) {
	integrationYaml := filepath.Join(i.composePath, "docker-compose.yml")

	cmdline := append([]string{"-f", integrationYaml}, args...)
	cmd := exec.Command("docker-compose", cmdline...)

	i.t.Log("Running", cmd.Env, cmd.Args)
	out, innerErr := cmd.Output()
	if exitErr, ok := innerErr.(*exec.ExitError); ok {
		fmt.Printf("stdout:\n%s\n", string(out))
		fmt.Printf("stderr:\n%s\n", string(exitErr.Stderr))
	}

	if innerErr != nil {
		i.t.Fatalf("Compose command failed: %v", innerErr)
	}
}

func (i *Test) prepareShutdownHandlers() {
	i.shutdownCalls = append(i.shutdownCalls,
		func() {
			i.handler.Close()
			i.server.Close()
			i.runComposeCommand("down", "-v")
		},
	)

	// Register cleanup handlers (on panic and ctrl+c) so the containers are
	// stopped even if ingestion or testing fails.
	i.t.Cleanup(i.Shutdown)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		i.Shutdown()
		os.Exit(int(syscall.SIGTERM))
	}()
}

// Shutdown stops the integration tests and destroys all its associated
// resources. It will be implicitly called when the calling test (i.e. the
// `testing.Test` passed to `New()`) is finished if it hasn't been explicitly
// called before.
func (i *Test) Shutdown() {
	i.shutdownOnce.Do(func() {
		// run them in the opposite order in which they where added
		for callI := len(i.shutdownCalls) - 1; callI >= 0; callI-- {
			i.shutdownCalls[callI]()
		}
	})
}

// Wait for core to be up and manually close the first ledger
func (i *Test) waitForCore() {
	i.t.Log("Waiting for core to be up...")
	for t := 30 * time.Second; t >= 0; t -= time.Second {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err := i.coreClient.Info(ctx)
		cancel()
		if err != nil {
			i.t.Logf("could not obtain info response: %v", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	i.UpgradeProtocol(hcnetCoreProtocolVersion)

	for t := 0; t < 5; t++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		info, err := i.coreClient.Info(ctx)
		cancel()
		if err != nil || !info.IsSynced() {
			i.t.Logf("Core is still not synced: %v %v", err, info)
			time.Sleep(time.Second)
			continue
		}
		i.t.Log("Core is up.")
		return
	}
	i.t.Fatal("Core could not sync after 30s")
}

func (i *Test) waitForAurora() {
	for t := 60; t >= 0; t -= 1 {
		time.Sleep(time.Second)

		i.t.Log("Waiting for ingestion and protocol upgrade...")
		root, err := i.auroraClient.Root()
		if err != nil {
			i.t.Logf("could not obtain root response %v", err)
			continue
		}

		if root.AuroraSequence < 3 ||
			int(root.AuroraSequence) != int(root.IngestSequence) {
			i.t.Logf("Aurora ingesting... %v", root)
			continue
		}

		if uint32(root.CurrentProtocolVersion) == hcnetCoreProtocolVersion {
			i.t.Logf("Aurora protocol version matches %d: %+v",
				root.CurrentProtocolVersion, root)
			return
		}
	}

	i.t.Fatal("Aurora not ingesting...")
}

// UpgradeProtocol arms Core with upgrade and blocks until protocol is upgraded.
func (i *Test) UpgradeProtocol(version uint32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	err := i.coreClient.Upgrade(ctx, int(version))
	cancel()
	if err != nil {
		i.t.Fatalf("could not upgrade protocol: %v", err)
	}

	for t := 0; t < 10; t++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		info, err := i.coreClient.Info(ctx)
		cancel()
		if err != nil {
			i.t.Logf("could not obtain info response: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if info.Info.Ledger.Version == int(version) {
			i.t.Logf("Protocol upgraded to: %d", info.Info.Ledger.Version)
			return
		}
		time.Sleep(time.Second)
	}

	i.t.Fatalf("could not upgrade protocol in 10s")
}

// Cluttering code with if err != nil is absolute nonsense.
func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// findDockerComposePath performs a best-effort attempt to find the project's
// Docker Compose files.
func findDockerComposePath() string {
	// Lets you check if a particular directory contains a file.
	directoryContainsFilename := func(dir string, filename string) bool {
		files, innerErr := ioutil.ReadDir(dir)
		panicIf(innerErr)

		for _, file := range files {
			if file.Name() == filename {
				return true
			}
		}

		return false
	}

	current, err := os.Getwd()
	panicIf(err)

	//
	// We have a primary and backup attempt for finding the necessary docker
	// files: via $GOPATH and via local directory traversal.
	//

	if gopath := os.Getenv("GOPATH"); gopath != "" {
		monorepo := filepath.Join(gopath, "src", "github.com", "hcnet", "go")
		if _, err = os.Stat(monorepo); !os.IsNotExist(err) {
			current = monorepo
		}
	}

	// In either case, we try to walk up the tree until we find "go.mod",
	// which we hope is the root directory of the project.
	for !directoryContainsFilename(current, "go.mod") {
		current, err = filepath.Abs(filepath.Join(current, ".."))

		// FIXME: This only works on *nix-like systems.
		if err != nil || filepath.Base(current)[0] == filepath.Separator {
			fmt.Println("Failed to establish project root directory.")
			panic(err)
		}
	}

	// Directly jump down to the folder that should contain the configs
	return filepath.Join(current, "exp", "services", "soroban-rpc", "internal", "test")
}
