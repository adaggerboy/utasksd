<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>uTasksd - Task Manager</title>

    <link rel="stylesheet" href="/static/index.css" />
    <link rel="stylesheet" href="/static/view.css" />
    <link rel="stylesheet" href="/static/tree_view.css" />
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
        <aside class="left-menu">
            <h3>Project tree view</h3>
            {{define "treeNode"}}
            {{range .}}

            <li class="{{ .VisibleClass }}">
                <details open>
                    <summary><a href="/web/task/{{ .ID }}">{{ .Name }}</a></summary>
                    {{if .Children}}
                    <ul>
                        {{template "treeNode" .Children}}
                    </ul>
                    {{end}}
                </details>
            </li>
            {{end}}
            {{end}}
            <ul class="tree">
                {{range .AvailableProjects}}
                <li class="project-node">
                    <details open>
                        <summary class="project-node-name"><a href="/web/project/{{ .ID }}/tasks">{{ .Name }}</a></summary>
                        {{if .Children}}
                            <ul>
                            {{template "treeNode" .Children}}
                            </ul>
                        {{end}}
                        </summary>
                    </details>
                </li>
                {{end}}
            </ul>
        </aside>
        <main class="central-content">
            <section class="view-block">
                <div class="list-view">
                    <h2>Tasks</h2>

                    {{ range .Tasks }} 
                    <div class="task-block" onclick="window.location.href= href='/web/task/{{ .ID }}'">
                        <div class="task-header">
                            {{if .FirePriority }}
                            <span class="priority-sign">!</span>
                            {{end}}
                            <a href="/web/project/{{ .ProjectID }}" class="project-name">{{ .ProjectName }}</a>
                            <span>:</span>
                            <a href="/web/task/{{ .ID }}" class="task-name">{{ .Name }}</a>
                            <a href="/web/task/{{ .ID }}" class="task-id">#{{ .ID }}</a>
                            <span class="task-status {{ .VisibleClass }}">{{ .Status }}</span>
                        </div>
                        <div class="task-description">
                            {{ .Description }}
                        </div>
                        <div class="assign-info">
                            <img src="{{ .Assigner.PathToAvatar }}" alt="" />
                            <span class="assigner-name">{{ .Assigner.Name }}</span>
                            <span>&rarr;</span>
                            {{range .Assignees}}
                            <img src="{{ .PathToAvatar }}" alt="" />
                            <span class="assigner-name">{{ .Name }}</span>
                            {{end}}
                        </div>
                    </div>
                    {{end}}
                </div>
            </section>
        </main>
        <aside class="right-menu">
            {{if .ProjectData.InProject }}
            <h3>Filters for tasks</h3>
            <div class="menu-item">
                <label for="filter-view">View</label>
                <select id="filter-view">
                    <option value="list">List</option>
                    <option value="kanban">Kanban</option>
                </select>
            </div>
            <div class="menu-item">
                <label for="filter-name">Name</label>
                <input type="text" id="filter-name" placeholder="Search by name" />
            </div>
            <div class="menu-item">
                <label for="filter-assigner-input">Assigner</label>
                <input type="text" id="filter-assigner-input" placeholder="Search by assigner" />
                <ul id="filter-assigner-suggestions" class="suggestions-list"></ul>
            </div>
            <div class="menu-item">
                <label for="filter-assignee-input">Assignee</label>
                <input type="text" id="filter-assignee-input" placeholder="Search by assignee" />
                <ul id="filter-assignee-suggestions" class="suggestions-list"></ul>
            </div>
            <div class="menu-item">
                <label for="filter-status">Status</label>
                <select id="filter-status">
                    {{range $tag, $value := .Core.Statuses}}
                    <option value="{{ $tag }}">{{ $value }}</option>
                    {{end}}
                    <option value="any" selected>any</option>
                </select>
            </div>
            <div class="menu-item">
                <label for="filter-priority">Priority</label>
                <select id="filter-priority">
                    {{range $tag, $value := .Core.Priorities}}
                    <option value="{{ $tag }}">{{ $value }}</option>
                    {{end}}
                    <option value="any" selected>any</option>
                </select>
            </div>
            <div class="menu-item">
                <button id="search" onclick="searchTasks()">Search</button>
            </div>
            {{end}}
        </aside>
        <script>
            window.jsGlobals.assigners = {{ .ProjectData.Assigners }};
            window.jsGlobals.assignees = {{ .ProjectData.Assignees }};
            window.jsGlobals.projectID = {{ .ProjectData.ID }};
        </script>
        <script src="/static/index.js"></script>
    </div>
</body>

</html>