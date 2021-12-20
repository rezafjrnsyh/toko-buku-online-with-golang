package repositories

import (
	"fmt"
	"log"
	"main/constant"
	domain "main/domain/model"
	"main/utils"
	"strings"

	"github.com/jmoiron/sqlx"
)

type memberRepo struct {
	db       *sqlx.DB
	BookRepo domain.IBookRepository
}

func (m *memberRepo) Find() ([]*domain.Member, utils.MessageErr) {
	members := make([]*domain.Member, 0)
	// Untuk format query
	query := fmt.Sprintf(`SELECT id, first_name, last_name, email, password, status FROM members`)

	// Eksekusi query
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, utils.ParseError(err)
	}
	defer rows.Close()

	for rows.Next() {
		member := &domain.Member{}
		getError := rows.Scan(&member.ID, &member.FirstName, &member.LastName, &member.Email, &member.Password, &member.Status)
		if err != nil {
			return nil, utils.NewInternalServerError(fmt.Sprintf("Error when trying to get message: %s", getError.Error()))
		}
		members = append(members, member)
	}
	if len(members) == 0 {
		return nil, utils.NewNotFoundError("no records found")
	}
	return members, nil
}

func (m *memberRepo) FindMemberById(id int) (*domain.Member, error) {
	fmt.Println("id", id)
	member := domain.Member{}
	getError := m.db.Get(&member, constant.FIND_MEMBER, id)
	if getError != nil {
		fmt.Println("this is the error man: ", getError)
		return nil, utils.ParseError(getError)
	}
	return &member, nil
}

func NewMemberRepo(db *sqlx.DB) domain.IMemberRepository {
	return &memberRepo{db: db, BookRepo: NewBookRepo(db)}
}

func (m *memberRepo) FindHistoryByMember(memberId int) ([]*domain.ResponseHistoryMember, utils.MessageErr) {
	query := fmt.Sprintf(`select first_name, title, price, Quantity, TotalPrice FROM MemberhasBooks as mhb 
			JOIN book as b on mhb.BookID = b.Id 
			JOIN members as m on mhb.MemberID = m.Id WHERE m.Id = ?;`)
	rows, err := m.db.Query(query, memberId)
	fmt.Println("ROWS", rows)
	if err != nil {
		return nil, utils.ParseError(err)
	}
	defer rows.Close()

	histories := make([]*domain.ResponseHistoryMember, 0)

	for rows.Next() {
		history := &domain.ResponseHistoryMember{}
		err = rows.Scan(&history.User, &history.Title, &history.Price, &history.Quantity, &history.TotalPrice)
		if err != nil {
			return nil, utils.NewInternalServerError(fmt.Sprintf("Error when trying to get message: %s", err.Error()))
		}
		histories = append(histories, history)
	}
	if len(histories) == 0 {
		return nil, utils.NewNotFoundError("no records found")
	}
	return histories, nil
}

func (m *memberRepo) AddBooks(purchases []domain.Purchase, memberId int) ([]domain.Purchase, utils.MessageErr) {
	newPurchases := make([]domain.Purchase, 0)
	trx, err := m.db.Begin()
	if err != nil {
		return nil, utils.ParseError(err)
	}
	query := fmt.Sprintf(`INSERT INTO MemberhasBooks(MemberID, BookID, Quantity, TotalPrice) VALUES(?,?,?,?)`)
	for _, purchase := range purchases {
		//fmt.Println("member id : ", memberId, "Book id :", purchase.Book.Id, purchase.Qty, purchase.TotalPrice)
		result, err := m.db.Exec(query, memberId, purchase.Book.ID, purchase.Qty, purchase.TotalPrice)
		if err != nil {
			s := strings.Split(err.Error(), ":")
			log.Println(s[1])
			trx.Rollback()
			return nil, utils.NewInternalServerError(fmt.Sprintf("error when trying to prepare user to save: %s", err.Error()))
		}

		id, err := result.LastInsertId()
		fmt.Println("in repo", id)
		if err != nil {
			return nil, utils.NewInternalServerError(fmt.Sprintf("error when trying to save message: %s", err.Error()))
		}

		purchase.Id = int(id)
		newPurchases = append(newPurchases, purchase)
	}
	return newPurchases, nil
}

func (m *memberRepo) FindByEmail(memberLogin *domain.MemberLogin) (*domain.Member, error) {
	//tx, err := m.db.Begin()
	//if err != nil {
	//	return nil,  utils.ParseError(err)
	//}
	memberNew := new(domain.Member)
	query := fmt.Sprintf("SELECT id, first_name, last_name, email, password, status FROM members WHERE email = ? AND password = ?")
	if getError := m.db.QueryRow(query, &memberLogin.Email, &memberLogin.Password).
		Scan(&memberNew.ID, &memberNew.FirstName, &memberNew.LastName, &memberNew.Email, &memberNew.Password, &memberNew.Status); getError != nil {
		fmt.Println("this is the error man: ", getError)
		//tx.Rollback()
		return nil, getError
	}
	return memberNew, nil
}

func (m *memberRepo) UpdateStatus(member *domain.Member) error {
	fmt.Println(member)
	res, _ := m.db.NamedExec(constant.UPDATE_STATUS_MEMBER, &member)
	RowsAffected, errRows := res.RowsAffected()
	if RowsAffected == 0 {
		return errRows
	}
	return nil
}

func (m *memberRepo) AddMember(member *domain.Member) (*domain.Member, utils.MessageErr) {
	lastInsertId := 0
	row := m.db.QueryRow(constant.INSERT_MEMBER, member.FirstName, member.LastName, member.Email, member.Password, member.Status).Scan(&lastInsertId)

	if row != nil {
		log.Println(row.Error())
		return nil, utils.NewInternalServerError("Internal Server Error")
	}
	member.ID = lastInsertId
	return member, nil
}
