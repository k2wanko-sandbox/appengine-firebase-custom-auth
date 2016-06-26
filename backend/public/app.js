function createUser(email, password) {
    const data = new FormData();
    data.append("email", email);
    data.append("password", password);
    return fetch("/register", {
        method: "POST",
        body: data
    }).then(x => x.json()).then(d => {
        const token = d.Token;
        return firebase.auth()
            .signInWithCustomToken(token)
            .then(user => {
                console.log("uid: ", user.uid);
                return user.updateEmail(email).then(() => user);
            });
    });
}

window.addEventListener("DOMContentLoaded", () => {
    const el_register = document.getElementById("register_submit");
    //const el_login = document.getElementById("login_submit");
    const el_uid = document.getElementById("uid");
    const el_email = document.getElementById("email");
    const el_password = document.getElementById("password");

    el_uid.innerText = "None";

    el_register.onclick = () => {
        const email = el_email.value;
        const password = el_password.value;
        console.log("email = %s, password = %s", email, password);
        createUser(email, password).then(user => {
            el_uid.innerText = user.uid;
        });
        return false;
    };
    
}, false);
