{{- define "cloud-provider-config-credentials" -}}
auth-url="{{ .Values.authUrl }}"
domain-name="{{ .Values.domainName }}"
tenant-name="{{ .Values.tenantName }}"
tenant-id="{{ .Values.tenantID }}"
username="{{ .Values.username }}"
password="{{ .Values.password }}"
region="{{ .Values.region }}"
{{- end -}}
