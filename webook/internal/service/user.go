package service

import (
	"awesomeProject/webook/internal/domain"
	"awesomeProject/webook/internal/repository"
	"awesomeProject/webook/internal/service/sms/aliyun"
	"context"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"golang.org/x/crypto/bcrypt"
	"os"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
var ErrInvalidUserNotFund = errors.New("不存在该用户")

type UserService interface {
	Login(ctx context.Context, email, password string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	SendSMS(ctx context.Context)
	EditUserProfile(ctx context.Context, id int64, u domain.User) error
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	//先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	//你要考虑加密放在哪些的问题
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
	//println(1111)
	//svc.repo.Create(ctx, u)

	//return svc.repo.Create(ctx, u)
}

func (svg *userService) EditUserProfile(ctx context.Context, id int64, u domain.User) error {
	return svg.repo.Edit(ctx, id, u)
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {

	u, err := svc.repo.FindById(ctx, id)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserNotFund
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (svc *userService) SendSMS(ctx context.Context) {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"))
	client, err := dysmsapi.NewClientWithOptions("cn-hangzhou", config, credential)

	if err != nil {
		panic(err)
	}

	SignName := "阿里云短信测试"
	TemplateCode := "SMS_154950909"
	PhoneNumbers := "1******"
	TemplateParam := "{\"code\":\"34561\"}"

	aliyService := aliyun.NewService(SignName, client)
	err = aliyService.Send(ctx, TemplateCode, TemplateParam, PhoneNumbers)
	if err != nil {
		panic(err)
	}

}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	//svc.repo.FindByPhone(ctx, phone)
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, nil
	}
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return u, err
	}
	return svc.repo.FindByPhone(ctx, phone)
}
