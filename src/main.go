package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"./db"

	"./config"
	"gopkg.in/mgo.v2/bson"

	"github.com/Jeffail/gabs"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

type Response struct {
	ReturnCode int    `json:"returnCode"`
	Message    string `json:"message"`
	//Data       interface{} `json:"data"`
}

func newResponse() Response {
	instance := Response{}
	instance.ReturnCode = 1
	return instance
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Ok!")
}

func InfoLoja(w http.ResponseWriter, r *http.Request) {

	response := newResponse()

	vars := mux.Vars(r)

	nomeLoja := vars["id"]

	l := db.GetMongo().C("configs")

	var infoLoja interface{}

	err := l.Find(bson.M{"id": nomeLoja}).Select(bson.M{"_id": 0}).One(&infoLoja)

	if err != nil {
		response.Message = "Loja não encontrada."
		response.ReturnCode = 0

		json, _ := json.Marshal(response)

		ResponseWithJSON(w, json, 404)
	} else {
		json, _ := json.Marshal(infoLoja)

		ResponseWithJSON(w, json, 200)
	}

}

func InsertForm(w http.ResponseWriter, r *http.Request) {

	response := newResponse()

	l := db.GetMongo().C("configs")

	body, _ := ioutil.ReadAll(r.Body)

	sBody := string(body)

	if sBody == "{}" {
		response.Message = "Adicione campos em seu JSON."
		response.ReturnCode = 1

		json, _ := json.Marshal(response)

		ResponseWithJSON(w, json, 404)
	} else {

		jsonParsed, err := gabs.ParseJSON([]byte(body))

		if err != nil {
			log.Println("Erro! jsonParsed falhou.")
			os.Exit(0)
		}

		id, ok := jsonParsed.Search("id").Data().(string)

		if ok == false || id == "" {
			response.Message = "Adicione um ID em seu JSON."
			response.ReturnCode = 2

			json, _ := json.Marshal(response)

			ResponseWithJSON(w, json, 404)
		} else {
			var v interface{}

			err = l.Find(bson.M{"id": id}).Select(bson.M{"_id": 0}).One(&v)

			if v == nil {
				err = json.Unmarshal(body, &v)

				err = l.Insert(v)

				if err != nil {
					fmt.Println("Deu erro no insert")
					os.Exit(0)
				}
			} else {
				err = json.Unmarshal(body, &v)

				_, err := l.Upsert(bson.M{"id": id}, v)

				if err != nil {
					log.Println("Deu erro no UPSERT")
					os.Exit(0)
				}

			}

		}

	}

}

func main() {

	config.Load()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
	})

	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/api/insertform", InsertForm).Methods("POST")
	router.HandleFunc("/api/infoloja/{id}", InfoLoja).Methods("GET")

	muxServer := http.NewServeMux()
	muxServer.Handle("/", router)

	n := negroni.Classic()
	n.Use(cors)
	n.UseHandler(muxServer)
	n.Run(":3000")
}
