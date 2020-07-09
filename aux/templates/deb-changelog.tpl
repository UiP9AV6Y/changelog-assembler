$PACKAGE ({{ .Version }}) {{ .GetAnnotation "distributions" "unstable" }};
{{- range $k, $v := .Annotations }}
  {{- if ne $k "distributions" }}
    {{- $k }}={{ $v }};
  {{- end }}
{{- end }}

{{ range .Entries }}
  * {{ .Title }} ({{ .Author }})
{{- end}}
-- {{ .Author }}  {{ .GetDateF "Mon, 2 Jan 2006 15:04:05 -0700"}}
