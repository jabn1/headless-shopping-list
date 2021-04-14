package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Item struct {
	Quantity int
	Status   string
	ETag     int `json:"-"`
}

type ItemDto struct {
	Name     string
	Quantity int
	Status   string
}

type ItemPatchDto struct {
	Name     *string
	Quantity *int
	Status   *string
}

type ShoppingList struct {
	Description string
	Date        string
	Items       map[string]*Item //key: item name
	ETag        int              `json:"-"`
}

type ShoppingListDto struct {
	Description *string
	Date        *string
	Items       map[string]*Item //key: item name
}

var shoppingLists map[int]*ShoppingList
var count int
var etagCount int
var port string
var listsEtag int

func main() {
	port = "5000"
	shoppingLists = map[int]*ShoppingList{}
	count = 1
	etagCount = 1
	listsEtag = 0
	r := registerRoutes()
	fmt.Println("Listening on port :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func registerRoutes() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/shoppinglists", getShoppingLists)
		r.Get("/shoppinglists/{id}", getShoppingList)
		r.Post("/shoppinglists", createShoppingList)
		r.Put("/shoppinglists/{id}", updateShoppingList)
		r.Patch("/shoppinglists/{id}", patchShoppingList)
		r.Delete("/shoppinglists/{id}", deleteShoppingList)
		r.Get("/shoppinglists/{id}/items", getItems)
		r.Get("/shoppinglists/{id}/items/{name}", getItem)
		r.Post("/shoppinglists/{id}/items", createItem)
		r.Put("/shoppinglists/{id}/items/{name}", updateItem)
		r.Delete("/shoppinglists/{id}/items/{name}", deleteItem)
		r.Head("/shoppinglists", getShoppingLists)
		r.Head("/shoppinglists/{id}", getShoppingList)
		r.Head("/shoppinglists/{id}/items", getItems)
		r.Head("/shoppinglists/{id}/items/{name}", getItem)

	})
	return r
}

func getItems(w http.ResponseWriter, r *http.Request) {
	shoppingListIdString := chi.URLParam(r, "id")
	ifNoneMatch := r.Header.Get("If-None-Match")
	status := r.URL.Query().Get("status")
	id, err := strconv.Atoi(shoppingListIdString)
	if err != nil {
		http.Error(w, "Invalid shopping list id", 400)
		return
	}

	if shoppingList, exists := shoppingLists[id]; exists {
		if ifNoneMatch == "" || ifNoneMatch != strconv.Itoa(shoppingLists[id].ETag) {
			w.Header().Set("Content-Type", "application/json")
			if status == "" {
				w.Header().Set("Etag", strconv.Itoa(shoppingLists[id].ETag))
				json.NewEncoder(w).Encode(shoppingList.Items)
			} else {
				resultItems := map[string]Item{}
				for key, element := range shoppingList.Items {
					if element.Status == status {
						resultItems[key] = *element
					}
				}
				w.Header().Set("Etag", strconv.Itoa(shoppingLists[id].ETag))
				json.NewEncoder(w).Encode(resultItems)
			}
		} else {
			w.WriteHeader(http.StatusNotModified)
		}

	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}
}

