package main

import (
	"github.com/qyrocloud/qyrodns/internal/app/qyrodns"
	"github.com/qyrocloud/qyrodns/internal/pkg/env"
)

func main() {
	qyrodns.NewServer(&qyrodns.ServerConfig{
		DNSHost:       env.GetOrDefault("DNS_HOST", "0.0.0.0"),
		DNSPort:       env.GetOrDefault("DNS_PORT", "5300"),
		AdminHost:     env.GetOrDefault("ADMIN_HOST", "0.0.0.0"),
		AdminPort:     env.GetOrDefault("ADMIN_PORT", "5301"),
		MongoEndpoint: env.GetOrDefault("MONGO_ENDPOINT", "mongodb://localhost:27017"),
		MongoDatabase: env.GetOrDefault("MONGO_DB", "qyrodns"),
		JwtSigningKey: env.GetOrDefault("JWT_SIGNING_KEY", "secret"),
		JwtIssuer:     env.GetOrDefault("JWT_ISSUER", "qyrodns"),
		JwtAudience:   env.GetOrDefault("JWT_AUDIENCE", "qyrodns"),
	}).Start()
}
