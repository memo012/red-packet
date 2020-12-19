package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/memo012/red-packet/resk/constant"
	"github.com/memo012/red-packet/resk/core/models/accounts"
	"github.com/memo012/red-packet/resk/infra/base"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var _ AccountService = new(AccountsService)

type AccountsService struct {
}

func (a *AccountsService) CreateAccount(dto AccountCreatedDTO) (*AccountDTO, error) {
	domain := accounts.AccountDomain{}
	// 验证输入参数
	err := base.Validate().Struct(&dto)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, v := range errs {
				logrus.Error(v)
			}
		}
		return nil, err
	}
	// 执行账户创建的业务逻辑
	amount, err := decimal.NewFromString(dto.Amount)
	if err != nil {
		return nil, err
	}
	account := AccountDTO{
		UserId:       dto.UserID,
		Username:     dto.UserName,
		AccountType:  dto.AccountType,
		CurrencyCode: dto.CurrencyCode,
		Status:       1,
		Balance:      amount,
	}
	rdto, err := domain.Create(account)
	return rdto, err
}

func (a *AccountsService) Transfer(dto AccountTransferDTO) (constant.TransferredStatus, error) {
	domain := accounts.AccountDomain{}
	// 验证输入参数
	err := base.Validate().Struct(&dto)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, v := range errs {
				logrus.Error(v)
			}
		}
		return constant.TransferredStatusFailure, err
	}
	// 执行转账逻辑
	amount, err := decimal.NewFromString(dto.AmountStr)
	if err != nil {
		return constant.TransferredStatusFailure, err
	}
	dto.Amount = amount
	if dto.ChangeFlag == constant.FlagTransferOut {
		if dto.ChangeType > 0 {
			return constant.TransferredStatusFailure, errors.New("changeType必须小于零")
		}
	} else {
		if dto.ChangeType < 0 {
			return constant.TransferredStatusFailure, errors.New("changeType必须大于零")
		}
	}
	status, err := domain.Transfer(dto)
	return status, err
}

func (a *AccountsService) StoreValue(dto AccountTransferDTO) (constant.TransferredStatus, error) {
	panic("implement me")
}

func (a *AccountsService) GetEnvelopeAccountByUserId(userId string) *AccountDTO {
	panic("implement me")
}
