package response

import (
	egocode "ego/src/boot/global/code"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type ApiSuccess struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type ApiOk struct {
	Code int `json:"code"`
}

type ApiError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type Responselink struct {
	Code    int         `json:"code"`
	Success bool        `json:"success""`
	Data    interface{} `json:"data"`
	Msg     string      `json:"msg"`
}

type ResponseTable struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Count int64       `json:"count""`
}

type ResponseMsg struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success""`
}

type ResponseEditor struct {
	Location string `json:"location"`
}

const (
	ERROR   = -1
	SUCCESS = 0

	CodeSuccess = 0  // 成功
	CodeError   = -1 // 通用失败
)

type ResponseE struct {
	Code int    `json:"code"` // 业务状态码
	Msg  string `json:"msg"`  // 提示信息
	Data any    `json:"data"` // 核心数据域（没有时返回 null 或空对象）
}

// 成功响应（支持不传数据，或传任意数据）
func OnSuccess(c *gin.Context, data ...any) {
	var resData any
	if len(data) > 0 {
		resData = data[0]
	} else {
		resData = map[string]any{}
	}

	c.JSON(http.StatusOK, ResponseE{
		Code: CodeSuccess,
		Msg:  "success",
		Data: resData,
	})
}

// 失败响应
func OnFailure(c *gin.Context, errMsg string) {
	c.JSON(http.StatusOK, ResponseE{
		Code: CodeError,
		Msg:  errMsg,
		Data: map[string]any{},
	})
}

func BackMsg(msg string, c *gin.Context) {
	var sucess bool
	if msg == "ok" {
		sucess = true
	} else {
		sucess = false
	}
	c.JSON(http.StatusOK, ResponseMsg{
		msg,
		sucess,
	})
}

func BackMsgOk(c *gin.Context) {
	BackMsg("ok", c)
}

func BackMsgErr(msg string, c *gin.Context) {
	BackMsg(msg, c)
}

func ResultTable(code int, data interface{}, msg string, count int64, c *gin.Context) {
	c.JSON(http.StatusOK, ResponseTable{
		code,
		data,
		msg,
		count,
	})
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func ResultEditor(src string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, ResponseEditor{
		src,
	})
}

func OkTableList(data interface{}, count int64, c *gin.Context) {
	ResultTable(SUCCESS, data, "", count, c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

func ErrWithPhone(message string, c *gin.Context) {
	Result(-1, nil, message, c)
}

func OkWithPhone(message string, data interface{}, c *gin.Context) {
	Result(1, data, message, c)
}

func ErrWithApp(code int, message string, c *gin.Context) {
	Result(code, nil, message, c)
}

func OkWithApp(message string, data interface{}, c *gin.Context) {
	Result(egocode.SUCESS, data, message, c)
}

func OkWithAppList(data interface{}, count int64, c *gin.Context) {
	ResultTable(egocode.SUCESS, data, "", count, c)
}

func OkWithLink(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Responselink{
		0,
		true,
		data,
		"操作成功",
	})
}

func Ok(c *gin.Context) {
	c.JSON(http.StatusOK, ApiOk{
		0,
	})
}

func OkApi(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, ApiSuccess{
		0,
		data,
	})
}

func ErrApi(msg string, c *gin.Context) {
	c.JSON(http.StatusOK, ApiError{
		-1,
		"Error:" + msg,
	})
}

func errHtmlBack(c *gin.Context, msg string, back bool) {
	c.HTML(200, "system/error.htm", gin.H{
		"msg":  msg,
		"back": back,
	})
}

func ErrHtml(c *gin.Context, msg string) {
	errHtmlBack(c, msg, false)
}

func OkHtml(c *gin.Context, msg string) {
	c.HTML(200, "system/success.htm", gin.H{
		"msg": msg,
	})
}
