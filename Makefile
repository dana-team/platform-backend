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
CAPP_RELEASE ?= v0.3.2
PREREQ_HELMFILE ?= $(shell pwd)/charts/platform_hub_prereq_helmfile.yaml
CLUSTER_PROXY_RBAC ?= $(shell pwd)/hack/managed-cluster-setup/cluster-proxy-rbac.yaml
CLUSTER_GATEWAY_RBAC ?= $(shell pwd)/hack/managed-cluster-setup/cluster-gateway-rbac.yaml
CNAME_RECORD_CRD ?= https://raw.githubusercontent.com/dana-team/provider-dns/main/package/crds/record.dns.crossplane.io_cnamerecords.yaml

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test -v $$(go list ./... | grep -v /e2e_tests) -coverprofile cover.out

.PHONY: test-e2e
test-e2e: ginkgo
	@test -n "${KUBECONFIG}" -o -r ${HOME}/.kube/config || (echo "Failed to find kubeconfig in ~/.kube/config or no KUBECONFIG set"; exit 1)
	echo "Running e2e tests"
	go clean -testcache
	$(LOCALBIN)/ginkgo -p --vv ./test/e2e_tests/... -coverprofile cover.out -timeout -- -platformURL=$(PLATFORM_URL) -clusterDomain=$(CLUSTER_DOMAIN)

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
	$(HELM) upgrade $(NAME) -n $(NAMESPACE) charts/$(NAME) --install --create-namespace \
		-f charts/$(NAME)/values.yaml \
		--set image.repository=$(IMG_REPO) \
		--set image.tag=$(IMG_TAG) \
		--set image.pullPolicy=$(IMG_PULL_POLICY) \
		--set config.cluster.name=$(CLUSTER_NAME) \
		--set config.cluster.domain=$(CLUSTER_DOMAIN)

.PHONY: undeploy
undeploy: helm ## Uninstall from the K8s cluster specified in ~/.kube/config.
	$(HELM) uninstall $(NAME) -n $(NAMESPACE)
	$(KUBECTL) delete ns $(NAMESPACE)

.PHONY: env-file
env-file: yq
	$(HELM) template -s templates/configmap.yaml charts/$(NAME) \
	--set config.cluster.name=${CLUSTER_NAME} \
	--set config.cluster.domain=${CLUSTER_DOMAIN} > $(ENV_FILE)_tmp
	$(YQ) eval -j $(ENV_FILE)_tmp | jq -r '.data | to_entries | .[] | "\(.key)=\(.value)"' > $(ENV_FILE)
	rm $(ENV_FILE)_tmp

.PHONY: install-prereq
install-prereq: deploy-capp install-cluster-proxy-role install-cluster-gateway-role

.PHONY: uninstall-prereq
uninstall-prereq: undeploy-capp uninstall-cluster-proxy-role uninstall-cluster-gateway-role

.PHONY: install-cluster-proxy-role
install-cluster-proxy-role:
	kubectl apply -f $(CLUSTER_PROXY_RBAC)

.PHONY: uninstall-cluster-proxy-role
uninstall-cluster-proxy-role:
	kubectl delete -f $(CLUSTER_PROXY_RBAC)

.PHONY: install-cluster-gateway-role
install-cluster-gateway-role:
	kubectl apply -f $(CLUSTER_GATEWAY_RBAC)

.PHONY: uninstall-cluster-gateway-role
uninstall-cluster-gateway-role:
	kubectl delete -f $(CLUSTER_GATEWAY_RBAC)

.PHONY: deploy-capp
deploy-capp: helm helm-plugins
	[ -d "container-app-operator" ] || git clone $(CAPP_REPO)

	make -C container-app-operator prereq-openshift \
		PROVIDER_DNS_REALM=${PROVIDER_DNS_REALM} \
		PROVIDER_DNS_KDC=${PROVIDER_DNS_KDC} \
		PROVIDER_DNS_POLICY=${PROVIDER_DNS_POLICY} \
		PROVIDER_DNS_NAMESERVER=${PROVIDER_DNS_NAMESERVER} \
		PROVIDER_DNS_USERNAME=${PROVIDER_DNS_USERNAME} \
		PROVIDER_DNS_PASSWORD=${PROVIDER_DNS_PASSWORD}

	$(HELM) upgrade --install container-app-operator container-app-operator/charts/container-app-operator \
      --wait --create-namespace --namespace capp-operator-system \
      --set image.manager.tag=$(CAPP_RELEASE) --set image.manager.pullPolicy=$(IMG_PULL_POLICY)

	rm -rf container-app-operator/

