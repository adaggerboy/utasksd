var apiBaseLink = window.jsCoreGlobals.api_base_link;

const handleFetchResponse = async (response, redirectUrl) => {
    if (response.ok) {
        window.location.href = redirectUrl;
    } else {
        const errorMessage = await response.text();
        throw new Error(errorMessage);
    }
};

const fetchApi = (endpoint, method, body = null) =>
    fetch(apiBaseLink + endpoint, {
        method,
        headers: { "Content-Type": "application/json" },
        body: body ? JSON.stringify(body) : null,
    });

function openDelModal() {
    let modal = document.getElementById("delModal");
    modal.style.display = "flex";
}

function closeDelModal() {
    let modal = document.getElementById("delModal");
    modal.style.display = "none";
}

function addFieldSuggestions(input, suggestions, list) {
    input.addEventListener("input", function () {
        const query = this.value.toLowerCase();
        suggestions.innerHTML = "";
        const filtered = Object.keys(list).filter((key) =>
            list[key].toLowerCase().includes(query)
        );

        filtered.length
            ? filtered.forEach((key) =>
                createSuggestionItem(list, key, input, suggestions)
            )
            : createNoMatchItem(suggestions);
    });

    input.addEventListener("keydown", function (e) {
        const selected = suggestions.querySelector(".selected");
        let next;
        switch (e.key) {
            case "ArrowDown":
                next = selected
                    ? selected.nextElementSibling
                    : suggestions.firstElementChild;
                break;
            case "ArrowUp":
                next = selected
                    ? selected.previousElementSibling
                    : suggestions.lastElementChild;
                break;
            case "Enter":
                if (selected)
                    assignValue(list, selected.dataset.key, input, suggestions);
                break;
            case "Escape":
                suggestions.innerHTML = "";
                break;
        }
        if (next) toggleSelectedClass(selected, next);
    });
}

function createSuggestionItem(list, key, input, suggestions) {
    const li = document.createElement("li");
    li.textContent = list[key];
    li.dataset.key = key;
    li.addEventListener("click", () =>
        assignValue(list, key, input, suggestions)
    );
    suggestions.appendChild(li);
}

function createNoMatchItem(suggestions) {
    const li = document.createElement("li");
    li.textContent = "No matches found";
    suggestions.appendChild(li);
}

function assignValue(list, key, input, suggestions) {
    input.value = list[key];
    input.keyVal = key;
    suggestions.innerHTML = "";
}

function toggleSelectedClass(selected, next) {
    if (selected) selected.classList.remove("selected");
    next.classList.add("selected");
}

const setupFieldSuggestions = (inputId, suggestionsId, list) => {
    const input = document.getElementById(inputId);
    const suggestions = document.getElementById(suggestionsId);
    if (input && suggestions) addFieldSuggestions(input, suggestions, list);
};

setupFieldSuggestions(
    "filter-registrar",
    "registrar-suggestions",
    window.jsGlobals.supports
);
setupFieldSuggestions(
    "filter-assigner-input",
    "filter-assigner-suggestions",
    window.jsGlobals.assigners
);
setupFieldSuggestions(
    "assigner-input2",
    "assigner-suggestions2",
    window.jsGlobals.assigners
);
setupFieldSuggestions(
    "filter-assignee-input",
    "filter-assignee-suggestions",
    window.jsGlobals.assignees
);
setupFieldSuggestions(
    "assignee-input2",
    "assignee-suggestions2",
    window.jsGlobals.assignees
);

let editMode = false;

function toggleEditMode(elementId) {
    if (editMode == true) {
        return;
    }
    editMode = true;
    const element = document.getElementById(elementId);
    const currentText = element.innerText;
    const elementRect = element.getBoundingClientRect();

    const textarea = document.createElement("textarea");
    textarea.value = currentText;
    textarea.classList.add("editable-textarea");
    textarea.style.width = elementRect.width + "px";
    textarea.style.height = elementRect.height + "px";

    element.innerHTML = "";
    element.appendChild(textarea);

    textarea.focus();

    textarea.addEventListener("blur", () => {
        saveChanges(element, textarea);
    });
}

