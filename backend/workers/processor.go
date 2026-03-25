package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"

	"lastsaas/modules/ai"
)

type Processor struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewProcessor(redisAddr string) *Processor {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	// Register handlers here
	// This acts as a bridge: when tasks are processed, it would normally call
	// the external python service out from here.
	mux.HandleFunc(ai.TypeProcessEdital, HandleProcessEdital)

	return &Processor{
		server: srv,
		mux:    mux,
	}
}

func (p *Processor) Start() error {
	return p.server.Start(p.mux)
}

func (p *Processor) Stop() {
	p.server.Stop()
}

func HandleProcessEdital(ctx context.Context, t *asynq.Task) error {
	var p ai.ProcessEditalPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	slog.Info("Processing task: Preparing to send edital to Python AI Service",
		"tenant_id", p.TenantID,
		"edital_id", p.EditalID,
		"file_url", p.FileURL,
	)

	// Here we would make an HTTP call to the python microservice
	// e.g. http.Post("http://ai-service/process", "application/json", body)

	return nil
}
