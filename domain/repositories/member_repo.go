package repositories

import (
	"database/sql"
	"fmt"
	"log"
	domain "main/domain/model"
	"main/utils"
	"strings"
)

type memberRepo struct {
	db *sql.DB
}

func (m *memberRepo) FindByEmail(memberLogin *domain.MemberLogin) (*domain.Member, utils.MessageErr) {
	//tx, err := m.db.Begin()
	//if err != nil {
	//	return nil,  utils.ParseError(err)
	//}
	memberNew := new(domain.Member)
	query := fmt.Sprintf("SELECT id, first_name, last_name, email, password, loggedIn FROM members WHERE email = ? AND password = ?")
	if getError := m.db.QueryRow(query, &memberLogin.Email, &memberLogin.Password).
		Scan(&memberNew.Id, &memberNew.FirstName, &memberNew.LastName, &memberNew.Email, &memberNew.Password, &memberNew.LoggedIn); getError != nil {
		fmt.Println("this is the error man: ", getError)
		//tx.Rollback()
		return nil,  utils.ParseError(getError)
	}
	return memberNew, nil
}

func (m *memberRepo) UpdateStatus(member *domain.Member, status int) utils.MessageErr {
	//tx, err := m.db.Begin()
	//if err != nil {
	//	return utils.ParseError(err)
	//}
	query := fmt.Sprintf("UPDATE members SET loggedIn = ? WHERE email = ?")
	_, updateErr := m.db.Exec(query, &status, &member.Email)
	if updateErr != nil {
		s := strings.Split(updateErr.Error(), ":")
		log.Println(s[1])
		if updateErr != nil {
			//tx.Rollback()
			return utils.ParseError(updateErr)
		}
	}
	return nil
}

func NewMemberRepo(db *sql.DB) domain.IMemberRepository {
	return &memberRepo{db: db}
}

func (m *memberRepo) AddMember(member *domain.Member) (*domain.Member, utils.MessageErr) {
	query := fmt.Sprintf(`INSERT INTO members(first_name, last_name, email, password) VALUES(?,?,?,?)`)
	result, err := m.db.Exec(query, &member.FirstName, &member.LastName, &member.Email, &member.Password)
	if err != nil {
		s := strings.Split(err.Error(), ":")
		log.Println(s[1])
		return nil, utils.NewInternalServerError(fmt.Sprintf("error when trying to prepare user to save: %s", err.Error()))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, utils.NewInternalServerError(fmt.Sprintf("error when trying to save message: %s", err.Error()))
	}

	member.Id = int(id)

	return member, nil
}
