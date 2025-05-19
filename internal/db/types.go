package db

// Settings represents the configurable settings that can be persisted
type Settings struct {
	// Crawler settings
	Domains       []string
	MaxDepth      int
	ThreadCount   int
	MaxQueueSize  int
	Parallelism   int
	DelayMS       int
	RandomDelayMS int

	// Anti-bot and proxy settings
	UserAgents      []string
	ProxyURL        string
	RotateUserAgent bool
}
