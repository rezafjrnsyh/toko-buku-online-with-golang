package services

import (
	"database/sql"
	domain "main/domain/model"
	"main/domain/repositories"
	"main/utils"
)

type memberService struct {
	db *sql.DB
	MemberRepo domain.IMemberRepository
	BookService domain.IBookService
}

func (m *memberService) BuyBook(id int, bookBuy *domain.ReqBuy) (*domain.Purchase, utils.MessageErr) {
	var purchase domain.Purchase
	book,err := m.BookService.FindBookById(id)
	if err != nil {
		return nil, err
	}
	m.BookService.ReduceStock(book, bookBuy, id)
	total := book.Price * bookBuy.Qty
	purchase.Book = book
	purchase.Qty = bookBuy.Qty
	purchase.TotalPrice = total

	return &purchase, nil
}

func (m *memberService) SignIn(memberLogin *domain.MemberLogin) utils.MessageErr {
	data, errFind := m.MemberRepo.FindByEmail(memberLogin)
	if errFind != nil {
		return errFind
	}
	err := m.MemberRepo.UpdateStatus(data, 1)
	if err != nil {
		return err
	}
	return nil
}

func NewMemberService(db *sql.DB) domain.IMemberService  {
	return &memberService{db: db, MemberRepo: repositories.NewMemberRepo(db), BookService: NewBookService(db)}
}

func (m *memberService) SignUp(member *domain.Member) (*domain.Member, utils.MessageErr) {
	member, err := m.MemberRepo.AddMember(member)
	if err != nil {
		return nil, err
	}
	return member, nil

}
