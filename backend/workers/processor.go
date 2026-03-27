package workers

import (
	"database/sql"

	"github.com/hibiken/asynq"
)

type Processor struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewProcessor(redisAddr string, db *sql.DB) *Processor {
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

	// Register handlers natively within the workers package
	mux.HandleFunc(TypeAIAnalysis, HandleAIAnalysisTask(db))

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
