package accounts

import (
	"errors"
	"github.com/memo012/red-packet/resk/constant"
	data2 "github.com/memo012/red-packet/resk/core/data"
	"github.com/memo012/red-packet/resk/core/service"
	"github.com/memo012/red-packet/resk/infra/base"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountDomain struct {
	account    data2.Account
	accountLog data2.AccountLog
}

// 创建logNo 的逻辑
func (domain *AccountDomain) createAccountLogNo() {
	domain.accountLog.LogNo = ksuid.New().Next().String()
}

// 创建accountNo 的逻辑
func (domain *AccountDomain) createAccountNo() {
	domain.account.AccountNo = ksuid.New().Next().String()
}

// 创建流水记录
func (domain *AccountDomain) createAccountLog() {
	// 通过account来创建流水 创建账户逻辑在前
	domain.accountLog = data2.AccountLog{}
	domain.createAccountLogNo()
	domain.accountLog.TradeNo = domain.accountLog.LogNo

	// 流水中的交易主体信息
	domain.accountLog.AccountNo = domain.account.AccountNo
	domain.accountLog.UserId = domain.account.UserId
	domain.accountLog.UserName = domain.account.UserName.String
	// 交易对象信息
	domain.accountLog.TargetAccountNo = domain.account.AccountNo
	domain.accountLog.TargetUserId = domain.account.UserId
	domain.accountLog.TargetUsername = domain.account.UserName.String
	// 交易金额
	domain.accountLog.Balance = domain.account.Balance
	domain.accountLog.Balance = domain.account.Balance
	// 交易变化属性
	domain.accountLog.Decs = "账户创建"
	domain.accountLog.ChangeType = constant.AccountCreated
	domain.accountLog.ChangeFlag = constant.FlagAccountCreated
}

// 账户创建的业务逻辑
func (domain *AccountDomain) Create(dto service.AccountDTO) (*service.AccountDTO, error) {
	// 创建账户持久化对象
	domain.account = data2.Account{}
	domain.account.FromDTO(&dto)
	domain.createAccountNo()
	domain.account.UserName.Valid = true
	// 创建对象流水持久化
	domain.createAccountLog()

	accountDao := AccountDao{}
	accountLogDao := AccountLogDao{}
	var rdto *service.AccountDTO
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		accountLogDao.runner = runner
		// 插入账户数据
		id, err := accountDao.Insert(&domain.account)
		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("创建账户失败")
		}
		// 如果插入成功 就插入流水数据
		id, err = accountLogDao.Insert(&domain.accountLog)
		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("创建账户流水失败")
		}
		domain.account = *accountDao.GetOne(domain.account.AccountNo)
		return nil
	})
	rdto = domain.account.ToDTO()
	return rdto, err
}

// 转账
func (a *AccountDomain) Transfer(dto service.AccountTransferDTO) (status constant.TransferredStatus, err error) {
	// 如果交易类型是支出 修正amount
	amount := dto.Amount
	if dto.ChangeFlag == constant.FlagTransferOut {
		amount = amount.Mul(decimal.NewFromFloat(-1))
	}

	// 创建流水记录
	a.accountLog = data2.AccountLog{}
	a.accountLog.FromTransferDTO(&dto)
	a.createAccountLogNo()
	// 检查余额是否足够：通过乐观锁来验证 更新余额的同时来验证余额是否足够
	err = base.Tx(func(runner *dbx.TxRunner) error {
		accountDao := AccountDao{runner: runner}
		acountLogDao := AccountLogDao{runner: runner}
		rows, err := accountDao.UpdateBalance(dto.TradeBody.AccountNo, amount)
		if err != nil {
			status = constant.TransferredStatusFailure
			return err
		}
		if rows <= 0 && dto.ChangeFlag == constant.FlagTransferOut {
			status = constant.TransferredStatusSufficientFunds
			return errors.New("余额不足")
		}
		account := accountDao.GetOne(dto.TradeBody.AccountNo)
		if account == nil {
			return errors.New("账户出错")
		}

		// 更新成功后 写入流水记录
		a.account = *account
		a.accountLog.Balance = a.account.Balance
		id, err := acountLogDao.Insert(&a.accountLog)
		if err != nil || id <= 0 {
			status = constant.TransferredStatusFailure
			return errors.New("账户流水创建失败")
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
	} else {
		status = constant.TransferredStatusSuccess
	}
	return status, err
}

// 根据账户编号查询账户信息
func (a *AccountDomain) GetAccount(accountNo string) *service.AccountDTO {
	accountDao := AccountDao{}
	var account *data2.Account
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		account = accountDao.GetOne(accountNo)
		return nil
	})
	if err != nil {
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

// 根据用户ID查询红包账户信息
func (a *AccountDomain) GetEnvelopeAccountByUserId(userId string) *service.AccountDTO {
	accountDao := AccountDao{}
	var account *data2.Account
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		account = accountDao.GetByUserId(userId, int(constant.EnvelopeAccountType))
		return nil
	})
	if err != nil {
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

// 根据流水ID来查询账户流水
func (a *AccountDomain) GetAccountLog(logNo string) *service.AccountLogDTO {
	dao := AccountLogDao{}
	var log *data2.AccountLog
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao.runner = runner
		log = dao.GetOne(logNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if log == nil {
		return nil
	}
	return log.ToDTO()
}

// 根据交易编号来查询账户流水
func (a *AccountDomain) GetAccountLogByTradeNo(tradeNo string) *service.AccountLogDTO {
	dao := AccountLogDao{}
	var log *data2.AccountLog
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao.runner = runner
		log = dao.GetByTradeNo(tradeNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if log == nil {
		return nil
	}
	return log.ToDTO()
}
