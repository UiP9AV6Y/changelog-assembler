package release

const (
	DefaultTemplate = `## {{ .Version }} ({{ .Date }})
{{ range .Entries.SortByReason }}
* [{{ .Reason.CapsString }}] {{- if .Component }} {{ .Component }}:{{ end }} {{ .Title }} {{- if .MergeRequest }} #{{ .MergeRequest }}{{ end }}
{{- end }}
`
	GroupedTemplate = `## {{ .Version }} ({{ .Date }})
{{ with .Entries.Authors }}
## Contributors
{{ range . }}
* {{ . }}
{{- end }}
{{- end }}
{{ range $reason, $entries := .Entries.GroupByReason }}
### {{ $reason.Description }}

{{ range $component, $changes := $entries.GroupByComponent }}
{{- if $component -}}
* {{ $component }}
{{- else -}}
* General changes
{{- end -}}
{{ range $changes }}
  * {{ .Title }} {{- if .MergeRequest }} #{{ .MergeRequest }}{{ end }}
{{- end }}
{{- end }}
{{ end }}
`
)
