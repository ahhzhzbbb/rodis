package server

const defaultPort = ":6379"

type Config struct {
	Port      string
	BatchSize int
	Expire    ExpireConfig
}

type ExpireConfig struct {
	CycleIntervalMs int
	SampleSize      int
	ExpireThreshold float64
	TimeBudgetMs    int
}

func DefaultConfig() Config {
	return Config{
		Port:      ":6379",
		BatchSize: 8,
		Expire: ExpireConfig{
			CycleIntervalMs: 1000,
			SampleSize:      20,
			ExpireThreshold: 0.25,
			TimeBudgetMs:    5,
		},
	}
}
