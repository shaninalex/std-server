package main

type contextKey string

const ContextUserIDKey contextKey = "userId"
const ContextDB contextKey = "db"
const ContextAppName contextKey = "appName"

const DatabaseConnectionString string = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

const EnvProduction = "production"
const EnvStaging = "staging"
const EnvDevelopment = "development"
const EnvTesting = "testing"
