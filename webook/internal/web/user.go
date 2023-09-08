package web

import (
	"awesomeProject/webook/internal/domain"
	"awesomeProject/webook/internal/service"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

const biz = "login"

//确保UserHandler实现了handlder接口
var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	codesvc     *service.CodeService
}

type UserCliams struct {
	jwt.RegisteredClaims
	//声明你自己的要加进去的token的数据
	Uid int64
	//自己随便加
	UserAgent string
}

func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		codesvc:     codeSvc,
	}
}

//func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
//	ug.GET("/profile", u.Profile)
//	ug.POST("/login", u.Login)
//	ug.POST("/edit", u.Edit)
//	ug.POST("signup", u.SignUp)
//}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	//ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	//ug.GET("/sms", u.SendSMS)
	ug.GET("/logout", u.Logout)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)

	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
}
func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//这面加校验，手机号之类的
	ok, err := u.codesvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}

	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 5,
		Msg:  "验证码校验通过",
	})

}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	err := u.codesvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
		return
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁请稍后再试",
		})
		return
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return

	}

}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	//Bind方法会根据Content-Type来解析你的数据到req里面
	//解析错了，就会直接写回一个400的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "您的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		//记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位,包括数字、特殊字符")
		return
	}
	//ctx.String(http.StatusOK, "注册成功")
	//fmt.Printf("%v", req)
	//这面就是数据库操作

	//调用一下svc的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:email`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	println(req.Email)
	println(req.Password)
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	fmt.Printf("%v", user)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或者密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//步骤2
	//这里登录成功了，设置session
	//sess := sessions.Default(ctx)
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		//Secure: true,
		HttpOnly: true,
		MaxAge:   120,
	})
	sess.Save()
	//sess.Save()
	ctx.String(http.StatusOK, "登录成功")

	return
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:email`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	println(req.Email)
	println(req.Password)
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	fmt.Printf("%v", user)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或者密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//步骤2
	//这里登录成功了，设置session
	//sess := sessions.Default(ctx)
	//sess := sessions.Default(ctx)
	//sess.Set("userId", user.Id)
	//sess.Save()
	//ctx.String(http.StatusOK, "登录成功")

	if err := u.setJWTToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	//fmt.Printf(user)
	ctx.String(http.StatusOK, "登录成功")
	//println(user)
	//fmt.Printf(user)
	return
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	//这里使用jwt设置登录态
	cliams := UserCliams{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, cliams)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil

}

// CustomValidator 是自定义的验证器结构
type CustomValidator struct {
	Validator *validator.Validate
}

// Engine 自定义验证器的实例化方法
func (cv *CustomValidator) Engine() interface{} {
	return cv.Validator
}

func (cv *CustomValidator) ValidateStruct(obj interface{}) error {
	if err := cv.Validator.Struct(obj); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		var errMsgs []string
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s: %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf(strings.Join(errMsgs, "; "))
	}

	return nil
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	claims, ok := c.(*UserCliams)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	println(claims.Uid)
	ctx.String(http.StatusOK, "你的profile")
}

func (u *UserHandler) SendSMS(ctx *gin.Context) {
	u.svc.SendSMS(ctx)
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	id := sess.Get("userId")
	value, ok := id.(int64)
	if !ok {
		//println(1111)
		return
	}
	user, err := u.svc.Profile(ctx, value)
	if err == service.ErrInvalidUserNotFund {
		ctx.String(http.StatusOK, "没有查询到该用户")
		return
	}
	//ctime = user.Ctime.Format("2006-01-02 15:04:05")
	//utime = user.Utime.Format("2006-01-02 15:04:05")
	//message := string[]{
	//	Id:       str(user.Id),
	//	Email:    user.Email,
	//	Password: user.Password,
	//
	//	//添加如下字段，用户昵称，生日和个人简介
	//	Nickname: user.Nickname,
	//	Birthday: user.Birthday,
	//	Abstract: user.Abstract,
	//	Ctime:    ctime,
	//	Utime:    utime,
	//}
	type UserReq struct {
		Id    int64
		Email string
		//Password string

		//添加如下字段，用户昵称，生日和个人简介
		Nickname string
		Birthday string
		Abstract string
		Ctime    string
		Utime    string
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	userReq := UserReq{
		Id:       user.Id,
		Email:    user.Email,
		Nickname: user.Nickname,
		Birthday: user.Birthday,
		Abstract: user.Abstract,
		Ctime:    user.Ctime.Format("2006-01-02 15:04:05"),
		Utime:    user.Utime.Format("2006-01-02 15:04:05"),
	}
	ctx.JSON(http.StatusOK, userReq)

}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditUserProfile struct {
		Nickname string `json:"nickname" binding:"required,customNicknameValid"`
		Birthday string `json:"birthday" binding:"required,customBirthdayValid"`
		Abstract string `json:"abstract" binding:"required,customAbstractValid"`
	}
	// 使用自定义验证器
	validatorInstance := validator.New()
	customValidator := &CustomValidator{Validator: validatorInstance}

	// 注册自定义验证规则
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("customNicknameValid", func(fl validator.FieldLevel) bool {
			// 验证昵称长度：昵称字符串长度小于10，英文字符和中文长度一样
			return utf8.RuneCountInString(fl.Field().String()) < 10
		})

		_ = v.RegisterValidation("customBirthdayValid", func(fl validator.FieldLevel) bool {
			// 验证生日格式：YYYY-MM-DD，例如 1992-01-01
			birthdayPattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`, 0)
			flag, _ := birthdayPattern.MatchString(fl.Field().String())
			return flag
		})

		_ = v.RegisterValidation("customAbstractValid", func(fl validator.FieldLevel) bool {
			// 验证个人简历：长度小于500，英文字符和中文长度一样
			return utf8.RuneCountInString(fl.Field().String()) <= 500
			//return len(fl.Field().String()) <= 500
		})
	}
	ctx.Set("custom_validator", customValidator)

	var req EditUserProfile
	sess := sessions.Default(ctx)
	//if err := ctx.Bind(&req); err != nil {
	//	//println(1111)
	//	return
	//}
	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		//println(1111)
		// 解析验证错误信息
		var errMsg string
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErr {
				switch e.Tag() {
				case "customNicknameValid":
					errMsg = "昵称字符串长度小于10，英文字符和中文长度一样"
				case "customBirthdayValid":
					errMsg = "生日格式：YYYY-MM-DD，例如 1992-01-01"
				case "customAbstractValid":
					errMsg = "验证个人简历长度小于500，英文字符和中文长度一样"
				}
			}
		}
		ctx.String(http.StatusOK, errMsg)
		//c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})

		return
	}
	//println(1111)
	id := sess.Get("userId")
	println(id)
	value, ok := id.(int64)
	if !ok {
		println(1111)
		return
	}
	//Bind方法会根据Content-Type来解析你的数据到req里面
	//解析错了，就会直接写回一个400的错误
	//if err := ctx.Bind(&req); err != nil {
	//	return
	//}
	//ok, err := u.emailExp.MatchString(req.Email)
	//todo,校验

	//调用一下svc的方法
	err := u.svc.EditUserProfile(ctx, value, domain.User{
		Nickname: req.Nickname,
		Birthday: req.Birthday,
		Abstract: req.Abstract,
	})

	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "修改个人信息成功")
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	// 我可以随便设置值了
	// 你要放在 session 里面的值
	sess.Options(sessions.Options{
		//Secure: true,
		//HttpOnly: true,
		MaxAge: -1,
	})
	sess.Save()
	println(111)
	ctx.String(http.StatusOK, "退出登录成功")
}

//
//func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
//	ug := server.GET("/users")
//	ug.GET("/profile", u.Profile)
//	//ug.POST("")
//}
//
//func (u *UserHandler) SignLogin(ctx *gin.Context) {
//
//}
//
//func (u *UserHandler) Login(ctx *gin.Context) {
//
//}
//func (u *UserHandler) Profile(ctx *gin.Context) {
//
//}
