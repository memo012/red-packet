package accounts

import (
	"github.com/memo012/red-packet/resk/constant"
	"github.com/memo012/red-packet/resk/core/service"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAccountDomain_Create(t *testing.T) {
	dto := service.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试用户",
		Balance:  decimal.NewFromFloat(0),
		Status:   1,
	}
	domain := new(accountDomain)
	Convey("账户创建", t, func() {
		rdto, err := domain.Create(dto)
		So(err, ShouldBeNil)
		So(rdto, ShouldNotBeNil)
		So(rdto.Balance.String(), ShouldEqual, dto.Balance.String())
		So(rdto.UserId, ShouldEqual, dto.UserId)
		So(rdto.Username, ShouldEqual, dto.Username)
		So(rdto.Status, ShouldEqual, dto.Status)
	})
}

func TestAccountDomain_Transfer(t *testing.T) {
	//2个账户，交易主体账户要有余额
	adto1 := &service.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试用户1",
		Balance:  decimal.NewFromFloat(100),
		Status:   1,
	}
	adto2 := &service.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试用户2",
		Balance:  decimal.NewFromFloat(100),
		Status:   1,
	}
	domain := accountDomain{}
	Convey("转账测试", t, func() {
		//账户1的准备
		dto1, err := domain.Create(*adto1)
		So(err, ShouldBeNil)
		So(dto1, ShouldNotBeNil)
		So(dto1.Balance.String(), ShouldEqual, adto1.Balance.String())
		So(dto1.UserId, ShouldEqual, adto1.UserId)
		So(dto1.Username, ShouldEqual, adto1.Username)
		So(dto1.Status, ShouldEqual, adto1.Status)
		adto1 = dto1
		//账户2的准备
		dto2, err := domain.Create(*adto2)
		So(err, ShouldBeNil)
		So(dto2, ShouldNotBeNil)
		So(dto2.Balance.String(), ShouldEqual, adto2.Balance.String())
		So(dto2.UserId, ShouldEqual, adto2.UserId)
		So(dto2.Username, ShouldEqual, adto2.Username)
		So(dto2.Status, ShouldEqual, adto2.Status)
		adto2 = dto2
		//转账操作验证
		// 1. 余额充足，金额转入其他账户
		Convey("余额充足，金额转入其他账户", func() {
			amount := decimal.NewFromFloat(1)
			body := service.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				UserName:  adto1.Username,
			}
			target := service.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId:    adto2.UserId,
				UserName:  adto2.Username,
			}
			dto := service.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				Amount:      amount,
				ChangeType:  constant.ChangeType(-1),
				ChangeFlag:  constant.FlagTransferOut,
				Decs:        "转账",
			}
			status, err := domain.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, constant.TransferredStatusSuccess)
			//实际余额更新后的预期值
			a2 := domain.GetAccount(adto1.AccountNo)
			So(a2, ShouldNotBeNil)
			So(a2.Balance.String(),
				ShouldEqual,
				adto1.Balance.Sub(amount).String())

		})
		// 2. 余额不足，金额转出
		Convey("余额不足，金额转出", func() {
			amount := adto1.Balance
			amount = amount.Add(decimal.NewFromFloat(200))
			body := service.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				UserName:  adto1.Username,
			}
			target := service.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId:    adto2.UserId,
				UserName:  adto2.Username,
			}
			dto := service.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				Amount:      amount,
				ChangeType:  constant.ChangeType(-1),
				ChangeFlag:  constant.FlagTransferOut,
				Decs:        "转账",
			}
			status, err := domain.Transfer(dto)
			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, constant.TransferredStatusSufficientFunds)
			//实际余额更新后的预期值
			a2 := domain.GetAccount(adto1.AccountNo)
			So(a2, ShouldNotBeNil)
			So(a2.Balance.String(),
				ShouldEqual,
				adto1.Balance.String())

		})
		// 3. 充值
		Convey("充值", func() {
			amount := decimal.NewFromFloat(11.1)
			body := service.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				UserName:  adto1.Username,
			}
			target := service.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId:    adto2.UserId,
				UserName:  adto2.Username,
			}
			dto := service.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				Amount:      amount,
				ChangeType:  constant.AccountStoreValue,
				ChangeFlag:  constant.FlagTransferIn,
				Decs:        "储值",
			}
			status, err := domain.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, constant.TransferredStatusSuccess)
			//实际余额更新后的预期值
			a2 := domain.GetAccount(adto1.AccountNo)
			So(a2, ShouldNotBeNil)
			So(a2.Balance.String(),
				ShouldEqual,
				adto1.Balance.Add(amount).String())
		})
	})
}
