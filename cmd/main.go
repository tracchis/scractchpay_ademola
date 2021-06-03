package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/scratchpay_ademola/internal/httputil"
	"github.com/scratchpay_ademola/internal/logger"
	"github.com/scratchpay_ademola/internal/os/process"

	"github.com/kelseyhightower/envconfig"

	"github.com/scratchpay_ademola/pkg/clinic"
	"go.uber.org/zap"
)

// These variables contain build information about the Version and a Buildstamp (timestamp, githash, etc.)
var (
	// injected by build tools but we aren't provisioning any here yet for this clinic search service
	AppName    = "scratchpay-clinic-search-service"
	Version    = "0.0.0"
	Buildstamp = "dev"
)

func main() {
	var cfg Config

	// load configuration
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(fmt.Errorf("error loading configuration: %s", err.Error()))
	}

	ctx, done := process.Init(AppName, Version, Buildstamp)
	defer done()

	// init logger
	log, err := logger.InitLogger(cfg.Env)
	if err != nil {
		panic(fmt.Errorf("error initialising logger: %s", err))
	}

	process.AtExit(func() { log.Sync() })

	// init mux
	ready := httputil.NewReady(
		httputil.TextHandler(http.StatusServiceUnavailable, "application/json", `"NOT READY"`),
	)

	mux := httputil.NewBaseMux(
		ready.Handler(httputil.TextHandler(http.StatusOK, "application/json", `"READY"`)),
	)

	clinicDataDownloader := clinic.NewDataDownloader(log)

	// init routes
	routes := initRoutes(clinicDataDownloader)

	mux.Handle("/", routes)

	// init HTTP Server for API
	httpServer := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", cfg.Port),
	}

	ready.Ready()
	log.Info("started clinic search service")

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal("Failed to start clinic search service", zap.Error(err))
	}

	// wait until the service should be stopped
	<-ctx.Done()

	// allow 15 seconds to shutdown everything gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown API server
	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Error("failed to shutdown clinic search server", zap.Error(err))
		process.Exit(1)
	}

	log.Debug("clinic search service shutdown successful")
}
