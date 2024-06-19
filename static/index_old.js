var apiBaseLink = window.jsCoreGlobals.api_base_link;


function openDelModal() {
    let modal = document.getElementById('delModal');
    modal.style.display = "flex";
}

function closeDelModal() {
    let modal = document.getElementById('delModal');
    modal.style.display = "none";
}


function addFieldSuggestions(assignerInput, assignerSuggestions, assigners) {

    assignerInput.addEventListener('input', function () {

        const query = this.value.toLowerCase();
        assignerSuggestions.innerHTML = '';

        if (query) {
            const filteredAssigners = Object.keys(assigners).filter(key =>
                assigners[key].toLowerCase().includes(query)
            );

            filteredAssigners.forEach(key => {
                const li = document.createElement('li');
                li.textContent = assigners[key];
                li.addEventListener('click', function () {
                    assignerInput.value = assigners[key];
                    assignerInput.keyVal = key;
                    assignerSuggestions.innerHTML = '';
                });
                assignerSuggestions.appendChild(li);
            });

            if (filteredAssigners.length === 0) {
                const li = document.createElement('li');
                li.textContent = 'No matches found';
                assignerSuggestions.appendChild(li);
            }
        }
    });

    assignerInput.addEventListener('keydown', function (e) {
        const selected = assignerSuggestions.querySelector('.selected');
        let next;

        switch (e.key) {
            case 'ArrowDown':
                if (selected) {
                    next = selected.nextElementSibling;
                } else {
                    next = assignerSuggestions.firstElementChild;
                }
                if (next) {
                    if (selected) selected.classList.remove('selected');
                    next.classList.add('selected');
                }
                break;
            case 'ArrowUp':
                if (selected) {
                    next = selected.previousElementSibling;
                } else {
                    next = assignerSuggestions.lastElementChild;
                }
                if (next) {
                    if (selected) selected.classList.remove('selected');
                    next.classList.add('selected');
                }
                break;
            case 'Enter':
                if (selected) {
                    const key = selected.dataset.key;
                    assignerInput.value = assigners[key];
                    assignerSuggestions.innerHTML = '';
                }
                break;
            case 'Escape':
                assignerSuggestions.innerHTML = '';
                break;
        }
    });
}


const assignerInput = document.getElementById('filter-assigner-input');
const assignerInput2 = document.getElementById('assigner-input2');
const assignerSuggestions = document.getElementById('filter-assigner-suggestions');
const assignerSuggestions2 = document.getElementById('assigner-suggestions2');

const assigneeInput = document.getElementById('filter-assignee-input');
const assigneeInput2 = document.getElementById('assignee-input2');
const assigneeSuggestions = document.getElementById('filter-assignee-suggestions');
const assigneeSuggestions2 = document.getElementById('assignee-suggestions2');

if(assigneeInput != null && assigneeSuggestions != null) {
    addFieldSuggestions(assigneeInput, assigneeSuggestions, window.jsGlobals.assignees);
}
if(assigneeInput2 != null && assigneeSuggestions2 != null) {
    addFieldSuggestions(assigneeInput2, assigneeSuggestions2, window.jsGlobals.assignees);
}
if(assignerInput != null && assignerSuggestions != null) {
    addFieldSuggestions(assignerInput, assignerSuggestions, window.jsGlobals.assigners);
}
if(assigneeInput2 != null && assignerSuggestions2 != null) {
    addFieldSuggestions(assignerInput2, assignerSuggestions2, window.jsGlobals.assigners);
}

let editMode = false;

function toggleEditMode(elementId) {
    if (editMode == true) {
        return;
    }
    editMode = true;
    const element = document.getElementById(elementId);
    const currentText = element.innerText;
    const elementRect = element.getBoundingClientRect();

    const textarea = document.createElement('textarea');
    textarea.value = currentText;
    textarea.classList.add('editable-textarea');
    textarea.style.width = elementRect.width + 'px';
    textarea.style.height = elementRect.height + 'px';

    element.innerHTML = '';
    element.appendChild(textarea);

    textarea.focus();

    textarea.addEventListener('blur', () => {
        saveChanges(element, textarea);
    });
}

