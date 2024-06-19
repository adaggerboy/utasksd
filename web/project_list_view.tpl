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
            <span class="account-name">{{ .Core.Account.Name }}</span>
            <button class="logout-button" onclick="logout()">Logout</button>
        </div>
    </header>
    <div class="container">
        <main class="central-content">
            <section class="view-block">
                <div class="list-view">
                    <h2>Projects</h2>

                    {{ range .Projects }} 
                    <div class="task-block" onclick="window.location.href= href='/web/issue/{{ .ID }}'">
                        <img src="{{ .Logo }}" alt="" class="project-logo"/>
                        <div class="project-over">
                            <div class="task-header">
                                <a href="/web/project/{{ .ID }}" class="project-name">{{ .Name }}</a>
                                <a href="/web/project/{{ .ID }}" class="task-id">#{{ .ID }}</a>
                                <span>|</span>
                                <span class="task-id">{{ .Role }}</a>
                            </div>
                            <div class="task-description">
                                {{ .Description }}
                            </div>
                            <div class="assign-info">
                                <img src="{{ .Owner.PathToAvatar }}" alt="" />
                                <span class="assigner-name">{{ .Owner.Name }}</span>
                            </div>
                        </div>
                    </div>
                    {{end}}
                </div>
            </section>
        </main>
        <script src="/static/index.js"></script>
    </div>
</body>

</html>