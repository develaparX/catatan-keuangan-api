package controller

import (
	"livecode-catatan-keuangan/models"
	"livecode-catatan-keuangan/models/dto"
	"livecode-catatan-keuangan/service"
	"livecode-catatan-keuangan/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
	rg      *gin.RouterGroup
}

func (u *UserController) loginHandler(ctx *gin.Context) {
	var payload dto.LoginDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.SendErrorResponse(ctx, 404)
		return
	}

	response, errors := u.service.Login(payload)
	if errors != nil {
		utils.SendErrorResponse(ctx, 404)
		return
	}
	utils.SendSingleResponse(ctx, response, 2000101)
}

func (u *UserController) registerHandler(ctx *gin.Context) {
	var payload models.User
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.SendErrorResponse(ctx, 401)
	}

	data, errors := u.service.CreateNew(payload)
	if errors != nil {
		utils.SendErrorResponse(ctx, 404)
	}
	utils.SendSingleResponse(ctx, data, 2000101)
}

func (u *UserController) Route() {
	router := u.rg.Group("/users")
	router.POST("/register", u.registerHandler)
	router.GET("/login", u.loginHandler)
}

func NewUserController(uS service.UserService, rg *gin.RouterGroup) *UserController {
	return &UserController{
		service: uS,
		rg:      rg,
	}
}
