package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (a *applicationDependences) serve() error {
	apiServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.config.port),
		Handler:      a.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(a.logger.Handler(), slog.LevelError),
	}

	//create a channel to kepp track of any errors during the shutdown process
	shutdownError := make(chan error)

	//crete a goroutine that runs in the background listining to the shutdown signals
	go func() {

		//receive the shutdown signal
		quit := make(chan os.Signal, 1)

		//signal occured
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		//blocks until a signal is received
		s := <-quit

		//message about shutdown in process
		a.logger.Info("shutting down server", "signal", s.String())

		//create a context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		//writing to error channel if error is found
		err := apiServer.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		//waiting for background tasks to finish
		a.logger.Info("completing background tasks", "address", apiServer.Addr)
		a.wg.Wait()
		shutdownError <- nil //sucessfull shutdown
	}()

	a.logger.Info("Starting Server", "address", apiServer.Addr, "environment", a.config.environment, "limiter-enabled", a.config.limiter)

	// something went wrong during shutdown if we don't get ErrServerClosed()
	// this only happens when we issue the shutdown command from our goroutine
	// otherwise our server keeps running as normal as it should.
	err := apiServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	//check the error channel to see if there were shutdown errors
	err = <-shutdownError
	if err != nil {
		return err
	}

	//gracefull shutdown was sucessfull
	a.logger.Info("server stoped", "address", apiServer.Addr)

	return nil
}
