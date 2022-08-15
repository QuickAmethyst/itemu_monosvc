package main

import (
	"context"
	"fmt"
	inventoryUC "github.com/QuickAmethyst/monosvc/module/inventory/usecase"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"

	"github.com/99designs/gqlgen/graphql"
	"github.com/QuickAmethyst/monosvc/graph"
	"github.com/QuickAmethyst/monosvc/graph/generated"
	inventorySql "github.com/QuickAmethyst/monosvc/module/inventory/repository/sql"
	sdkGraphql "github.com/QuickAmethyst/monosvc/stdlibgo/graphql"
	"go.uber.org/zap"
	"os"

	"github.com/QuickAmethyst/monosvc/stdlibgo/config"
	sdkGrace "github.com/QuickAmethyst/monosvc/stdlibgo/grace"
	"github.com/QuickAmethyst/monosvc/stdlibgo/http"
	"github.com/QuickAmethyst/monosvc/stdlibgo/httpserver"
	sdkLogger "github.com/QuickAmethyst/monosvc/stdlibgo/logger"
	"log"
	"syscall"
)

const defaultPort = "8080"

var (
	err                error
	port               string
	logger             sdkLogger.Logger
	fileConf           config.File
	conf               Config
	rest               http.Http
	stdLog             *log.Logger
	graphES            graphql.ExecutableSchema
	inventorySqlClient sql.PostgresSQL
	resolver           graph.Resolver
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
	if inventorySqlClient, err = sql.NewPostgresSQL(conf.InventoryDatabase); err != nil {
		logger.Fatal(err.Error())
	}

	if err = inventorySqlClient.Master().PingContext(context.Background()); err != nil {
		logger.Fatal(err.Error())
	}
}

func initResolver() {
	inventorySqlInstance := inventorySql.New(&inventorySql.Options{
		MasterDB: inventorySqlClient.Master(),
		SlaveDB:  inventorySqlClient.Slave(),
		Logger:   logger,
	})

	resolver = graph.Resolver{
		Logger: logger,
		InventoryUsecase: inventoryUC.New(&inventoryUC.Options{InventorySQL: inventorySqlInstance}),
	}
}

func initGraph() {
	graphES = generated.NewExecutableSchema(generated.Config{Resolvers: &resolver})
	fmt.Printf("%+v", conf.HttpCors)
	rest = http.New(http.Options{Cors: &conf.HttpCors})

	graphqlH, playgroundH := sdkGraphql.New(graphES, sdkGraphql.Options{Development: conf.Development})

	rest.Handle(http.MethodGet, "/graphql", playgroundH)
	rest.Handle(http.MethodPost, "/graphql/query", graphqlH)
}

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	initLogger()
	initConf()
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