function saveChanges(element, textarea) {
    element.innerHTML = textarea.value;
    editMode = false;
}

function getTaskInformation() {
    const taskName = document.getElementById("taskName").innerText.trim();
    const taskDescription = document
        .getElementById("taskDescription")
        .innerText.trim();
    const status = document.getElementById("status").value;
    const priority = document.getElementById("priority").value;
    const assigner = parseInt(document.getElementById("assigner-input2").keyVal);
    const assignee = parseInt(document.getElementById("assignee-input2").keyVal);
    const startDate = document.getElementById("start-date").value.trim();
    const dueDate = document.getElementById("due-date").value.trim();

    const dependentTasks = Array.from(
        document.querySelectorAll(".dependent-tasks .task")
    ).map((taskElement) => {
        const taskId = parseInt(taskElement.querySelector(".task-id").innerText);
        const taskName = taskElement.querySelector(".task-name").innerText.trim();
        const dependencyType = taskElement.querySelector(".dependency-type").value;
        return { id: taskId, name: taskName, dependency_type: dependencyType };
    });

    const linkedIssues = Array.from(
        document.querySelectorAll(".linked-issues .issue")
    ).map((issueElement) => {
        const issueId = parseInt(issueElement.querySelector(".issue-id").innerText);
        return issueId;
    });

    const attachments = Array.from(
        document.querySelectorAll(".task-block .image-wrapper img")
    ).map((imgElement) => {
        const url = new URL(imgElement.src);
        return url.pathname; // Get the file path from the src attribute
    });

    const task = {
        name: taskName || undefined,
        description: taskDescription || undefined,
        status: status || undefined,
        project: window.jsGlobals.projectID,
        priority: priority || undefined,
        assigner: assigner || undefined,
        assignees: assignee ? [assignee] : undefined, // If assignee is not empty, create array
        start_date: startDate || undefined,
        due_date: dueDate || undefined,
        dependent_tasks: dependentTasks,
        linked_issues: linkedIssues,
        attachments: attachments,
    };

    Object.keys(task).forEach(
        (key) => task[key] === undefined && delete task[key]
    );

    return task;
}

function addDepTask() {
    document.getElementById("depModal").style.display = "block";
}

function closeDepModal() {
    document.getElementById("depModal").style.display = "none";
}

function addDepTaskToList(taskId, taskName, dependencyType) {
    const dependentTasksList = document.querySelector(".dependent-tasks");
    const taskElement = document.createElement("div");
    taskElement.classList.add("task");
    taskElement.setAttribute("data-task-id", taskId);
    taskElement.innerHTML = `
        <a href="#"><span class="task-id">${taskId}</span><span class="task-name">${taskName}</span></a>
        <select class="dependency-type"><option value="${dependencyType}" selected>${dependencyType}</option></select>
        <button class="delete-task delete-button" onclick="deleteDepTask(${taskId})">x</button>`;
    dependentTasksList.appendChild(taskElement);
}

function openuplModal(id) {
    document.getElementById("uplModal").style.display = "block";
}

function closeUplModal(id) {
    document.getElementById("uplModal").style.display = "none";
}

function handleUpload(id, handler) {
    let element = document.getElementById(id)
    if(element == null) return;
    element.addEventListener("click", function () {
        var fileInput = document.getElementById('fileInput');
        var file = fileInput.files[0];
        var formData = new FormData();
        formData.append('file', file);
    
        fetch(apiBaseLink + '/files/upload', {
            method: 'POST',
            body: formData
        })
            .then(response => {
                if (response.ok) {
                    return response.json();
                } else {
                    return response.text().then(errorMessage => {
                        throw new Error(errorMessage);
                    });
                }
            })
            .then(data => {
                var filePath = data.path;
                handler(filePath);
                closeUplModal();
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Upload failed: ' + error.message);
            });
    });
}

