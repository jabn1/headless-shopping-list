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

type ItemDto struct {
	Name         string
	Quantity     int
	Status       string
	ShoppingList int
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
		r.Post("/shoppinglist", createShoppingList)
		r.Put("/shoppinglist/{id}", updateShoppingList)
		r.Delete("/shoppinglist/{id}", deleteShoppingList)
		r.Get("/item", getItems)
		r.Get("/item/{name}", getItem)
		r.Post("/item", createItem)
		r.Put("/item", updateItem)
		r.Delete("/item/{name}", deleteItem)
	})
	return r
}

func getItems(w http.ResponseWriter, r *http.Request) {
	shoppingListIdString := r.URL.Query().Get("shoppingListId")
	id, err := strconv.Atoi(shoppingListIdString)
	if err != nil {
		http.Error(w, "Invalid shopping list id", 400)
		return
	}

	if shoppingList, exists := shoppingLists[id]; exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shoppingList.Items)
	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}
}

func getItem(w http.ResponseWriter, r *http.Request) {
	shoppingListIdString := r.URL.Query().Get("shoppingListId")
	itemName := chi.URLParam(r, "name")
	id, err := strconv.Atoi(shoppingListIdString)
	if err != nil {
		http.Error(w, "Invalid shopping list id", 400)
		return
	}

	if shoppingList, exists := shoppingLists[id]; exists {
		if item, exists := shoppingList.Items[itemName]; exists {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]Item{itemName: item})
		} else {
			http.Error(w, "Shopping list id does not exist", 404)
			return
		}

	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}
}

func createItem(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Empty request body", 400)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newItem ItemDto
	err := decoder.Decode(&newItem)
	if err != nil {
		http.Error(w, "Could not parse request body", 400)
		return
	}
	if shoppingList, exists := shoppingLists[newItem.ShoppingList]; exists {
		if _, exists := shoppingList.Items[newItem.Name]; exists {
			http.Error(w, "An item with that name already exists", 409)
			return
		} else {
			shoppingLists[newItem.ShoppingList].Items[newItem.Name] = Item{Quantity: newItem.Quantity, Status: newItem.Status}
		}

	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Empty request body", 400)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newItem ItemDto
	err := decoder.Decode(&newItem)
	if err != nil {
		http.Error(w, "Could not parse request body", 400)
		return
	}
	if shoppingList, exists := shoppingLists[newItem.ShoppingList]; exists {
		if _, exists := shoppingList.Items[newItem.Name]; exists {
			shoppingLists[newItem.ShoppingList].Items[newItem.Name] = Item{Quantity: newItem.Quantity, Status: newItem.Status}
		} else {
			http.Error(w, "Item does not exist", 404)
			return
		}

	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	shoppingListIdString := r.URL.Query().Get("shoppingListId")
	itemName := chi.URLParam(r, "name")
	id, err := strconv.Atoi(shoppingListIdString)
	if err != nil {
		http.Error(w, "Invalid shopping list id", 400)
		return
	}
	if shoppingList, exists := shoppingLists[id]; exists {
		if _, exists := shoppingList.Items[itemName]; exists {
			delete(shoppingLists[id].Items, itemName)
		} else {
			http.Error(w, "Item does not exist", 404)
			return
		}

	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}
	w.WriteHeader(http.StatusOK)
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
