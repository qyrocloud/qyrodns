package qyrodns

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"github.com/qyrocloud/qyrodns/internal/app/qyrodns/admin"
	"github.com/qyrocloud/qyrodns/internal/app/qyrodns/apikey"
	"github.com/qyrocloud/qyrodns/internal/app/qyrodns/deletion"
	dnsLib "github.com/qyrocloud/qyrodns/internal/app/qyrodns/dns"
	"github.com/qyrocloud/qyrodns/internal/app/qyrodns/health"
	"github.com/qyrocloud/qyrodns/internal/app/qyrodns/namespace"
	"github.com/qyrocloud/qyrodns/internal/pkg/auth"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	config *ServerConfig
}

func NewServer(config *ServerConfig) *Server {
	return &Server{config: config}
}

type ServerConfig struct {
	DNSHost       string
	DNSPort       string
	AdminHost     string
	AdminPort     string
	MongoEndpoint string
	MongoDatabase string
	JwtSigningKey string
	JwtIssuer     string
	JwtAudience   string
}

func (s *Server) Start() {
	// Database setup

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(s.config.MongoEndpoint))

	if err != nil {
		log.Fatal("error while connecting to mongo database: ", err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("error while pinging mongo database: ", err)
	}

	mongoDatabase := client.Database(s.config.MongoDatabase)

	// Services setup

	apiKeysCollection := mongoDatabase.Collection("api_keys")
	authenticator := auth.NewAuthenticator(s.config.JwtSigningKey, s.config.JwtIssuer, s.config.JwtAudience, apiKeysCollection)
	adminService := admin.NewService(mongoDatabase.Collection("admins"), authenticator)
	apiKeyService := apikey.NewService(apiKeysCollection)
	namespaceService := namespace.NewService(mongoDatabase.Collection("namespaces"))
	apiKeyAccessService := namespace.NewApiKeyAccessService(mongoDatabase.Collection("api_key_accesses"), namespaceService, apiKeyService)
	recordService := dnsLib.NewRecordService(mongoDatabase.Collection("records"), namespaceService)

	// DNS server setup

	dnsHandler := dnsLib.NewHandler(recordService)
	dns.HandleFunc(".", dnsHandler.Handle)

	dnsServer := &dns.Server{
		Addr: fmt.Sprintf("%s:%s", s.config.DNSHost, s.config.DNSPort),
		Net:  "udp",
	}

	go func() {
		log.Printf("starting DNS server on %s:%s", s.config.DNSHost, s.config.DNSPort)
		err = dnsServer.ListenAndServe()

		if err != nil {
			log.Fatal("error while starting DNS server: ", err)
		}
	}()

	// Admin server setup

	router := gin.Default()
	health.NewCheckHandler(router).Register()
	admin.NewHandler(router, authenticator, adminService).Register()
	apikey.NewHandler(router, authenticator, apiKeyService).Register()
	namespace.NewHandler(router, authenticator, namespaceService).Register()
	deletion.NewNamespaceDeletionHandler(router, authenticator, namespaceService, apiKeyAccessService, recordService).Register()
	namespace.NewApiKeyAccessHandler(router, authenticator, apiKeyAccessService).Register()
	dnsLib.NewRecordAdminHandler(router, authenticator, recordService).Register()
	dnsLib.NewRecordHandler(router, authenticator, apiKeyAccessService, recordService).Register()

	log.Printf("starting admin server on %s:%s", s.config.AdminHost, s.config.AdminPort)

	err = router.Run(fmt.Sprintf("%s:%s", s.config.AdminHost, s.config.AdminPort))

	if err != nil {
		log.Fatal("error starting admin server: ", err)
	}
}
