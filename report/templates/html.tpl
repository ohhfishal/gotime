<!DOCTYPE html>
<html>
<head>
<title>Daily Report</title>
<style>
  body {
    font-family: sans-serif;
  }
  table {
    border-collapse: collapse;
    width: 100%;
    margin-bottom: 1em;
  }
  th, td {
    border: 1px solid #ddd;
    padding: 8px;
    text-align: left;
  }
  th {
    background-color: #f2f2f2;
  }
</style>
</head>
<body>

  <h1>{{ .Time | date }}</h1>

  {{ if .Schedule | not }}
    <p>No work logged today</p>
  {{ else }}
    <h2>Schedule</h2>
    <table>
      <thead>
        <tr>
          <th>Category</th>
          <th>Start</th>
          <th>Length</th>
          <th>Note</th>
        </tr>
      </thead>
      <tbody>
        {{range .Schedule}}
          <tr>
            <td>{{ truncRight 12 "..." .Category }}</td>
            <td>{{ time .Time }}</td>
            <td>{{ duration .Duration }}</td>
            <td>{{ truncRight 40 "..." .Note }}</td>
          </tr>
        {{end}}
      </tbody>
    </table>

    <h2>Summary</h2>
    <table>
      <thead>
        <tr>
          <th>Category</th>
          <th>Length</th>
        </tr>
      </thead>
      <tbody>
        {{range $category, $total := .Categories}}
          <tr>
            <td>{{ truncRight 32 "..." $category }}</td>
            <td>{{ duration $total }}</td>
          </tr>
        {{end}}
        <tr>
          <td>{{ truncRight 32 "..." "total" }}</td>
          <td>{{ duration .Total }}</td>
        </tr>
      </tbody>
    </table>
  {{ end }}

</body>
</html>
