#!/bin/bash

## Initialize variables
capprelease=""

initialize_capp_release() {
    if test -s "$1"; then
        capprelease="$1"
    else
        capprelease="main"
    fi
}

initialize_capp_release "$1"

cappimage="ghcr.io/dana-team/container-app-operator:${capprelease}"
capprepo=https://github.com/dana-team/container-app-operator
branch=main

git clone ${capprepo} -b ${branch}
make -C container-app-operator prereq-openshift
make -C container-app-operator deploy IMG="${cappimage}"
kubectl wait --for=condition=ready pods -l control-plane=controller-manager -n capp-operator-system

kubectl get configmap dns-zone -n capp-operator-system &> /dev/null
if [ $? -ne 0 ]; then
    kubectl create configmap dns-zone --from-literal=zone=capp-zone. -n capp-operator-system
fi

rm -rf container-app-operator/