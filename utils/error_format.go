package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func ParseError(err error) MessageErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	log.Println("NUMBER:", sqlErr.Number)
	if !ok {
		if strings.Contains(err.Error(), "no rows in result set") {
			return NewNotFoundError("no record matching given id")
		}
		return NewInternalServerError(fmt.Sprintf("error when trying to save message: %s", err.Error()))
	}
	switch sqlErr.Number {
	case 1062:
		fmt.Println("masuk 1")
		return NewInternalServerError("title already taken")
	case 1064:
		fmt.Println("masuk")
		return NewInternalServerError("Internal Server Error")
	}
	return NewInternalServerError(fmt.Sprintf("error when processing request: %s", err.Error()))
}
