package domain

import (
	"main/utils"
)

type Member struct {
	ID        int     `gorm:"column:id;size:36;auto_increment;primaryKey;"`
	FirstName string  `db:"first_name" json:"firstname"`
	LastName  string  `db:"last_name" json:"lastname"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Status    int     `json:"status"`
	Books     []*Book `gorm:"many2many:members_books;ForeignKey:id"`
}

func (c *Member) TableName() string {
	return "mst_members"
}

type MemberLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ReqBuy struct {
	Qty int `json:"qty"`
}

type Purchase struct {
	Id         int `json:"id"`
	Book       *Book
	Qty        int `json:"qty"`
	TotalPrice int `json:"total_price"`
}

type Buy struct {
	BookID int `json:"bookId"`
	Qty    int `json:"qty"`
}

type RequestBuyBooks struct {
	Buys []Buy `json:"buys"`
}

type ResponseHistoryMember struct {
	User       string `json:"user"`
	Title      string `json:"title"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
	TotalPrice int    `json:"totalPrice"`
}
type IMemberRepository interface {
	Find() ([]*Member, utils.MessageErr)
	AddMember(member *Member) (*Member, utils.MessageErr)
	UpdateStatus(member *Member) error
	FindByEmail(memberLogin *MemberLogin) (*Member, error)
	//note pointer
	AddBooks(purchases []Purchase, memberId int) ([]Purchase, utils.MessageErr)
	FindHistoryByMember(memberId int) ([]*ResponseHistoryMember, utils.MessageErr)
	FindMemberById(id int) (*Member, error)
}

type IMemberService interface {
	FindMembers() ([]*Member, utils.MessageErr)
	SignUp(member *Member) (*Member, utils.MessageErr)
	SignIn(member *MemberLogin) (*Member, error)
	ActivatedMember(memberId int) error
	//BuyBook(id int, bookBuy *Buy) (*Purchase, utils.MessageErr)
	Buys(buys []Buy, memberId int) ([]Purchase, utils.MessageErr)
	GetHistoryTrxMember(memberId int) ([]*ResponseHistoryMember, utils.MessageErr)
}
