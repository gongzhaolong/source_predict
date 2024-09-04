package errorcode

type ResponseCode int64

const (
	CodeSuccess        ResponseCode = 0 + iota
	CodeInvalidParam                // 1
	CodeInvalidMatch                // 2
	CodeInvalidData                 // 3
	CodeAnalysisFailed              // 4
	CodeServeBusy                   // 5
)

// 这样做是为了不把真正的错误返回给用户，只返回错误提示信息，真正的错误在终端或者日志中自己看
var codeMsgMap = map[ResponseCode]string{
	CodeSuccess:        "success",
	CodeInvalidParam:   "请求参数错误",
	CodeInvalidMatch:   "传入时间和数据长度不一致",
	CodeInvalidData:    "传入数据不能为空",
	CodeAnalysisFailed: "数据预测失败",
}

func (c ResponseCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServeBusy] // 根据 c 没查到提示信息，就返回个服务繁忙。。。
	}
	return msg
}
