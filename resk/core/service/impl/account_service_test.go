package impl

import (
	"github.com/memo012/red-packet/resk/constant"
	service2 "github.com/memo012/red-packet/resk/core/service"
	_ "github.com/memo012/red-packet/resk/testx"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAccountsService_CreateAccount(t *testing.T) {
	dto := service2.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户",
		Amount:       "100",
		AccountName:  "测试账户",
		AccountType:  2,
		CurrencyCode: "CNY",
	}
	service := new(AccountsService)
	Convey("账户创建", t, func() {
		rdto, err := service.CreateAccount(dto)
		So(err, ShouldBeNil)
		So(rdto, ShouldNotBeNil)
		So(rdto.Balance.String(), ShouldEqual, dto.Amount)
		So(rdto.UserId, ShouldEqual, dto.UserId)
		So(rdto.Username, ShouldEqual, dto.Username)
		So(rdto.Status, ShouldEqual, 1)
	})
}

//转账业务应用服务层测试用例
func TestAccountsService_Transfer(t *testing.T) {

	Convey("转账", t, func() {
		//准备2个账户
		a1 := service2.AccountCreatedDTO{
			UserId:       ksuid.New().Next().String(),
			Username:     "测试用户1",
			Amount:       "100",
			AccountName:  "测试账户1",
			AccountType:  2,
			CurrencyCode: "CNY",
		}
		a2 := service2.AccountCreatedDTO{
			UserId:       ksuid.New().Next().String(),
			Username:     "测试用户2",
			Amount:       "100",
			AccountName:  "测试账户2",
			AccountType:  2,
			CurrencyCode: "CNY",
		}
		service := new(AccountsService)
		adto1, err := service.CreateAccount(a1)
		So(err, ShouldBeNil)
		So(adto1, ShouldNotBeNil)
		adto2, err := service.CreateAccount(a2)
		So(err, ShouldBeNil)
		So(adto2, ShouldNotBeNil)

		//从账户1转入账户2一定金额，其中账户1的余额是足够的
		//Convey("余额足够，从账户1转入账户2一定金额", func() {
		//	//转账过程需要2个交易的参与者：交易主体和交易对象
		//	body := service2.TradeParticipator{
		//		AccountNo: adto1.AccountNo,
		//		UserId:    adto1.UserId,
		//		UserName:  adto1.Username,
		//	}
		//	target := service2.TradeParticipator{
		//		AccountNo: adto2.AccountNo,
		//		UserId:    adto2.UserId,
		//		UserName:  adto2.Username,
		//	}
		//	amount := decimal.NewFromFloat(10)
		//	dto := service2.AccountTransferDTO{
		//		TradeBody:   body,
		//		TradeTarget: target,
		//		TradeNo:     ksuid.New().Next().String(),
		//		AmountStr:   amount.String(),
		//		ChangeType:  constant.ChangeType(-1),
		//		ChangeFlag:  constant.FlagTransferOut,
		//		Decs:        "转出",
		//	}
		//	status, err := service.Transfer(dto)
		//	So(err, ShouldBeNil)
		//	So(status, ShouldEqual, constant.TransferredStatusSuccess)
		//
		//	ra1 := service.GetAccount(adto1.AccountNo)
		//	So(ra1, ShouldNotBeNil)
		//	So(ra1.Balance.String(), ShouldEqual, adto1.Balance.Sub(amount).String())
		//
		//})
		////从账户1转入账户2一定金额，但余额不足，转账会失败
		//Convey("余额不足，从账户1转入账户2一定金额", func() {
		//	//转账过程需要2个交易的参与者：交易主体和交易对象
		//	body := service2.TradeParticipator{
		//		AccountNo: adto1.AccountNo,
		//		UserId:    adto1.UserId,
		//		UserName:  adto1.Username,
		//	}
		//	target := service2.TradeParticipator{
		//		AccountNo: adto2.AccountNo,
		//		UserId:    adto2.UserId,
		//		UserName:  adto2.Username,
		//	}
		//	amount := adto1.Balance.Add(decimal.NewFromFloat(200))
		//	dto := service2.AccountTransferDTO{
		//		TradeBody:   body,
		//		TradeTarget: target,
		//		TradeNo:     ksuid.New().Next().String(),
		//		AmountStr:   amount.String(),
		//		ChangeType:  constant.ChangeType(-1),
		//		ChangeFlag:  constant.FlagTransferOut,
		//		Decs:        "转出",
		//	}
		//	status, err := service.Transfer(dto)
		//	So(err, ShouldNotBeNil)
		//	So(status, ShouldEqual, constant.TransferredStatusSufficientFunds)
		//
		//	ra1 := service.GetAccount(adto1.AccountNo)
		//	So(ra1, ShouldNotBeNil)
		//	So(ra1.Balance.String(), ShouldEqual, adto1.Balance.String())
		//})
		//给账户1储值
		Convey("给账户1储值", func() {
			//转账过程需要2个交易的参与者：交易主体和交易对象
			body := service2.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				UserName:  adto1.Username,
			}
			target := body
			amount := decimal.NewFromFloat(10)
			dto := service2.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				AmountStr:   amount.String(),
				ChangeType:  constant.AccountStoreValue,
				ChangeFlag:  constant.FlagTransferIn,
				Decs:        "储值",
			}
			status, err := service.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, constant.TransferredStatusSuccess)

			ra1 := service.GetAccount(adto1.AccountNo)
			So(ra1, ShouldNotBeNil)
			So(ra1.Balance.String(), ShouldEqual, adto1.Balance.Add(amount).String())

		})
	})

}
