package domain

type Environment string

var (
	EnvTest       Environment = "test"
	EnvStaging    Environment = "staging"
	EnvProduction Environment = "production"
)
