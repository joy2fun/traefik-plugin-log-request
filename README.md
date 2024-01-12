# Traefik plugin to log requests

file provider example:

```yml
http:
  middlewares:
    my-plugin:
      plugin:
        log-request:
          ResponseBody: false # also including response body
          RequestIDHeaderName: X-Request-Id
```

crd example:

```yml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: log-request
spec:
  plugin:
    log-request:
      ResponseBody: false
      RequestIDHeaderName: X-Request-Id
```

helm chart values example (local plugin mode):

```yaml
additionalArguments:
  - >-
    --experimental.localplugins.log-request.modulename=github.com/joy2fun/traefik-plugin-log-request
additionalVolumeMounts:
  - mountPath: /plugins-local
    name: plugins
deployment:
  additionalVolumes:
    - hostPath:
        path: /data/plugins-local
      name: plugins
```