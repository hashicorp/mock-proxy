#! /usr/bin/env bash
set -euo pipefail

# Always work from the root of the repo.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR"/.. && pwd)"

# Operate in the certs directory.
cd "$ROOT_DIR/certs"

if ! command -v openssl > /dev/null 2>&1; then
  echo "This script requires openssl to run correctly"
  exit 1
fi

if ! grep -q "v3_ca" /etc/ssl/openssl.cnf; then
  echo "You are missing a section in /etc/ssl/openssl.cnf"
  echo "To correct this error, add the following section to your openssl.cnf config"
  echo ""
  echo "[ v3_ca ]"
  echo "basicConstraints = critical,CA:TRUE"
  echo "subjectKeyIdentifier = hash"
  echo "authorityKeyIdentifier = keyid:always,issuer:always"
  echo ""
  echo "Then rerun this command."
  exit 1
fi

echo "Okay, here we go, let's generate a certificate"
echo ""
openssl req -new \
  -newkey rsa:2048 \
  -sha256 \
  -days 365 \
  -nodes \
  -x509 \
  -extensions v3_ca \
  -keyout ca.pem \
  -out ca.pem \
  -subj "/C=US/ST=California/L=San Francisco/O=HashiCorp/OU=Engineering Services/CN=example.com"
echo ""

echo "And the public key thereof"
echo ""
openssl x509 -in ca.pem -outform DER -out ca.crt
echo ""

echo "Okay! That looks good, good certificating."
