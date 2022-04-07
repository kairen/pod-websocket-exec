ORG := github.com
OWNER := kubedev
REPOPATH ?= $(ORG)/$(OWNER)/pod-websocket-exec

GOOS ?= $(shell go env GOOS)

.PHONY: build
build: k8s-ws-exec

.PHONY: k8s-ws-exec
k8s-ws-exec: depend
	GOOS=$(GOOS) go build -a -o $@ .

.PHONY: depend
depend:
	@dep ensure

.PHONY: clean
clean:
	@rm -rf k8s-ws-exec
