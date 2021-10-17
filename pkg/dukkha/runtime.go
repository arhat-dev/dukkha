package dukkha

// RuntimeOptions for task execution
type RuntimeOptions struct {
	FailFast            bool
	ColorOutput         bool
	TranslateANSIStream bool
	RetainANSIStyle     bool
	Workers             int
}

func (c *dukkhaContext) SetRuntimeOptions(opts RuntimeOptions) {
	c.runtimeOpts = opts
}
