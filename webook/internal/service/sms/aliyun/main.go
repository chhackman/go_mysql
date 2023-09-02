package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"golang.org/x/net/context"
	"os"
)

func main() {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"))
	client, err := dysmsapi.NewClientWithOptions("cn-hangzhou", config, credential)

	if err != nil {
		panic(err)
	}

	SignName := "阿里云短信测试"
	TemplateCode := "SMS_154950909"
	PhoneNumbers := "18526282359"
	TemplateParam := "{\"code\":\"1234\"}"

	aliyService := NewService(SignName, client)
	aliyService.Send(context.Background(), TemplateCode, TemplateParam, PhoneNumbers)

}
