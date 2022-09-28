package actions

import "github.com/hcnet/go/services/aurora/internal/corestate"

type CoreStateGetter interface {
	GetCoreState() corestate.State
}
