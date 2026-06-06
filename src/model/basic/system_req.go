package model

var SystemReqDel struct {
	Id int `json:"id" binding:"required,gt=0"`
}
