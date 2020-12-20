package api

import (
	"github.com/gin-gonic/gin"
	"github.com/memo012/red-packet/resk/constant"
	service2 "github.com/memo012/red-packet/resk/core/service"
	"github.com/memo012/red-packet/resk/core/service/impl"
	"github.com/memo012/red-packet/resk/infra"
	"github.com/memo012/red-packet/resk/infra/base"
	"net/http"
)

func init() {
	infra.RegisterApi(new(AccountApi))
}

type AccountApi struct {
	service service2.AccountService
}

func (a *AccountApi) Init() {
	a.service = new(impl.AccountsService)
	groupRouter := base.Gin().Group("/v1/account")
	groupRouter.POST("/create", a.createHandler)
	groupRouter.POST("/transfer", a.transferHandler)
	groupRouter.GET("/envelope/get", a.getEnvelopeAccountHandler)
	//groupRouter.Get("/get", a.getAccountHandler)
}

// 账户创建的接口：/v1/account/create
func (a *AccountApi) createHandler(context *gin.Context) {
	var account service2.AccountCreatedDTO
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err := context.ShouldBindJSON(&account); err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		context.JSON(http.StatusOK, r)
		return
	}
	// 执行账户的代码
	res, err := a.service.CreateAccount(account)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		context.JSON(http.StatusOK, r)
		return
	}
	r.Data = res
	context.JSON(http.StatusOK, r)
}

// 转账接口：/v1/account/transfer
func (a *AccountApi) transferHandler(context *gin.Context) {
	//获取请求参数，
	var account service2.AccountTransferDTO
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err := context.ShouldBindJSON(&account); err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		context.JSON(http.StatusOK, r)
		return
	}
	//执行转账逻辑
	status, err := a.service.Transfer(account)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if status != constant.TransferredStatusSuccess {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		context.JSON(http.StatusOK, r)
	}
	r.Data = status
	context.JSON(http.StatusOK, r)
}

// 查询红包账户接口：/v1/account/envelope/get
//查询红包账户的web接口: /v1/account/envelope/get
func (a *AccountApi) getEnvelopeAccountHandler(ctx *gin.Context) {
	userId := ctx.Param("userId")
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if userId == "" {
		r.Code = base.ResCodeRequestParamsError
		r.Message = "用户ID不能为空"
		ctx.JSON(http.StatusOK, r)
		return
	}
	account := a.service.GetEnvelopeAccountByUserId(userId)
	r.Data = account
	ctx.JSON(http.StatusOK, r)
}

// 查询账户信息接口：/v1/account/get
func (a *AccountApi) getAccountHandler(context *gin.Context) {
	accountNo := context.Param("accountNo")
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if accountNo == "" {
		r.Code = base.ResCodeRequestParamsError
		r.Message = "账户编号不能为空"
		context.JSON(http.StatusOK, r)
		return
	}
	account := a.service.GetAccount(accountNo)
	r.Data = account
	context.JSON(http.StatusOK, r)
}
