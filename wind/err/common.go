package merr

var (
	ParamsError           = NewError(1001, "参数错误", Abnormal)
	FindOneResultNilError = NewError(1201, "findOne result is nil", Abnormal)
	//NoDocumentError       = NewError(1202, "no document err", Abnormal)
	TransObjectIdError       = NewError(1203, "TransObjectIdError", Abnormal)
	DomainOptionError       = NewError(1203, "查询或修改条件错误", Abnormal)
)
