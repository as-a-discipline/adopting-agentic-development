{{- define "pulse.labels" -}}
app.kubernetes.io/part-of: pulse
{{- end -}}

{{- define "pulse.service.selectorLabels" -}}
app.kubernetes.io/name: {{ .Values.service.name }}
{{- end -}}

{{- define "pulse.web.selectorLabels" -}}
app.kubernetes.io/name: {{ .Values.web.name }}
{{- end -}}
