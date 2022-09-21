package utils

var configMapTemplate = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.Name}}
  {{if ne .Namespace "" }}
  namespace: {{.Namespace}}
  {{end}}
data:
{{if ne .File ""}}
swagger: {{.SwaggerContent}}
{{else}}
swagger: {{.DefaultSwagger}}
{{end}}
`
