package xerr

import merr "mond/wind/err"

var (
	LoginErr = merr.NewError(19001001, "登录失败，token验证不合法", merr.Abnormal)
)
