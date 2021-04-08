package main

import (
	"encoding/json"
	"fmt"
	"strconv"

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
	Description string
	Date        string
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
		r.Get("/shoppinglist/{id}", getShoppingList)
		// r.Get("/shoppinglist/{msgId}", getMessages)
		r.Post("/shoppinglist", createShoppingList)
		r.Put("/shoppinglist/{id}", updateShoppingList)
		r.Delete("/shoppinglist/{id}", deleteShoppingList)
	})
	return r
}

func getShoppingLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shoppingLists)
}

func getShoppingList(w http.ResponseWriter, r *http.Request) {
	msgId := chi.URLParam(r, "id")
	if msgId == "" {
		http.Error(w, "Empty message id", 400)
		return
	}
	id, err := strconv.Atoi(msgId)
	if err != nil {
		http.Error(w, "Invalid shopping list id", 400)
		return
	}

	if shoppingList, exists := shoppingLists[id]; exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[int]ShoppingList{id: shoppingList})
	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}
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

func updateShoppingList(w http.ResponseWriter, r *http.Request) {
	msgId := chi.URLParam(r, "id")
	if msgId == "" {
		http.Error(w, "Empty message id", 400)
		return
	}

	if r.Body == nil {
		http.Error(w, "Empty request body", 400)
		return
	}
	id, err := strconv.Atoi(msgId)
	if err != nil {
		http.Error(w, "Invalid message id", 400)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newShoppingList ShoppingList
	err = decoder.Decode(&newShoppingList)
	if err != nil {
		http.Error(w, "Could not parse request body", 400)
		return
	}

	if _, exists := shoppingLists[id]; exists {
		shoppingLists[id] = newShoppingList
	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}
}

func deleteShoppingList(w http.ResponseWriter, r *http.Request) {
	msgId := chi.URLParam(r, "id")
	if msgId == "" {
		http.Error(w, "Empty message id", 400)
		return
	}
	id, err := strconv.Atoi(msgId)
	if err != nil {
		http.Error(w, "Invalid shopping list id", 400)
		return
	}

	if _, exists := shoppingLists[id]; exists {
		delete(shoppingLists, id)
	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}
}
