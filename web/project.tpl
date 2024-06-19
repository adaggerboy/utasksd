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
            <span class="account-name">{{ .Core.Account.Name }}</span>
            <button class="logout-button" onclick="logout()">Logout</button>
        </div>
    </header>
    <div class="container">
        <main class="central-content">
            <section class="task-block">
                <div class="project-info">
                    <img src="{{ .Project.Logo }}" id="logo" alt="Project Logo" class="project-logo">
                    <h1>{{ .Project.Name }}</h1>
                    <p id="projectDescription" onclick="toggleEditMode('projectDescription')">{{ .Project.Description }}</p>
                </div>

                {{if .IsOwner}}
                <div class="add-user">
                    <input type="text" placeholder="Enter user ID" id="user-id-input">
                    <button onclick="addUser()">Add User</button>
                </div>
                <div class="table-container">
                    <table id="user-table-body">
                        <thead>
                            <tr>
                                <th>User</th>
                                <th>Role</th>
                                <th>Action</th>
                            </tr>
                        </thead>
                        <tbody>
                            {{range .Project.Users}}
                            <tr data-user-id="{{ .ID }}">
                                <td class="table-avatar"><img src="{{ .PathToAvatar }}" alt=""><span>{{ .Name }}</span></td>
                                <td>
                                    <select value="{{ .Role }}">
                                        <option value="manager" {{ if eq .Role "manager" }}selected{{end}}>Manager</option>
                                        <option value="worker" {{ if eq .Role "worker" }}selected{{end}}>Worker</option>
                                        <option value="support" {{ if eq .Role "support" }}selected{{end}}>Support agent</option>
                                    </select>
                                </td>
                                <td><button class="delete-button" onclick="deleteUser(this)">Delete</button></td>
                            </tr>
                            {{ end }}
                        </tbody>
                    </table>
                </div>
                {{ end }}
            </section>
        </main>
        <aside class="right-menu">
            {{if .IsOwner}}
            <h3>Project preferences</h3>
            <div class="menu-item">
                <label for="project-name">Project name:</label>
                <input type="text" id="project-name" value="{{ .Project.Name }}">
            </div>
            <div class="menu-item image-selector">
                <label for="project-image">Project image</label>
                <button onclick="openuplModal()">Upload image</button>
            </div>
            <div class="menu-item">
                <button id="save-task" onclick="updateProject()">Save project</button>
                <button class="delete-button" onclick="openDelModal()">Delete task</button>
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
                <button class="confirm" onclick="deleteProject(); closeDelModal()">Yes</button>
                <button class="cancel" onclick="closeDelModal('taskModal')">No</button>
            </div>
        </div>
        <script>
            window.jsGlobals.projectID = {{ .Project.ID }};
        </script>
        <script src="/static/index.js"></script>
    </div>
</body>

</html>