# emx

Real-time proxy interceptor which supports application layer and database.

![emx demo](./_img.gif)

```sql
SELECT req.method,
       res.status_code,
       req.host,
       req.request_uri,
       req.body,
       res.body
    FROM requests req,
         responses res
    WHERE req.session_id = res.session_id
    ORDER BY req.created_at DESC;
```

- GET /api/v1/log
- GET /api/v1/log/:id
- GET /api/v1/request/:id
- GET /api/v1/response/:id

## Installation

Required binaries

- `go` https://go.dev/doc/install
- `openssl` https://wiki.openssl.org/index.php/Binaries

```sh
go build -o ./emx . && mv emx /usr/local/bin
```

## Usage

To start emx you need to type following on terminal:

- Set configuration for emx

```sh
export PROXY_HOST=0.0.0.0
export PROXY_PORT=8080
export APP_HOST=127.0.0.1
export APP_PORT=8888
export PROXY_DECRYPT_CERT_FILE=_certs/ca.crt
export PROXY_DECRYPT_KEY_FILE=_certs/ca.key
export TLS_MODE=off
export SQLITE_DSN=:memory:?cache=shared
```

`TLS_MODE` is required only for application. Proxy server will decrypt https using `PROXY_DECRYPT_CERT_FILE` and `PROXY_DECRYPT_KEY_FILE` files.

- Run

```
emx
```


or 

```
chmod +x ./bin && ./bin
```

Visit http://127.0.0.1:8888/ui to see proxy logs.

### Environment

List of environment variables:

- TLS_MODE (`off`, `tls`, `mutual_tls`) - required
- TLS_CERT_FILE (`path_to_file`) - optional for `off`
- TLS_KEY_FILE (`path_to_file`) - optional for `off`
- TLS_CA_FILE (`path_to_file`) - optional for `off`
- PROXY_HOST (`0.0.0.0`) - required
- PROXY_PORT (`8080`) - required
- PROXY_DECRYPT_CERT_FILE (`path_to_file`) - required
- PROXY_DECRYPT_KEY_FILE (`path_to_file`) - required
- PROXY_VERBOSE (`bool`) - optional (default: `false`)
- APP_HOST (`0.0.0.0`) - required
- APP_PORT (`8443`) - required
- SQL_VERBOSE  (`bool`) - optional (default: `false`)
- SQLITE_DSN (`file:emxdb.sqlite?cache=shared`) - optional (default: `:memory:?cache=shared`)
- MOCK_FILE (`path_to_file`) optional (example: `mock.yaml`)

## Mocking

You can easily mock api response using following like configuration file:

mock.yaml
```yaml
patterns:
  - method: "GET"
    host: "example.com"
    path: "/"
    response:
      status_code: 200
      headers: |
        Content-Type: application/json
      body: |
        {
          "message": "Hello, World!"
        }

  - method: "GET"
    host: "www.google.com"
    path: "/"
    response:
      status_code: 200
      headers: |
        Content-Type: application/json
      body: |
        {
          "message": "Hello, Google!"
        }
```

```sh
$ curl --proxy 127.0.0.1:8080 --ssl-no-revoke https://example.com
{
  "message": "Hello, World!"
}
```

```sh
$ curl --proxy 127.0.0.1:8080 --ssl-no-revoke https://www.google.com
{
  "message": "Hello, Google!"
}
```

## Intercepting HTTPS

To intercept HTTPS requests you need to trust CA cert file. If you did not generate CA cert file, you can run following command on terminal:

```sh
chmod +x ./openssl_gen.sh && ./openssl_gen.sh
```

Distribute generated `_certs/ca.crt` file to your device, and trust this certificate.

How to trust certificate:
- Windows  https://learn.microsoft.com/en-us/windows-hardware/drivers/install/trusted-root-certification-authorities-certificate-store
- MacOS https://support.apple.com/en-gb/guide/keychain-access/kyca11871/mac
- Iphone https://support.apple.com/en-us/102390
- Android https://developer.android.com/privacy-and-security/security-ssl#Pinning
