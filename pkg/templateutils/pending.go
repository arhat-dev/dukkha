package templateutils

// pendingNS to define placeholder funcs
type pendingNS struct{}

// Var used as template variables
func (pendingNS) Var() map[string]any { return nil }

// Include other template by name, execute with data
func (pendingNS) Include(name String, data any) (_ string, _ error) { return }