function addImageToAttachments(filePath) {
    const taskBlock = document.getElementById("task-central-block");
    const imgWrapper = document.createElement("div");
    imgWrapper.classList.add("image-wrapper");

    const img = document.createElement("img");
    img.src = filePath;
    img.setAttribute("onclick", "deleteAttachment(this)");

    

    const deleteBtn = document.createElement("button");
    deleteBtn.textContent = "Delete";
    deleteBtn.classList.add("delete-btn");
    deleteBtn.classList.add("delete-button");
    imgWrapper.appendChild(img);
    imgWrapper.appendChild(deleteBtn);
    taskBlock.appendChild(imgWrapper);
}

function deleteAttachment(buttonElement) {
    const imageWrapper = buttonElement.closest(".image-wrapper");
    if (imageWrapper) {
        imageWrapper.remove();
    }
}

function logout() {
    fetch(apiBaseLink + "/api/v1/auth/session", {
        method: "DELETE",
    }).then((nil) => {
        window.location.href = "/web/index";
    });
}

function openTrModal() {
    document.getElementById("trModal").style.display = "block";
}

function closeTrModal() {
    document.getElementById("trModal").style.display = "none";
}

function getIssueInformation() {
    const taskName = document.getElementById("taskName").innerText.trim();
    const status = document.getElementById("status").value;
    const taskDescription = document
        .getElementById("taskDescription")
        .innerText.trim();
    const reporter = parseInt(document.getElementById("reporter").keyVal);

    const attachments = Array.from(
        document.querySelectorAll(".task-block .image-wrapper img")
    ).map((imgElement) => {
        const url = new URL(imgElement.src);
        return url.pathname; // Get the file path from the src attribute
    });

    const issue = {
        name: taskName || undefined,
        description: taskDescription || undefined,
        status: status,
        project: window.jsGlobals.projectID,
        reporter: reporter.length > 0 ? reporter : undefined,
        attachments: attachments,
    };

    Object.keys(issue).forEach(
        (key) => issue[key] === undefined && delete issue[key]
    );

    return issue;
}

function updateTask() {
    const taskID = window.jsGlobals.taskID;
    fetchApi(`/api/v1/task/${taskID}`, "PUT", getTaskInformation())
        .then((response) => handleFetchResponse(response, `/web/task/${taskID}`))
        .catch((error) => {
            console.error("Error:", error);
            alert("Update task failed: " + error.message);
        });
}

function deleteTask() {
    const taskID = window.jsGlobals.taskID;
    fetchApi(`/api/v1/task/${taskID}`, "DELETE")
        .then((response) => handleFetchResponse(response, "/web/index"))
        .catch((error) => {
            console.error("Error:", error);
            alert("Delete task failed: " + error.message);
        });
}

function createTask() {
    fetchApi("/api/v1/task/", "POST", getTaskInformation())
        .then((response) => {
            if (response.ok) {
                return response
                    .json()
                    .then((data) => {
                        if (data.id) {
                            window.location.href = "/web/task/" + data.id;
                        } else {
                            throw new Error("No task id returned from server");
                        }
                    })
                    .catch((errorMessage) => {
                        throw new Error(errorMessage);
                    });
            } else {
                return response.text().then((errorMessage) => {
                    throw new Error(errorMessage);
                });
            }
        })
        .catch((error) => {
            console.error("Error:", error);
            alert("Create task failed: " + error.message);
        });
}

function trackRecord() {
    const taskID = window.jsGlobals.taskID;
    const trackedTime = document.getElementById("trackTime").value;
    const note = document.getElementById("trackNote").value;
    const data = { text: note, duration: trackedTime };

    fetchApi(`/api/v1/task/${taskID}/tracking-record`, "POST", data)
        .then((response) => handleFetchResponse(response, `/web/task/${taskID}`))

        .catch((error) => {
            console.error("Error:", error);
            alert("Track record failed: " + error.message);
        });
}

