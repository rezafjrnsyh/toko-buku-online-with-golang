package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	domain "main/domain/model"
	"main/services"
	"main/utils"
	"net/http"
	"strconv"
)

type categoryController struct {
	CategoryService domain.ICategoryService
}

const (
	CATEGORY_LIST_PATH  = "/category/list"
	CATEGORY_CREATE_PATH = "/category"
	CATEGORY_GET_BY_ID_PATH = "/category/:id"
	CATEGORY_DELETE_PATH = "/category/:id"
	CATEGORY_UPDATE_PATH = "/category/:id"
)

func NewCategoryController(db *sql.DB, r *gin.RouterGroup)  {
	Controller := categoryController{CategoryService: services.NewCategoryService(db)}
	r.GET(CATEGORY_LIST_PATH, Controller.lstCategory )
	r.POST(CATEGORY_CREATE_PATH, Controller.AddCategory)
	r.GET(CATEGORY_GET_BY_ID_PATH, Controller.GetCategoryById)
	r.PUT(CATEGORY_UPDATE_PATH, Controller.UpdateCategory)
	r.DELETE(CATEGORY_DELETE_PATH, Controller.DeleteCategory)
}

func (s *categoryController) lstCategory(c *gin.Context) {
	categories, err := s.CategoryService.FindCategory()
	fmt.Print("err",err)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.Response(http.StatusOK, "ok", categories))
}

func (s *categoryController) AddCategory(c *gin.Context) {
	var category domain.Category
	err := c.ShouldBindJSON(&category)
	if err != nil {
		theErr := utils.NewUnprocessibleEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
		//s := strings.Split(err.Error(), "'")
		//errField := fmt.Errorf("field %s can't be empty", s[3])
		//c.JSON(http.StatusBadRequest, gin.H{"message": errField.Error(), "code": 400})
	}

	newCategory, error := s.CategoryService.CreateCategory(&category)
	if err != nil {
		c.JSON(error.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, utils.Response(http.StatusCreated, "Category successfully created", newCategory))
}

func (s *categoryController) GetCategoryById(c *gin.Context) {
	param := c.Param("id")
	id,err := strconv.Atoi(param)
	if err != nil {
		log.Println("Failed to converted to int")
		c.JSON(http.StatusInternalServerError, gin.H{"code" : 500, "message" : "Internal Server Error"})
	}
	category, er := s.CategoryService.FindCategoryById(id)
	if er != nil {
		log.Println(er.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code" : 400, "message" : "data not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok", "data": category})
	}
}

func (s *categoryController) UpdateCategory(c *gin.Context) {
	param := c.Param("id")
	id,errparse := strconv.Atoi(param)
	if errparse != nil {
		log.Println("Failed to converted to int")
		c.JSON(http.StatusInternalServerError, gin.H{"code" : 500, "message" : "Internal Server Error"})
	}

	var category domain.Category
	err := c.ShouldBindJSON(&category)
	if err != nil {
		theErr := utils.NewUnprocessibleEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
		//s := strings.Split(err.Error(), "'")
		//errField := fmt.Errorf("field %s can't be empty", s[3])
		//c.JSON(http.StatusBadRequest, gin.H{"message": errField.Error(), "code": 400})
	}

	newCategory, error := s.CategoryService.UpdateCategory(&category, id)
	if err != nil {
		c.JSON(error.Status(), err)
		return
	}
	c.JSON(http.StatusOK, utils.Response(http.StatusCreated, "Category updated successfully", newCategory))
}

func (s *categoryController) DeleteCategory(c *gin.Context) {
	param := c.Param("id")
	id,err := strconv.Atoi(param)
	if err != nil {
		log.Println("Failed to converted to int")
		c.JSON(http.StatusInternalServerError, gin.H{"code" : 500, "message" : "Internal Server Error"})
	}
	result, err := s.CategoryService.DeleteCategory(id)
	log.Println("rows:",result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code" : 500, "message" : "Internal server error"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Data deleted successfully", "data": result})
	}
}
