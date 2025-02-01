| Category    | Start | Length    | Note
| --------    | ----- | ------    | ----
{{range . -}}
| {{if (len .Category) | gt 12 -}}
    {{- printf "%-11s | " .Category -}}
  {{- else -}}
    {{- trunc 8 .Category | printf "%s... | " -}}
  {{- end -}}
 {{- time .Time }} | {{ duration .Duration | printf "%9s | " }}{{ trunc 30 .Note }}
{{end}}
