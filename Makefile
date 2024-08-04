# Image URL to use all building/pushing image targets
IMG ?= backend:latest
NAME ?= platform-backend
NAMESPACE ?= platform-backend-system

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker
platformUrl ?=

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test:  fmt vet ## Run tests.
	go test -v $$(go list ./... | grep -v /e2e_tests) -coverprofile cover.out

.PHONY: test-e2e
test-e2e: ginkgo
	@test -n "${KUBECONFIG}" -o -r ${HOME}/.kube/config || (echo "Failed to find kubeconfig in ~/.kube/config or no KUBECONFIG set"; exit 1)
	echo "Running e2e tests"
	go clean -testcache
	$(LOCALBIN)/ginkgo -p --vv ./test/e2e_tests/... -coverprofile cover.out -timeout -- -platformUrl=$(platformUrl)

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter & yamllint
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix

.PHONY: build
build: fmt vet ## Build binary.
	go build -o bin/backend main.go

.PHONY: run
run: fmt vet ## Run on your host.
	go run ./cmd/main.go

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	$(CONTAINER_TOOL) build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) push ${IMG}

.PHONY: deploy
deploy: install-helm ## Deploy to the K8s cluster specified in ~/.kube/config.
	cd chart && $(HELM) upgrade $(NAME) -n $(NAMESPACE) . --install --create-namespace \
	-f values.yaml \
	--set image=$(IMG) \
	--set clusterName=$(CLUSTER_NAME) \
	--set clusterDomain=$(CLUSTER_DOMAIN)

.PHONY: undeploy
undeploy: install-helm ## Deploy to the K8s cluster specified in ~/.kube/config.
	helm uninstall $(NAME) -n $(NAMESPACE)
	$(KUBECTL) delete ns $(NAMESPACE)

.PHONY: env
env:
	bash hack/create-env.sh ${CLUSTER_NAME} ${CLUSTER_DOMAIN} .env

.PHONY: deploy-capp
deploy-capp: ## Run the deploy-capp script
	$(shell pwd)/hack/deploy-capp.sh $(CAPP_RELEASE)

.PHONY: undeploy-capp
undeploy-capp: ## Run the uninstall-prereq script
	$(shell pwd)/hack/undeploy-capp.sh $(CAPP_RELEASE)

##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUBECTL ?= kubectl
HELM ?= $(LOCALBIN)/helm
KUSTOMIZE ?= $(LOCALBIN)/kustomize-$(KUSTOMIZE_VERSION)
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
GINKGO ?= $(LOCALBIN)/ginkgo
HELM_URL ?= https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3

## Tool Versions
KUSTOMIZE_VERSION ?= v5.3.0
GOLANGCI_LINT_VERSION ?= v1.54.2

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5,$(KUSTOMIZE_VERSION))

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,${GOLANGCI_LINT_VERSION})

.PHONY: ginkgo
ginkgo: $(GINKGO) ## Download ginkgo locally if necessary.
$(GINKGO): $(LOCALBIN)
	test -s $(LOCALBIN)/ginkgo || GOBIN=$(LOCALBIN) go install github.com/onsi/ginkgo/v2/ginkgo@latest

.PHONY: install-helm
install-helm: $(HELM) ## Install helm on the local machine
$(HELM): $(LOCALBIN)
	wget -O $(LOCALBIN)/get-helm.sh $(HELM_URL)
	chmod 700 $(LOCALBIN)/get-helm.sh
	HELM_INSTALL_DIR=$(LOCALBIN) $(LOCALBIN)/get-helm.sh

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary (ideally with version)
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f $(1) ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv "$$(echo "$(1)" | sed "s/-$(3)$$//")" $(1) ;\
}
endef
