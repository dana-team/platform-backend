#!/bin/bash

HTFILE=e2e-tests-htpass-secret
USERNAME=e2e-user
PASSWORD=e2e-password

###

# Function to generate bcrypt hash
generate_bcrypt_hash() {
    local password=$1
    openssl passwd -apr1 "$password"
}

# Generate bcrypt hash of the password
HASH=$(generate_bcrypt_hash "$PASSWORD")

# Create or update the .htpasswd file
if [ -f "$HTFILE" ]; then
    # Update existing file: remove existing entry for the username, if any
    grep -v "^$USERNAME:" "$HTFILE" > "$HTFILE.tmp" && mv "$HTFILE.tmp" "$HTFILE"
fi

# Append the new username:hashed_password entry
echo "$USERNAME:$HASH" >> "$HTFILE"