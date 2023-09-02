package aliyun

import (
	"context"
	"fmt"
	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type Service struct {
	signName string
	client   *dysmsapi.Client
}

func NewService(signName string, client *dysmsapi.Client) *Service {
	return &Service{
		signName: signName,
		client:   client,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args string, numbers ...string) error {
	req := dysmsapi.CreateSendSmsRequest()
	req.Scheme = "https"
	req.SignName = s.signName
	req.TemplateCode = tpl
	req.TemplateParam = args
	if len(numbers) > 1 {
		return fmt.Errorf("由于阿里云模板只能发送单人，请写一个收件人")
	}
	req.PhoneNumbers = numbers[0]
	response, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	fmt.Print("response is %#v\n", response)
	return nil

}