function saveChanges(element, textarea) {
    const newText = textarea.value;
    element.innerHTML = newText;
    editMode = false;
}

function getTaskInformation() {
    const taskName = document.getElementById('taskName').innerText.trim();
    const taskDescription = document.getElementById('taskDescription').innerText.trim();
    const status = document.getElementById('status').value;
    const priority = document.getElementById('priority').value;
    const assigner = parseInt(document.getElementById('assigner-input2').keyVal);
    const assignee = parseInt(document.getElementById('assignee-input2').keyVal);
    const startDate = document.getElementById('start-date').value.trim();
    const dueDate = document.getElementById('due-date').value.trim();

    const dependentTasks = Array.from(document.querySelectorAll('.dependent-tasks .task')).map(taskElement => {
        const taskId = parseInt(taskElement.querySelector('.task-id').innerText);
        const taskName = taskElement.querySelector('.task-name').innerText.trim();
        const dependencyType = taskElement.querySelector('.dependency-type').value;
        return { id: taskId, name: taskName, dependency_type: dependencyType };
    });

    const linkedIssues = Array.from(document.querySelectorAll('.linked-issues .issue')).map(issueElement => {
        const issueId = parseInt(issueElement.querySelector('.issue-id').innerText);
        const issueName = issueElement.querySelector('.issue-name').innerText.trim();
        return { id: issueId, name: issueName };
    });

    const attachments = Array.from(document.querySelectorAll('.task-block .image-wrapper img')).map(imgElement => {
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
        attachments: attachments
    };

    Object.keys(task).forEach(key => task[key] === undefined && delete task[key]);

    return task;
}

function updateTask() {
    const taskID = window.jsGlobals.taskID;
    const userData = getTaskInformation();

    fetch(apiBaseLink + "/api/v1/task/" + taskID, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData)
    }).then(response => {
        if (response.ok) {
            window.location.href = "/web/task/" + taskID;
        } else {
            return response.text().then(errorMessage => {
                throw new Error(errorMessage);
            });
        }
    }).catch(error => {
        console.error('Error:', error);
        alert('Update failed: ' + error.message);
    });
}

function deleteTask() {
    const taskID = window.jsGlobals.taskID;

    fetch(apiBaseLink + "/api/v1/task/" + taskID, {
        method: 'DELETE',
    }).then(response => {
        if (response.ok) {
            window.location.href = "/web/index";
        } else {
            return response.text().then(errorMessage => {
                throw new Error(errorMessage);
            });
        }
    }).catch(error => {
        console.error('Error:', error);
        alert('Delete failed: ' + error.message);
    });
}

function createTask() {
    const userData = getTaskInformation();

    fetch(apiBaseLink + "/api/v1/task/", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData)
    }).then(response => {
        if (response.ok) {
            return response.json().then(data => {
                if (data.id) {
                    window.location.href = "/web/task/" + data.id;
                } else {
                    throw new Error("No task id returned from server");
                }
            }).catch(errorMessage => {
                throw new Error(errorMessage);
            });
        } else {
            return response.text().then(errorMessage => {
                throw new Error(errorMessage);
            });
        }
    }).catch(error => {
        console.error('Error:', error);
        alert('Update failed: ' + error.message);
    });
}
function searchTasks() {
    const viewOption = document.getElementById('filter-view').value;
    const nameFilter = document.getElementById('filter-name').value.trim();
    let assigneeFilter = null;
    if (document.getElementById('filter-assignee-input') !== null) {
        assigneeFilter = parseInt(document.getElementById('filter-assignee-input').keyVal);
    }
    let assignerFilter = null;
    if (document.getElementById('filter-assigner-input') !== null) {
        assignerFilter = parseInt(document.getElementById('filter-assigner-input').keyVal);
    }
    const statusFilter = document.getElementById('filter-status').value;
    const priorityFilter = document.getElementById('filter-priority').value;
    // const isRelevantFilter = document.getElementById('filter-is-relevant').checked;
    // const isWarningFilter = document.getElementById('filter-is-warning').checked;

    const queryParams = new URLSearchParams();
    if (nameFilter !== '') queryParams.append('name', nameFilter);
    if (assigneeFilter !== null && !isNaN(assigneeFilter)) queryParams.append('assignee', assigneeFilter);
    if (assignerFilter !== null && !isNaN(assignerFilter)) queryParams.append('assigner', assignerFilter);
    if (statusFilter !== 'any') queryParams.append('status', statusFilter);
    if (priorityFilter !== 'any') queryParams.append('priority', priorityFilter);
    queryParams.append('project', window.jsGlobals.projectID)
    // if (isRelevantFilter) queryParams.append('isRelevant', isRelevantFilter);
    // if (isWarningFilter) queryParams.append('isWarning', isWarningFilter);

    const endpoint = `/web/search_tasks/${viewOption}?${queryParams.toString()}`;

    window.location.href = endpoint;
}


