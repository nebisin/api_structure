package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (s *server) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.port),
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sign := <-quit

		s.logger.WithField("signal", sign.String()).Info("shutting down the server")

		ctx,cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		s.logger.WithField("addr", srv.Addr).Info("completing background tasks")

		s.wg.Wait()
		shutdownError <- nil
	}()

	s.logger.WithField("addr", srv.Addr).Info("starting the server")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <- shutdownError
	if err != nil {
		return err
	}

	s.logger.WithField("addr", srv.Addr).Info("stopped the server")

	return nil
}