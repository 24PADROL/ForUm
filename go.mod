module engine

go 1.23.0

toolchain go1.23.4

require (
	github.com/go-sql-driver/mysql v1.9.0
	github.com/google/uuid v1.6.0
	golang.org/x/crypto v0.36.0
)

require filippo.io/edwards25519 v1.1.0 // indirect

// Run 'go mod tidy' to update dependencies
