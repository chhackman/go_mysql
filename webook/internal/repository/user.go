package repository

import (
	"awesomeProject/webook/internal/domain"
	"awesomeProject/webook/internal/repository/cache"
	"awesomeProject/webook/internal/repository/dao"
	"database/sql"
	"fmt"
	"golang.org/x/net/context"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, Phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	Edit(ctx context.Context, id int64, u domain.User) error
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, Phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, Phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

//先从cache里面找
//再从dao里面找
//找到了回写cache

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		//必然有数据
		return u, nil
	}
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	//ctime := changeTime(ue.Ctime)
	//utime := changeTime(ue.Utime)
	//u = domain.User{
	//	Id:    ue.Id,
	//	Email: ue.Email,
	//	//Password: u.Password,
	//	Nickname: ue.Nickname,
	//	Birthday: ue.Birthday,
	//	Abstract: ue.Abstract,
	//	Ctime:    ctime,
	//	Utime:    utime,
	//}
	u = r.entityToDomain(ue)
	_ = r.cache.Set(ctx, u)
	//if err != nil {
	//	//打日志监控
	//}

	//go func() {
	//	err = r.cache.Set(ctx, u)
	//	if err != nil {
	//		//打日志监控
	//	}
	//
	//}()
	return u, nil
}

func changeTime(srcTime int64) string {
	millisecondTimestamp := int64(srcTime) // 假设这是一个毫秒时间戳示例

	// 将毫秒时间戳除以1000得到秒级时间戳
	secondTimestamp := millisecondTimestamp / 1000

	// 将秒级时间戳转换为time.Time类型
	timestampTime := time.Unix(secondTimestamp, 0)

	// 载入中国时区
	chinaTimezone, _ := time.LoadLocation("Asia/Shanghai")

	// 将时间调整为中国时区
	timestampTimeInChina := timestampTime.In(chinaTimezone)

	// 格式化为日期时间字符串
	formattedDate := timestampTimeInChina.Format("2006-01-02 15:04:05") // 格式可以根据需求进行调整

	fmt.Println("Formatted Date in China Timezone:", formattedDate)
	return formattedDate

}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	println(1111)
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

//func (r *CachedUserRepository) FindById(int64) {
//	//先从cache里面找
//	//再从dao里面找
//	//找到了回写cache
//}

func (r *CachedUserRepository) Edit(ctx context.Context, id int64, u domain.User) error {
	return r.dao.EditUserProfile(ctx, id, dao.User{
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		Abstract: u.Abstract,
	})
}

func (r *CachedUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			// 我确实有手机号
			Valid: u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
	}
}

func (r *CachedUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}