function addDepTask() {
    document.getElementById('depModal').style.display = 'block';
}

function closeDepModal() {
    document.getElementById('depModal').style.display = 'none';
}

async function fetchTaskName(id) {
    return fetch(apiBaseLink + `/api/v1/task/${id}`)
    .then(response => response.json())
    .then(data => {
        return data.name;
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Get task name failed: ' + error.message);
    });
}

function addDepTaskToList(taskId, taskName, dependencyType) {
    const dependentTasksList = document.querySelector('.dependent-tasks');
    const taskElement = document.createElement('div');
    taskElement.classList.add('task');
    taskElement.setAttribute('data-task-id', taskId);
    taskElement.innerHTML = `
        <a href="#">
            <span class="task-id">${taskId}</span>
            <span class="task-name">${taskName}</span>
        </a>
        <select class="dependency-type">
            <option value="${dependencyType}" selected>${dependencyType}</option>
        </select>
        <button class="delete-task delete-button" onclick="deleteDepTask(${taskId})">x</button>
    `;
    dependentTasksList.appendChild(taskElement);
}

function addDependentTask() {
    const taskId = document.getElementById('depTaskID').value;
    const dependencyType = document.getElementById('dependencyType').value;
    fetchTaskName(taskId).then(name => {
        addDepTaskToList(taskId, name, dependencyType);
        closeModal();
    })
    
}

function deleteDepTask(taskId) {
    const taskElement = document.querySelector(`.task[data-task-id="${taskId}"]`);
    if (taskElement) {
        taskElement.remove();
    } else {
        console.error('Task element not found.');
    }
}



function openuplModal() {
    document.getElementById('uplModal').style.display = 'block';
}

function closeUplModal() {
    document.getElementById('uplModal').style.display = 'none';
}

// Upload function
document.getElementById("uploadBtn").addEventListener("click", function(){
    var fileInput = document.getElementById('fileInput');
    var file = fileInput.files[0];
    var formData = new FormData();
    formData.append('file', file);

    fetch(apiBaseLink + '/files/upload', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if(response.ok) {
            return response.json();
        } else {
            return response.text().then(errorMessage => {
                throw new Error(errorMessage);
            });
        }
    })
    .then(data => {
        var filePath = data.path;
        addImageToAttachments(filePath);
        closeUplModal();
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Upload failed: ' + error.message);
    });
});


function addImageToAttachments(filePath) {
    const taskBlock = document.getElementById('task-central-block')
    const imgWrapper = document.createElement('div');
    imgWrapper.classList.add('image-wrapper');

    const img = document.createElement('img');
    img.src = filePath;

    const deleteBtn = document.createElement('button');
    deleteBtn.textContent = 'Delete';
    deleteBtn.classList.add('delete-btn');
    deleteBtn.classList.add('delete-button');
    imgWrapper.appendChild(img);
    imgWrapper.appendChild(deleteBtn);
    taskBlock.appendChild(imgWrapper);
}

function deleteAttachment(buttonElement) {
    const imageWrapper = buttonElement.closest('.image-wrapper');
    if (imageWrapper) {
        imageWrapper.remove();
    }
}


function uploadComment() {
    const taskID = window.jsGlobals.taskID;
    const newCommentText = document.getElementById('new-comment-text').value.trim();
    if (newCommentText) {
        fetch('/api/v1/comment', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                task_id: taskID,
                text: newCommentText
            })
        })
        .then(response => {
            if(response.ok) {
                window.location.reload();
            } else {
                return response.text().then(errorMessage => {
                    throw new Error(errorMessage);
                });
            }
        }).catch(error => {
            console.error('Error:', error);
            alert('Upload failed: ' + error.message);
        });
    }
}

