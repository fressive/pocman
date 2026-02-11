#!/bin/bash

usage() {
    echo "usage: $0 <agent_name>"
    exit 1
}

if [ $# -lt 1 ]; then
    usage
fi

CERT_DIR="./cert"
mkdir -p $CERT_DIR

AGENT_NAME=$1

openssl genrsa -out $CERT_DIR/agent_$AGENT_NAME.key 2048
openssl req -new -key $CERT_DIR/agent_$AGENT_NAME.key -out $CERT_DIR/agent_$AGENT_NAME.csr -subj "/CN=pocman-agent-$AGENT_NAME"
openssl x509 -req -in $CERT_DIR/agent_$AGENT_NAME.csr -CA $CERT_DIR/ca.pem -CAkey $CERT_DIR/ca.key -CAcreateserial -out $CERT_DIR/agent_$AGENT_NAME.pem -days 3650