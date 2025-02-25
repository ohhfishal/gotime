{{- if .Schedule | not -}} 
 No work logged today {{printf "\n"}}
{{- else }} Category       Start   Length      Note
 --------       -----   ------      ----
{{range .Schedule -}}
  {{- truncRight 12 "..." .Category | printf " %-12s   " -}} 
  {{- time .Time | printf "%s   " -}}
  {{- duration .Duration | printf "%-9s   " -}}
  {{- truncRight 40 "..." .Note | printf "%s\n" -}}
{{- end }}
 Category                           Length  
 --------                           ------
{{range $category, $total := .Categories -}}
  {{- truncRight 32 "..." $category | printf " %-32s   " -}}
  {{- duration $total | printf "%-9s" }}
  {{- printf "\n" -}}
{{- end -}}
  {{- truncRight 32 "..." "total" | printf " %-32s   " -}}
  {{- duration .Total | printf "%-9s" }}
  {{- printf "\n" -}}
{{- end -}}
