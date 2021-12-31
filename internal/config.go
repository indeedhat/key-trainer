package internal

// RunnerConfig that can be passe to NewRunner in order
// to setup the test config
type RunnerConfig struct {
	Contains    string
	ContainsAny string
	TotalWords  uint
	TimeLimit   uint
}
