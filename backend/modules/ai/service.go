package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeProcessEdital = "ai:process_edital"
)

type ProcessEditalPayload struct {
	TenantID uuid.UUID
	EditalID uuid.UUID
	FileURL  string
}

type Service struct {
	client *asynq.Client
}

func NewService(redisAddr string) *Service {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &Service{
		client: client,
	}
}

func (s *Service) EnqueueEditalProcessing(ctx context.Context, tenantId, editalId uuid.UUID, fileUrl string) error {
	payload, err := json.Marshal(ProcessEditalPayload{
		TenantID: tenantId,
		EditalID: editalId,
		FileURL:  fileUrl,
	})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeProcessEdital, payload)

	info, err := s.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %v", err)
	}
	_ = info
	return nil
}
