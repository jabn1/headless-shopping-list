package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

var messages map[int]string
var count int

func main() {
	port := "5000"
	messages = map[int]string{
		1: "Hola",
		2: "Mundo",
	}
	count = 3
	r := registerRoutes()
	fmt.Println("Listening on port :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func registerRoutes() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/messages", getMessages)
		r.Post("/messages", createMessage)
		r.Put("/messages/{msgId}", updateMessage)
		r.Delete("/messages/{msgId}", deleteMessage)
	})
	return r
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(messages)
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid message", 400)
		return
	}
	if len(body) == 0 {
		http.Error(w, "Empty request body", 400)
		return
	}

	messages[count] = string(body)
	count += 1
}

func updateMessage(w http.ResponseWriter, r *http.Request) {
	msgId := chi.URLParam(r, "msgId")
	body, err := ioutil.ReadAll(r.Body)
	if msgId == "" {
		http.Error(w, "Empty message id", 400)
		return
	}
	if err != nil {
		http.Error(w, "Invalid message", 400)
		return
	}
	if len(body) == 0 {
		http.Error(w, "Empty request body", 400)
		return
	}
	id, err := strconv.Atoi(msgId)
	if err != nil {
		http.Error(w, "Invalid message id", 400)
		return
	}

	if _, exists := messages[id]; exists {
		messages[id] = string(body)
	} else {
		http.Error(w, "Message id does not exist", 400)
		return
	}
}

func deleteMessage(w http.ResponseWriter, r *http.Request) {
	msgId := chi.URLParam(r, "msgId")
	if msgId == "" {
		http.Error(w, "Empty message id", 400)
		return
	}
	id, err := strconv.Atoi(msgId)
	if err != nil {
		http.Error(w, "Invalid message id", 400)
		return
	}

	if _, exists := messages[id]; exists {
		delete(messages, id)
	} else {
		http.Error(w, "Message id does not exist", 400)
		return
	}
}
