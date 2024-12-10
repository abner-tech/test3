function newUser() {
    fetch("http://localhost:4000/api/v1/register/user", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: 'abner',
            email: 'abner@example.com',
            password: 'password'
        })
    }).then(function(response) {
        setTimeout(() => {
            response.text().then(function(text) {
                document.getElementById("fetch-result").innerHTML = text;
            });
        }, 4000);
    }, function(err) {
        setTimeout(() => {
            document.getElementById("fetch-result").innerHTML = err;
        }, 4000);
    });
}

function activate() {
    fetch("http://localhost:4000/api/v1/users/activated", {
        method: "PUT",
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            token: 'HLFGJCTPATRSW6MSGYSE4DEFPQ'
        })
    }).then(function(response) {
        setTimeout(() => {
            response.text().then(function(text) {
                document.getElementById("fetch-result").innerHTML = text;
            });
        }, 4000);
    }, function(err) {
        setTimeout(() => {
            document.getElementById("fetch-result").innerHTML = err;
        }, 4000);
    });
}

function resetPassReq() {
    fetch("http://localhost:4000/api/v1/tokens/password-reset", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            email: 'abner@example.com',
        })
    }).then(function(response) {
        setTimeout(() => {
            response.text().then(function(text) {
                document.getElementById("fetch-result").innerHTML = text;
            });
        }, 4000);
    }, function(err) {
        setTimeout(() => {
            document.getElementById("fetch-result").innerHTML = err;
        }, 4000);
    });
}

function PassReset() {
    fetch("http://localhost:4000/api/v1/tokens/password-reset", {
        method: "PUT",
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            token: '6O6LWJKLEKFHMDHGCMQTE62CUY',
            newPassword: 'PASSWORD_CAPS'
        })
    }).then(function(response) {
        setTimeout(() => {
            response.text().then(function(text) {
                document.getElementById("fetch-result").innerHTML = text;
            });
        }, 4000);
    }, function(err) {
        setTimeout(() => {
            document.getElementById("fetch-result").innerHTML = err;
        }, 4000);
    });
}
