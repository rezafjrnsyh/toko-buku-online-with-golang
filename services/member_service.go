package services

import (
	"database/sql"
	"fmt"
	"log"
	domain "main/domain/model"
	"main/domain/repositories"
	"main/utils"
)

type memberService struct {
	db *sql.DB
	MemberRepo domain.IMemberRepository
	BookService domain.IBookService
}

func (m *memberService) FindMembers() ([]*domain.Member, utils.MessageErr) {
	members,err := m.MemberRepo.Find()
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (m *memberService) ActivatedMember(memberId int) utils.MessageErr {
	member, err := m.MemberRepo.FindMemberById(memberId)
	fmt.Println(member)
	if err != nil {
		log.Println("Service member :", err)
		return err
	}
	err = m.MemberRepo.UpdateStatus(member, 1)
	if err != nil {
		log.Println("Service member :", err)
		return err
	}
	return nil
}

func (m *memberService) GetHistoryTrxMember(memberId int) ([]*domain.ResponseHistoryMember, utils.MessageErr) {
	histories, err := m.MemberRepo.FindHistoryByMember(memberId)
	if err != nil {
		log.Println("Service member :", err)
		return nil, err
	}
	return histories, nil
}

func (m *memberService) Buys(buys []domain.Buy, memberId int) ([]domain.Purchase, utils.MessageErr) {
	fmt.Println("Service : ", buys)
	//var books []*domain.Book
	var purchase domain.Purchase
	var purchases []domain.Purchase

	for _,buy := range buys {
		book, _ := m.BookService.FindBookById(buy.BookID)
		if book.Stock < 1 || buy.Qty > book.Stock {
			return nil, utils.NewBadRequestError("Insufficient stock available book " + book.Title)
		}
		//m.BookService.ReduceStock(book, &buy)
		purchase.Book = book
		purchase.Qty = buy.Qty
		purchase.TotalPrice = book.Price * buy.Qty
		fmt.Println("in Loop :", purchase)
		purchases = append(purchases, purchase)
	}

	data, err := m.MemberRepo.AddBooks(purchases, memberId)
	if err != nil {
		return nil, err
	}

	newPurchases := m.BookService.ReduceStock(data)
	newPurchases2 := m.BookService.AddPurchaseAmountBook(newPurchases)
	return newPurchases2, nil
}

//func (m *memberService) BuyBook(id int, bookBuy *domain.Buy) (*domain.Purchase, utils.MessageErr) {
//	var purchase domain.Purchase
//	book,err := m.BookService.FindBookById(id)
//	if err != nil {
//		return nil, err
//	}
//	newBook := m.BookService.ReduceStock(book, bookBuy)
//	total := book.Price * bookBuy.Qty
//	purchase.Book = newBook
//	purchase.Qty = bookBuy.Qty
//	purchase.TotalPrice = total
//
//	return &purchase, nil
//}

func (m *memberService) SignIn(memberLogin *domain.MemberLogin) (*domain.Member,  error) {
	fmt.Println("credential", memberLogin)
	member, errFind := m.MemberRepo.FindByEmail(memberLogin)
	fmt.Println("test",member)
	fmt.Println(errFind)

	if errFind != nil {
		return nil, errFind
	}

	//if member.Email != memberLogin.Email || member.Password != memberLogin.Password {
	//	fmt.Println("err status")
	//	return nil, utils.NewBadRequestError("Please check your Email or Password is wrong!")
	//}

	return member, nil
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
