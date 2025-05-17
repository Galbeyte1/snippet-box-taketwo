# snippet-box-taketwo

### Generating TLS Certificates for Local Development

To run the project locally with HTTPS, generate a self-signed TLS certificate using Go’s built-in `generate_cert.go` tool:

```bash
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost --cert=tls/cert.pem --key=tls/key.pem