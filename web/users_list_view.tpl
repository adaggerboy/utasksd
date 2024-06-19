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
                <table>
                    <thead>
                        <tr>
                            <th></th>
                            <th>ID</th>
                            <th>Username</th>
                            <th>First Name</th>
                            <th>Last Name</th>
                            <th>New Password</th>
                            <th>Director</th>
                            <th>Admin</th>
                            <th>Active</th>
                            {{if .IsAdmin}}
                            <th>Action</th>
                            {{end}}
                        </tr>
                    </thead>
                    <tbody>
                        {{$isAdmin := .IsAdmin}}{{range .Users}}
                        <tr>
                            <td><img class="avatar" src="{{.AvatarPath}}"/></td>
                            <td id="id">{{.ID}}</td>
                            <td><a href="/web/user/{{.ID}}">{{.Username}}</a></td>
                            <td>{{.Firstname}}</td>
                            <td>{{.Lastname}}</td>
                            
                            <td>{{if $isAdmin}}<input type="password" placeholder="New Password">{{end}}</td>
                            <td><input type="checkbox" {{if .Director}}checked{{end}} {{if not $isAdmin}}disabled{{end}}> </td>
                            <td><input type="checkbox" {{if .Admin}}checked{{end}} {{if not $isAdmin}}disabled{{end}}></td>
                            <td><input type="checkbox" {{if .Active}}checked{{end}} {{if not $isAdmin}}disabled{{end}}></td>
                            {{if $isAdmin}}
                            <td>
                                <button onclick="saveUser(this)">Save User</button>
                            </td>
                            {{end}}
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </section>
        </main>
        <script>
        </script>
        <script src="/static/index.js"></script>
    </div>
</body>

</html>