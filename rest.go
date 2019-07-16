package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"html"
)

type Routes []Route

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var indexContent string

func Rest() {
	result := "crowd-mind\n\n"
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
		result += fmt.Sprintf("%s %s %s\n", route.Method, route.Pattern, route.Name)
	}
	indexContent = result
	log.Fatal(http.ListenAndServe(":8080", router))
}

var routes = Routes{
	Route{"Index", "GET", "/", Index},
	Route{"Categories", "GET", "/categories", CategoryList},
	Route{"List by category", "GET", "/list/{category}", ContentByCategory},
	Route{"Retrieve by id", "GET", "/{id}", ContentById},
	Route{"Accept entry", "POST", "/accept/{id}", Accept},
	Route{"Reject entry", "POST", "/reject/{id}", Reject},
	Route{"Update existing category {keyword: ..., content: ...}", "POST", "/update/{id}", Update},
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", html.EscapeString(indexContent))
}

func CategoryList(w http.ResponseWriter, r *http.Request) {

}

func ContentByCategory(w http.ResponseWriter, r *http.Request) {

}

func ContentById(w http.ResponseWriter, r *http.Request) {

}

func Accept(w http.ResponseWriter, r *http.Request) {

}

func Reject(w http.ResponseWriter, r *http.Request) {

}

func Update(w http.ResponseWriter, r *http.Request) {

}
