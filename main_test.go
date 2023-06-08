package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

func TestLogin(t *testing.T) {
	// Create a request body with valid credentials
	creds := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "admin",
		Password: "password",
	}
	body, _ := json.Marshal(creds)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and server
	router := mux.NewRouter()
	router.HandleFunc("/login", login).Methods("POST")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Login handler returned wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the token exists
	if _, ok := response["token"]; !ok {
		t.Errorf("Login handler didn't return a token in the response body")
	}
}

func TestGetAllBooks(t *testing.T) {
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and server
	router := mux.NewRouter()
	router.HandleFunc("/books", getAllBooks).Methods("GET")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("GetAllBooks handler returned wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// Check the response body
	var books []Book
	err = json.Unmarshal(rr.Body.Bytes(), &books)
	if err != nil {
		t.Fatal(err)
	}

	// Check if at least one book is returned
	if len(books) == 0 {
		t.Errorf("GetAllBooks handler returned an empty list of books")
	}
}

func TestCreateBook(t *testing.T) {
	book := Book{
		Title:         "Test Book",
		PublishedYear: "2023",
		ISBN:          1234567890,
	}
	body, _ := json.Marshal(book)

	req, err := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and server
	router := mux.NewRouter()
	router.HandleFunc("/books", createBook).Methods("POST")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("CreateBook handler returned wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// Check the response body
	var createdBook Book
	err = json.Unmarshal(rr.Body.Bytes(), &createdBook)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the created book matches the sent book
	if createdBook.Title != book.Title || createdBook.PublishedYear != book.PublishedYear || createdBook.ISBN != book.ISBN {
		t.Errorf("CreateBook handler didn't create the book correctly")
	}
}

func TestGetBook(t *testing.T) {
	// Create a book ID to retrieve
	bookID := 1

	req, err := http.NewRequest("GET", "/books/"+strconv.Itoa(bookID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and server
	router := mux.NewRouter()
	router.HandleFunc("/books/{id}", getBook).Methods("GET")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("GetBook handler returned wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// Check the response body
	var book Book
	err = json.Unmarshal(rr.Body.Bytes(), &book)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the retrieved book ID matches the requested book ID
	if book.ID != bookID {
		t.Errorf("GetBook handler didn't retrieve the correct book")
	}
}

func TestUpdateBook(t *testing.T) {
	// Create a book ID to update
	bookID := 1

	book := Book{
		Title:         "Updated Book",
		PublishedYear: "2022",
		ISBN:          9876543210,
	}
	body, _ := json.Marshal(book)

	req, err := http.NewRequest("PUT", "/books/"+strconv.Itoa(bookID), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and server
	router := mux.NewRouter()
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("UpdateBook handler returned wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// Check the response body
	var updatedBook Book
	err = json.Unmarshal(rr.Body.Bytes(), &updatedBook)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the updated book matches the sent book
	if updatedBook.Title != book.Title || updatedBook.PublishedYear != book.PublishedYear || updatedBook.ISBN != book.ISBN {
		t.Errorf("UpdateBook handler didn't update the book correctly")
	}
}

func TestDeleteBook(t *testing.T) {
	// Create a book ID to delete
	bookID := 1

	req, err := http.NewRequest("DELETE", "/books/"+strconv.Itoa(bookID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and server
	router := mux.NewRouter()
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusNoContent {
		t.Errorf("DeleteBook handler returned wrong status code: got %d, expected %d", rr.Code, http.StatusNoContent)
	}
}

func TestGetAllAuthors(t *testing.T) {
	req, err := http.NewRequest("GET", "/authors", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and server
	router := mux.NewRouter()
	router.HandleFunc("/authors", getAllAuthors).Methods("GET")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("GetAllAuthors handler returned wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// Check the response body
	var authors []Author
	err = json.Unmarshal(rr.Body.Bytes(), &authors)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the number of authors is greater than zero
	if len(authors) == 0 {
		t.Errorf("GetAllAuthors handler didn't retrieve any authors")
	}
}

func TestCreateAuthor(t *testing.T) {
	author := Author{
		Name:    "John Doe",
		Country: "United States",
	}
	body, _ := json.Marshal(author)

	req, err := http.NewRequest("POST", "/authors", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and server
	router := mux.NewRouter()
	router.HandleFunc("/authors", createAuthor).Methods("POST")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("CreateAuthor handler returned wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// Check the response body
	var createdAuthor Author
	err = json.Unmarshal(rr.Body.Bytes(), &createdAuthor)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the created author matches the sent author
	if createdAuthor.Name != author.Name || createdAuthor.Country != author.Country {
		t.Errorf("CreateAuthor handler didn't create the author correctly")
	}
}

func TestGetAuthor(t *testing.T) {
	// Create an author ID to retrieve
	authorID := 1

	req, err := http.NewRequest("GET", "/authors/"+strconv.Itoa(authorID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and server
	router := mux.NewRouter()
	router.HandleFunc("/authors/{id}", getAuthor).Methods("GET")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("GetAuthor handler returned wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// Check the response body
	var author Author
	err = json.Unmarshal(rr.Body.Bytes(), &author)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the retrieved author ID matches the requested author ID
	if author.ID != authorID {
		t.Errorf("GetAuthor handler didn't retrieve the correct author")
	}
}

func TestUpdateAuthor(t *testing.T) {
	// Create an author ID to update
	authorID := 1

	author := Author{
		Name:    "Updated Author",
		Country: "United Kingdom",
	}
	body, _ := json.Marshal(author)

	req, err := http.NewRequest("PUT", "/authors/"+strconv.Itoa(authorID), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and server
	router := mux.NewRouter()
	router.HandleFunc("/authors/{id}", updateAuthor).Methods("PUT")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("UpdateAuthor handler returned wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// Check the response body
	var updatedAuthor Author
	err = json.Unmarshal(rr.Body.Bytes(), &updatedAuthor)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the updated author matches the sent author
	if updatedAuthor.Name != author.Name || updatedAuthor.Country != author.Country {
		t.Errorf("UpdateAuthor handler didn't update the author correctly")
	}
}

func TestDeleteAuthor(t *testing.T) {
	// Create an author ID to delete
	authorID := 1

	req, err := http.NewRequest("DELETE", "/authors/"+strconv.Itoa(authorID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and server
	router := mux.NewRouter()
	router.HandleFunc("/authors/{id}", deleteAuthor).Methods("DELETE")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusNoContent {
		t.Errorf("DeleteAuthor handler returned wrong status code: got %d, expected %d", rr.Code, http.StatusNoContent)
	}
}

// TestCreateAuthorBook tests the CreateAuthorBook function
func TestCreateAuthorBook(t *testing.T) {
	// Create a new request body
	authorBook := AuthorBook{
		AuthorID: 1,
		BookID:   1,
	}
	body, err := json.Marshal(authorBook)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new request
	req, err := http.NewRequest("POST", "/authorbooks", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Set up a mock router
	router := mux.NewRouter()
	router.HandleFunc("/authorbooks", CreateAuthorBook).Methods("POST")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}

	// Check the response body
	expected := authorBook
	var responseAuthorBook AuthorBook
	err = json.Unmarshal(rr.Body.Bytes(), &responseAuthorBook)
	if err != nil {
		t.Fatal(err)
	}
	if responseAuthorBook.AuthorID != expected.AuthorID || responseAuthorBook.BookID != expected.BookID {
		t.Errorf("Expected author book %+v, but got %+v", expected, responseAuthorBook)
	}
}

// TestGetAuthorBook tests the GetAuthorBook function
func TestGetAuthorBook(t *testing.T) {
	// Create a new request
	req, err := http.NewRequest("GET", "/authorbooks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Set up a mock router
	router := mux.NewRouter()
	router.HandleFunc("/authorbooks/{id}", GetAuthorBook).Methods("GET")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}

	// Check the response body
	// Add assertions for the expected response body
}

// TestUpdateAuthorBook tests the UpdateAuthorBook function
func TestUpdateAuthorBook(t *testing.T) {
	// Create a new request body
	authorBook := AuthorBook{
		AuthorBookID: 1,
		AuthorID:     1,
		BookID:       2,
	}
	body, err := json.Marshal(authorBook)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new request
	req, err := http.NewRequest("PUT", "/authorbooks/1", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Set up a mock router
	router := mux.NewRouter()
	router.HandleFunc("/authorbooks/{id}", UpdateAuthorBook).Methods("PUT")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}

	// Check the response body
	// Add assertions for the expected response body
}

// TestDeleteAuthorBook tests the DeleteAuthorBook function
func TestDeleteAuthorBook(t *testing.T) {
	// Create a new request
	req, err := http.NewRequest("DELETE", "/authorbooks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Set up a mock router
	router := mux.NewRouter()
	router.HandleFunc("/authorbooks/{id}", DeleteAuthorBook).Methods("DELETE")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}
}
