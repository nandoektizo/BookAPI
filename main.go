package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// Book represents a book in the library
type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	PublishedYear string `json:"published_year"`
	ISBN          int    `json:"isbn"`
}

// Author represents an author of a book
type Author struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

type AuthorBook struct {
	AuthorBookID int `json:"author_book_id"`
	AuthorID     int `json:"author_id"`
	BookID       int `json:"book_id"`
}

// JWT claims struct
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var (
	db              *sql.DB
	secretKey       = []byte("your-secret-key")
	authorizedUsers = map[string]string{
		"admin": "password",
		"user":  "password",
	}
)

func main() {
	var err error
	db, err = sql.Open("mysql", "username:password@tcp(localhost:3306)/library")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/books", getAllBooks).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	router.HandleFunc("/authors", getAllAuthors).Methods("GET")
	router.HandleFunc("/authors", createAuthor).Methods("POST")
	router.HandleFunc("/authors/{id}", getAuthor).Methods("GET")
	router.HandleFunc("/authors/{id}", updateAuthor).Methods("PUT")
	router.HandleFunc("/authors/{id}", deleteAuthor).Methods("DELETE")
	router.HandleFunc("/authorbooks", CreateAuthorBook).Methods("POST")
	router.HandleFunc("/authorbooks/{id}", GetAuthorBook).Methods("GET")
	router.HandleFunc("/authorbooks/{id}", UpdateAuthorBook).Methods("PUT")
	router.HandleFunc("/authorbooks/{id}", DeleteAuthorBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

// Login handles the user login and generates a JWT token
func login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	password, ok := authorizedUsers[creds.Username]
	if !ok || password != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Username: creds.Username,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// Middleware to validate JWT token
func validateToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CRUD operations for books
// CRUD operations for books

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, published_year, isbn FROM books")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		books := []Book{}
		for rows.Next() {
			var book Book
			err := rows.Scan(&book.ID, &book.Title, &book.PublishedYear, &book.ISBN)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			books = append(books, book)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	})(w, r)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		var book Book
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := db.Exec("INSERT INTO books (title, published_year, isbn) VALUES (?, ?,?)", book.Title, book.PublishedYear, book.ISBN)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ID, _ := result.LastInsertId()
		book.ID = int(ID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	})(w, r)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var book Book
		err := db.QueryRow("SELECT id, title, published_year, isbn FROM books WHERE id = ?", id).Scan(&book.ID, &book.Title, &book.PublishedYear, &book.ISBN)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	})(w, r)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var book Book
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = db.Exec("UPDATE books SET title = ?, published_year = ? , isbn = ?WHERE id = ?", book.Title, book.PublishedYear, book.ISBN, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		book.ID, _ = strconv.Atoi(id)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	})(w, r)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		_, err := db.Exec("DELETE FROM books WHERE id = ?", id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})(w, r)
}

func getAllAuthors(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, counttry FROM authors")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		authors := []Author{}
		for rows.Next() {
			var author Author
			err := rows.Scan(&author.ID, &author.Name, &author.Country)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			authors = append(authors, author)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(authors)
	})(w, r)
}

func createAuthor(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		var author Author
		err := json.NewDecoder(r.Body).Decode(&author)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := db.Exec("INSERT INTO authors (name, country) VALUES (?, ?)", author.Name, author.Country)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ID, _ := result.LastInsertId()
		author.ID = int(ID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(author)
	})(w, r)
}

func getAuthor(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var author Author
		err := db.QueryRow("SELECT id, name, country FROM books WHERE id = ?", id).Scan(&author.ID, &author.Name, &author.Country)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(author)
	})(w, r)
}

func updateAuthor(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var author Author
		err := json.NewDecoder(r.Body).Decode(&author)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = db.Exec("UPDATE authors SET name = ?, country = ?WHERE id = ?", author.Name, author.Country, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		author.ID, _ = strconv.Atoi(id)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(author)
	})(w, r)
}

func deleteAuthor(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		_, err := db.Exec("DELETE FROM authors WHERE id = ?", id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})(w, r)
}

// CreateAuthorBook creates a new author book relationship
func CreateAuthorBook(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		var authorBook AuthorBook
		err := json.NewDecoder(r.Body).Decode(&authorBook)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check if the author and book exist
		var authorExists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM authors WHERE id = ?)", authorBook.AuthorID).Scan(&authorExists)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !authorExists {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Author does not exist"})
			return
		}

		var bookExists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM books WHERE id = ?)", authorBook.BookID).Scan(&bookExists)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !bookExists {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Book does not exist"})
			return
		}

		result, err := db.Exec("INSERT INTO author_books (author_id, book_id) VALUES (?, ?)", authorBook.AuthorID, authorBook.BookID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ID, _ := result.LastInsertId()
		authorBook.AuthorBookID = int(ID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(authorBook)
	})(w, r)
}

// GetAuthorBook retrieves a specific author book relationship
func GetAuthorBook(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var authorBook AuthorBook
		err := db.QueryRow("SELECT author_book_id, author_id, book_id FROM author_books WHERE author_book_id = ?", id).Scan(&authorBook.AuthorBookID, &authorBook.AuthorID, &authorBook.BookID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(authorBook)
	})(w, r)
}

// UpdateAuthorBook updates an author book relationship
func UpdateAuthorBook(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var authorBook AuthorBook
		err := json.NewDecoder(r.Body).Decode(&authorBook)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check if the author and book exist
		var authorExists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM authors WHERE id = ?)", authorBook.AuthorID).Scan(&authorExists)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !authorExists {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Author does not exist"})
			return
		}

		var bookExists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM books WHERE id = ?)", authorBook.BookID).Scan(&bookExists)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !bookExists {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Book does not exist"})
			return
		}

		_, err = db.Exec("UPDATE author_books SET author_id = ?, book_id = ? WHERE author_book_id = ?", authorBook.AuthorID, authorBook.BookID, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		authorBook.AuthorBookID, _ = strconv.Atoi(id)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(authorBook)
	})(w, r)
}

// DeleteAuthorBook deletes an author book relationship
func DeleteAuthorBook(w http.ResponseWriter, r *http.Request) {
	// Token validation middleware
	validateToken(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		_, err := db.Exec("DELETE FROM author_books WHERE author_book_id = ?", id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})(w, r)
}
