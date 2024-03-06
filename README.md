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
          MaxLineSize: 16384
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
      MaxLineSize: 16384
```

configMap via helm chart

```yml
apiVersion: v1
kind: ConfigMap
metadata:
  name: traefik-plugin-log-request
data:
{{ (.Files.Glob "plugin-log-request/*").AsConfig | indent 2 }}
```

traefik helm chart values example (local plugin mode):

```yaml
additionalArguments:
  - >-
    --experimental.localplugins.log-request.modulename=github.com/joy2fun/traefik-plugin-log-request
additionalVolumeMounts:
  - mountPath: /plugins-local/src/github.com/joy2fun/traefik-plugin-log-request
    name: plugins
deployment:
  additionalVolumes:
    - configMap:
        name: traefik-plugin-log-request
        items: 
          - key: dot.traefik.yml
            path: .traefik.yml
          - key: go.mod
            path: go.mod
          - key: main.go
            path: main.go
      name: plugins
```
