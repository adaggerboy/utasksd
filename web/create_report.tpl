<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Request report</title>
    <link rel="stylesheet" href="/static/auth.css">
    <link rel="stylesheet" href="/static/fonts.css" />
    
    <base href="/">

    <style>
        :root {
        {{range $key, $value := .Core.CSSVariables}}
            --{{$key}}: {{$value}};
        {{end}}
        }
    </style>
</head>
<body>
    <div class="container">
        <form id="createReportForm" method="get" action="/web/report" class="form">
            <h2>Request project</h2>
            <div class="break"></div>
            <div class="form-group">
                <label for="name">Project:</label>
                <input type="text" id="name" name="name" value="{{.Name}}" disabled>
                <span class="error" id="nameError"></span>
            </div>
            <div class="form-group">
                <label for="startDate">Start Date:</label>
                <input type="month" id="startDate" name="startDate">
            </div>
            <div class="form-group">
                <label for="dueDate">Due Date:</label>
                <input type="month" id="dueDate" name="dueDate">
            </div>
                <input type="hidden" id="hiddenParam1" name="project" value="{{.ID}}">
            <div class="break"></div>
            <button type="submit">Create report</button>
        </form>
    </div>
</body>
</html>
