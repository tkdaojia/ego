package model

type J_App_Link struct {
	Data string `json:"data"`
}

type J_App_Link_Info struct {
	Appid     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}

type J_App_Link_User struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Truename string `json:"truename"`
	Role     string `json:"role"`
	Sex      int    `json:"sex"`
	Status   int    `json:"status"`
}
