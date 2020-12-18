package service

import (
	"github.com/memo012/red-packet/resk/constant"
	"time"
)

type AccountService interface {
	// 创建账户
	CreateAccount(dto AccountCreatedDTO) (*AccountDTO, error)
	// 转账交易
	Transfer(dto AccountTransferDTO) (constant.TransferredStatus, error)
	// 储值
	StoreValue(dto AccountTransferDTO) (constant.TransferredStatus, error)
	// 账户查询
	GetEnvelopeAccountByUserId(userId string) *AccountDTO
}

// 账户交易的参与者
type TradeParticipator struct {
	AccountNo string
	UserId    string
	UserName  string
}

// 账户转账
type AccountTransferDTO struct {
	TradeNo     string
	TradeBody   TradeParticipator
	TradeTarget TradeParticipator
	AmountStr   string
	ChangeType  constant.ChangeType
	ChangeFlag  constant.ChangeFlag
	Desc        string
}

// 账户创建
type AccountCreatedDTO struct {
	UserID       string
	UserName     string
	AccountName  string
	AccountType  int
	CurrencyCode string
	Amount       string
}

type AccountDTO struct {
	AccountCreatedDTO
	AccountNo   string
	CreatedTime time.Time
}
