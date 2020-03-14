// +build jsoniter

package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"

	api "github.com/kobeHub/Pegasus-engine/api/v1"
	"github.com/kobeHub/Pegasus-engine/config"
	route "github.com/kobeHub/Pegasus-engine/pkg/common/router"
	"github.com/kobeHub/Pegasus-engine/pkg/common/logmw"
)

func main() {
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

	router := route.New()
	api.Register(router)
	routerH := logmw.Middleware(router)
	srv := http.Server{Addr: *srvAddr, Handler: routerH}
	srvc := make(chan struct{})

	go func() {
		log.WithFields(log.Fields{
			"Listen":  *srvAddr,
			"Version": viper.GetString("Version"),
		}).Info("Welcome to pegasus-engine")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
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
		hup = make(chan os.Signal, 1)
		//_hupReady = make(chan bool)
		term = make(chan os.Signal, 1)
	)
	signal.Notify(hup, syscall.SIGHUP)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	/*
		go func() {
			<-hupReady
			for {
				select {
				case <-hup:

				}
			}
		}()*/
	for {
		select {
		case <-term:
			log.Info("Received SIGTERM, exiting gracefully....")
			os.Exit(0)
		case <-srvc:
			log.Info("Error exist.")
			os.Exit(1)
		}
	}
}
