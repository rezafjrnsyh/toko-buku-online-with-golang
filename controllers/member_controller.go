package controllers

import (
	"database/sql"
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
	BUY_BOOK_PATH = "/member/book/buy/:id"
)

func NewMemberController(db *sql.DB, r *gin.RouterGroup) {
	Controller := memberController{MemberService: services.NewMemberService(db)}
	r.POST(SIGN_UP_PATH, Controller.SignUpMember)
	r.POST(SIGN_IN_PATH, Controller.SignInMember)
	r.PUT(BUY_BOOK_PATH, Controller.Buy)
}

func (m *memberController) Buy(c *gin.Context) {
	param := c.Param("id")
	id,err := strconv.Atoi(param)
	if err != nil {
		log.Println("Failed to converted to int")
		c.JSON(http.StatusInternalServerError, gin.H{"code" : 500, "message" : "Internal Server Error"})
	}

	var qty *domain.ReqBuy
	errBind := c.ShouldBindJSON(&qty)
	if errBind != nil {
		theErr := utils.NewUnprocessibleEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}
	trx, errBuy := m.MemberService.BuyBook(id, qty)
	if errBuy != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal server error",
		})
	}
	c.JSON(http.StatusOK, utils.Response(200, "Your transaction was successful", trx))
}

func (m *memberController) SignInMember(c *gin.Context)  {
	var member domain.MemberLogin
	errBind := c.ShouldBindJSON(&member)
	if errBind != nil {
		theErr := utils.NewUnprocessibleEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}
	err := m.MemberService.SignIn(&member)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Please validate email or password correct!",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "Verified"})

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