function updateIssue() {
    const issueID = window.jsGlobals.issueID;
    const userData = getIssueInformation();

    fetchApi(`/api/v1/issue/${issueID}`, "PUT", userData)
        .then((response) => handleFetchResponse(response, `/web/issue/${issueID}`))

        .catch((error) => {
            console.error("Error:", error);
            alert("Update issue failed: " + error.message);
        });
}

function deleteIssue() {
    const issueID = window.jsGlobals.issueID;

    fetchApi(`/api/v1/issue/${issueID}`, "DELETE")
        .then((response) => handleFetchResponse(response, "/web/index"))
        .catch((error) => {
            console.error("Error:", error);
            alert("Delete issue failed: " + error.message);
        });
}

function createIssue() {
    const userData = getIssueInformation();

    fetchApi("/api/v1/issue/", "POST", userData)
        .then((response) => {
            if (response.ok) {
                return response
                    .json()
                    .then((data) => {
                        if (data.id) {
                            window.location.href = "/web/issue/" + data.id;
                        } else {
                            throw new Error("No task id returned from server");
                        }
                    })
                    .catch((errorMessage) => {
                        throw new Error(errorMessage);
                    });
            } else {
                return response.text().then((errorMessage) => {
                    throw new Error(errorMessage);
                });
            }
        })
        .catch((error) => {
            console.error("Error:", error);
            alert("Create issue failed: " + error.message);
        });
}
function uploadComment() {
    const taskID = window.jsGlobals.taskID;
    const newCommentText = document
        .getElementById("new-comment-text")
        .value.trim();

    if (newCommentText) {
        fetchApi("/api/v1/comment", "POST", {
            task_id: taskID,
            text: newCommentText,
        })
        .then((response) => handleFetchResponse(response, `/web/task/${taskID}`))
        .catch((error) => {
            console.error("Error:", error);
            alert("Upload failed: " + error.message);
        });
    }
}

function deleteComment(buttonElement) {
    const taskID = window.jsGlobals.taskID;
    const commentElement = buttonElement.closest(".comment");
    const commentID = commentElement.getAttribute("data-comment-id");

    fetchApi(`/api/v1/comment/${commentID}`, "DELETE")
        .then((response) => handleFetchResponse(response, `/web/task/${taskID}`))
        .catch((error) => {
            console.error("Error:", error);
            alert("Delete comment failed: " + error.message);
        });
}

function searchTasks() {
    const getFilterValue = (id) => document.getElementById(id)?.value.trim();
    const viewOption = getFilterValue("filter-view");
    const nameFilter = getFilterValue("filter-name");
    const assigneeFilter = parseInt(
        document.getElementById("filter-assignee-input")?.keyVal
    );
    const assignerFilter = parseInt(
        document.getElementById("filter-assigner-input")?.keyVal
    );
    const statusFilter = getFilterValue("filter-status");
    const priorityFilter = getFilterValue("filter-priority");

    const queryParams = new URLSearchParams({
        ...(nameFilter && { name: nameFilter }),
        ...(assigneeFilter &&
            !isNaN(assigneeFilter) && { assignee: assigneeFilter }),
        ...(assignerFilter &&
            !isNaN(assignerFilter) && { assigner: assignerFilter }),
        ...(statusFilter !== "any" && { status: statusFilter }),
        ...(priorityFilter !== "any" && { priority: priorityFilter }),
        project: window.jsGlobals.projectID,
    }).toString();

    window.location.href = `/web/search_tasks/${viewOption}?${queryParams}`;
}


const fetchIssueName = async (id) =>
    fetchApi(`/api/v1/issue/${id}`, "GET")
        .then((response) => response.json())
        .then((data) => data.name);

