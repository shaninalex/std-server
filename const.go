package main

type contextKey string

const ContextUserIDKey contextKey = "user_id"
const ContextDB contextKey = "db"

const DatabaseConnectionString string = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
