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
            <a class="account-name" href="/web/user/{{ .Core.Account.ID }}">{{ .Core.Account.Name }}</a>
            <button class="logout-button" onclick="logout()">Logout</button>
        </div>
    </header>
    <div class="container">
        <aside class="left-menu">
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
        </aside>
        <main class="central-content">
            <section class="task-block" id="task-central-block">
                <div class="task-header">
                    <span class="task-name"  id="taskName" {{if .Task.EditCapability  }}onclick="toggleEditMode('taskName')"{{end}}>{{ .Task.Name }}</span>
                    <span class="task-id">#{{ .Task.ID }}</span>
                </div>
                <div class="task-description"  id="taskDescription" {{if .Task.ChangeCapability  }}onclick="toggleEditMode('taskDescription')"{{end}}>
                    {{ .Task.Description }}
                </div>
                {{range $item := .Task.Attachments}}
                <div class="image-wrapper">
                    <img src="{{$item}}" />
                    <button class="delete-btn delete-button" onclick="deleteAttachment(this)">Delete</button>
                </div>
                {{end}}
            </section>
            {{if .Task.Created }}
            <section class="comments-block">
                <div class="comment-upload">
                    <textarea id="new-comment-text" placeholder="Write a comment..."></textarea>
                    <button id="upload-comment-btn" onclick="uploadComment()">Upload Comment</button>
                </div>
                {{range .Task.Comments }}
                    <div class="comment" data-comment-id="{{ .ID }}">
                        <img src="{{ .Avatar }}" alt="Avatar" class="avatar" />
                        <div class="comment-details">
                            <span class="comment-name">{{ .Name }}</span>
                            <span class="comment-time">{{ .TimeAgo }} ago</span>
                            {{if .CanRemove}}<button class="delete-btn delete-button" onclick="deleteComment(this)">Delete</button>{{end}}
                        </div>
                        <div class="comment-text">
                            {{ .Text }}
                        </div>
                        
                    </div>
                {{end}}
            </section>
            {{end}}
        </main>
        <aside class="right-menu">
            <div class="menu-item">
                <label for="status">Status</label>
                <select id="status" {{if not .Task.ChangeCapability }}disabled{{end}}>
                    {{range $tag, $value := .Core.Statuses}}
                    <option value="{{ $tag }}" {{if eq $tag $.Task.Status }}selected{{ end }}>{{ $value }}</option>
                    {{end}}
                </select>
            </div>
            <div class="menu-item">
                <label for="status">Priority</label>
                <select id="priority" {{if not .Task.EditCapability  }}disabled{{end}}>
                    {{range $tag, $value := .Core.Priorities}}
                    <option value="{{ $tag }}"  {{if eq $tag $.Task.Priority }}selected{{ end }}>{{ $value }}</option>
                    {{end}}
                </select>
            </div>
            <div class="menu-item">
                <label for="assigner">Assigner:</label>
                <input type="text" id="assigner-input2" name="assigner" {{if not .Task.ManageCapability }}disabled{{end}} value="{{ .Task.Assigner }}" keyVal="{{ .Task.AssignerID }}"/>
                <ul id="assigner-suggestions2" class="suggestions-list"></ul>
            </div>
            <div class="menu-item">
                <label for="assignee">Assignee:</label>
                <input type="text" id="assignee-input2" name="assignee" {{if not .Task.EditCapability  }}disabled{{end}} value="{{ .Task.Assignee }}" keyVal="{{ .Task.AssigneeID }}"/>
                <ul id="assignee-suggestions2" class="suggestions-list"></ul>
            </div>
            <div class="menu-item">
                <label for="start-date">Start Date:</label>
                <input type="date" id="start-date" name="start-date" value="{{ .Task.StartDate }}" {{if not .Task.EditCapability  }}disabled{{end}}/>
            </div>
            <div class="menu-item">
                <label for="due-date">Due Date:</label>
                <input type="date" id="due-date" name="due-date" value="{{ .Task.DueDate }}" {{if not .Task.EditCapability  }}disabled{{end}}/>
            </div>

            <div class="menu-item">
                <button id="attach-button" {{if not .Task.EditCapability }}disabled{{end}} onclick="openuplModal()">Attach</button>
            </div>
            {{if .Task.Created }}
            <div class="tracking">
                <h3>Tracking</h3>
                <div class="menu-item">
                    <button id="track-work" {{if not .Task.ChangeCapability }}disabled{{end}} onclick="openTrModal()">Track work</button>
                </div>
                <table class="tracking-table">
                    {{range .Task.TrackingRecords}}
                    <tr>
                        <td class="left notes">{{ .Text }}</td>
                        <td class="right">{{ .Duration }}</td>
                    </tr>
                    {{end}}
                </table>
                <div class="menu-item">
                    <p class="tracked">Tracked: {{.Task.OverallDuration}}</p>
                </div>
            </div>
            <div class="linked-issues">
                <h3>Linked Issues</h3>
                {{$task := .Task}}{{range .Task.LinkedIssues}}
                <div class="issue" data-issue-id={{.ID}}>
                    <a href="/web/issue/{{.ID}}">
                        <span class="issue-id">{{ .ID }}</span>
                        <span class="issue-name">{{ .Name }}</span>
                    </a>
                    <button class="delete-issue delete-button" onclick="deleteLinkedIssue({{ .ID }})" {{ if not $task.EditCapability }}disabled{{ end }}>x</button>
                </div>
                {{end}}
                <br><button class="add-issue" onclick="linkIssue()" {{if not .Task.EditCapability  }}disabled{{end}}>+ Add More</button><br>
            </div>

            <div class="dependent-tasks">
                <h3>Dependent Tasks</h3>
                {{$core := .Core}}{{$task := .Task}}{{range .Task.DependentTasks}}
                <div class="task" data-task-id={{.ID}}>
                    <a href="/web/task/{{.ID}}">
                        <span class="task-id">{{.ID}}</span>
                        <span class="task-name">{{.Name}}</span>
                    </a>
                    <select class="dependency-type" {{if not $task.EditCapability  }}disabled{{end}}>
                        {{$depTask := .}}{{range $tag, $value := $core.TaskDependencyTypes}}
                        <option value="{{ $tag }}" {{if eq $tag $depTask.DependencyType }}selected{{ end }}>{{ $value }}</option>
                        {{end}}
                    </select>
                    <button class="delete-task delete-button" onclick="deleteDepTask({{ .ID }})" {{ if not $task.EditCapability }}disabled{{ end }}>x</button>
                </div>
                {{end}}
                <br><button class="add-task" onclick="addDepTask()" {{if not .Task.EditCapability  }}disabled{{end}}>+ Add More</button><br>
            </div>
            {{end}}
            <div class="menu-item">
                <button id="save-task" onclick="{{ if .Task.Created }}updateTask();{{else}}createTask();{{end}}">Save task</button>
                {{if .Task.Created}}
                <button class="delete-button" onclick="openDelModal()">Delete task</button>
                {{end}}
            </div>
        </aside>

        <div id="trModal" class="modal">
            <div class="modal-content">
                <span class="close" onclick="closeTrModal()">&times;</span>
                <h2>Track Time</h2>
                <label for="trackTime">Tracked Time:</label>
                <input type="text" id="trackTime" name="trackedTime" placeholder="2h">
                <br>
                <label for="trackNote">Note:</label>
                <textarea id="trackNote" name="note"></textarea>
                <br>
                <button id="trackBtn" onclick="trackRecord()">Upload</button>
            </div>
        </div>

        <div id="uplModal" class="modal">
            <div class="modal-content">
                <span class="close" onclick="closeUplModal()">&times;</span>
                <h2>Upload File</h2>
                <input type="file" id="fileInput">
                <button id="uploadBtn">Upload</button>
            </div>

        </div>


        <div id="depModal" class="modal">
            <div class="modal-content">
                <span class="close" onclick="closeDepModal()">&times;</span>
                <p>Add dependent task</p>
                <label for="depTaskID">Task ID:</label>
                <input type="text" id="depTaskID" placeholder="Enter task ID">
                <br>
                <label for="dependencyType">Dependency type:</label>
                <select id="dependencyType">
                {{range $tag, $value := .Core.TaskDependencyTypes}}
                    <option value="{{ $tag }}" >{{ $value }}</option>
                {{end}}
                </select>
                <button class="accept" onclick="addDependentTask()">Add</button>
            </div>
        </div>

        <div id="lnModal" class="modal">
            <div class="modal-content">
                <span class="close" onclick="closeLnModal()">&times;</span>
                <p>Link issue</p>
                <label for="issueID">Issue ID:</label>
                <input type="text" id="issueID" placeholder="Enter issue ID">
                <button class="accept" onclick="addLinkedIssue()">Link</button>
            </div>
        </div>

        <div id="delModal" class="modal">
            <div class="modal-content">
                <p>Are you sure?</p>
                <button class="confirm" onclick="deleteTask(); closeDelModal()">Yes</button>
                <button class="cancel" onclick="closeDelModal('taskModal')">No</button>
            </div>
        </div>
        <script>
            window.jsGlobals.assigners = {{ .ProjectData.Assigners }};
            window.jsGlobals.assignees = {{ .ProjectData.Assignees }};
            window.jsGlobals.projectID = {{ .ProjectData.ID }};
            window.jsGlobals.taskID = {{ .Task.ID }};
        </script>
        <script src="/static/index.js"></script>
    </div>
</body>

</html>