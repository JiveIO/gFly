# Folder resources/tls/

https://gist.github.com/ethicka/27c36c975a5c2cbbd1874bc78bab61c4

## Step 1. Configuration

```bash
vi localhost.conf

[req]
default_bits = 1024
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
```

## Step 2. Generate file

```bash
openssl genrsa -out localhost.key 2048
openssl rsa -in localhost.key -out localhost.key.rsa
openssl req -new -key localhost.key.rsa -subj /CN=localhost -out localhost.csr -config localhost.conf
openssl x509 -req -extensions v3_req -days 3650 -in localhost.csr -signkey localhost.key.rsa -out localhost.crt -extfile localhost.conf
```

## Step 3. Import to Mac

```bash
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain localhost.crt
```

## Step 4. Setup

```.env
SERVER_TLS_CERT=./resources/tls/localhost.crt
SERVER_TLS_KEY=./resources/tls/localhost.key
```

## Step 5. Verify

On Chrome
    
    Open dev mode and check tab Security.

One Safari

    Click on SSL icon on HTTP bar of website.