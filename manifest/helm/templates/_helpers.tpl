{{- define "cassata.fullname" -}}
{{- .Release.Name }}-{{ .Chart.Name }}
{{- end -}}
{{- define "cassata.selectorLabels" -}}
app: {{ .Chart.Name }}
release: {{ .Release.Name }}
{{- end -}}
{{- define "cassata.labels" -}}
app: {{ .Chart.Name }}
release: {{ .Release.Name }}
{{- end -}}
