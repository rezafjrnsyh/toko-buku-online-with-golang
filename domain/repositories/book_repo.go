package repositories

import (
	"database/sql"
	"fmt"
	"log"
	domain "main/domain/model"
	"main/utils"
	"strconv"
	"strings"
)

type bookRepo struct {
	db *sql.DB
}

// NewBookRepo is a constructor
func NewBookRepo(db *sql.DB) domain.IBookRepository {
	return &bookRepo{db: db}
}

func (c *bookRepo) UpdatePurchaseAmount(book *domain.Book) (*domain.Book, error) {
	fmt.Println("In Book Repo : " , &book.PurchaseAmount, &book.Id )
	query := fmt.Sprintf("UPDATE book SET purchase_amount = ? WHERE id = ?")
	// Eksekusi query
	_, updateErr := c.db.Exec(query, &book.PurchaseAmount, strconv.Itoa(book.Id))
	if updateErr != nil {
		return nil, utils.ParseError(updateErr)
	}
	return book,nil
}

func (c *bookRepo) UpdateStock(book *domain.Book) (*domain.Book, error) {
	fmt.Println("In Book Repo : " , &book.Stock, &book.Id )
	query := fmt.Sprintf("UPDATE book SET stock = ? WHERE id = ?")
	// Eksekusi query
	_, updateErr := c.db.Exec(query, &book.Stock, strconv.Itoa(book.Id))
	if updateErr != nil {
		//s := strings.Split(updateErr.Error(), ":")
		//log.Println(s[1])
		return nil, utils.ParseError(updateErr)
	}
	return book,nil
}

func (c *bookRepo) Find() ([]*domain.Book, error) {
	// Membuat object slice category
	books := make([]*domain.Book, 0)
	// Untuk format query
	query := fmt.Sprintf(`SELECT id, title, description, year, pages, language, publisher, price, stock, purchase_amount FROM book`)

	// Eksekusi query
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, utils.ParseError(err)
	}
	defer rows.Close()

	for rows.Next() {
		book := &domain.Book{}
		getError := rows.Scan(&book.Id,&book.Title, &book.Description, &book.Year, &book.Pages, &book.Language,
			&book.Publisher, &book.Price, &book.Stock, &book.PurchaseAmount)
		if err != nil {
			return nil, utils.NewInternalServerError(fmt.Sprintf("Error when trying to get message: %s", getError.Error()))
		}
		books = append(books, book)
	}
	if len(books) == 0 {
		return nil, utils.NewNotFoundError("no records found")
	}
	return books, nil
}

func (c *bookRepo) Create(book *domain.Book) (*domain.Book, error) {
	query := fmt.Sprintf(`INSERT INTO book(title, description, year, pages, language, publisher, price, stock) VALUES(?,?,?,?,?,?,?,?)`)
	result, err := c.db.Exec(query, &book.Title, &book.Description, &book.Year, &book.Pages, &book.Language,
		&book.Publisher, &book.Price, &book.Stock)
	if err != nil {
		s := strings.Split(err.Error(), ":")
		log.Println(s[1])
		return nil, utils.NewInternalServerError(fmt.Sprintf("error when trying to prepare user to save: %s", err.Error()))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, utils.NewInternalServerError(fmt.Sprintf("error when trying to save message: %s", err.Error()))
	}

	book.Id = int(id)

	return book, nil
}

func (c *bookRepo) FindById(id int) (*domain.Book, error) {
	book := new(domain.Book)

	query := fmt.Sprintf(`SELECT id, title, description, year, pages, language, publisher, price, stock FROM book WHERE id=?`)
	if getError := c.db.QueryRow(query, id).
		Scan(&book.Id, &book.Title, &book.Description, &book.Year, &book.Pages, &book.Language,
		&book.Publisher, &book.Price, &book.Stock); getError != nil {
		fmt.Println("this is the error man: ", getError)
		return nil,  utils.ParseError(getError)
	}
	return book, nil
}

func (c *bookRepo) Update(book *domain.Book) (*domain.Book, error) {
	query := fmt.Sprintf("UPDATE book SET title = ?, description = ?, year = ?, pages = ?, language = ?, publisher = ?, price = ?, stock = ? WHERE id = ?")
	_, updateErr := c.db.Exec(query, &book.Title, &book.Description, &book.Year, &book.Pages, &book.Language,
		&book.Publisher, &book.Price, &book.Stock, &book.Id, )
	if updateErr != nil {
		s := strings.Split(updateErr.Error(), ":")
		log.Println(s[1])
		if updateErr != nil {
			return nil,  utils.ParseError(updateErr)
		}
	}

	return book, nil
}

func (c *bookRepo) Delete(id int) (int64, error) {
	query := fmt.Sprintf("DELETE FROM book WHERE id = ?")
	result, err := c.db.Exec(query, id)
	fmt.Println("repo", err.Error())
	if err != nil {
		return 0, utils.ParseError(err)
	}
	RowsAffected, errRows := result.RowsAffected()
	if errRows != nil {
		return 0, utils.ParseError(errRows)
	}
	return RowsAffected,nil
}
