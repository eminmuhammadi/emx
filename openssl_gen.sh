#!/usr/bin/env bash

# exit when any command fails
set -ex

CERT_FOLDER=_certs

# Generate X509 certificate and private key without prompting for any input (using extension file)
rm -rf $CERT_FOLDER && mkdir -p $CERT_FOLDER

# Check if openssl is installed
if ! command -v openssl &>/dev/null; then
    echo "openssl could not be found, please install it first."
    exit
fi

openssl genrsa -aes256 -passout pass:1 -out $CERT_FOLDER/ca.key 4096
openssl rsa -passin pass:1 -in $CERT_FOLDER/ca.key -out $CERT_FOLDER/ca.key.tmp
mv $CERT_FOLDER/ca.key.tmp $CERT_FOLDER/ca.key

# Create server-config
cat >$CERT_FOLDER/ca.cnf <<EOF
[ ca ]
default_ca	= CA_default

[ CA_default ]
default_md	= sha256

[ v3_ca ]
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid:always,issuer
basicConstraints = critical,CA:true

[ req ]
distinguished_name	= req_distinguished_name
prompt = no

[ req_distinguished_name ]
countryName			   = IL
stateOrProvinceName	   = Tel Aviv
localityName           = Center
0.organizationName     = EMX Certificate Authority
organizationalUnitName = Digital Lab
commonName			   = EMX Certificate Authority Root
emailAddress           = ca@emxlab
EOF

# Create CA certificate
openssl req -config $CERT_FOLDER/ca.cnf - -key $CERT_FOLDER/ca.key -new -x509 -days 7300 -sha256 -extensions v3_ca -out $CERT_FOLDER/ca.crt


# Create server-config
cat >$CERT_FOLDER/server.cnf <<EOF
[ req ]
default_bits       = 4096
default_keyfile    = server.key
distinguished_name = req_distinguished_name
req_extensions     = req_ext
x509_extensions    = v3_ca
string_mask        = utf8only
prompt             = no

[ req_distinguished_name ]
countryName			   = IL
stateOrProvinceName	   = Tel Aviv
localityName           = Center
0.organizationName     = EMX Certificate Authority
organizationalUnitName = Digital Lab
commonName			   = EMX Lab Server
emailAddress           = server@emxlab

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost

[ v3_ca ]
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid:always,issuer
basicConstraints = critical,CA:false
keyUsage = critical,digitalSignature,keyEncipherment
extendedKeyUsage = serverAuth,clientAuth
subjectAltName = @alt_names
EOF

# Generate server certificate
openssl genrsa -out $CERT_FOLDER/server.key 4096
openssl req -new -key $CERT_FOLDER/server.key -out $CERT_FOLDER/server.csr -config $CERT_FOLDER/server.cnf
openssl x509 -req -in $CERT_FOLDER/server.csr -CA $CERT_FOLDER/ca.crt -CAkey $CERT_FOLDER/ca.key -CAcreateserial -out $CERT_FOLDER/server.crt -days 365 -extensions req_ext -extfile $CERT_FOLDER/server.cnf

# Create client-config
cat >$CERT_FOLDER/client.cnf <<EOF
[ req ]
default_bits       = 4096
default_keyfile    = client.key
distinguished_name = req_distinguished_name
req_extensions     = req_ext
x509_extensions    = v3_ca
string_mask        = utf8only
prompt             = no

[ req_distinguished_name ]
countryName			   = IL
stateOrProvinceName	   = Tel Aviv
localityName           = Center
0.organizationName     = EMX Certificate Authority
organizationalUnitName = Digital Lab
commonName			   = EMX Lab Client
emailAddress           = client@emxlab

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost

[ v3_ca ]
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid:always,issuer
basicConstraints = critical,CA:false
keyUsage = critical,digitalSignature,keyEncipherment
extendedKeyUsage = serverAuth,clientAuth
subjectAltName = @alt_names
EOF

# Generate client certificate
openssl genrsa -out $CERT_FOLDER/client.key 4096
openssl req -new -key $CERT_FOLDER/client.key -out $CERT_FOLDER/client.csr -config $CERT_FOLDER/client.cnf
openssl x509 -req -in $CERT_FOLDER/client.csr -CA $CERT_FOLDER/ca.crt -CAkey $CERT_FOLDER/ca.key -CAcreateserial -out $CERT_FOLDER/client.crt -days 365 -extensions req_ext -extfile $CERT_FOLDER/client.cnf

# Verify certificates
openssl verify -CAfile $CERT_FOLDER/ca.crt $CERT_FOLDER/server.crt
openssl verify -CAfile $CERT_FOLDER/ca.crt $CERT_FOLDER/client.crt