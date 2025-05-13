# {{ .Time | date }}
{{ if .Schedule | not -}} 
 No work logged today {{printf "\n"}}
{{ else -}} 
## Schedule
| Category | Start | Length | Note                     |
| -------- | ----- | ------ | ------------------------ |
{{range .Schedule -}}
  {{- truncRight 12 "..." .Category | printf "| %s" -}} 
  {{- time .Time | printf " | %s" -}}
  {{- duration .Duration | printf " | %s" -}}
  {{- truncRight 40 "..." .Note | printf " | %s |\n" -}}
{{- end }}
## Summary
| Category | Length |
| -------- | ------ |
{{range $category, $total := .Categories -}}
  {{- truncRight 32 "..." $category | printf "| %s" -}}
  {{- duration $total | printf " | %s |" }}
  {{- printf "\n" -}}
{{- end -}}
  {{- truncRight 32 "..." "total" | printf "| %s" -}}
  {{- duration .Total | printf " | %s |" }}
  {{- printf "\n" -}}
{{- end -}}
