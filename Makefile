# Image URL to use all building/pushing image targets
IMG ?= backend:$(IMG_TAG)
IMG_REPO ?= ghcr.io/dana-team/$(NAME)
IMG_TAG ?= main
IMG_PULL_POLICY ?= Always
NAME ?= platform-backend
NAMESPACE ?= platform-backend-system

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

PLATFORM_URL ?=
ENV_FILE ?= .env
CAPP_REPO ?= https://github.com/dana-team/container-app-operator
CAPP_RELEASE ?= main

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
	$(LOCALBIN)/ginkgo -p --vv ./test/e2e_tests/... -coverprofile cover.out -timeout -- -platformUrl=$(PLATFORM_URL)

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
deploy: helm ## Deploy to the K8s cluster specified in ~/.kube/config.
	REQUIRED_VARS="\
		CLUSTER_NAME \
		CLUSTER_DOMAIN \
	"
	for var in $${REQUIRED_VARS}; do \
		if [ -z "$${!var}" ]; then \
			echo "Error: Variable $$var is not set."; \
			exit 1; \
		fi; \
	done

	$(HELM) upgrade $(NAME) -n $(NAMESPACE) charts/$(NAME) --install --create-namespace \
		-f charts/$(NAME)/values.yaml \
		--set image.repository=$(IMG_REPO) \
		--set image.tag=$(IMG_TAG) \
		--set image.pullPolicy=$(IMG_PULL_POLICY) \
		--set config.cluster.name=$(CLUSTER_NAME) \
		--set config.cluster.domain=$(CLUSTER_DOMAIN)

.PHONY: undeploy
undeploy: helm ## Deploy to the K8s cluster specified in ~/.kube/config.
	$(HELM) uninstall $(NAME) -n $(NAMESPACE)
	$(KUBECTL) delete ns $(NAMESPACE)

.PHONY: env-file
env-file:
	REQUIRED_VARS="\
		CLUSTER_NAME \
		CLUSTER_DOMAIN \
	"
	for var in $${REQUIRED_VARS}; do \
		if [ -z "$${!var}" ]; then \
			echo "Error: Variable $$var is not set."; \
			exit 1; \
		fi; \
	done

	$(HELM) template -s templates/configmap.yaml charts/$(NAME) \
	--set config.cluster.name=${CLUSTER_NAME} \
	--set config.cluster.domain=${CLUSTER_DOMAIN} > $(ENV_FILE)_tmp
	$(YQ) eval -j $(ENV_FILE)_tmp | jq -r '.data | to_entries | .[] | "\(.key)=\(.value)"' > $(ENV_FILE)
	rm $(ENV_FILE)_tmp

.PHONY: install-prereq
install-prereq: deploy-capp

.PHONY: uninstall-prereq
uninstall-prereq: undeploy-capp

.PHONY: deploy-capp
deploy-capp: helm helm-plugins
	REQUIRED_VARS="\
		PROVIDER_DNS_REALM \
		PROVIDER_DNS_KDC \
		PROVIDER_DNS_POLICY \
		PROVIDER_DNS_NAMESERVER \
		PROVIDER_DNS_USERNAME \
		PROVIDER_DNS_PASSWORD \
	"
	for var in $${REQUIRED_VARS}; do \
		if [ -z "$${!var}" ]; then \
			echo "Error: Variable $$var is not set."; \
			exit 1; \
		fi; \
	done

	[ -d "container-app-operator" ] || git clone $(CAPP_REPO)

	make -C container-app-operator prereq-openshift \
		PROVIDER_DNS_REALM=${PROVIDER_DNS_REALM} \
		PROVIDER_DNS_KDC=${PROVIDER_DNS_KDC} \
		PROVIDER_DNS_POLICY=${PROVIDER_DNS_POLICY} \
		PROVIDER_DNS_NAMESERVER=${PROVIDER_DNS_NAMESERVER} \
		PROVIDER_DNS_USERNAME=${PROVIDER_DNS_USERNAME} \
		PROVIDER_DNS_PASSWORD=${PROVIDER_DNS_PASSWORD}

	$(HELM) upgrade --install capp-operator container-app-operator/charts/container-app-operator \
      --wait --create-namespace --namespace capp-operator-system \
      --set image.manager.tag=$(CAPP_RELEASE)

	rm -rf container-app-operator/

.PHONY: undeploy-capp
undeploy-capp: helm helm-plugins
	[ -d "container-app-operator" ] || git clone $(CAPP_REPO)
	make -C container-app-operator uninstall-prereq-openshift
	$(HELM) uninstall capp-operator --namespace capp-operator-system
	rm -rf container-app-operator/

.PHONY: doc-chart
doc-chart: helm-docs helm
	$(HELM_DOCS) charts/

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
YQ ?= $(LOCALBIN)/yq
HELM_DOCS ?= $(LOCALBIN)/helm-docs-$(HELM_DOCS_VERSION)
HELM_URL ?= https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3

## Tool Versions
KUSTOMIZE_VERSION ?= v5.3.0
GOLANGCI_LINT_VERSION ?= v1.60.1
HELM_DOCS_VERSION ?= v1.14.2

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

.PHONY: helm
helm: $(HELM) ## Install helm on the local machine
$(HELM): $(LOCALBIN)
	wget -O $(LOCALBIN)/get-helm.sh $(HELM_URL)
	chmod 700 $(LOCALBIN)/get-helm.sh
	HELM_INSTALL_DIR=$(LOCALBIN) $(LOCALBIN)/get-helm.sh

.PHONY: helm-plugins
helm-plugins: ## Install helm plugins on the local machine
	@if ! $(HELM) plugin list | grep -q 'diff'; then \
		$(HELM) plugin install https://github.com/databus23/helm-diff; \
	fi
	@if ! $(HELM) plugin list | grep -q 'git'; then \
		$(HELM) plugin install https://github.com/aslafy-z/helm-git; \
	fi
	@if ! $(HELM) plugin list | grep -q 's3'; then \
		$(HELM) plugin install https://github.com/hypnoglow/helm-s3; \
	fi
	@if ! $(HELM) plugin list | grep -q 'secrets'; then \
		$(HELM) plugin install https://github.com/jkroepke/helm-secrets; \
	fi

.PHONY: helm-docs
helm-docs: $(HELM_DOCS)
$(HELM_DOCS): $(LOCALBIN)
	$(call go-install-tool,$(HELM_DOCS),github.com/norwoodj/helm-docs/cmd/helm-docs,$(HELM_DOCS_VERSION))

.PHONY: yq
yq: $(YQ) ## Download yq locally if necessary.
$(YQ): $(LOCALBIN)
	test -s $(LOCALBIN)/yq || GOBIN=$(LOCALBIN) go install github.com/mikefarah/yq/v4@latest

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