#!/bin/bash

USERNAME=e2e-user
PASSWORD=e2e-password

generate_bcrypt_hash() {
    local password=$1
    openssl passwd -apr1 "$password"
}

HASH=$(generate_bcrypt_hash "$PASSWORD")
echo "$USERNAME:$HASH"