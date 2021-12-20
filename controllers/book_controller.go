package controllers

import (
	"fmt"
	"log"
	domain "main/domain/model"
	"main/services"
	"main/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type bookController struct {
	BookService domain.IBookService
}

const (
	BOOK_LIST_PATH  = "/"
	BOOK_CREATE_PATH = "/"
	BOOK_GET_BY_ID_PATH = "/:id"
	BOOK_DELETE_PATH = "/:id"
	BOOK_UPDATE_PATH = "/:id"
	ADD_STOCK_BOOK_PATH = "/:id/stock"
)

func NewBookController(db *sqlx.DB, r *gin.Engine)  {
	// var data map[string]string
	// data["username"] = "root"
	Controller := bookController{BookService: services.NewBookService(db)}
	bookRoutes := r.Group("/books") 
	// bookRoutes.Use(gin.BasicAuth(data))
	{
		bookRoutes.GET(BOOK_LIST_PATH, Controller.lstBook)
		bookRoutes.POST(BOOK_CREATE_PATH, Controller.AddBook)
		bookRoutes.GET(BOOK_GET_BY_ID_PATH, Controller.GetBookById)
		bookRoutes.PUT(ADD_STOCK_BOOK_PATH, Controller.addStockBook)
		bookRoutes.PUT(BOOK_UPDATE_PATH, Controller.UpdateBook)
		bookRoutes.DELETE(BOOK_DELETE_PATH, Controller.DeleteBook)
		r.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"Code" : http.StatusNotFound, "Message": "PAGE_NOT_FOUND"})
		})
	}
	// r.Any("/*any", func (c *gin.Context)  {
	// 	c.String(404,"PAGE_NOT_FOUND",)
	// })
}

func (b *bookController) addStockBook(c *gin.Context) {
	param := c.Param("id")
	id,err := strconv.Atoi(param)
	if err != nil {
		log.Println("Failed to converted to int")
		c.JSON(http.StatusInternalServerError, gin.H{"code" : 500, "message" : "Internal Server Error"})
	}
	var stock domain.Book
	errBind := c.ShouldBindJSON(&stock)
	if errBind != nil {
		theErr := utils.NewUnprocessibleEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
		//s := strings.Split(err.Error(), "'")
		//errField := fmt.Errorf("field %s can't be empty", s[3])
		//c.JSON(http.StatusBadRequest, gin.H{"message": errField.Error(), "code": 400})
	}

	errStock := b.BookService.AddStock(stock.Stock, id)
	if errStock != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code" : 500, "message" : "Internal Server Error"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code" : 200, "message" : "Add stock book successfully"})
	}
}

func (b *bookController) lstBook(c *gin.Context) {
	books, err := b.BookService.FindBook()
	fmt.Print("err",err)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.Response(http.StatusOK, "ok", books))
}

func (b *bookController) AddBook(c *gin.Context) {
	var book domain.Book
	err := c.ShouldBindJSON(&book)
	if err != nil {
		theErr := utils.NewUnprocessibleEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}

	newBook, error := b.BookService.CreateBook(&book)
	log.Println("controller ", newBook)
	if error != nil {
		c.JSON(http.StatusInternalServerError, utils.NewInternalServerError(error.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.Response(http.StatusCreated, "Book created successfully", newBook))
}

func (b *bookController) GetBookById(c *gin.Context) {
	param := c.Param("id")
	id,err := strconv.Atoi(param)
	if err != nil {
		log.Println("Failed to converted to int")
		c.JSON(http.StatusInternalServerError, gin.H{"code" : 500, "message" : "Internal Server Error"})
	}
	book, er := b.BookService.FindBookById(id)
	if er != nil {
		log.Println(er.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code" : 400, "message" : er.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok", "data": book})
	}
}

func (b *bookController) UpdateBook(c *gin.Context) {
	param := c.Param("id")
	id,errparse := strconv.Atoi(param)
	if errparse != nil {
		log.Println("Failed to converted to int")
		c.JSON(http.StatusInternalServerError, gin.H{"code" : 500, "message" : "Internal Server Error"})
	}

	var book domain.Book
	err := c.ShouldBindJSON(&book)
	if err != nil {
		theErr := utils.NewUnprocessibleEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
		//s := strings.Split(err.Error(), "'")
		//errField := fmt.Errorf("field %s can't be empty", s[3])
		//c.JSON(http.StatusBadRequest, gin.H{"message": errField.Error(), "code": 400})
	}

	newBook, error := b.BookService.UpdateBook(&book, id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, utils.NewInternalServerError("Internal Server Error"))
	}
	c.JSON(http.StatusOK, utils.Response(http.StatusCreated, "Book updated successfully", newBook))
}

func (b *bookController) DeleteBook(c *gin.Context) {
	param := c.Param("id")
	id,err := strconv.Atoi(param)
	if err != nil {
		log.Println("Failed to converted to int")
		c.JSON(http.StatusInternalServerError, gin.H{"code" : 500, "message" : "Internal Server Error"})
	}
	result, err := b.BookService.DeleteBook(id)
	log.Println("rows:",result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code" : 500, "message" : "Internal server error"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Data deleted successfully", "data": result})
	}
}
