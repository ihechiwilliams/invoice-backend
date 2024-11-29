package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutdownDuration = 30 * time.Second

func HandleSignals(ctx context.Context, cancelCtx context.CancelFunc, callback func()) {
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownDuration)

		go func() {
			<-shutdownCtx.Done()

			if shutdownCtx.Err() == context.DeadlineExceeded {
				panic("graceful shutdown timed out.. forcing exit.")
			}
		}()

		callback()

		cancel()
		cancelCtx()
	}()
}
