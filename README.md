# BookAPI
Using golang , mysql to create CRUD and implement jwt token also unit testing


before you running the project please create table on mysql following this struct


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
