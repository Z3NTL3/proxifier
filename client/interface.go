package client

type (
	Context struct {
		IP   string
		Port int
	}

	ProxyCtx  = Context
	TargetCtx = Context
)
