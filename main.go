package main

import (
	"Source_Predict/entity"
	"Source_Predict/errorcode"
	"Source_Predict/function"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// 路由
	router.GET("/api/v1/resource/predict", AnalysisResourceInfo)
	// 运行服务器
	router.Run(":8080")
}

func AnalysisResourceInfo(ctx *gin.Context) {
	req := &entity.DataForAnalysisReq{} //从ctx中解析参数
	if err := ctx.ShouldBindJSON(req); err != nil {
		entity.ResponseError(ctx, errorcode.CodeInvalidParam)
		return
	}
	//数据校验
	code := function.DataCheck(*req)
	if code != errorcode.CodeSuccess {
		entity.ResponseError(ctx, code)
		return
	}
	//数据预测
	err, resp := function.BatchReadAndAnalysis(req)
	if err != nil {
		entity.ResponseError(ctx, errorcode.CodeAnalysisFailed)
		return
	}
	entity.ResponseSuccess(ctx, resp)
}
