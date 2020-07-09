* {{ .DateF "Mon Jan  2 2006" }} {{ .Author }} - {{ .Version }}
{{ range .Entries }}
- {{ .Title }}
{{- end }}
