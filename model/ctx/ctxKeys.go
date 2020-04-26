package ctx

type CtxKey string

const (
	ReqIDKey CtxKey = "reqID"
	TxKey    CtxKey = "tx"
)
