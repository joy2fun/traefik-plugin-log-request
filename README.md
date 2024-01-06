# Traefik plugin to log requests

file provider example:

```yml
http:
  middlewares:
    my-plugin:
      plugin:
        log-request:
          responseBody: true # also including response body
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
      responseBody: true
```

helm chart values example:

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