function deleteComment(buttonElement) {
    const commentElement = buttonElement.closest('.comment');
    const commentID = commentElement.getAttribute('data-comment-id');

    fetch(apiBaseLink + `/api/v1/comment/${commentID}`, {
        method: 'DELETE'
    })
    .then(response => {
        if(response.ok) {
            window.location.reload();
        } else {
            return response.text().then(errorMessage => {
                throw new Error(errorMessage);
            });
        }
    }).catch(error => {
        console.error('Error:', error);
        alert('Upload failed: ' + error.message);
    });
}

function logout() {
    fetch(apiBaseLink + "/api/v1/auth/session", {
        method: 'DELETE',
    }).then(nil => {
        window.location.href = "/web/index"
    })
}

function openTrModal() {
    document.getElementById("trModal").style.display = "block";
}

function closeTrModal() {
    document.getElementById("trModal").style.display = "none";
}

function trackRecord() {
    let taskID = window.jsGlobals.taskID;
    var trackedTime = document.getElementById("trackTime").value;
    var note = document.getElementById("trackNote").value;
    var data = {
        "text": note,
        "duration": trackedTime
    };
    fetch(apiBaseLink + "/api/v1/task/" + taskID + "/tracking-record", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if(response.ok) {
            window.location.reload();
        } else {
            return response.text().then(errorMessage => {
                throw new Error(errorMessage);
            });
        }
    }).catch(error => {
        console.error('Error:', error);
        alert('Upload failed: ' + error.message);
    });
}


function updateIssue() {
    const taskID = window.jsGlobals.issueID;
    const userData = getIssueInformation();

    fetch(apiBaseLink + "/api/v1/issue/" + taskID, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData)
    }).then(response => {
        if (response.ok) {
            window.location.href = "/web/issue/" + taskID;
        } else {
            return response.text().then(errorMessage => {
                throw new Error(errorMessage);
            });
        }
    }).catch(error => {
        console.error('Error:', error);
        alert('Update failed: ' + error.message);
    });
}

function deleteIssue() {
    const taskID = window.jsGlobals.issueID;

    fetch(apiBaseLink + "/api/v1/issue/" + taskID, {
        method: 'DELETE',
    }).then(response => {
        if (response.ok) {
            window.location.href = "/web/index";
        } else {
            return response.text().then(errorMessage => {
                throw new Error(errorMessage);
            });
        }
    }).catch(error => {
        console.error('Error:', error);
        alert('Delete failed: ' + error.message);
    });
}

function createIssue() {
    const userData = getIssueInformation();

    fetch(apiBaseLink + "/api/v1/issue/", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData)
    }).then(response => {
        if (response.ok) {
            return response.json().then(data => {
                if (data.id) {
                    window.location.href = "/web/issue/" + data.id;
                } else {
                    throw new Error("No issue id returned from server");
                }
            }).catch(errorMessage => {
                throw new Error(errorMessage);
            });
        } else {
            return response.text().then(errorMessage => {
                throw new Error(errorMessage);
            });
        }
    }).catch(error => {
        console.error('Error:', error);
        alert('Create failed: ' + error.message);
    });
}


function getIssueInformation() {
    const taskName = document.getElementById('taskName').innerText.trim();
    const taskDescription = document.getElementById('taskDescription').innerText.trim();
    const reporter = parseInt(document.getElementById('reporter').keyVal);

    const attachments = Array.from(document.querySelectorAll('.task-block .image-wrapper img')).map(imgElement => {
        const url = new URL(imgElement.src);
        return url.pathname; // Get the file path from the src attribute
    });

    const issue = {
        name: taskName || undefined,
        description: taskDescription || undefined,
        project: window.jsGlobals.projectID,
        reporter: reporter.length > 0 ? reporter : undefined,
        attachments: attachments
    };

    Object.keys(issue).forEach(key => issue[key] === undefined && delete issue[key]);

    return issue;
}
