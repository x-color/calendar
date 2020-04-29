package ctx

type CtxKey string

const (
	ReqIDKey  CtxKey = "reqID"
	UserIDKey CtxKey = "userID"
	TxKey     CtxKey = "tx"
)
