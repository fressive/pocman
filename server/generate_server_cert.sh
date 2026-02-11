#!/bin/bash

if [ $# -eq 0 ]; then
    echo "usage: $0 <Server IP or Domain 1> [Server IP or Domain 2] ..."
    echo "example: $0 127.0.0.1 localhost 192.168.1.50 my-server.com"
    exit 1
fi

CERT_DIR="./cert"
mkdir -p $CERT_DIR

ALT_NAMES=""
IP_COUNT=1
DNS_COUNT=1

for arg in "$@"; do
    if [[ $arg =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        ALT_NAMES="${ALT_NAMES}IP.${IP_COUNT} = ${arg}\n"
        ((IP_COUNT++))
    else
        ALT_NAMES="${ALT_NAMES}DNS.${DNS_COUNT} = ${arg}\n"
        ((DNS_COUNT++))
    fi
done

cat > $CERT_DIR/openssl.cnf <<EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
x509_extensions = v3_ca
prompt = no

[req_distinguished_name]
CN = gRPC-Service

[v3_req]
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth, clientAuth
subjectAltName = @alt_names

[v3_ca]
basicConstraints = critical, CA:true
keyUsage = critical, digitalSignature, cRLSign, keyCertSign

[alt_names]
$(echo -e "$ALT_NAMES")
EOF

openssl genrsa -out $CERT_DIR/ca.key 2048
openssl req -x509 -new -nodes -key $CERT_DIR/ca.key -sha256 -days 3650 \
    -out $CERT_DIR/ca.pem -config $CERT_DIR/openssl.cnf -extensions v3_ca -subj "/CN=MyRootCA"
openssl genrsa -out $CERT_DIR/server.key 2048
openssl req -new -key $CERT_DIR/server.key -out $CERT_DIR/server.csr -config $CERT_DIR/openssl.cnf
openssl x509 -req -in $CERT_DIR/server.csr -CA $CERT_DIR/ca.pem -CAkey $CERT_DIR/ca.key -CAcreateserial \
    -out $CERT_DIR/server.pem -days 3650 -sha256 -extfile $CERT_DIR/openssl.cnf -extensions v3_req