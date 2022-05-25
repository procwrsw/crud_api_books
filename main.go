package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     string `json:"id"`
	Isbn   string `json:"isbn"`
	Name   string `json:"name"`
	Writer *Writer
}

type Writer struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var books []Book

func getBooks(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func deleteBook(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(req.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(100))
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, req *http.Request) {
	// at first we delete the book
	// next we create the book
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)

			var book Book
			_ = json.NewDecoder(req.Body).Decode(&book)
			book.ID = strconv.Itoa(rand.Intn(100))
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
		}
	}
}

func main() {

	r := mux.NewRouter()

	books = append(books, Book{ID: "1", Isbn: "234456", Name: "Book one",
		Writer: &Writer{Firstname: "Monica", Lastname: "Goths"}})
	books = append(books, Book{ID: "2", Isbn: "123098", Name: "Book two",
		Writer: &Writer{Firstname: "Lukas", Lastname: "Tomson"}})
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
