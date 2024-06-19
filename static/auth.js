var apiBaseLink = window.jsCoreGlobals.api_base_link;


let registrationForm = document.getElementById('registrationForm')
if (registrationForm != null) {
    registrationForm.addEventListener('submit', function (event) {
        event.preventDefault();

        var isValidForm = true;
        inputs.forEach(function (input) {
            if (!validateInput(input)) {
                isValidForm = false;
            }
        });

        if (!isValidForm) {
            alert("Data is invalid")
            return
        }

        const form = event.target;
        const userData = {
            user: {
                username: form.username.value,
                firstname: form.firstname.value,
                lastname: form.lastname.value,
                email: form.email.value,
                avatar_path: "null"
            },
            secret: form.password.value
        };

        fetch(apiBaseLink + "/api/v1/auth/with-password", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        }).then(response => {
            if (response.ok) {
                window.location.href = 'index';
            } else {
                return response.text().then(errorMessage => {
                    throw new Error(errorMessage);
                });
            }
        }).catch(error => {
            console.error('Error:', error);
            alert('Registration failed: ' + error.message);
        });
    });
}

let loginForm = document.getElementById('loginForm')
if (loginForm != null) {
    loginForm.addEventListener('submit', function (event) {
        event.preventDefault();

        var isValidForm = true;
        inputs.forEach(function (input) {
            if (!validateInput(input)) {
                isValidForm = false;
            }
        });

        if (!isValidForm) {
            alert("Data is invalid")
            return
        }

        const form = event.target;
        const userData = {
            username: form.username.value,
            secret: form.password.value
        };

        fetch(apiBaseLink + "/api/v1/auth/session", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        }).then(response => {
            if (response.ok) {
                window.location.href = 'index';
            } else {
                return response.text().then(errorMessage => {
                    throw new Error(errorMessage);
                });
            }
        }).catch(error => {
            console.error('Error:', error);
            alert('Login failed: ' + error.message);
        });
    });
}


let createProjectForm = document.getElementById('createProjectForm')
if (createProjectForm != null) {
    createProjectForm.addEventListener('submit', function (event) {
        event.preventDefault();

        var isValidForm = true;
        inputs.forEach(function (input) {
            if (!validateInput(input)) {
                isValidForm = false;
            }
        });

        if (!isValidForm) {
            alert("Data is invalid")
            return
        }

        const form = event.target;
        const userData = {
            logo_path: 'null',
            description: form.description.value,
            name: form.name.value,
        };

        fetch(apiBaseLink + "/api/v1/project", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        }).then(response => {
            if (response.ok) {
                return response.json().then(data => {
                    window.location.href = '/web/project/' + data.id;
                }).catch(error => {
                    console.error('Error:', error);
                    alert('Creation failed: ' + error.message);
                })
            } else {
                return response.text().then(errorMessage => {
                    throw new Error(errorMessage);
                });
            }
        }).catch(error => {
            console.error('Error:', error);
            alert('Login failed: ' + error.message);
        });
    });
}


var inputs = document.querySelectorAll('input');
inputs.forEach(function (input) {
    input.addEventListener('input', function () {
        validateInput(this);
    });
});


let createReportForm = document.getElementById('createReportForm')
if (createReportForm != null) {
    createReportForm.addEventListener('submit', function (event) {
        event.preventDefault();

        var isValidForm = true;
        inputs.forEach(function (input) {
            if (!validateInput(input)) {
                isValidForm = false;
            }
        });

        if (!isValidForm) {
            alert("Data is invalid")
            return
        }

        const form = event.target;
        const userData = {
            logo_path: 'null',
            description: form.description.value,
            name: form.name.value,
        };

        fetch(apiBaseLink + "/api/v1/project", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        }).then(response => {
            if (response.ok) {
                return response.json().then(data => {
                    window.location.href = '/web/project/' + data.id;
                }).catch(error => {
                    console.error('Error:', error);
                    alert('Creation failed: ' + error.message);
                })
            } else {
                return response.text().then(errorMessage => {
                    throw new Error(errorMessage);
                });
            }
        }).catch(error => {
            console.error('Error:', error);
            alert('Login failed: ' + error.message);
        });
    });
}


var inputs = document.querySelectorAll('input');
inputs.forEach(function (input) {
    input.addEventListener('input', function () {
        validateInput(this);
    });
});


function validateInput(input) {
    var errorSpan = input.nextElementSibling;
    var inputValue = input.value.trim();
    var isValid = true;

    if (inputValue === '') {
        errorSpan.textContent = 'Required field';
        isValid = false;
    } else if (input.id === 'email' && !isValidEmail(inputValue)) {
        errorSpan.textContent = 'Invalid email';
        isValid = false;
    } else if (input.id === 'password' && inputValue.length < 8) {
        errorSpan.textContent = 'Password too short';
        isValid = false;
    } else {
        errorSpan.textContent = '';
    }

    input.style.borderColor = isValid ? '' : 'red';

    return isValid;
}


function isValidEmail(email) {
    var emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    return emailRegex.test(email);
}
