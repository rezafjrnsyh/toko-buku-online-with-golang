package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	domain "main/domain/model"
	"main/services"
	"main/utils"
	"net/http"
	"strconv"
)

type memberController struct {
	MemberService domain.IMemberService
}

const (
	SIGN_UP_PATH = "/member/signup"
	SIGN_IN_PATH = "/member/signin"
	BUYS_BOOK_PATH = "/member/book/buys/:memberId"
	GET_HISTORY_PATH = "/member/history/:id"
	ACTIVATED_PATH = "/member/activated/:memberId"
)

func NewMemberController(db *sql.DB, r *gin.RouterGroup) {
	Controller := memberController{MemberService: services.NewMemberService(db)}
	r.POST(SIGN_UP_PATH, Controller.SignUpMember)
	r.POST(SIGN_IN_PATH, Controller.SignInMember)
	r.POST(BUYS_BOOK_PATH, Controller.Buys)
	r.GET(GET_HISTORY_PATH, Controller.HistoryTrx)
	r.PUT(ACTIVATED_PATH, Controller.ActivatedMember)
}

func (m *memberController) HistoryTrx(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Println("Failed to converted to int")
		c.JSON(http.StatusInternalServerError, utils.NewInternalServerError("Internal server error"))
	}

	histories, errget := m.MemberService.GetHistoryTrxMember(id)
	if errget != nil {
		c.JSON(http.StatusInternalServerError, utils.NewInternalServerError("Internal server error"))
		return
	}
	c.JSON(http.StatusOK, utils.Response(http.StatusOK,"Success", histories))
}

func (m *memberController) Buys(c *gin.Context) {
	param := c.Param("memberId")
	id,err := strconv.Atoi(param)
	fmt.Println("id", id)
	if err != nil {
		log.Println("Failed to converted to int")
		c.JSON(http.StatusInternalServerError, gin.H{"code" : 500, "message" : "Internal Server Error"})
	}
	var requestBook domain.RequestBuyBooks

	errBind := c.ShouldBindJSON(&requestBook)
	if errBind != nil {
		theErr := utils.NewUnprocessibleEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}

	purchases, errPurchase := m.MemberService.Buys(requestBook.Buys, id)
	if errPurchase != nil {
		c.JSON(http.StatusUnauthorized, errPurchase)
	} else {
		c.JSON(http.StatusOK, utils.Response(http.StatusOK, "Success", purchases))
	}
}

//func (m *memberController) Buy(c *gin.Context) {
//	param := c.Param("id")
//	id,err := strconv.Atoi(param)
//	if err != nil {
//		log.Println("Failed to converted to int")
//		c.JSON(http.StatusInternalServerError, gin.H{"code" : 500, "message" : "Internal Server Error"})
//	}
//
//	var qty *domain.Buy
//	errBind := c.ShouldBindJSON(&qty)
//	if errBind != nil {
//		theErr := utils.NewUnprocessibleEntityError("invalid json body")
//		c.JSON(theErr.Status(), theErr)
//		return
//	}
//	trx, errBuy := m.MemberService.BuyBook(id, qty)
//	if errBuy != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"status":  http.StatusInternalServerError,
//			"message": "Internal server error",
//		})
//	}
//	c.JSON(http.StatusOK, utils.Response(200, "Your transaction was successful", trx))
//}

func (m *memberController) SignInMember(c *gin.Context)   {
	var member domain.MemberLogin
	errBind := c.ShouldBindJSON(&member)
	if errBind != nil {
		theErr := utils.NewUnprocessibleEntityError("invalid json body")
		c.Error(theErr)
		return
	}
	_,err := m.MemberService.SignIn(&member)
	fmt.Println("controller",err)
	if err != nil{
		fmt.Println("notfound")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Fail"})
		return
	}
	//	fmt.Println("success")
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK,"message": "Verified"})
}

func (m *memberController) SignUpMember(c *gin.Context)  {
	var member domain.Member
	errBind := c.ShouldBindJSON(&member)
	if errBind != nil {
		theErr := utils.NewUnprocessibleEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}

	MemberNew, err := m.MemberService.SignUp(&member)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, utils.Response(http.StatusCreated, "Member registration successfully", MemberNew))
}

func (m *memberController) ActivatedMember(c *gin.Context)  {
	param := c.Param("memberId")
	id,errParse := strconv.Atoi(param)
	if errParse != nil {
		log.Println("Failed to converted to int")
		c.JSON(http.StatusInternalServerError, gin.H{"code" : 500, "message" : "Internal Server Error"})
	}

	err := m.MemberService.ActivatedMember(id)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Member activated successfully"})
}
