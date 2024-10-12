package forms

type PasswordLoginForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"` // 手机号码格式有规矩可循，自定义validator
	Password string `json:"password" form:"password" binding:"required,min=3,max=10"`
}

type RegisterForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"`
	Password string `json:"password" form:"password" binding:"required,min=3,max=20"`
	Code     string `json:"code" form:"code" binding:"required,min=6,max=6"`
}
