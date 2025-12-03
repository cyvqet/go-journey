package web

import (
	"errors"
	"log"
	"net/http"
	"webook/internal/domain"
	"webook/internal/service"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

const (
	emailRegex    = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	passwordRegex = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (u *UserHandler) RegisterRouter(r *gin.Engine) {
	ug := r.Group("/user")

	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.POST("/profile", u.Profile)
}

func (u *UserHandler) Signup(c *gin.Context) {

	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效请求"})
		return
	}

	ok, err := ValidateEmail(req.Email)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "邮箱格式错误"})
		return
	}

	ok, err = ValidatePassword(req.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "密码格式错误"})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"message": "两次输入的密码不一致"})
		return
	}

	err = u.svc.SignUp(c.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		c.JSON(http.StatusOK, gin.H{"message": "邮箱冲突"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统异常"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func (u *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效请求"})
		return
	}

	err := u.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvaildUserOrPassword) {
			c.JSON(http.StatusOK, gin.H{"message": "用户名/密码错误"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登陆成功"})
}

func (u *UserHandler) Edit(c *gin.Context) {}

func (u *UserHandler) Profile(c *gin.Context) {}

func ValidatePassword(password string) (bool, error) {
	re := regexp2.MustCompile(passwordRegex, 0)
	return re.MatchString(password)
}

func ValidateEmail(email string) (bool, error) {
	re := regexp2.MustCompile(emailRegex, 0)
	return re.MatchString(email)
}
