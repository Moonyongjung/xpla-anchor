package types

type XGoError struct {
	errCode uint64
	desc    string
}

// Return error code and message generating on the anchor.
var (
	ErrParseConfig   = new(101, "error parsing config file")
	ErrParseApp      = new(102, "error parsing app file")
	ErrGenXplaClient = new(103, "error generating XPLA client")
	ErrInit          = new(104, "error init")
	ErrExecute       = new(105, "error execute")
	ErrKey           = new(106, "error key")
	ErrContract      = new(107, "error contract")
	ErrGw            = new(108, "error gateway")
	ErrBlockMng      = new(109, "error block management")
	ErrQuery         = new(110, "error query")
	ErrAccount       = new(111, "error account")
)

func new(errCode uint64, desc string) XGoError {
	var xErr XGoError
	xErr.errCode = errCode
	xErr.desc = desc

	return xErr
}

func (x XGoError) ErrCode() uint64 {
	return x.errCode
}

func (x XGoError) Desc() string {
	return x.desc
}
