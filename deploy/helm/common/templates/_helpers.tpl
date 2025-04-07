{{/*
Common labels
*/}}
{{- define "common.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Values.image.tag | default .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "common.selectorLabels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "common.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "common.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "common.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "common.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common annotations
*/}}
{{- define "common.annotations" -}}
helm.sh/chart: {{ include "common.chart" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Pod annotations
*/}}
{{- define "common.podAnnotations" -}}
{{- if .Values.podAnnotations }}
{{- toYaml .Values.podAnnotations | nindent 8 }}
{{- end }}
prometheus.io/scrape: "true"
prometheus.io/port: "{{ .Values.metrics.port | default "8080" }}"
prometheus.io/path: "{{ .Values.metrics.path | default "/metrics" }}"
{{- end }}

{{/*
Common environment variables
*/}}
{{- define "common.env" -}}
- name: APP_ENV
  value: {{ .Values.env.environment | quote }}
- name: APP_LOG_LEVEL
  value: {{ .Values.env.logLevel | quote }}
- name: APP_LOG_FORMAT
  value: {{ .Values.env.logFormat | quote }}
- name: APP_METRICS_ENABLED
  value: {{ .Values.metrics.enabled | quote }}
- name: APP_METRICS_PORT
  value: {{ .Values.metrics.port | quote }}
- name: APP_METRICS_PATH
  value: {{ .Values.metrics.path | quote }}
- name: APP_TRACING_ENABLED
  value: {{ .Values.tracing.enabled | quote }}
- name: APP_TRACING_PROVIDER
  value: {{ .Values.tracing.provider | quote }}
- name: APP_TRACING_ENDPOINT
  value: {{ .Values.tracing.endpoint | quote }}
- name: APP_TRACING_SAMPLING_RATE
  value: {{ .Values.tracing.samplingRate | quote }}
- name: APP_TRACING_PROPAGATION
  value: {{ .Values.tracing.propagation | quote }}
{{- if .Values.tracing.serviceName }}
- name: APP_TRACING_SERVICE_NAME
  value: {{ .Values.tracing.serviceName | quote }}
{{- else }}
- name: APP_TRACING_SERVICE_NAME
  value: {{ .Chart.Name | quote }}
{{- end }}
{{- if .Values.tracing.environment }}
- name: APP_TRACING_ENVIRONMENT
  value: {{ .Values.tracing.environment | quote }}
{{- else }}
- name: APP_TRACING_ENVIRONMENT
  value: {{ .Release.Namespace | quote }}
{{- end }}
{{- if eq .Values.tracing.provider "zipkin" }}
- name: APP_TRACING_ZIPKIN_ENDPOINT
  value: {{ .Values.tracing.zipkin.endpoint | quote }}
{{- end }}
{{- if eq .Values.tracing.provider "otlp" }}
- name: APP_TRACING_OTLP_ENDPOINT
  value: {{ .Values.tracing.otlp.endpoint | quote }}
- name: APP_TRACING_OTLP_INSECURE
  value: {{ .Values.tracing.otlp.insecure | quote }}
- name: APP_TRACING_OTLP_TIMEOUT
  value: {{ .Values.tracing.otlp.timeout | quote }}
{{- end }}
{{- end }}

{{/*
Database environment variables
*/}}
{{- define "common.databaseEnv" -}}
{{- if .Values.database.enabled }}
- name: APP_DATABASE_HOST
  valueFrom:
    secretKeyRef:
      name: {{ .Release.Name }}-db-credentials
      key: host
- name: APP_DATABASE_PORT
  valueFrom:
    secretKeyRef:
      name: {{ .Release.Name }}-db-credentials
      key: port
- name: APP_DATABASE_USER
  valueFrom:
    secretKeyRef:
      name: {{ .Release.Name }}-db-credentials
      key: username
- name: APP_DATABASE_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ .Release.Name }}-db-credentials
      key: password
- name: APP_DATABASE_NAME
  valueFrom:
    secretKeyRef:
      name: {{ .Release.Name }}-db-credentials
      key: database
- name: APP_DATABASE_SSLMODE
  value: {{ .Values.database.sslMode | quote }}
{{- end }}
{{- end }}

{{/*
Cache environment variables
*/}}
{{- define "common.cacheEnv" -}}
{{- if .Values.cache.enabled }}
- name: APP_CACHE_ENABLED
  value: "true"
- name: APP_CACHE_TYPE
  value: {{ .Values.cache.type | quote }}
{{- if eq .Values.cache.type "redis" }}
- name: APP_CACHE_REDIS_HOST
  valueFrom:
    secretKeyRef:
      name: {{ .Release.Name }}-redis-credentials
      key: host
- name: APP_CACHE_REDIS_PORT
  valueFrom:
    secretKeyRef:
      name: {{ .Release.Name }}-redis-credentials
      key: port
- name: APP_CACHE_REDIS_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ .Release.Name }}-redis-credentials
      key: password
{{- end }}
{{- end }}
{{- end }}

{{/*
JWT environment variables
*/}}
{{- define "common.jwtEnv" -}}
{{- if .Values.jwt.enabled }}
- name: APP_JWT_SECRET
  valueFrom:
    secretKeyRef:
      name: {{ .Release.Name }}-jwt
      key: secret
- name: APP_JWT_EXPIRY
  value: {{ .Values.jwt.expiry | quote }}
{{- end }}
{{- end }}

{{/*
Common probes
*/}}
{{- define "common.probes" -}}
livenessProbe:
  httpGet:
    path: {{ .Values.probes.liveness.path | default "/health/live" }}
    port: http
  initialDelaySeconds: {{ .Values.probes.liveness.initialDelaySeconds | default 30 }}
  periodSeconds: {{ .Values.probes.liveness.periodSeconds | default 10 }}
  timeoutSeconds: {{ .Values.probes.liveness.timeoutSeconds | default 5 }}
  failureThreshold: {{ .Values.probes.liveness.failureThreshold | default 3 }}
readinessProbe:
  httpGet:
    path: {{ .Values.probes.readiness.path | default "/health/ready" }}
    port: http
  initialDelaySeconds: {{ .Values.probes.readiness.initialDelaySeconds | default 5 }}
  periodSeconds: {{ .Values.probes.readiness.periodSeconds | default 10 }}
  timeoutSeconds: {{ .Values.probes.readiness.timeoutSeconds | default 5 }}
  failureThreshold: {{ .Values.probes.readiness.failureThreshold | default 3 }}
{{- end }}

{{/*
Common security context
*/}}
{{- define "common.securityContext" -}}
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
  allowPrivilegeEscalation: false
  capabilities:
    drop:
    - ALL
  seccompProfile:
    type: RuntimeDefault
{{- end }}

{{/*
Common container security context
*/}}
{{- define "common.containerSecurityContext" -}}
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  allowPrivilegeEscalation: false
  capabilities:
    drop:
    - ALL
  readOnlyRootFilesystem: true
{{- end }}
