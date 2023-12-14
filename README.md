### About plugin
This plugin allows to generate uuid V4 correlation ID for every request and transmit them to next request.
Its allow to send pure correlationId key from client by send header name for this correlation or generate them by plugin with default header correlation-id name.
### Configure (default)
```yaml
...
    correlation:
      headerName: "correlation-id" # default correlation ID header name if client does not sent them
...
```
If client sent this header name by correlation, then plugin catch this value and transmit them to next
```shell
curl -XPOST https://<some domain>/<some api uri> -H 'Correlation-Id: <somw uuid correlation value or another unique value>'
```

### Install
For install this middleware plugin you can configure them like this

#### Configure traefik-ingress for kubernates
First step is to add experimental plugin
```yaml
# anywhere/traefik.yml
experimental:
  plugins:
    correlation-id:
      moduleName: github.com/killer-djon/traefik-correlation
      version: v1.3.0
```
Next step is to add customDefenition for middleware
```yaml
...
-   apiVersion: traefik.containo.us/v1alpha1
    kind: Middleware
    metadata:
      name: my-traefik-correlation
      namespace: default
    spec:
      plugin:
        correlation-id:
          headerName: "correlation-id" # this is headerName to catch them
...
```
And as the third step is to add ingress annotation for this middleware (if need to be at single middleware for single service) or set this plugin for avery services as entryPoint (web, websecure ...)
1. At ingress annotation
```yaml
...
traefik.ingress.kubernetes.io/router.middlewares: default-my-traefik-correlation@kubernetescrd
```
2. At general entrypoint of traefik instance
```yaml
ports:
  ...
  websecure:
    port: 8443
    expose: true
    exposedPort: 443
    protocol: TCP
    appProtocol: https
    middlewares:
    - default-my-traefik-correlation@kubernetescrd
# in this case you installed plugin will works for every request
```