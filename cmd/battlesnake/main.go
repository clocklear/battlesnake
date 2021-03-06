package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clocklear/battlesnake/lib/gamerecorder"
	v1 "github.com/clocklear/battlesnake/lib/v1"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/go-kit/kit/log"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Addr         string        `default:":8080" split_words:"true"`
	ReadTimeout  time.Duration `default:"5s" required:"true" split_words:"true"`
	WriteTimeout time.Duration `default:"5s" required:"true" split_words:"true"`
	Recorder     struct {
		Enabled           bool          `default:"true" split_words:"true"`
		OutputPath        string        `default:"" split_words:"true"`
		MaxAgeBeforePrune time.Duration `default:"2m" split_words:"true"`
		PruneInterval     time.Duration `default:"1m" split_words:"true"`
	} `split_words:"true"`
	SolveOption struct {
		Lookahead                bool `default:"true" split_words:"true"`
		ConsiderOpponentNextMove bool `default:"true" split_words:"true"`
		UseSingleBestOption      bool `default:"false" split_words:"true"`
		FoodReward               int  `default:"20" split_words:"true"`
		HazardPenalty            int  `default:"40" split_words:"true"`
	} `split_words:"true"`
	Logger struct {
		Enabled bool `default:"true" split_words:"true"`
	}
	NewRelic struct {
		LicenseKey string `split_words:"true"`
	} `split_words:"true"`
}

func main() {
	// Create our logger
	base := log.NewJSONLogger(os.Stdout)
	base = log.WithPrefix(base, "date", log.DefaultTimestampUTC)
	l := logger{base: base}

	// Parse our configuration and make sure we have everything that we need.
	var c config
	err := envconfig.Process("", &c)
	if err != nil {
		l.Fatal("could not process env", "err", err.Error())
	}

	if !c.Logger.Enabled {
		l = logger{
			base: log.NewNopLogger(),
		}
	}

	// Create new relic agent
	var nr *newrelic.Application
	if c.NewRelic.LicenseKey != "" {
		nr, err = newrelic.NewApplication(
			newrelic.ConfigFromEnvironment(),
		)
		if err != nil {
			l.Fatal("failed to create new relic", "err", err.Error())
		}
	}

	// Create gamerecorder
	var gr gamerecorder.GameRecorder
	if c.Recorder.Enabled {
		gr = gamerecorder.NewFileArchive(
			c.Recorder.OutputPath,
			c.Recorder.PruneInterval,
			c.Recorder.MaxAgeBeforePrune)
	} else {
		gr = gamerecorder.NoopGameRecorder{}
	}

	// Create handler
	h := handler{
		l:   l,
		rec: gr,
		nr:  nr,
		so:  v1.SolveOptions(c.SolveOption),
	}

	// Create http server
	appServer := http.Server{
		Addr:         c.Addr,
		Handler:      router(&h),
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	}

	// Create a channel to listen for http shutdown errors
	errs := make(chan error, 1)
	go func() {
		l.Info("starting battlesnake server", "addr", c.Addr)
		errs <- appServer.ListenAndServe()
	}()

	// Listen for stopping signals, and attempt to shut down gracefully.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errs:
		l.Fatal("received error", "err", err.Error())
		os.Exit(1)
	case s := <-osSignals:
		l.Info("received signal", "signal", s)

		// Shutdown supporting constructs
		nr.Shutdown(time.Second * 5)
		if fa, isFileArchive := gr.(*gamerecorder.FileArchive); isFileArchive {
			_ = fa.Shutdown()
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		l.Info("stopping battlesnake server")
		if err := appServer.Shutdown(ctx); err != nil {
			l.Error("could not shutdown battlesnake server", "err", err.Error())
			if err := appServer.Close(); err != nil {
				l.Error("could not close battlesnake server", "err", err.Error())
			}
		}
		cancel()
		os.Exit(0)
	}
}
