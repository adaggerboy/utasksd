<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>uTasksd - Task Manager</title>

    <base href="/">

    <link rel="stylesheet" href="/static/index.css" />
    <link rel="stylesheet" href="/static/fonts.css" />
    <link rel="stylesheet" href="/static/task.css" />
    <link rel="stylesheet" href="/static/project.css" />

    <style>
        :root {
        {{range $key, $value := .Core.CSSVariables}}
            --{{$key}}: {{$value}};
        {{end}}
        }
    </style>
    <script>
        window.jsCoreGlobals = {{.Core.JSVariablesJSON}};
        window.jsGlobals = {};
    </script>
</head>

<body>
    <header>
        <img class="logo" src="/static/logo.png" />
        <div class="search-bar">
            <input type="text" placeholder="Search...">
            <button>Search</button>
        </div>
        <nav class="menu">
            {{range  .Core.Options}}
            <a href="{{ .Href }}">{{ .Label }}</a>
            {{end}}
        </nav>
        <div class="account-section">
            <img src="{{ .Core.Account.PathToAvatar }}" alt="" />
            <a class="account-name" href="/web/user/{{ .Core.Account.ID }}">{{ .Core.Account.Name }}</a>
            <button class="logout-button" onclick="logout()">Logout</button>
        </div>
    </header>
    <div class="container">
        <main class="central-content">
            <section class="task-block">
                <div class="project-info">
                    <img src="{{ .User.AvatarPath }}" id="logo" alt="User Logo" class="project-logo">
                    <h1>{{ .User.Firstname }} {{ .User.Lastname }}</h1>
                </div>

                <p>Username: {{ .User.Username}}</p>
                <p>Email: {{ .User.Email}}</p>
                <p>ID: {{ .User.ID }}</p>
            </section>
        </main>
        <aside class="right-menu">
            {{if .ItsMe}}
            <h3>User preferences</h3>
            <div class="menu-item image-selector">
                <label for="project-image">User avatar</label>
                <button onclick="openuplModal()">Upload image</button>
            </div>
            <div class="menu-item">
                <button id="save-task" onclick="updateUser()">Save user</button>
                <button class="delete-button" onclick="openDelModal()">Deactivate user</button>
            </div>
            {{end}}
        </aside>

        <div id="uplModal" class="modal">
            <div class="modal-content">
                <span class="close" onclick="closeUplModal()">&times;</span>
                <h2>Upload File</h2>
                <input type="file" id="fileInput">
                <button id="avaBtn">Upload</button>
            </div>

        </div>

        <div id="delModal" class="modal">
            <div class="modal-content">
                <p>Are you sure?</p>
                <button class="confirm" onclick="deactivateUser(); closeDelModal()">Yes</button>
                <button class="cancel" onclick="closeDelModal('taskModal')">No</button>
            </div>
        </div>
        <script>
            window.jsGlobals.userID = {{ .User.ID }};
        </script>
        <script src="/static/index.js"></script>
    </div>
</body>

</html>