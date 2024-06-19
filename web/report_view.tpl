<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Report</title>
    <link rel="stylesheet" href="/static/fonts.css" />
    <style>
        * {
            font-size: 12pt;
            font-family: 'Roboto';
        }
        table {
            border-collapse: collapse;
            width: 1000px;
        }

        th, td {
            border: 1px solid #000;
            text-align: left;
            height:20px;
        }

        th {
            font-weight: 500;
        }

    </style>
</head>
<body>

<h2>Time efficiency report</h2>
<table>
  <tr>
    <th>ID</th>
    <th>Username</th>
    <th>First Name</th>
    <th>Last Name</th>
    <th>Tracking Records</th>
    <th>Summary</th>
    <th>Average Per Day</th>
    <th>Average Per Month</th>
  </tr>
  {{range .TimeEfficiency}}
  <tr>
    <td>{{.ID}}</td>
    <td>{{.Username}}</td>
    <td>{{.Firstname}}</td>
    <td>{{.Lastname}}</td>
    <td>{{.TrackedRecords}}</td>
    <td>{{.SummaryHours}}</td>
    <td>{{.AveragePerDay}}</td>
    <td>{{.AveragePerMonth}}</td>
  </tr>
  {{end}}
</table>

<h2>Tasks report</h2>
<table>
  <tr>
    <th>ID</th>
    <th>Username</th>
    <th>First Name</th>
    <th>Last Name</th>
    <th>Tasks</th>
    <th>Done Tasks</th>
    <th>Issues</th>
    <th>Done Issues</th>
  </tr>
  {{range .Tasks}}
  <tr>
    <td>{{.ID}}</td>
    <td>{{.Username}}</td>
    <td>{{.Firstname}}</td>
    <td>{{.Lastname}}</td>
    <td>{{.Tasks}}</td>
    <td>{{.DoneTasks}}</td>
    <td>{{.Issues}}</td>
    <td>{{.ClosedIssues}}</td>
  </tr>
  {{end}}
</table>

</body>
</html>
