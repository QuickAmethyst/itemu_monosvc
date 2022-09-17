package main

import (
	"context"
	"github.com/QuickAmethyst/monosvc/graph/generated"
	accountUC "github.com/QuickAmethyst/monosvc/module/account/usecase"
	accountingUC "github.com/QuickAmethyst/monosvc/module/accounting/usecase"
	inventoryUC "github.com/QuickAmethyst/monosvc/module/inventory/usecase"
	sdkAuth "github.com/QuickAmethyst/monosvc/stdlibgo/auth"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"
	"github.com/go-redis/redis/v9"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/QuickAmethyst/monosvc/graph"
	accountSql "github.com/QuickAmethyst/monosvc/module/account/repository/sql"
	accountingSql "github.com/QuickAmethyst/monosvc/module/accounting/repository/sql"
	inventorySql "github.com/QuickAmethyst/monosvc/module/inventory/repository/sql"
	"github.com/QuickAmethyst/monosvc/stdlibgo/config"
	sdkGrace "github.com/QuickAmethyst/monosvc/stdlibgo/grace"
	sdkGraphql "github.com/QuickAmethyst/monosvc/stdlibgo/graphql"
	"github.com/QuickAmethyst/monosvc/stdlibgo/http"
	"github.com/QuickAmethyst/monosvc/stdlibgo/httpserver"
	sdkLogger "github.com/QuickAmethyst/monosvc/stdlibgo/logger"
	"go.uber.org/zap"
	"log"
	"syscall"
)

var (
	err                 error
	logger              sdkLogger.Logger
	fileConf            config.File
	conf                Config
	rest                http.Http
	stdLog              *log.Logger
	graphES             graphql.ExecutableSchema
	redisClient         redis.UniversalClient
	accountSqlClient    sql.PostgresSQL
	inventorySqlClient  sql.PostgresSQL
	accountingSqlClient sql.PostgresSQL
	resolver            graph.Resolver
	auth                sdkAuth.Auth
)

func initConf() {
	fileConf, err = config.NewFile(config.Option{
		Path:            "./config.yml",
		Type:            "yaml",
		Logger:          logger,
		RestartOnChange: true,
	})

	if err != nil {
		logger.Fatal(err.Error())
	}

	fileConf.ReadAndWatch(&conf)
}

func initLogger() {
	if logger, err = sdkLogger.New(sdkLogger.Option{Development: true}); err != nil {
		log.Fatal("Failed to create logger", err)
	}

	stdLog, err = zap.NewStdLogAt(logger, log.LstdFlags)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func initDB() {
	if accountSqlClient, err = sql.NewPostgresSQL(conf.AccountDatabase); err != nil {
		logger.Fatal(err.Error())
	}

	if inventorySqlClient, err = sql.NewPostgresSQL(conf.InventoryDatabase); err != nil {
		logger.Fatal(err.Error())
	}

	if accountingSqlClient, err = sql.NewPostgresSQL(conf.AccountingDatabase); err != nil {
		logger.Fatal(err.Error())
	}
}

func initRedis() {
	redisClient = redis.NewUniversalClient(&conf.Redis)
}

func initResolver() {
	accountSQLRepo := accountSql.New(&accountSql.Options{
		MasterDB: accountSqlClient.Master(),
		SlaveDB:  accountSqlClient.Slave(),
		Logger:   logger,
	})

	inventorySQLRepo := inventorySql.New(&inventorySql.Options{
		MasterDB: inventorySqlClient.Master(),
		SlaveDB:  inventorySqlClient.Slave(),
		Logger:   logger,
	})

	accountingSQLRepo := accountingSql.New(&accountingSql.Options{
		MasterDB: accountingSqlClient.Master(),
		SlaveDB:  accountingSqlClient.Slave(),
		Logger:   logger,
	})

	auth, err = sdkAuth.New(&sdkAuth.Options{
		Redis:                redisClient,
		PublicKeyPath:        "etc/rsa/public.pem",
		PrivateKeyPath:       "etc/rsa/private.pem",
		AccessTokenDuration:  1 * time.Hour,
		RefreshTokenDuration: 30 * 24 * time.Hour,
	})

	if err != nil {
		logger.Fatal(err.Error())
	}

	resolver = graph.Resolver{
		Logger:            logger,
		InventoryUsecase:  inventoryUC.New(&inventoryUC.Options{InventorySQL: inventorySQLRepo}),
		AccountUsecase:    accountUC.New(&accountUC.Options{AccountSQL: accountSQLRepo, Auth: auth}),
		AccountingUsecase: accountingUC.New(&accountingUC.Options{AccountingSQL: accountingSQLRepo}),
	}
}

func initGraph() {
	c := generated.Config{Resolvers: &resolver}
	c.Directives.Authenticated = sdkGraphql.AuthenticatedDirective(auth)

	graphES = generated.NewExecutableSchema(c)
	rest = http.New(http.Options{Cors: &conf.HttpCors})
	graphqlH, playgroundH := sdkGraphql.New(graphES, sdkGraphql.Options{Development: conf.Development})

	rest.Handle(http.MethodGet, "/graphql", playgroundH)
	rest.Handle(http.MethodPost, "/graphql/query", graphqlH)
}

func init() {
	initLogger()
	initConf()
	initRedis()
	initDB()
	initResolver()
	initGraph()
}

func main() {
	server := httpserver.New(conf.HttpServer, rest.Handler(), stdLog.Writer())
	grace, err := sdkGrace.New(logger, conf.Grace)
	if err != nil {
		logger.Fatal(err.Error())
	}

	defer func() {
		grace.Stop()
	}()

	go func() {
		grace.ListenForUpgrade(syscall.SIGHUP)
	}()

	grace.Serve(context.Background(), server)
}
