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
        <aside class="left-menu">
            <h3>Filters for issues</h3>
            <div class="menu-item">
                <label for="filter-name">Name</label>
                <input type="text" id="filter-name" placeholder="Search by name" />
            </div>
            <div class="menu-item">
                <label for="filter-registrar">Registrar</label>
                <input type="text" id="filter-registrar" placeholder="Search by registrar" />
                <ul id="registrar-suggestions" class="suggestions-list"></ul>
            </div>
            <div class="menu-item">
                <label for="filter-status">Status</label>
                <select id="filter-status">
                    {{range $tag, $value := .Core.IssueStatuses}}
                    <option value="{{ $tag }}">{{ $value }}</option>
                    {{end}}
                    <option value="any" selected>any</option>
                </select>
            </div>
            <div class="menu-item">
                <button id="search" onclick="searchIssues()">Search</button>
            </div>
        </aside>
        <main class="central-content">
            <section class="task-block" id="task-central-block">
                <div class="task-header">
                    <span class="task-name"  id="taskName" {{if .Issue.EditCapability  }}onclick="toggleEditMode('taskName')"{{end}}>{{ .Issue.Name }}</span>
                    <span class="task-id">#{{ .Issue.ID }}</span>
                </div>
                <div class="task-description"  id="taskDescription" {{if .Issue.EditCapability  }}onclick="toggleEditMode('taskDescription')"{{end}}>
                    {{ .Issue.Description }}
                </div>
                {{range $item := .Issue.Attachments}}
                <div class="image-wrapper">
                    <img src="{{$item}}" />
                    <button class="delete-btn delete-button" onclick="deleteAttachment(this)">Delete</button>
                </div>
                {{end}}
            </section>
        </main>
        <aside class="right-menu">
            <div class="menu-item">
                <label for="status">Status:</label>
                <input type="text" id="status" name="status" disabled value="{{ .Issue.Status }}"/>
            </div>
            <div class="menu-item">
                <label for="registrar">Registrar:</label>
                <input type="text" id="registrar" name="registrar" disabled value="{{ .Issue.Registrar }}"/>
            </div>
            <div class="menu-item">
                <label for="reporter">Reporter:</label>
                <input type="text" id="reporter" name="reporter" {{if .Issue.Created }}disabled{{end}} value="{{ .Issue.Reporter }}"/>
            </div>
            <div class="menu-item">
                <label for="start-date">Start Date:</label>
                <input type="date" id="start-date" name="start-date" value="{{ .Issue.StartDate }}" disabled/>
            </div>
            <div class="menu-item">
                <label for="due-date">Close Date:</label>
                <input type="date" id="due-date" name="due-date" value="{{ .Issue.DueDate }}" {{if not .Issue.EditCapability  }}disabled{{end}} disabled/>
            </div>
            <div class="menu-item">
                <button id="attach-button" {{if not .Issue.EditCapability }}disabled{{end}} onclick="openuplModal()">Attach</button>
            </div>
            <div class="menu-item">
                <button id="save-issue" onclick="{{ if .Issue.Created }}updateIssue();{{else}}createIssue();{{end}}">Save issue</button>
                {{if .Issue.Created }}
                {{if .Issue.Closed}}
                <button id="reopen-issue" onclick="reopenIssue()" {{if not .Issue.EditCapability }}disabled{{end}}>Reopen</button>
                {{else}}
                <button id="close-issue" onclick="closeIssue()" {{if not .Issue.EditCapability }}disabled{{end}}>Close</button>
                {{end}}
                <button class="delete-button" onclick="openDelModal()">Delete issue</button>
                {{end}}
            </div>
        </aside>

        <div id="uplModal" class="modal">
            <div class="modal-content">
                <span class="close" onclick="closeUplModal()">&times;</span>
                <h2>Upload File</h2>
                <input type="file" id="fileInput">
                <button id="uploadBtn">Upload</button>
            </div>

        </div>
        <div id="delModal" class="modal">
            <div class="modal-content">
                <p>Are you sure?</p>
                <button class="confirm" onclick="deleteIssue(); closeDelModal()">Yes</button>
                <button class="cancel" onclick="closeDelModal('taskModal')">No</button>
            </div>
        </div>
        <script>
            window.jsGlobals.supports = {{ .ProjectData.Supports }};
            window.jsGlobals.projectID = {{ .ProjectData.ID }};
            window.jsGlobals.issueID = {{ .Issue.ID }};
        </script>
        <script src="/static/index.js"></script>
    </div>
</body>

</html>