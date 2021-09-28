package services

import (
	"database/sql"
	"fmt"
	domain "main/domain/model"
	"main/domain/repositories"
)

type bookService struct {
	db          *sql.DB
	BookRepo domain.IBookRepository
}

func (b *bookService) AddPurchaseAmountBook(purchases []domain.Purchase) []domain.Purchase {
	for _,purchase := range purchases{
		purchase.Book.PurchaseAmount += purchase.Qty
		_, _ = b.BookRepo.UpdatePurchaseAmount(purchase.Book)
		//newPurchases = append(newPurchases, newBook)
	}
	return purchases
}

func (b *bookService) ReduceStock(purchases []domain.Purchase) []domain.Purchase {
	//var newPurchases []*domain.Book
	for _,purchase := range purchases{
		purchase.Book.Stock -= purchase.Qty
		_, _ = b.BookRepo.UpdateStock(purchase.Book)
		//newPurchases = append(newPurchases, newBook)
	}
	return purchases
}

func (b *bookService) ReduceStockWithId(book *domain.Book, buy *domain.Buy, id int) *domain.Book{
	book.Stock = book.Stock - buy.Qty
	newBook,_ := b.BookRepo.UpdateStock(book)
	return newBook
}

func NewBookService(db *sql.DB) domain.IBookService {
	return &bookService{db: db, BookRepo: repositories.NewBookRepo(db)}
}

func (b *bookService) AddStock(stock int, id int) error {
	book, err := b.FindBookById(id)
	if err != nil {
		return err
	}

	book.Stock = book.Stock + stock
	_, errStock := b.BookRepo.UpdateStock(book)
	if errStock != nil {
		return errStock
	}
	return nil
}

func (b bookService) FindBook() ([]*domain.Book, error) {
	books,err := b.BookRepo.Find()
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (b *bookService) CreateBook(book *domain.Book) (*domain.Book, error) {
	book, err := b.BookRepo.Create(book)
	fmt.Println("Service :", book)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (b *bookService) FindBookById(id int) (*domain.Book, error) {
	book, err := b.BookRepo.FindById(id)
	if err != nil{
		return nil, err
	}

	return book, nil
}

func (b *bookService) UpdateBook(book *domain.Book, id int) (*domain.Book, error) {
	data,errfound := b.FindBookById(id)
	if errfound != nil {
		return nil, errfound
	}

	data.Title = book.Title
	data.Description = book.Description
	data.Year = book.Year
	data.Pages = book.Pages
	data.Language = book.Language
	data.Publisher = book.Publisher
	data.Price = book.Price
	data.Stock = book.Stock

	NewBook, err := b.BookRepo.Update(data)

	if err != nil{
		return nil, err
	}

	return NewBook, nil
}

func (b *bookService) DeleteBook(id int) (int64, error) {
	book,errfound := b.FindBookById(id)
	if errfound != nil{
		return 0, errfound
	}

	result, err := b.BookRepo.Delete(book.Id)
	if err != nil {
		return result, err
	}
	return result, nil
}
