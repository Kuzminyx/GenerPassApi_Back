package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	ghandlers "github.com/GeneratePassAPI/bin/handlers"
	utility "github.com/GeneratePassAPI/bin/utility"
)

type apihandler func(w http.ResponseWriter, r *http.Request, js ghandlers.JSONclient)

func main() {

	router := http.NewServeMux()

	//Site routes
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("/", ghandlers.OpenIndex)            //Ввод логина и пароля
	router.Handle("/main", apihandler(ghandlers.OpenMain)) //ЛК пользователя
	router.HandleFunc("/signIn", ghandlers.SignIn)         //Обертка для проверки учетки с веб клиента

	//API routes
	router.HandleFunc("/getid", ghandlers.Getid)
	router.HandleFunc("/adduser", ghandlers.CreateNewUser)
	router.Handle("/watchid", apihandler(ghandlers.Watchid))
	router.Handle("/write", apihandler(ghandlers.Write))
	router.Handle("/list", apihandler(ghandlers.List))
	router.Handle("/change", apihandler(ghandlers.Change))
	router.Handle("/del", apihandler(ghandlers.Delete))
	router.Handle("/gener", apihandler(ghandlers.Gener))

	server := &http.Server{
		Addr:    ":9901",
		Handler: router,
	}

	go server.ListenAndServe()

	// Wait for an interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

}

func (api apihandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	clientdata := json.NewDecoder(r.Body)

	clientreq := ghandlers.JSONclient{}
	result := ""

	alert := clientdata.Decode(&clientreq)
	if alert != nil {
		log.Println(alert.Error())
		resp := utility.Response{}
		resp.IsArray = false
		resp.Msg = alert.Error()
		resp.Status = "error"
		utility.SendJSON(w, resp, 500)
	}

	if clientreq.ID != "" {
		result = ghandlers.Findid(w, clientreq.ID)
	}

	if result != "" {
		clientreq.User = result
		api(w, r, clientreq)
	} else {
		msg := "Нет данных об авторизации, необходимо получить новый ID"
		log.Println(msg)
		resp := utility.Response{}
		resp.IsArray = false
		resp.Msg = msg
		resp.Status = "error"
		utility.SendJSON(w, resp, 404)
	}
}
