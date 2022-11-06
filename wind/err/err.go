package merr

import (
	"encoding/json"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

// ErrorCode is type for error code.
type ErrorCode uint32

type ErrorLabel string

const (
	UnknownError ErrorCode = 100
	MaxSystemErr           = 1000
	
	Normal   ErrorLabel = "normal"   //常规错误
	Abnormal ErrorLabel = "abnormal" //异常错误
)

var (
	UnknownErr                 = NewError(100, "unknown error", Abnormal)
	SysErrTimeoutErr           = NewError(101, "system timeout", Abnormal)
	MongoContextTimeoutErr     = NewError(110, "mongo context timeout", Abnormal)
	MongoFindOneResultNilErr   = NewError(111, "mongo find one result nil", Abnormal)
	HttpRequestErr             = NewError(120, "http request err", Abnormal)
	RedisLockTimeoutErr        = NewError(501, "redis lock timeout", Abnormal)
	ProducerNotInitErr         = NewError(511, "producer not init", Abnormal)
	ProducerStopErr            = NewError(512, "producer stopped", Abnormal)
	HttpStatusErr              = NewError(513, "http status err", Abnormal)
	SysErrSentinelBreaker      = NewError(601, "服务熔断", Abnormal)
	ResourceErrSentinelBreaker = NewError(602, "外部资源熔断", Abnormal)
)

type MetaError struct {
	Code  ErrorCode
	Msg   string
	Label ErrorLabel
	Data  string //可以携带更多信息 如果需要的话
}

func NewError(code ErrorCode, msg string, label ErrorLabel) MetaError {
	return MetaError{Code: code, Msg: msg, Label: label}
}

func (n MetaError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s, label: %s data: %s", n.Code, n.Msg, n.Label, n.Data)
}

func ParseMetaErrToStatusErr(t MetaError) error {
	s1 := status.New(codes.Code(t.Code), t.Msg)
	
	errInfo := errdetails.ErrorInfo{
		Metadata: map[string]string{
			"code":  fmt.Sprintf("%d", t.Code),
			"msg":   t.Msg,
			"label": string(t.Label),
			"data":  t.Data,
		},
	}
	s2, err := s1.WithDetails(&errInfo)
	if err != nil {
		return err
	}
	return s2.Err()
}

func ParseErrToMetaErr(err error) MetaError {
	if e, ok := err.(MetaError); ok {
		return e
	}
	status, ok := status.FromError(err)
	if ok {
		e := MetaError{Code: ErrorCode(status.Code()), Msg: status.Message(), Label: Abnormal}
		for _, v := range status.Details() {
			switch d := v.(type) {
			case *errdetails.ErrorInfo:
				errCode := d.Metadata["code"]
				errMsg := d.Metadata["msg"]
				errLabel := d.Metadata["label"]
				if errCode != "" {
					_code, _ := strconv.ParseInt(errCode, 10, 64)
					e.Code = ErrorCode(_code)
				}
				if errMsg != "" {
					e.Msg = errMsg
				}
				if d.Metadata["data"] != "" {
					e.Data = d.Metadata["data"]
				}
				if errLabel != "" {
					e.Label = ErrorLabel(errLabel)
				}
			default:
			}
		}
		return e
	}
	return MetaError{Code: UnknownError, Msg: err.Error(), Label: Abnormal}
}

func (n *MetaError) SetCode(code int) error {
	return MetaError{
		Code:  ErrorCode(code),
		Msg:   n.Msg,
		Label: n.Label,
		Data:  n.Data,
	}
}

func (n *MetaError) SetMsg(msg string) error {
	return MetaError{
		Code:  n.Code,
		Msg:   msg,
		Label: n.Label,
		Data:  n.Data,
	}
}

func (n *MetaError) WithMsg(msg string) error {
	return MetaError{
		Code:  n.Code,
		Msg:   fmt.Sprintf("%s %s", n.Msg, msg),
		Label: n.Label,
		Data:  n.Data,
	}
}

func (n *MetaError) SetData(data map[string]interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return MetaError{
		Code:  n.Code,
		Msg:   n.Msg,
		Label: n.Label,
		Data:  string(bytes),
	}
}
