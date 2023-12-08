# Middlewares

## CORS
Support 6 access controls `Access-Control-Allow-Origin`, `Access-Control-Allow-Headers`, `Access-Control-Allow-Methods`, `Access-Control-Allow-Credentials`, `Access-Control-Expose-Headers`, `Access-Control-Max-Age`. 

### How to use
```go
// Add global middlewares
app.Use(middleware.CORS(map[string]string{
    gfly.HeaderAccessControlAllowOrigin: "*",
}))
```
### Access controls:
- Access-Control-Allow-Origin: Accept all domains `*` (default)
- Access-Control-Allow-Headers: Accept all header parameters `Authorization, Content-Type, x-requested-with, origin, true-client-ip, X-Correlation-ID` (default)
- Access-Control-Allow-Methods: List supported methods `PUT`, `POST`, `GET`, `DELETE`, `OPTIONS`, `PATCH` (default)
- Access-Control-Allow-Credentials: N/A
- Access-Control-Expose-Headers: N/A
- Access-Control-Max-Age: N/A
