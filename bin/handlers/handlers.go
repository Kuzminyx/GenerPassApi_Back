package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	gsql "github.com/GeneratePassAPI/bin/bd"
	"github.com/GeneratePassAPI/bin/utility"
)

//User - объект для работы с пользователями с помощью JSON ...
type User struct {
	UID          string
	Login        string `json:"user"`
	Pass         string `json:"pass"`
	Email        string `json:"mail"`
	ActivateTime time.Time
}

//Gdata - объект для хранения клиентских вложенных данных ...
type Gdata struct {
	Login   string `json:"login"`
	Pass    string `json:"pass"`
	Email   string `json:"mail"`
	Address string `json:"address"`
}

//JSONclient - объект для работы с JSON ...
type JSONclient struct {
	User      string
	ID        string `json:"id"`
	Operation string `json:"proc"`
	Client    Gdata  `json:"data"`
}

var active = make(map[string]User)

//Getid - получить новый ИД ...
func Getid(w http.ResponseWriter, r *http.Request) {

	resp := utility.Response{Status: "Ok", IsArray: false}
	user := User{}
	jsonquery := json.NewDecoder(r.Body)

	alert := jsonquery.Decode(&user)

	fmt.Println(user)

	if alert != nil {
		resp.Msg = alert.Error()
		resp.Status = "error"
		utility.SendJSON(w, resp, 500)
		return
	}

	userdb := gsql.NewUser()
	userdb.Login = user.Login
	userdb.Pass = user.Pass

	uid, err := userdb.GetUser()
	if err != nil {
		resp.Msg = err.Error()
		resp.Status = "error"
		utility.SendJSON(w, resp, 500)
		return
	}

	newid := utility.GenerateUID()

	user.UID = uid
	user.ActivateTime = time.Now()
	active[newid] = user

	resp.Msg = newid
	utility.SendJSON(w, resp, 200)

}

//Findid - ведем поиск и проверку актуальности по УИД ...
func Findid(w http.ResponseWriter, id string) string {

	result := ""

	user, istrue := active[id]

	if istrue != true {
		return result
	}

	activatetime := user.ActivateTime

	delta := activatetime.Sub(time.Now()) / 86400000000000

	if delta > 3 {
		delete(active, id)
		return result
	}

	result = user.UID

	return result

}

//Watchid - найти и посмотреть текущий ID ...
func Watchid(w http.ResponseWriter, r *http.Request, clientdata JSONclient) {
	resp := utility.Response{Status: "Ok", IsArray: false, Msg: clientdata.ID}
	utility.SendJSON(w, resp, 200)
}

//Write - запишем новые/измененные данные в БД ...
func Write(w http.ResponseWriter, r *http.Request, clientdata JSONclient) {

	resp := utility.Response{Status: "Ok", IsArray: false}
	code := 200

	dbdata := gsql.DBselector{}
	initSelector(clientdata, &dbdata)

	err := dbdata.Add()
	if err != nil {
		resp.Msg = err.Error()
		resp.Status = "error"
		code = 500
		utility.SendJSON(w, resp, code)
		return
	}

	resp.Msg = "Ok"
	utility.SendJSON(w, resp, code)

}

//CreateNewUser - создаем нового пользователя БД ...
func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	resp := utility.Response{Status: "Ok", IsArray: false}

	js := json.NewDecoder(r.Body)

	alert := js.Decode(&user)
	if alert != nil {
		resp.Msg = alert.Error()
		resp.Status = "error"
		utility.SendJSON(w, resp, 500)
		return
	}

	dbUser := gsql.NewUser()

	dbUser.Login = user.Login
	dbUser.Pass = user.Pass
	dbUser.Email = user.Email

	uid, err := dbUser.CreateUser()
	if err != nil {
		resp.Msg = err.Error()
		resp.Status = "error"
		utility.SendJSON(w, resp, 500)
		return
	}

	user.UID = uid
	user.ActivateTime = time.Now()

	newid := utility.GenerateUID()
	active[newid] = user

	resp.Msg = newid
	utility.SendJSON(w, resp, 200)

}

//List - получаем и выводим список сохраненных данных из БД ...
func List(w http.ResponseWriter, r *http.Request, clientdata JSONclient) {
	resp := utility.Response{}

	dbdata := gsql.DBselector{}
	initSelector(clientdata, &dbdata)

	arraydb, err := dbdata.Sel()
	if err != nil {
		resp.IsArray = false
		resp.Status = "error"
		resp.Msg = err.Error()
		code := 500
		utility.SendJSON(w, resp, code)
		return
	}

	resp.IsArray = true
	resp.Arraymsg = arraydb
	resp.Status = "Ok"
	utility.SendJSON(w, resp, 200)
}

//Change - Хэндлер для записи изменений в БД ...
func Change(w http.ResponseWriter, r *http.Request, clientdata JSONclient) {
	dbdata := gsql.DBselector{}
	initSelector(clientdata, &dbdata)
	alert := dbdata.Upd()
	if alert != nil {
		utility.SendJSON(w, utility.Response{Status: "error", IsArray: false, Arraymsg: nil, Msg: alert.Error()}, 405)
		return
	}
	utility.SendJSON(w, utility.Response{Status: "Ok", IsArray: false, Arraymsg: nil, Msg: "Ok"}, 200)
}

//Delete - Хэндлер для записи изменений в БД ...
func Delete(w http.ResponseWriter, r *http.Request, clientdata JSONclient) {
	dbdata := gsql.DBselector{}
	initSelector(clientdata, &dbdata)
	alert := dbdata.Del()
	if alert != nil {
		utility.SendJSON(w, utility.Response{Status: "error", IsArray: false, Arraymsg: nil, Msg: alert.Error()}, 405)
		return
	}
	utility.SendJSON(w, utility.Response{Status: "Ok", IsArray: false, Arraymsg: nil, Msg: "Ok"}, 200)
}

func initSelector(cd JSONclient, dbsel *gsql.DBselector) {
	dbsel.User = cd.User
	dbsel.Login = cd.Client.Login
	dbsel.Pass = cd.Client.Pass
	dbsel.Email = cd.Client.Email
	dbsel.Address = cd.Client.Address
}

//Gener - создаем уникальный код ...
func Gener(w http.ResponseWriter, r *http.Request, clientdata JSONclient) {
	code, alert := utility.GenerateUniqCode()
	if alert != nil {
		utility.SendJSON(w, utility.Response{Status: "error", IsArray: false, Arraymsg: nil, Msg: alert.Error()}, 500)
		return
	}
	utility.SendJSON(w, utility.Response{Status: "Ok", IsArray: false, Arraymsg: nil, Msg: code}, 200)
}

//OpenIndex = обработаем роут на главную страницу ...
func OpenIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, alert := template.ParseFiles("template/signin.html")
	if alert != nil {
		fmt.Printf(alert.Error())
	}
	tmpl.Execute(w, nil)
}

//OpenMain - открываем основную страницу ...
func OpenMain(w http.ResponseWriter, r *http.Request, clientdata JSONclient) {
	tmpl, alert := template.ParseFiles("template/main.html")
	if alert != nil {
		fmt.Printf(alert.Error())
	}
	tmpl.Execute(w, nil)
}

//SignIn - Логинимся ...
func SignIn(w http.ResponseWriter, r *http.Request) {
	Getid(w, r)
}
