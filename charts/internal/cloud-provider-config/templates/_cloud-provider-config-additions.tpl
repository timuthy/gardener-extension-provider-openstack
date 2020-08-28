{{- define "cloud-provider-config-additions" -}}
[OVH]
cluster-namespace={{ required ".Values.clusterName is required" .Values.clusterName }}
ovhapis-url=https://api.ovh.com
{{- end -}}