const fetchTaskName = async (id) =>
    fetchApi(`/api/v1/task/${id}`, "GET")
        .then((response) => response.json())
        .then((data) => data.name);

function addDependentTask() {
    const taskId = document.getElementById("depTaskID").value;
    const dependencyType = document.getElementById("dependencyType").value;
    fetchTaskName(taskId)
        .then((name) => {
            addDepTaskToList(taskId, name, dependencyType);
            closeDepModal();
        })
        .catch((error) => {
            console.error("Error:", error);
            alert("Fetch task name failed: " + error.message);
        });
}

function deleteDepTask(taskId) {
    document
        .querySelector(`.dependent-tasks .task[data-task-id="${taskId}"]`)
        .remove();
}

function linkIssue() {
    document.getElementById("lnModal").style.display = "block";
}

function closeLnModal() {
    document.getElementById("lnModal").style.display = "none";
}

function addLnIssToList(issueID, issueName, dependencyType) {
    const list = document.querySelector(".linked-issues");
    const taskElement = document.createElement("div");
    taskElement.classList.add("issue");
    taskElement.setAttribute("data-issue-id", issueID);
    taskElement.innerHTML = `
        <a href="/web/issue/${issueID}"><span class="issue-id">${issueID}</span><span class="issue-name">${issueName}</span></a>
        <button class="delete-issue delete-button" onclick="deleteLinkedIssue(${issueID})">x</button>`;
    list.appendChild(taskElement);
}

function addLinkedIssue() {
    const issueID = document.getElementById("issueID").value;
    fetchIssueName(issueID).then((name) => {
        addLnIssToList(issueID, name);
        closeDepModal();
    });
}


function deleteLinkedIssue(issueID) {
    document
        .querySelector(`.linked-issues .issue[data-issue-id="${issueID}"]`)
        .remove();
}

function closeIssue() {
    document.getElementById("status").value = 'closed';
    updateIssue();
}
function reopenIssue() {
    document.getElementById("status").value = 'reopened';
    updateIssue();
}


function searchIssues() {
    const getFilterValue = (id) => document.getElementById(id)?.value.trim();
    const nameFilter = getFilterValue("filter-name");
    const registrarFilter = parseInt(
        document.getElementById("filter-registrar")?.keyVal
    );
    const statusFilter = getFilterValue("filter-status");
    const queryParams = new URLSearchParams({
        ...(nameFilter && { name: nameFilter }),
        ...(registrarFilter &&
            !isNaN(registrarFilter) && { registrar: registrarFilter }),
        ...(statusFilter !== "any" && { status: statusFilter }),
        project: window.jsGlobals.projectID,
    }).toString();

    window.location.href = `/web/search_issues?${queryParams}`;
}


function saveUser(cell) {
    let foundRow = cell.closest('tr');

    let id = parseInt(foundRow.querySelector('td:nth-child(2)').innerHTML);

    const newPassword = foundRow.querySelector('td:nth-child(7) input[type="password"]').value;
    const isDirector = foundRow.querySelector('td:nth-child(8) input[type="checkbox"]').checked;
    const isAdmin = foundRow.querySelector('td:nth-child(9) input[type="checkbox"]').checked;
    const isActive = foundRow.querySelector('td:nth-child(10) input[type="checkbox"]').checked;

    if (newPassword !== '') {

        fetchApi(`/api/v1/auth/${id}/credentials`, "PUT", { secret: newPassword })
        .then((response) => handleFetchResponse(response, "/web/users"))
        .catch((error) => {
            console.error("Error:", error);
            alert("User credentials update failed: " + error.message);
        });
    }

    fetchApi(`/api/v1/auth/${id}/permissions`, "PUT", {
        is_active: isActive,
        is_director: isDirector,
        is_admin: isAdmin
    })
    .then((response) => handleFetchResponse(response, "/web/users"))
    .catch((error) => {
        console.error("Error:", error);
        alert("User permissions update failed: " + error.message);
    });

}


