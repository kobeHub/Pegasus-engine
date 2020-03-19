// +build jsoniter

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"

	api "github.com/kobeHub/Pegasus-engine/api/v1"
	"github.com/kobeHub/Pegasus-engine/config"
	"github.com/kobeHub/Pegasus-engine/pkg/common/k8s"
	"github.com/kobeHub/Pegasus-engine/pkg/common/logmw"
	route "github.com/kobeHub/Pegasus-engine/pkg/common/router"
)

func main() {
	root := context.Background()
	os.Exit(run(root))
}

func run(ctx context.Context) int {
	var environ string
	if environ = os.Getenv("pegasus_env"); environ == "" {
		environ = "local"
	}
	config.Init(environ)
	var (
		srvAddr = kingpin.Flag("listen-address",
			"Cluster listen address, empty string to disable HA mode").Default(
			viper.GetString("DefaultClusterAddr")).String()
	)
	kingpin.Version(viper.GetString("version"))
	kingpin.Parse()
	k8s.Init(environ)

	router := route.New()
	api.Register(router)
	routerH := logmw.Middleware(router)
	srv := http.Server{Addr: *srvAddr, Handler: routerH}
	srvc := make(chan bool)

	go func() {
		log.WithFields(log.Fields{
			"Listen":  *srvAddr,
			"Version": viper.GetString("Version"),
		}).Info("Welcome to pegasus-engine")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Server listen error")
			close(srvc)
		}
		defer func() {
			if err := srv.Close(); err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Error("Error on closing the server")
			}
		}()

	}()

	var (
		term = make(chan os.Signal, 1)
	)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	// Watch tasks
	go k8s.WatchPods(ctx, srvc)

	for {
		select {
		case <-term:
			log.Info("Received SIGTERM, exiting gracefully....")
			ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Error("Gracefully shutdown error")
			}
			// catching ctx.Done(). timeout of 1 seconds.
			select {
			case <-ctx.Done():
				log.Info("Server timeout 1 seconds...")
			}

			return 0
		case <-srvc:
			log.Info("Error exist")
			return 1
		}
	}
}
