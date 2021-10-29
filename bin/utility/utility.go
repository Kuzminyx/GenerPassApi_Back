package utility

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	gsql "github.com/GeneratePassAPI/bin/bd"
)

//Response - структура ответа для клиента ...
type Response struct {
	Status   string
	IsArray  bool
	Arraymsg []gsql.DBselector
	Msg      string
}

//GenerateUID - генерация уникального идентификатора ...
func GenerateUID() string {
	uuid := ""
	tempB := make([]byte, 16)
	_, err := rand.Read(tempB)
	if err != nil {
		log.Panic(err)
	}
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", tempB[0:4], tempB[4:6], tempB[6:8], tempB[8:10], tempB[10:])
	return uuid
}

//SendJSON - форимруем и отправляем JSON на клиент ...
func SendJSON(w http.ResponseWriter, msg Response, code int) {
	log.Println(msg.Status)
	w.WriteHeader(code)
	encode := json.NewEncoder(w)
	encode.Encode(msg)
}

//GenerateUniqCode - генерим уникальный рандомный код ...
func GenerateUniqCode() (string, error) {

	code := ""

	tempByte := make([]byte, 16)

	_, alert := rand.Read(tempByte)
	if alert != nil {
		return code, alert
	}

	code = fmt.Sprintf("%x", tempByte[0:8])

	return code, nil
}
