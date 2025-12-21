package handler

import (
	"lendbook/internal/entity"
	"lendbook/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	usecase usecase.BookUsecase
}

func NewBookHandler(u usecase.BookUsecase) *BookHandler {
	return &BookHandler{
		usecase: u,
	}
}

func (h *BookHandler) AddBook(c echo.Context) error {
	userid := c.Get("userID").(int)
	var req entity.AddBookParams
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	err = h.usecase.AddBook(c.Request().Context(), userid, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "book added successfully"})
}

func (h *BookHandler) ListBooks(c echo.Context) error {
	books, err := h.usecase.ListBooks(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": books})
}

func (h *BookHandler) GetBookDetails(c echo.Context) error {
	bookID, err := strconv.Atoi(c.Param("id")) // convert string parameter to int ID
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid book id"})
	}

	details, err := h.usecase.GetBookDetails(c.Request().Context(), bookID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": details})
}

func (h *BookHandler) DeleteBook(c echo.Context) error {
	userID := c.Get("userID").(int)

	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid book id"})
	}

	if err := h.usecase.DeleteBook(c.Request().Context(), userID, bookID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "book deleted successfully"})
}

func (h *BookHandler) BorrowBook(c echo.Context) error {
	userID := c.Get("userID").(int)

	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid book id"})
	}

	if err := h.usecase.BorrowBook(c.Request().Context(), userID, bookID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "book borrowed successfully"})
}
func (h *BookHandler) ReturnBook(c echo.Context) error {
	userID := c.Get("userID").(int)

	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid book id"})
	}

	if err := h.usecase.ReturnBook(c.Request().Context(), userID, bookID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "book returned successfully"})
}
