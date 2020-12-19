package accounts

import (
	data "github.com/memo012/red-packet/resk/core/data/accounts"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountLogDao struct {
	runner *dbx.TxRunner
}

//通过流水编号查询流水记录
func (dao *AccountLogDao) GetOne(logNo string) *data.AccountLog {
	a := &data.AccountLog{LogNo: logNo}
	ok, err := dao.runner.GetOne(a)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a

}

//通过交易编号来查询流水记录
func (dao *AccountLogDao) GetByTradeNo(tradeNo string) *data.AccountLog {
	sql := "select * from account_log where trade_no=?"
	out := &data.AccountLog{}
	ok, err := dao.runner.Get(out, sql, tradeNo)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return out

}

//流水记录的写入
func (dao *AccountLogDao) Insert(l *data.AccountLog) (id int64, err error) {
	rs, err := dao.runner.Insert(l)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}
