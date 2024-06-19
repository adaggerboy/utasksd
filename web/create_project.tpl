<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Create project</title>
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
    <script>
        window.jsCoreGlobals = {{.Core.JSVariablesJSON}};
        window.jsGlobals = {{.JSVariablesJSON}};
    </script>
</head>
<body>
    <div class="container">
        <form id="createProjectForm" class="form">
            <h2>Create project</h2>
            <div class="break"></div>
            <div class="form-group">
                <label for="name">Name:</label>
                <input type="text" id="name" name="name">
                <span class="error" id="nameError"></span>
            </div>
            <div class="form-group">
                <label for="description">Description:</label>
                <textarea id="description" name="description"></textarea>
            </div>
            <div class="break"></div>
            <button type="submit">Create project</button>
        </form>
    </div>
    <script src="/static/auth.js"></script>
</body>
</html>
