package outbox

import (
	"context"
	"log"
	"time"

	"example.com/m/internal/usecase/event"
)

type relayWorker struct {
	usecase  event.EventRelayUseCaseInterface
	interval time.Duration
}

func NewRelayWorker(
	uc event.EventRelayUseCaseInterface,
	interval time.Duration,
) event.RelayWorker {
	return &relayWorker{
		usecase:  uc,
		interval: interval,
	}
}

func (w *relayWorker) Start(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("Starting EventRelay Worker (interval: %v)", w.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Event Relay Worker...")
			return ctx.Err()
		case <-ticker.C:
			if err := w.usecase.ProcessOutbox(ctx); err != nil {
				log.Printf("Event Relay process failed: %v", err)
			}
		}
	}
}
