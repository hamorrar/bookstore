package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hamorrar/bookstore/internal/database"
)

// createBook creates a book
// @Summary		creates a book
// @Description	creates a book
// @Tags		book
// @Accept		json
// @Produce		json
// @Param		book body database.Book true "new book to add to db"
// @Success		201	{object} database.Book "successfully created a book"
// @Failure 403 {object} gin.H "wrong role"
// @Failure 400 {object} gin.H "error binding JSON"
// @Failure 500 {object} gin.H "error creating books"
// @Router			/api/v1/books [post]
// @Security CookieAuth
func (app *application) createBook(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unathorized to create book"})
		return
	}

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

// getPageOfBooks gets a page of books
// @Summary		gets a page of books
// @Description	gets a page of books
// @Tags		book
// @Accept		json
// @Produce		json
// @Param		page query int false "page number to request"
// @Param limit query int false "max number of books to return per page"
// @Success		200	{array} database.Book "successfully got a page of books"
// @Failure 500 {object} gin.H "error getting a page"
// @Router			/api/v1/books [get]
func (app *application) getPageOfBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	books, err := app.models.Books.GetPageOfBooks(limit, page)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get books on page %d and limit %d", page, limit), "error msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

// getAllBooks gets all books
// @Summary		gets all books
// @Description	gets all books
// @Tags		book
// @Accept		json
// @Produce		json
// @Success		200	{array} database.Book "successfully got all books"
// @Failure 500 {object} gin.H "error getting all books"
// @Router			/api/v2/books/all [get]
// @Security CookieAuth
func (app *application) getAllBooks(c *gin.Context) {

	user := app.GetUserFromContext(c)
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to get all books"})
		return
	}

	limit := 3
	page := 1
	var allBooks []*database.Book
	for {
		books, err := app.models.Books.GetPageOfBooks(limit, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get all books page by page."})
			return
		}

		allBooks = append(allBooks, books...)
		if len(books) < limit {
			break
		}
		page++

	}

	c.JSON(http.StatusOK, allBooks)
}

// getBook get one book
// @Summary		get one book
// @Description	get one book by id
// @Tags		book
// @Accept		json
// @Produce		json
// @Param		id query int true "id of book to get"
// @Success		200	{object} database.Book "successfully got a book"
// @Failure 400 {object} gin.H "invalid book id"
// @Failure 404 {object} gin.H "book not found with this id"
// @Failure 500 {object} gin.H "error getting book"
// @Router			/api/v1/books/:id [get]
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

// deleteBook delete a book
// @Summary		delete book
// @Description	delete a book by id
// @Tags		book
// @Accept		json
// @Produce		json
// @Param		id query int true "id of book to delete"
// @Success		204	"successfully deleted"
// @Failure 403 {object} gin.H "wrong role"
// @Failure 400 {object} gin.H "invalid id"
// @Failure 500 {object} gin.H "error deleting book"
// @Router			/api/v1/books/:id [delete]
// @Security CookieAuth
func (app *application) deleteBook(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to delete a book"})
		return
	}

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

// updateBook updates a book
// @Summary		update a book
// @Description	update a book by id
// @Tags		book
// @Accept		json
// @Produce		json
// @Param		id query int true "id of book to update"
// @Param book body database.Book true "updated book data"
// @Success 200	{object} database.Book "successfully updated a book"
// @Failure 403 {object} gin.H "wrong role"
// @Failure 400 {object} gin.H "invalid id"
// @Failure 500 {object} gin.H "error getting book"
// @Failure 404 {object} gin.H "book to update not found"
// @Failure 400 {object} gin.H "error binding JSON"
// @Failure 500 {object} gin.H "failed to update book"
// @Router			/api/v1/books/:id [put]
// @Security CookieAuth
func (app *application) updateBook(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update book"})
		return
	}

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