.PHONY: undeploy-capp
undeploy-capp: helm helm-plugins
	[ -d "container-app-operator" ] || git clone $(CAPP_REPO)
	make -C container-app-operator uninstall-prereq-openshift
	$(HELM) uninstall container-app-operator --namespace capp-operator-system
	rm -rf container-app-operator/

.PHONY: install-cnamerecord-crd
install-cnamerecord-crd:
	kubectl apply -f $(CNAME_RECORD_CRD)

.PHONY: uninstall-cnamerecord-crd
uninstall-cnamerecord-crd:
	kubectl delete -f $(CNAME_RECORD_CRD)

.PHONY: setup-hub
setup-hub: helmfile install-capp-crds clusteradm ## Setup hub cluster with RCS
	kubectl get namespace ${PLACEMENTS_NAMESPACE} || kubectl create namespace ${PLACEMENTS_NAMESPACE}
	$(CLUSTERADM) create clusterset ${MANAGED_CLUSTER_NAME}
	$(CLUSTERADM) clusterset set ${MANAGED_CLUSTER_NAME} --clusters ${MANAGED_CLUSTER_NAME}
	$(CLUSTERADM) clusterset bind ${MANAGED_CLUSTER_NAME} --namespace ${PLACEMENTS_NAMESPACE}
	$(CLUSTERADM) create placement ${PLACEMENT_NAME} --namespace ${PLACEMENTS_NAMESPACE} --clustersets=${MANAGED_CLUSTER_NAME}

	$(HELMFILE) apply -f $(PREREQ_HELMFILE) \
	--state-values-set placementName=${PLACEMENT_NAME} \
	--state-values-set placementsNamespace=${PLACEMENTS_NAMESPACE}
	rm -rf container-app-operator/
	$(MAKE) install-cnamerecord-crd

.PHONY: cleanup-hub
cleanup-hub: helmfile  ## cleanup hub cluster.
	$(HELMFILE) -f $(PREREQ_HELMFILE) destroy
	kubectl delete placements ${PLACEMENT_NAME} --namespace ${PLACEMENTS_NAMESPACE} --ignore-not-found
	$(CLUSTERADM) clusterset unbind ${MANAGED_CLUSTER_NAME} --namespace ${PLACEMENTS_NAMESPACE}
	$(CLUSTERADM) delete clusterset ${MANAGED_CLUSTER_NAME}
	kubectl delete ns ${PLACEMENTS_NAMESPACE} --ignore-not-found
	$(MAKE) uninstall-cnamerecord-crd

.PHONY: doc-chart
doc-chart: helm-docs helm
	$(HELM_DOCS) charts/

.PHONY: install-capp-crds
install-capp-crds:
	[ -d "container-app-operator" ] || git clone $(CAPP_REPO)
	make -C container-app-operator install
	rm -rf container-app-operator/

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
HELMFILE ?= $(LOCALBIN)/helmfile-$(HELMFILE_VERSION)
CLUSTERADM ?= $(LOCALBIN)/clusteradm

HELM_URL ?= https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
HELMFILE_URL ?= https://github.com/helmfile/helmfile/releases/download/v${HELMFILE_VERSION}/helmfile_${HELMFILE_VERSION}_linux_amd64.tar.gz
CLUSTERADM_URL ?= https://raw.githubusercontent.com/open-cluster-management-io/clusteradm/main/install.sh

## Tool Versions
KUSTOMIZE_VERSION ?= v5.3.0
GOLANGCI_LINT_VERSION ?= v1.60.3
HELM_DOCS_VERSION ?= v1.14.2
HELMFILE_VERSION ?= 0.167.1

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

.PHONY: helmfile
helmfile: $(HELMFILE) ## Install helmfile on the local machine
$(HELMFILE): $(LOCALBIN)
	wget -O $(LOCALBIN)/helmfile.tar.gz $(HELMFILE_URL)
	tar -xzvf $(LOCALBIN)/helmfile.tar.gz -C $(LOCALBIN)
	rm $(LOCALBIN)/helmfile.tar.gz $(LOCALBIN)/*.md $(LOCALBIN)/LICENSE
	mv $(LOCALBIN)/helmfile $(LOCALBIN)/helmfile-$(HELMFILE_VERSION)

.PHONY: clusteradm
clusteradm: $(CLUSTERADM) ## Download clusteradm locally if necessary.
$(CLUSTERADM): $(LOCALBIN)
	test -s $(LOCALBIN)/clusteradm || curl -L $(CLUSTERADM_URL) | sed 's|/usr/local/bin|$(LOCALBIN)|g' | bash

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