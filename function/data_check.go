package function

import (
	"Source_Predict/entity"
	"Source_Predict/errorcode"
)

func DataCheck(req entity.DataForAnalysisReq) (errcode errorcode.ResponseCode) {

	if len(req.DataAll) == 0 || len(req.Time) == 0 {
		return errorcode.CodeInvalidData
	}
	if len(req.DataAll) != len(req.Time) {
		return errorcode.CodeInvalidMatch
	}
	return errorcode.CodeSuccess

}
