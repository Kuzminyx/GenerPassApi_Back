function signinUser() {

    var login = $('#login').val();
    var pass = $('#password').val();

    //var txt = window.btoa(pass + "/" + login + "/" + hour);
    //var dataVal = JSON.stringify({ "codetxt": txt });
    var dataVal = JSON.stringify({ "user": login, "pass": pass });
    sendMsgOnServer("/signIn", dataVal, "POST");
}

function sendMsgOnServer(url, dataform, typeform){
    $.ajax({
        url: url,
        async: true,
        contentType: "application/json",
        data: dataform,
        type: typeform,
        complete: completesend(msg, status) 
    });
}

function completesend(msg, status) { 
    alert(msg); 
    if (status == "success") 
        { window.location.href = "/main"; } 
    else 
        { $("#alert").removeClass("d-none")}
}