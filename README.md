# Loyalty System

## Swagger

```azure
http://{server_host}:{server_port}/swagger/index.html
```
./gophermarttest-darwin-arm64 \
-test.v -test.run=^TestGophermart$ \
-gophermart-binary-path=cmd/gophermart/gophermart \
-gophermart-host=localhost \
-gophermart-port=8080 \
-gophermart-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable" \
-accrual-binary-path=cmd/accrual/accrual_darwin_arm64 \
-accrual-host=localhost \
-accrual-port=8081 \
-accrual-database-uri="postgresql://postgres:postgres@5432/praktikum?sslmode=disable"
