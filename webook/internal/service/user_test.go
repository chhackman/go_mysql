package service

import (
	"awesomeProject/webook/internal/domain"
	"awesomeProject/webook/internal/repository"
	repomocks "awesomeProject/webook/internal/repository/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func Test_userService_Login(t *testing.T) {
	//做成一个每个测试用例都用到的时间
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		//输入
		//ctx      context.Context
		email    string
		password string

		//输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功", //用户和密码是对的
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1@qq.com").
					Return(domain.User{
						Email:    "1@qq.com",
						Password: "$2a$10$3N4HZCQMaQ0IaB/dnPuAHuqnihBKGWGAs0gHrpiJWqb1am9BE8Coi",
						Phone:    "111111111",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "1@qq.com",
			password: "Huawei@123",

			wantUser: domain.User{
				Email:    "1@qq.com",
				Password: "$2a$10$3N4HZCQMaQ0IaB/dnPuAHuqnihBKGWGAs0gHrpiJWqb1am9BE8Coi",
				Phone:    "111111111",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在", //用户和密码是对的
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "1@qq.com",
			password: "Huawei@123",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1@qq.com").
					Return(domain.User{}, errors.New("mock db 错误"))
				return repo
			},
			email:    "1@qq.com",
			password: "Huawei@123",

			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1@qq.com").
					Return(domain.User{
						Email:    "1@qq.com",
						Password: "$2a$10$3N4HZCQMaQ0IaB/dnPuAHuqnihBKGWGAs0gHrpiJWqb1am9BE8Coi11111",
						Phone:    "111111111",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "1@qq.com",
			password: "Huawei@1w22223",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			//具体测试代码
			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("Huawei@123"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
