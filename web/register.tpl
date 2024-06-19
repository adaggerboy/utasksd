<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Registration form</title>
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
        <form id="registrationForm" class="form">
            <h2>Register</h2>
            <div class="break"></div>
            <div class="form-group">
                <label for="username">Username:</label>
                <input type="text" id="username" name="username">
                <span class="error" id="usernameError"></span>
            </div>
            <div class="form-group">
                <label for="firstname">First Name:</label>
                <input type="text" id="firstname" name="firstname">
                <span class="error" id="firstnameError"></span>
            </div>
            <div class="form-group">
                <label for="lastname">Last Name:</label>
                <input type="text" id="lastname" name="lastname">
                <span class="error" id="lastnameError"></span>
            </div>
            <div class="form-group">
                <label for="email">Email:</label>
                <input type="email" id="email" name="email">
                <span class="error" id="emailError"></span>
            </div>
            <div class="form-group">
                <label for="password">Password:</label>
                <input type="password" id="password" name="password">
                <span class="error" id="passwordError"></span>
            </div>
            <div class="break"></div>
            <button type="submit">Register</button>
        </form>
    </div>
    <script src="/static/auth.js"></script>
</body>
</html>
