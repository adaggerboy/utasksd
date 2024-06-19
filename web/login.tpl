<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login form</title>
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
        <form id="loginForm" class="form">
            <h2>Login</h2>
            <div class="break"></div>
            <div class="form-group">
                <label for="username">Username:</label>
                <input type="text" id="username" name="username" required>
                <span class="error" id="firstnameError"></span>
            </div>
            <div class="form-group">
                <label for="password">Password:</label>
                <input type="password" id="password" name="password" required>
                <span class="error" id="firstnameError"></span>
            </div>
            <div class="break"></div>
            <button type="submit">Login</button>
            <button type="button" onclick="window.location.href='/web/register'" class="side-button">Register</button>
        </form>
    </div>
    <script src="/static/auth.js"></script>
</body>
</html>
