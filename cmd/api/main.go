package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"syscall"
	"time"

	"github.com/jesseokeya/go-rest-api-template/config"
	"github.com/jesseokeya/go-rest-api-template/cronjob"
	"github.com/jesseokeya/go-rest-api-template/data"
	"github.com/jesseokeya/go-rest-api-template/lib/connect"
	"github.com/jesseokeya/go-rest-api-template/lib/session"
	"github.com/jesseokeya/go-rest-api-template/server"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/zenazn/goji/graceful"
)

var (
	flags    = flag.NewFlagSet("template", flag.ExitOnError)
	confFile = flags.String("config", "", "path to config file")
)

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal().Err(err).Msgf("invalid flags")
	}

	conf, err := config.NewFromFile(*confFile, os.Getenv("CONFIG"))
	if err != nil {
		log.Fatal().Err(err).Msgf("invalid config file")
	}

	// initialize seed
	rand.Seed(time.Now().Unix())

	// [db]
	db, err := data.NewDB(conf.DB)
	if err != nil {
		log.Fatal().Msgf("database: connection failed: %v", err)
	}

	// [connect]
	connect.Configure(conf.Connect)

	// [jwt]
	tokAuth := session.Setup(conf.JWT)

	// [cronjob]
	var c *cron.Cron
	if conf.Environment != "development" {
		c = cron.New()
		c.AddFunc("@every 15m", cronjob.Run)
		c.Start()

		for _, v := range c.Entries() {
			// for debugging purposes, print out the func name and next time
			// it's scheduled to run
			log.Info().Msgf("%v: next run @ %v", runtime.FuncForPC(reflect.ValueOf(v.Job).Pointer()).Name(), v.Next)
		}
	}

	graceful.AddSignal(syscall.SIGINT, syscall.SIGTERM)
	graceful.Timeout(10 * time.Second) // Wait timeout for handlers to finish.
	graceful.PreHook(func() {
		// prehook
		if c != nil {
			ctx := c.Stop()
			log.Info().Msg("waiting for crons")
			select {
			case <-ctx.Done():
				log.Info().Msg("done.")
			}
		}
	})
	graceful.PostHook(func() {
		// finishing up
	})

	// new web hander
	h, err := server.New(tokAuth, db, server.Debug(conf.Environment != "production"))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start web server")
	}

	// bind.
	bind := conf.Bind
	if bind == "" && conf.Port != "" {
		bind = fmt.Sprintf(":%s", conf.Port)
	}

	log.Info().Msgf("[%s] API starting on %s", conf.Environment, bind)
	if err := graceful.ListenAndServe(bind, h.Routes()); err != nil {
		log.Fatal().Err(err).Msg("cannot bind to host")
	}
	graceful.Wait()
}
