package entity

import (
	"Source_Predict/errorcode"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 响应
type ResponseData struct {
	RetCode errorcode.ResponseCode `json:"retcode"`
	Msg     interface{}            `json:"message"` // 定义空接口是因为 message 的数据类型有很多种，这样定义限制宽松
	Data    interface{}            `json:"data,omitempty"`
}

// 响应Data
type ResourceResp struct {
	Metric  string  `json:"metric"`
	Period  uint    `json:"period"`
	Predict float64 `json:"predict"`
}

// ResponseError 错误返回响应
func ResponseError(c *gin.Context, code errorcode.ResponseCode) {
	c.JSON(http.StatusOK, &ResponseData{
		RetCode: code,
		Msg:     code.Msg(),
		Data:    nil,
	})
}

// ResponseSuccess 返回成功响应
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		RetCode: errorcode.CodeSuccess,
		Msg:     errorcode.CodeSuccess.Msg(),
		Data:    data,
	})
}
