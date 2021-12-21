package domain

type Book struct {
	ID	int `gorm:"column:id;size:36;auto_increment;primaryKey;"`
	Title string `json:"title"`
	Description string `json:"desc"`
	Price int	`json:"price"`
	Stock	int	`json:"stock"`
	PurchaseAmount int	`db:"purchase_amount" json:"purchaseAmount"`
	Members []*Member `gorm:"many2many:members_books"`
}

func (c *Book) TableName() string {
	return "mst_books"
}

type IBookRepository interface {
	Find() ([]*Book, error)
	Create(book *Book) (*Book, error)
	FindById(id int) (*Book, error)
	Update(book *Book) (*Book, error)
	Delete(id int) (int64, error)
	UpdateStock(book *Book) (*Book , error)
	UpdatePurchaseAmount(book *Book) (*Book , error)
}

type IBookService interface {
	FindBook() ([]*Book, error)
	CreateBook(book *Book) (*Book, error)
	FindBookById(id int) (*Book, error)
	UpdateBook(book *Book, id int) (*Book, error)
	DeleteBook(id int) (int64,error)
	AddStock(stock int, id int) error
	ReduceStock(book []Purchase) []Purchase
	AddPurchaseAmountBook(book []Purchase) []Purchase
}
