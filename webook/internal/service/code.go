package service

import (
	"awesomeProject/webook/internal/repository"
	"awesomeProject/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

var (
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
)

const tpl = "SMS_154950909"

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smSvc sms.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smSvc,
	}
}

//发送验证码
//biz区别业务场景
func (svc *CodeService) Send(ctx context.Context,
	////biz区别业务场景
	biz string,
	//这个码谁来管,谁来生成
	phone string) error {
	//两个步骤，生成一个验证码
	code := svc.generateCode()
	//塞进去Redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		//有问题
		return err
	}
	//发送出去
	err = svc.smsSvc.Send(ctx, tpl, code, phone)
	return err

}
func (svc *CodeService) generateCode() string {
	//六位数,num在0,99999之间，包含0和999999
	num := rand.Intn(1000000)

	//不够6位的，加上前导0
	//000001
	return fmt.Sprintf("%06d", num)
}

func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}
