```shell
# 生成私钥
openssl genrsa -out server.key 2048

# 生成证书签名请求 (CSR)
openssl req -new -key server.key -out server.csr

# 自签名证书（有效期 365 天）
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
```