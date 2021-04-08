package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"

	// "strconv"

	"github.com/go-chi/chi"
)

type Item struct {
	Quantity int
	Status   string
}

type ShoppingList struct {
	Title       string
	Description string
	Fecha       string
	Items       map[string]Item //key: item name
}

var shoppingLists map[int]ShoppingList
var count int

func main() {
	port := "5000"
	shoppingLists = map[int]ShoppingList{}
	count = 1
	r := registerRoutes()
	fmt.Println("Listening on port :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func registerRoutes() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/shoppinglist", getShoppingLists)
		// r.Get("/shoppinglist/{msgId}", getMessages)
		r.Post("/shoppinglist", createShoppingList)
		// r.Put("/shoppinglist/{msgId}", updateMessage)
		// r.Delete("/shoppinglist/{msgId}", deleteMessage)
	})
	return r
}

func getShoppingLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shoppingLists)
}

func createShoppingList(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Empty request body", 400)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newShoppingList ShoppingList
	err := decoder.Decode(&newShoppingList)
	if err != nil {
		http.Error(w, "Could not parse request body", 400)
		return
	}

	shoppingLists[count] = newShoppingList
	count += 1
	w.WriteHeader(http.StatusCreated)

}

// func updateMessage(w http.ResponseWriter, r *http.Request) {
// 	msgId := chi.URLParam(r, "msgId")
// 	body, err := ioutil.ReadAll(r.Body)
// 	if msgId == "" {
// 		http.Error(w, "Empty message id", 400)
// 		return
// 	}
// 	if err != nil {
// 		http.Error(w, "Invalid message", 400)
// 		return
// 	}
// 	if len(body) == 0 {
// 		http.Error(w, "Empty request body", 400)
// 		return
// 	}
// 	id, err := strconv.Atoi(msgId)
// 	if err != nil {
// 		http.Error(w, "Invalid message id", 400)
// 		return
// 	}

// 	if _, exists := messages[id]; exists {
// 		messages[id] = string(body)
// 	} else {
// 		http.Error(w, "Message id does not exist", 400)
// 		return
// 	}
// }

// func deleteMessage(w http.ResponseWriter, r *http.Request) {
// 	msgId := chi.URLParam(r, "msgId")
// 	if msgId == "" {
// 		http.Error(w, "Empty message id", 400)
// 		return
// 	}
// 	id, err := strconv.Atoi(msgId)
// 	if err != nil {
// 		http.Error(w, "Invalid message id", 400)
// 		return
// 	}

// 	if _, exists := messages[id]; exists {
// 		delete(messages, id)
// 	} else {
// 		http.Error(w, "Message id does not exist", 400)
// 		return
// 	}
// }
