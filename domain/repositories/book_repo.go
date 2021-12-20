package repositories

import (
	"fmt"
	"log"
	"main/constant"
	domain "main/domain/model"
	"main/utils"

	"github.com/jmoiron/sqlx"
)

type bookRepo struct {
	db *sqlx.DB
}

// NewBookRepo is a constructor
func NewBookRepo(db *sqlx.DB) domain.IBookRepository {
	return &bookRepo{db: db}
}

func (c *bookRepo) UpdatePurchaseAmount(book *domain.Book) (*domain.Book, error) {
	fmt.Println("In Book Repo : ", &book.PurchaseAmount, &book.ID)
	query := fmt.Sprintf("UPDATE book SET purchase_amount = ? WHERE id = ?")
	// Eksekusi query
	_, updateErr := c.db.Exec(query, &book.PurchaseAmount, book.ID)
	if updateErr != nil {
		return nil, utils.ParseError(updateErr)
	}
	return book, nil
}

func (c *bookRepo) UpdateStock(book *domain.Book) (*domain.Book, error) {
	fmt.Println("In Book Repo : ", &book.Stock, &book.ID)
	query := fmt.Sprintf("UPDATE book SET stock = ? WHERE id = ?")
	// Eksekusi query
	_, updateErr := c.db.Exec(query, &book.Stock, book.ID)
	if updateErr != nil {
		//s := strings.Split(updateErr.Error(), ":")
		//log.Println(s[1])
		return nil, utils.ParseError(updateErr)
	}
	return book, nil
}

func (c *bookRepo) Find() ([]*domain.Book, error) {
	// Membuat object slice category
	books := []*domain.Book{}

	// Eksekusi query
	err := c.db.Select(&books, constant.FIND_BOOKS)
	if err != nil {
		fmt.Println("err repo :", err)
		return nil, utils.ParseError(err)
	}
	// defer rows.Close()

	// for rows.Next() {
	// 	book := &domain.Book{}
	// 	getError := rows.Scan(&book.ID, &book.Title, &book.Description, &book.Price, &book.Stock, &book.PurchaseAmount)
	// 	if err != nil {
	// 		return nil, utils.NewInternalServerError(fmt.Sprintf("Error when trying to get message: %s", getError.Error()))
	// 	}
	// 	books = append(books, book)
	// }
	if len(books) == 0 {
		return nil, utils.NewNotFoundError("no records found")
	}
	return books, nil
}

func (c *bookRepo) Create(book *domain.Book) (*domain.Book, error) {
	lastInsertId := 0
	row := c.db.QueryRow(constant.INSERT_BOOK, book.Title,book.Description, book.Price, book.Stock).Scan(&lastInsertId)
	if row != nil {
		log.Println(row.Error())
		return nil, utils.NewInternalServerError("Internal Server Error")
	}
	book.ID = lastInsertId
	return book, nil
}

func (c *bookRepo) FindById(id int) (*domain.Book, error) {
	book := domain.Book{}
	getError := c.db.Get(&book, constant.FIND_BOOK_BY_ID, id)
	if getError != nil {
		fmt.Println("err repo ", getError)
		return nil, utils.ParseError(getError)
	}
	return &book, nil
}

func (c *bookRepo) Update(book *domain.Book) (*domain.Book, error) {
	_, updateErr := c.db.NamedExec(constant.UPDATE_BOOK, book)
	if updateErr != nil {
		log.Println("err repo ", updateErr)
		return nil, utils.ParseError(updateErr)
	}

	return book, nil
}

func (c *bookRepo) Delete(id int) (int64, error) {
	result:= c.db.MustExec(constant.DELETE_BOOK, id)
	RowsAffected, errRows := result.RowsAffected()
	if errRows != nil {
		return 0, utils.ParseError(errRows)
	}
	return RowsAffected, nil
}
