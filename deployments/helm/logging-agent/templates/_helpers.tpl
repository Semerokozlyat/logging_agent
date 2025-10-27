{{/*
Expand the name of the chart.
*/}}
{{- define "logging-agent.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "logging-agent.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
K8s namespace name
*/}}
{{- define "logging-agent.namespace" -}}
{{- default "logging-system" .Values.kubernetes.customNamespaceName }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "logging-agent.labels" -}}
helm.sh/chart: {{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{ include "logging-agent.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: monitoring
{{- end }}

{{/*
Selector labels
*/}}
{{- define "logging-agent.selectorLabels" -}}
app.kubernetes.io/name: {{ include "logging-agent.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