function deleteProject() {
    const projectID = window.jsGlobals.projectID;
    fetchApi(`/api/v1/project/${projectID}`, "DELETE")
        .then((response) => handleFetchResponse(response, "/web/index"))
        .catch((error) => {
            console.error("Error:", error);
            alert("Delete project failed: " + error.message);
        });
}

function deleteUser(button) {
    const row = button.closest('tr');
    row.remove();
}

const fetchUserData = async (id) =>
    fetchApi(`/api/v1/user/${id}`, "GET")
        .then((response) => response.json());


function addUserToTable(id, avatar, firstName, lastName, email) {
    let tableBody = document.querySelectorAll('#user-table-body tbody')[0];
    let row = document.createElement('tr');
    row.setAttribute('data-user-id', id);
    row.innerHTML = `
        <td class="table-avatar"><img src="${avatar}" alt="User Avatar"><span>${firstName} ${lastName}</span></td>
        <td>${email}</td>
        <td>
            <select value="worker">
                <option value="manager">Manager</option>
                <option value="worker" selected>Worker</option>
                <option value="support">Support agent</option>
            </select>
        </td>
        <td><button class="delete-button" onclick="deleteUser(this)">Delete</button></td>
    `;
    tableBody.appendChild(row);
}

async function addUser() {
    const userIdInput = document.getElementById('user-id-input');
    const userId = userIdInput.value.trim();

    if (!userId) {
        alert('Please enter a user ID.');
        return;
    }

    const existingUser = document.querySelector(`tr[data-user-id="${userId}"]`);
    if (existingUser) {
        alert('User already exists in the table.');
        return;
    }

    fetchUserData(userId).then((userData) => {
        addUserToTable(userId, userData.avatar_path, userData.firstname, userData.lastname, userData.email);
    })
    .catch((error) => {
        console.error("Error:", error);
        alert("Fetch user name failed: " + error.message);
    });
}

function updateProject() {
    const projectID = window.jsGlobals.projectID;

    const projectName = document.getElementById('project-name').value;
    const projectLogo = document.getElementById('logo').src;
    const projectDescription = document.querySelector('.project-info p').textContent;

    const managers = [];
    const workers = [];
    const supports = [];

    document.querySelectorAll('#user-table-body tbody tr').forEach(row => {
        const userId = parseInt(row.getAttribute('data-user-id'));
        const role = row.querySelector('select').value;

        if (role === 'manager') {
            managers.push(userId);
        } else if (role === 'worker') {
            workers.push(userId);
        } else if (role === 'support') {
            supports.push(userId);
        }
    });

    fetchApi(`/api/v1/project/${projectID}`, "PUT", {
        name: projectName,
        logo_path: new URL(projectLogo).pathname,
        description: projectDescription,
        managers: managers,
        workers: workers,
        support_agents: supports,
    })
    .then((response) => handleFetchResponse(response, `/web/project/${projectID}`))
    .catch((error) => {
        console.error("Error:", error);
        alert("Update project failed: " + error.message);
    });
    
}

function updateUser() {
    const userID = window.jsGlobals.userID;

    const userAva = document.getElementById('logo').src;

    fetchApi(`/api/v1/user`, "PUT", {
        avatar_path: new URL(userAva).pathname,
    })
    .then((response) => handleFetchResponse(response, `/web/user/${userID}`))
    .catch((error) => {
        console.error("Error:", error);
        alert("Update user failed: " + error.message);
    });
    
}

function changeAva(filePath) {
    document.getElementById("logo").src = filePath
}

handleUpload("uploadBtn", addImageToAttachments);
handleUpload("avaBtn", changeAva);


function deactivateUser() {
    fetchApi(`/api/v1/user`, "DELETE")
        .then((response) => handleFetchResponse(response, "/web/index"))
        .catch((error) => {
            console.error("Error:", error);
            alert("Delete user failed: " + error.message);
        });
}