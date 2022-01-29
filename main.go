package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID       string `json: "id"`
	Title    string `json: "title"`
	Author   string `json: "author"`
	Quantity int16  `json: "quantity"`
}

var books = []Book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, books)
}

func createBooks(ctx *gin.Context) {
	var newBook Book

	if err:= ctx.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	ctx.IndentedJSON(http.StatusCreated, newBook)
}


func getBookById(ctx *gin.Context) {
	id := ctx.Param("id")
	book, err := searchBookById(id)

	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "The requested book could not be found."})
		return
	}

	ctx.IndentedJSON(http.StatusOK, book)
}

func searchBookById(id string) (*Book, error) {
	for idx, book := range books {
		if book.ID == id {
			return &books[idx], nil
		}
	}

	return nil, errors.New("invalid id: book not found")
}

func checkoutBook(ctx *gin.Context) {
	id, ok := ctx.GetQuery("id")

	if !ok {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing query parameter id"})
		return
	}

	book, err := searchBookById(id)

	if err != nil || book.Quantity <= 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "The requested book could not be found."})
		return
	}

	book.Quantity -= 1

	ctx.IndentedJSON(http.StatusOK, book)
}

func returnBook(ctx *gin.Context) {
	id, ok := ctx.GetQuery("id")

	if !ok {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing query parameter id"})
		return
	}
	book, err := searchBookById(id)

	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book could not be found"})
		return
	}
	book.Quantity += 1
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}

func main() {
	router := gin.Default()

	router.GET("/books", getBooks)
	router.POST("/books", createBooks)
	router.GET("/books/:id", getBookById)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}