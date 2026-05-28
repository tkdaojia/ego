package msg

const (
	SqlFindErr     = "数据不存在"
	SqlFindListErr = "数据列表异常"
	SqlUpdateErr   = "数据更新失败"
	SqlCountErr    = "统计数量异常"
	SqlCreateErr   = "数据新增异常"
	SqlDeleteErr   = "数据删除异常"

	SqlErr = "SQL异常"

	QueryParamErr = "缺少必要的查询参数"

	NameExistsErr = "名称重复已存在"

	ReqParamErr = "参数格式错误或必填项缺失"

	IdInvalidErr = "ID数据异常"

	UserDataLost      = "用户信息丢失"
	UserPasswordEmpty = "用户密码不能唯恐"

	PasswordEncryptErr = "密码加密失败，请稍后重试"
)
