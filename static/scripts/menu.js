function signinUser() {
    var datenow = new Date;
    var hour = datenow.getHours();
    var login = $('#login').val();
    var pass = $('#password').val();

    var txt = window.btoa(pass + "/" + login + "/" + hour);
    var dataVal = JSON.stringify({ "codetxt": txt });
    sendMsgonServer("/signIn", dataVal, "GET");
}

function sendMsgonServer(url, dataform, typeform){
    $.ajax({
        url: url,
        async: true,
        contentType: "application/json",
        data: dataform,
        type: typeform,
        complete: function completesend(msg, status) { alert(msg); if (status == "success") { window.location.href = "/main"; } else { $("#alert").removeClass("d-none") } }
    });
}