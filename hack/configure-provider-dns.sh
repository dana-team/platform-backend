#!/bin/bash

# Initialize variables with defaults
domaincontrollerip=""
domaincontrollerpassword=""
domaincontrollerserver=""
domaincontrollerrealm=""
domaincontrolleruser=""

# Functions to initialize variables
initialize_variable() {
    local var_name="$1"
    local input_value="$2"
    local default_value="$3"

    if [ -n "$input_value" ]; then
        eval "$var_name=\"$input_value\""
    else
        eval "$var_name=\"$default_value\""
    fi
}

# Initialize variables based on input
initialize_variable domaincontrollerip "$1" "127.0.0.1"
initialize_variable domaincontrollerpassword "$2" "passw0rd"
initialize_variable domaincontrollerserver "$3" "dana-wdc-1.dana-dev.com"
initialize_variable domaincontrollerrealm "$4" "DANA-DEV.COM"
initialize_variable domaincontrolleruser "$5" "dana"

# Convert realm to lowercase
domaincontrollerrealm_lowercase=$(echo "$domaincontrollerrealm" | tr '[:upper:]' '[:lower:]')

# Create krb5.conf file
krb5File=krb5.conf

cat <<EOL > $krb5File
includedir /etc/krb5.conf.d/

[logging]
    default = FILE:/var/log/krb5libs.log
    kdc = FILE:/var/log/krb5kdc.log
    admin_server = FILE:/var/log/kadmind.log

[libdefaults]
    dns_lookup_realm = false
    ticket_lifetime = 24h
    renew_lifetime = 7d
    forwardable = true
    rdns = false
    pkinit_anchors = FILE:/etc/pki/tls/certs/ca-bundle.crt
    spake_preauth_groups = edwards25519
    default_realm = $domaincontrollerrealm
    default_ccache_name = KEYRING:persistent:%{uid}

[realms]
 $domaincontrollerrealm = {
     kdc = $domaincontrollerserver
     admin_server = $domaincontrollerserver
     default_domain = $domaincontrollerrealm_lowercase
 }

[domain_realm]
 .$domaincontrollerrealm_lowercase = $domaincontrollerrealm
 $domaincontrollerrealm_lowercase = $domaincontrollerrealm
EOL

# Create ConfigMap if it doesn't exist
kubectl get configmap krb5-config -n crossplane-system &> /dev/null
if [ $? -ne 0 ]; then
    kubectl create configmap krb5-config --from-file=krb5.conf=krb5.conf -n crossplane-system
else
    kubectl create configmap krb5-config --from-file=krb5.conf=krb5.conf -n crossplane-system -o yaml --dry-run=client | kubectl replace -f -
fi

rm $krb5File

# Create DeploymentRuntimeConfig object
cat <<EOF | kubectl apply -f -
apiVersion: pkg.crossplane.io/v1beta1
kind: DeploymentRuntimeConfig
metadata:
  name: dns-config
spec:
  deploymentTemplate:
    spec:
      selector:
        matchLabels:
          pkg.crossplane.io/provider: provider-dns
      template:
        spec:
          containers:
          - args:
            - --debug
            name: package-runtime
            volumeMounts:
            - mountPath: /etc/krb5.conf
              name: krb5-config
              readOnly: true
              subPath: krb5.conf
          volumes:
          - configMap:
              name: krb5-config
            name: krb5-config
          dnsConfig:
            nameservers:
              - $domaincontrollerip
          dnsPolicy: None
EOF

# Create Secret and ProviderConfig objects
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: dns-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    {
      "rfc": "3645",
      "server": "$domaincontrollerserver",
      "realm": "$domaincontrollerrealm",
      "username": "$domaincontrolleruser",
      "password": "$domaincontrollerpassword"
    }
EOF

cat <<EOF | kubectl apply -f -
apiVersion: dns.dns.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: dns-default
spec:
  credentials:
    source: Secret
    secretRef:
      name: dns-creds
      namespace: crossplane-system
      key: credentials
EOF

# Patch the Provider to use the new runtime config
kubectl patch Provider dana-team-provider-dns --type='json' -p='[{"op": "replace", "path": "/spec/runtimeConfigRef/name", "value": "dns-config"}]'