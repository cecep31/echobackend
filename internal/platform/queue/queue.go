package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"echobackend/config"

	"github.com/hibiken/asynq"
)

// SkipRetry marks a task error as permanent.
var SkipRetry = asynq.SkipRetry

// HandlerFunc handles a queued task payload.
type HandlerFunc func(ctx context.Context, payload []byte) error

// TaskOptions controls how a task is enqueued.
type TaskOptions struct {
	Queue    string
	Timeout  time.Duration
	MaxRetry int
}

// Service owns the shared Asynq client and worker server.
type Service struct {
	client       *asynq.Client
	server       *asynq.Server
	mux          *asynq.ServeMux
	defaultQueue string
	maxRetry     int
	started      bool
}

// NewService creates a shared Asynq queue service. Empty Redis URL disables it.
func NewService(cfg config.QueueConfig) *Service {
	service := &Service{
		mux:          asynq.NewServeMux(),
		defaultQueue: cfg.DefaultQueue,
		maxRetry:     cfg.MaxRetry,
	}

	if cfg.RedisURL == "" {
		slog.Warn("queue: QUEUE_REDIS_URL/REDIS_URL is empty, background jobs disabled")
		return service
	}

	redisOpt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		slog.Warn("queue: invalid Redis URL, background jobs disabled", "error", err)
		return service
	}
	redisOpt = withRedisTimeouts(redisOpt, cfg.ConnectTimeout)

	service.client = asynq.NewClient(redisOpt)

	// Verify broker connectivity in the background so startup is not blocked.
	go func() {
		if err := service.client.Ping(); err != nil {
			slog.Warn("queue: failed to connect to Redis, background jobs disabled", "error", err)
			_ = service.client.Close()
			service.client = nil
			return
		}
		slog.Info("queue: Asynq enabled", "queue", cfg.DefaultQueue, "concurrency", cfg.Concurrency)
	}()

	service.server = asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: cfg.Concurrency,
		Queues: map[string]int{
			cfg.DefaultQueue: 1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			slog.Error("queue: task failed", "error", err, "task_type", task.Type())
		}),
	})

	slog.Info("queue: server configured", "queue", cfg.DefaultQueue, "concurrency", cfg.Concurrency)
	return service
}

// IsConfigured reports whether the queue broker and worker are available.
func (s *Service) IsConfigured() bool {
	return s != nil && s.client != nil && s.server != nil
}

// Handle registers a task handler.
func (s *Service) Handle(taskType string, handler HandlerFunc) {
	if s == nil || s.mux == nil || taskType == "" || handler == nil {
		return
	}

	s.mux.HandleFunc(taskType, func(ctx context.Context, task *asynq.Task) error {
		return handler(ctx, task.Payload())
	})
}

// Start begins processing registered task handlers.
func (s *Service) Start() {
	if !s.IsConfigured() || s.started {
		return
	}
	s.started = true

	go func() {
		if err := s.server.Run(s.mux); err != nil {
			slog.Error("queue: Asynq server stopped", "error", err)
		}
	}()
}

// Close stops the worker and closes the client.
func (s *Service) Close() error {
	if s == nil {
		return nil
	}

	if s.server != nil {
		s.server.Shutdown()
	}

	if s.client != nil {
		return s.client.Close()
	}

	return nil
}

// EnqueueJSON marshals payload as JSON and enqueues it as a background task.
func (s *Service) EnqueueJSON(taskType string, payload any, opts TaskOptions) error {
	if !s.IsConfigured() {
		return fmt.Errorf("queue service not configured")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	queueName := opts.Queue
	if queueName == "" {
		queueName = s.defaultQueue
	}

	maxRetry := opts.MaxRetry
	if maxRetry == 0 {
		maxRetry = s.maxRetry
	}

	taskOptions := []asynq.Option{
		asynq.Queue(queueName),
		asynq.MaxRetry(maxRetry),
	}
	if opts.Timeout > 0 {
		taskOptions = append(taskOptions, asynq.Timeout(opts.Timeout))
	}

	task := asynq.NewTask(taskType, body)
	_, err = s.client.Enqueue(task, taskOptions...)
	return err
}

func withRedisTimeouts(opt asynq.RedisConnOpt, timeout time.Duration) asynq.RedisConnOpt {
	if timeout <= 0 {
		return opt
	}

	switch v := opt.(type) {
	case asynq.RedisClientOpt:
		v.DialTimeout = timeout
		v.ReadTimeout = timeout
		v.WriteTimeout = timeout
		return v
	case asynq.RedisFailoverClientOpt:
		v.DialTimeout = timeout
		v.ReadTimeout = timeout
		v.WriteTimeout = timeout
		return v
	case asynq.RedisClusterClientOpt:
		v.DialTimeout = timeout
		v.ReadTimeout = timeout
		v.WriteTimeout = timeout
		return v
	default:
		return opt
	}
}
