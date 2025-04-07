{{/*
Common labels
*/}}
{{- define "kyc-service.labels" -}}
app: {{ .Chart.Name }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "kyc-service.selectorLabels" -}}
app: {{ .Chart.Name }}
{{- end }}