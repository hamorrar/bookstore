package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hamorrar/bookstore/internal/database"
)

func (app *application) createBook(c *gin.Context) {
	var book database.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := app.models.Books.CreateBook(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (app *application) getAllBooks(c *gin.Context) {
	books, err := app.models.Books.GetAllBooks()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

func (app *application) getBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	book, err := app.models.Books.GetBook(id)
	if book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get book"})
		return
	}

	c.JSON(http.StatusOK, book)

}

func (app *application) deleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := app.models.Books.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (app *application) updateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	existingBook, err := app.models.Books.GetBook(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get book"})
		return
	}

	if existingBook == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	updatedBook := &database.Book{}
	updatedBook.Id = id

	if err := c.ShouldBindJSON(updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedBook.Id = id

	if err := app.models.Books.UpdateBook(updatedBook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book"})
		return
	}

	c.JSON(http.StatusOK, updatedBook)

}
