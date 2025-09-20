package pkg

type contextKey string

const (
	ContextSession = contextKey("session")
	ContextUser    = contextKey("user")

	ContextUserIDKey = contextKey("userId")
	ContextDB        = contextKey("db")
	ContextAppName   = contextKey("appName")
	EnvProduction    = contextKey("production")
	EnvStaging       = contextKey("staging")
	EnvDevelopment   = contextKey("development")
	EnvTesting       = contextKey("testing")
)
