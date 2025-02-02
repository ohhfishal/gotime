| Category     | Start | Length    | Note
| --------     | ----- | ------    | ----
{{range .Schedule -}}
| {{ truncRight 12 "..." .Category | printf "%-12s" }} | {{ time .Time }} | {{ duration .Duration | printf "%-9s | " }}{{ trunc 30 .Note }}
{{end}}

 Category        Length  
 --------        ------
{{range $category, $total := .Categories -}}
  {{- truncRight 12 "..." $category | printf "| %-12s | " -}}
  {{- duration $total | printf "%-9s" }}
{{end}}
