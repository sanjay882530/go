# Docker build targets use an optional "TAG" environment
# variable can be set to use custom tag name. For example:
#   TAG=my-registry.example.com/keystore:dev make keystore
XDRS = xdr/Hcnet-SCP.x \
xdr/Hcnet-ledger-entries.x \
xdr/Hcnet-ledger.x \
xdr/Hcnet-overlay.x \
xdr/Hcnet-transaction.x \
xdr/Hcnet-types.x

XDRGEN_COMMIT=3f6808cd161d72474ffbe9eedbd7013de7f92748

.PHONY: xdr xdr-clean xdr-update

keystore:
	$(MAKE) -C services/keystore/ docker-build

ticker:
	$(MAKE) -C services/ticker/ docker-build

friendbot:
	$(MAKE) -C services/friendbot/ docker-build

webauth:
	$(MAKE) -C exp/services/webauth/ docker-build

recoverysigner:
	$(MAKE) -C exp/services/recoverysigner/ docker-build

regulated-assets-approval-server:
	$(MAKE) -C services/regulated-assets-approval-server/ docker-build

gxdr/xdr_generated.go: $(XDRS)
	go run github.com/xdrpp/goxdr/cmd/goxdr -p gxdr -enum-comments -o $@ $(XDRS)
	go fmt $@

xdr/%.x:
	curl -Lsf -o $@ https://raw.githubusercontent.com/hcnet/hcnet-core/master/src/protocol-curr/$@

xdr/xdr_generated.go: $(XDRS)
	docker run -it --rm -v $$PWD:/wd -w /wd ruby /bin/bash -c '\
		gem install specific_install -v 0.3.7 && \
		gem specific_install https://github.com/hcnet/xdrgen.git -b $(XDRGEN_COMMIT) && \
		xdrgen \
			--language go \
			--namespace xdr \
			--output xdr/ \
			$(XDRS)'
	go fmt $@

xdr: gxdr/xdr_generated.go xdr/xdr_generated.go

xdr-clean:
	rm xdr/*.x || true

xdr-update: xdr-clean xdr
