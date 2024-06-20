<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>uTasksd - Task Manager</title>

    <link rel="stylesheet" href="/static/index.css" />
    <link rel="stylesheet" href="/static/view.css" />
    <link rel="stylesheet" href="/static/fonts.css" />
    <link rel="stylesheet" href="/static/project.css" />

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
        <main class="central-content">
            <section class="view-block">
                <div class="table-container">
                    <table id="user-table-body">
                        <thead>
                            <tr>
                                <th></th>
                                <th>ID</th>
                                <th>Username</th>
                                <th>First Name</th>
                                <th>Last Name</th>
                                <th>Email</th>

                                {{if .IsAdmin}}
                                <th>New Password</th>
                                <th>Director</th>
                                <th>Admin</th>
                                <th>Active</th>
                                <th>Action</th>
                                {{end}}
                            </tr>
                        </thead>
                        <tbody>
                            {{$isAdmin := .IsAdmin}}{{range .Users}}
                            <tr>
                                <td class="table-avatar"><img src="{{.AvatarPath}}"/></td>
                                <td id="id">{{.ID}}</td>
                                <td><a href="/web/user/{{.ID}}">{{.Username}}</a></td>
                                <td>{{.Firstname}}</td>
                                <td>{{.Lastname}}</td>
                                <td>{{.Email}}</td>
                                
                                {{if $isAdmin}}
                                <td><input type="password" placeholder="New Password"></td>
                                <td><input type="checkbox" {{if .Director}}checked{{end}} {{if not $isAdmin}}disabled{{end}}> </td>
                                <td><input type="checkbox" {{if .Admin}}checked{{end}} {{if not $isAdmin}}disabled{{end}}></td>
                                <td><input type="checkbox" {{if .Active}}checked{{end}} {{if not $isAdmin}}disabled{{end}}></td>
                                <td>
                                    <button onclick="saveUser(this)">Save User</button>
                                </td>
                                {{end}}
                            </tr>
                            {{end}}
                        </tbody>
                    </table>
                </div>
            </section>
        </main>
        <script>
        </script>
        <script src="/static/index.js"></script>
    </div>
</body>

</html>