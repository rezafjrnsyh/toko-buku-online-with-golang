package utils

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"strings"
)

func ParseError(err error) MessageErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), "no rows in result set") {
			return NewNotFoundError("no record matching given id")
		}
		return NewInternalServerError(fmt.Sprintf("error when trying to save message: %s", err.Error()))
	}
	switch sqlErr.Number {
	case 1062:
		return NewInternalServerError("title already taken")
	}
	return NewInternalServerError(fmt.Sprintf("error when processing request: %s", err.Error()))
}
