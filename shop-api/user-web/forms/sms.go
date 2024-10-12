package forms

type SendSmsForm struct {
	Mobile string `json:"mobile" binding:"required,mobile"`
	Type   uint   `json:"type" binding:"required,oneof=1 2"`
}
