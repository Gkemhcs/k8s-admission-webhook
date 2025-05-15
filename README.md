# k8s-admission-webhook
Custom  K8s Admission Webhook
### Steps to Setup
- **Create SSL certs using openssl**
```bash
CERTS_DIRECTORY="certs"
mkdir -p $CERTS_DIRECTORY
cd $CERTS_DIRECTORY
#create root ca  private key
openssl genrsa -out rootCA.key 4096
#create root ca cert
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 3650 -out rootCA.crt -subj "/CN=Root-CA"
#create a webhook server private key
openssl genrsa -out webhook.key 4096
# create a webhook server csr 
openssl req -new -key webhook.key -out webhook.csr -subj "/CN=webhook-service.webhook-ns.svc"
cat <<EOF> webhook.ext
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = webhook-service.webhook-ns.svc
DNS.2 = webhook-service.webhook-ns.svc.cluster.local
EOF

#sign the csr using root ca private key
openssl x509 -req -in webhook.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out webhook.crt -days 365 -sha256 -extfile webhook.ext

```
- **Replace the caBundle with current root ca cert* in validating and mutating webhook configuration*
```bash
export BASE64_ENCODED_CERT=$(cat $CERTS_DIRECTORY/rootCA.crt|base64 -w 0)
sed -i  "s/CA_BUNDLE/$BASE64_ENCODED_CERT/" validation-webhook-configuration.yaml

```

- deploy the webhook server onto cluster

```bash
kubectl create ns webhook-ns
kubectl apply -k .
```


