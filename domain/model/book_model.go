package domain

import "main/utils"

type Book struct {
	Id	int `json:"id"`
	Title string `json:"title"`
	Description string `json:"desc"`
	Year int	`json:"year"`
	Pages int	`json:"pages"`
	Language string	`json:"language"`
	Publisher	string `json:"publisher"`
	Price int	`json:"price"`
	Stock	int	`json:"stock"`
	PurchaseAmount int	`json:"purchaseAmount"`
}

type IBookRepository interface {
	Find() ([]*Book, utils.MessageErr)
	Create(book *Book) (*Book, utils.MessageErr)
	FindById(id int) (*Book, utils.MessageErr)
	Update(book *Book) (*Book, utils.MessageErr)
	Delete(id int) (int64,utils.MessageErr)
	UpdateStock(book *Book) (*Book , utils.MessageErr)
	UpdatePurchaseAmount(book *Book) (*Book , utils.MessageErr)
}

type IBookService interface {
	FindBook() ([]*Book, utils.MessageErr)
	CreateBook(book *Book) (*Book, utils.MessageErr)
	FindBookById(id int) (*Book, utils.MessageErr)
	UpdateBook(book *Book, id int) (*Book, utils.MessageErr)
	DeleteBook(id int) (int64,utils.MessageErr)
	AddStock(stock int, id int) utils.MessageErr
	ReduceStock(book []Purchase) []Purchase
	AddPurchaseAmountBook(book []Purchase) []Purchase
}