func getItem(w http.ResponseWriter, r *http.Request) {
	shoppingListIdString := chi.URLParam(r, "id")
	ifNoneMatch := r.Header.Get("If-None-Match")
	itemName := chi.URLParam(r, "name")
	id, err := strconv.Atoi(shoppingListIdString)
	if err != nil {
		http.Error(w, "Invalid shopping list id", 400)
		return
	}

	if shoppingList, exists := shoppingLists[id]; exists {
		if item, exists := shoppingList.Items[itemName]; exists {
			if ifNoneMatch == "" || ifNoneMatch != strconv.Itoa(shoppingList.Items[itemName].ETag) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Etag", strconv.Itoa(shoppingLists[id].Items[itemName].ETag))
				json.NewEncoder(w).Encode(map[string]Item{itemName: *item})
			} else {
				w.WriteHeader(http.StatusNotModified)
			}
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
	shoppingListIdString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(shoppingListIdString)
	if err != nil {
		http.Error(w, "Invalid shopping list id", 400)
		return
	}

	if r.Body == nil {
		http.Error(w, "Empty request body", 400)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newItem ItemDto
	err = decoder.Decode(&newItem)
	if err != nil {
		http.Error(w, "Could not parse request body", 400)
		return
	}
	if shoppingList, exists := shoppingLists[id]; exists {
		if _, exists := shoppingList.Items[newItem.Name]; exists {
			http.Error(w, "An item with that name already exists", 409)
			return
		} else {
			shoppingLists[id].Items[newItem.Name] = &Item{Quantity: newItem.Quantity, Status: newItem.Status, ETag: etagCount}
			etagCount += 1
			shoppingLists[id].ETag = etagCount
			etagCount += 1
			listsEtag = etagCount
			etagCount += 1
		}

	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}
	w.Header().Set("Location", "http://localhost:"+port+"/shoppinglists/"+strconv.Itoa(id)+"/"+newItem.Name)
	w.Header().Set("Etag", strconv.Itoa(shoppingLists[id].Items[newItem.Name].ETag))
	w.WriteHeader(http.StatusCreated)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	shoppingListIdString := chi.URLParam(r, "id")
	ifMatch := r.Header.Get("If-Match")
	itemName := chi.URLParam(r, "name")
	id, err := strconv.Atoi(shoppingListIdString)
	if err != nil {
		http.Error(w, "Invalid shopping list id", 400)
		return
	}

	if r.Body == nil {
		http.Error(w, "Empty request body", 400)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newItem Item
	err = decoder.Decode(&newItem)
	if err != nil {
		http.Error(w, "Could not parse request body", 400)
		return
	}
	if shoppingList, exists := shoppingLists[id]; exists {
		if _, exists := shoppingList.Items[itemName]; exists {
			if ifMatch == strconv.Itoa(shoppingList.Items[itemName].ETag) {
				shoppingLists[id].Items[itemName] = &Item{Quantity: newItem.Quantity, Status: newItem.Status, ETag: etagCount}
				etagCount += 1
				shoppingLists[id].ETag = etagCount
				etagCount += 1
				listsEtag = etagCount
				etagCount += 1
				w.Header().Set("Etag", strconv.Itoa(shoppingLists[id].Items[itemName].ETag))
			} else {
				w.WriteHeader(http.StatusConflict)
			}
		} else {
			http.Error(w, "Item does not exist", 404)
			return
		}

	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}

}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	shoppingListIdString := chi.URLParam(r, "id")
	itemName := chi.URLParam(r, "name")
	id, err := strconv.Atoi(shoppingListIdString)
	if err != nil {
		http.Error(w, "Invalid shopping list id", 400)
		return
	}
	if shoppingList, exists := shoppingLists[id]; exists {
		if _, exists := shoppingList.Items[itemName]; exists {
			delete(shoppingLists[id].Items, itemName)
			shoppingLists[id].ETag = etagCount
			etagCount += 1
			listsEtag = etagCount
			etagCount += 1
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
	ifNoneMatch := r.Header.Get("If-None-Match")
	if ifNoneMatch == "" || ifNoneMatch != strconv.Itoa(listsEtag) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Etag", strconv.Itoa(listsEtag))
		json.NewEncoder(w).Encode(shoppingLists)
	} else {
		w.WriteHeader(http.StatusNotModified)
	}
}

func getShoppingList(w http.ResponseWriter, r *http.Request) {
	msgId := chi.URLParam(r, "id")
	ifNoneMatch := r.Header.Get("If-None-Match")
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
		if ifNoneMatch == "" || ifNoneMatch != strconv.Itoa(shoppingLists[id].ETag) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Etag", strconv.Itoa(shoppingLists[id].ETag))
			json.NewEncoder(w).Encode(map[int]ShoppingList{id: *shoppingList})
		} else {
			w.WriteHeader(http.StatusNotModified)
		}

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
	for key := range newShoppingList.Items {
		newShoppingList.Items[key].ETag = etagCount
		etagCount += 1
	}
	newShoppingList.ETag = etagCount
	etagCount += 1
	listsEtag = etagCount
	etagCount += 1
	shoppingLists[count] = &newShoppingList

	w.Header().Set("Etag", strconv.Itoa(shoppingLists[count].ETag))
	w.Header().Set("Location", "http://localhost:"+port+"/shoppinglists/"+strconv.Itoa(count))
	count += 1

	w.WriteHeader(http.StatusCreated)

}

func updateShoppingList(w http.ResponseWriter, r *http.Request) {
	msgId := chi.URLParam(r, "id")
	ifMatch := r.Header.Get("If-Match")
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
		if ifMatch == strconv.Itoa(shoppingLists[id].ETag) {
			newShoppingList.ETag = etagCount
			etagCount += 1
			for key := range newShoppingList.Items {
				newShoppingList.Items[key].ETag = etagCount
				etagCount += 1
			}
			listsEtag = etagCount
			etagCount += 1
			shoppingLists[id] = &newShoppingList
			w.Header().Set("Etag", strconv.Itoa(shoppingLists[id].ETag))
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusConflict)
		}

	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}
}

func patchShoppingList(w http.ResponseWriter, r *http.Request) {
	msgId := chi.URLParam(r, "id")
	ifMatch := r.Header.Get("If-Match")
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
	var newShoppingList ShoppingListDto
	err = decoder.Decode(&newShoppingList)
	if err != nil {
		http.Error(w, "Could not parse request body", 400)
		return
	}

	if _, exists := shoppingLists[id]; exists {
		if ifMatch == strconv.Itoa(shoppingLists[id].ETag) {
			if newShoppingList.Date != nil {
				shoppingLists[id].Date = *newShoppingList.Date
			}
			if newShoppingList.Description != nil {
				shoppingLists[id].Description = *newShoppingList.Description
			}
			if newShoppingList.Items != nil {
				shoppingLists[id].Items = *&newShoppingList.Items
			}
			shoppingLists[id].ETag = etagCount
			etagCount += 1
			for key := range newShoppingList.Items {
				newShoppingList.Items[key].ETag = etagCount
				etagCount += 1
			}
			listsEtag = etagCount
			etagCount += 1
			w.Header().Set("Etag", strconv.Itoa(shoppingLists[id].ETag))
		} else {
			w.WriteHeader(http.StatusConflict)
		}

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
		listsEtag = etagCount
		etagCount += 1
	} else {
		http.Error(w, "Shopping list id does not exist", 404)
		return
	}
}
