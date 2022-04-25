package taskchain

type TaskChainFactory struct {
	Chains      []*TaskChain
	ChainMap    map[string]*TaskChain //{name}-{version}
	LatestChain map[string]*TaskChain //{name}
}

type TaskChain struct {
	Name    string
	Version int
	Stages  []string
	Failure []string
}
