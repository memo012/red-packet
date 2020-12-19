package service

import (
	"github.com/memo012/red-packet/resk/constant"
	"github.com/shopspring/decimal"
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

//转账对象
type AccountTransferDTO struct {
	TradeNo     string              `validate:"required"`         //交易单号 全局不重复字符或数字，唯一性标识
	TradeBody   TradeParticipator   `validate:"required"`         //交易主体
	TradeTarget TradeParticipator   `validate:"required"`         //交易对象
	AmountStr   string              `validate:"required,numeric"` //交易金额,该交易涉及的金额
	Amount      decimal.Decimal     ``                            //交易金额,该交易涉及的金额
	ChangeType  constant.ChangeType `validate:"required,numeric"` //流水交易类型，0 创建账户，>0 为收入类型，<0 为支出类型，自定义
	ChangeFlag  constant.ChangeFlag `validate:"required,numeric"` //交易变化标识：-1 出账 1为进账，枚举
	Decs        string              ``                            //交易描述
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

//账户
type AccountDTO struct {
	AccountNo    string          //账户编号,账户唯一标识
	AccountName  string          //账户名称,用来说明账户的简短描述,账户对应的名称或者命名，比如xxx积分、xxx零钱
	AccountType  int             //账户类型，用来区分不同类型的账户：积分账户、会员卡账户、钱包账户、红包账户
	CurrencyCode string          //货币类型编码：CNY人民币，EUR欧元，USD美元 。。。
	UserId       string          //用户编号, 账户所属用户
	Username     string          //用户名称
	Balance      decimal.Decimal //账户可用余额
	Status       int             //账户状态，账户状态：0账户初始化，1启用，2停用
	CreatedAt    time.Time       //创建时间
	UpdatedAt    time.Time       //更新时间
}