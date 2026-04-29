package cycle

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/platform/sse"
)

// EventPayload is the JSON payload sent via Postgres NOTIFY and forwarded to
// SSE clients.
type EventPayload struct {
	Type      string `json:"type"`
	CycleID   string `json:"cycle_id"`
	UCID      string `json:"uc_id,omitempty"`
	JobType   string `json:"job_type,omitempty"`
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
	Timestamp string `json:"timestamp"`
}

// CycleEventBroker listens on a Postgres channel (billing_cycle_events) for
// job progress notifications and fans them out to SSE subscribers filtered by
// cycle_id.
type CycleEventBroker struct {
	broker  *sse.Broker
	conn    *pgx.Conn
	channel string
	logger  *slog.Logger
	closeCh chan struct{}
}

// NewCycleEventBroker creates a new CycleEventBroker, issues LISTEN on the
// Postgres channel, and starts a background goroutine that forwards incoming
// notifications to matched subscribers.
func NewCycleEventBroker(listenConn *pgx.Conn, logger *slog.Logger) *CycleEventBroker {
	b := &CycleEventBroker{
		broker:  sse.NewBroker(),
		conn:    listenConn,
		channel: "billing_cycle_events",
		logger:  logger,
		closeCh: make(chan struct{}),
	}
	go b.listen()
	return b
}

func (b *CycleEventBroker) listen() {
	ctx := context.Background()

	_, err := b.conn.Exec(ctx, "LISTEN "+b.channel)
	if err != nil {
		b.logger.Error("sse_listen_failed", "channel", b.channel, "error", err)
		return
	}
	b.logger.Info("sse_listening", "channel", b.channel)

	for {
		select {
		case <-b.closeCh:
			return
		default:
		}

		notification, err := b.conn.WaitForNotification(ctx)
		if err != nil {
			select {
			case <-b.closeCh:
				return
			default:
				b.logger.Error("sse_wait_notification", "error", err)
				time.Sleep(1 * time.Second)
				continue
			}
		}

		var payload EventPayload
		if err := json.Unmarshal([]byte(notification.Payload), &payload); err != nil {
			b.logger.Warn("sse_invalid_payload", "payload", notification.Payload)
			continue
		}

		if payload.CycleID == "" {
			continue
		}

		// Fan out to subscribers for this specific cycle_id
		b.broker.Publish(payload.CycleID, []byte(notification.Payload))
	}
}

// Subscribe returns a channel that receives raw JSON event payloads for the
// given cycle_id.
func (b *CycleEventBroker) Subscribe(cycleID string) chan []byte {
	return b.broker.Subscribe(cycleID)
}

// Unsubscribe removes a subscription channel for the given cycle_id.
func (b *CycleEventBroker) Unsubscribe(cycleID string, ch chan []byte) {
	b.broker.Unsubscribe(cycleID, ch)
}

// Close stops the listener goroutine and releases the dedicated Postgres
// connection.
func (b *CycleEventBroker) Close() error {
	select {
	case <-b.closeCh:
	default:
		close(b.closeCh)
	}
	b.broker.Close()
	return b.conn.Close(context.Background())
}

// NotifyJobEvent sends a job progress notification via Postgres NOTIFY on the
// "billing_cycle_events" channel. This is a best-effort operation; errors are
// logged but not returned.
func NotifyJobEvent(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger, job *SyncJob, status, message string) {
	cycleID, ucID := getJobEventIDs(job)
	if cycleID == "" {
		return // cannot route without cycle_id
	}
	payload, err := json.Marshal(EventPayload{
		Type:      "job_update",
		CycleID:   cycleID,
		UCID:      ucID,
		JobType:   job.Type,
		Status:    status,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		logger.Warn("notify_job_event_marshal", "error", err)
		return
	}
	_, err = pool.Exec(ctx, "SELECT pg_notify('billing_cycle_events', $1)", string(payload))
	if err != nil {
		logger.Warn("notify_job_event_failed", "error", err)
	}
}

// getJobEventIDs extracts cycle_id and uc_id from a SyncJob's payload.
func getJobEventIDs(job *SyncJob) (cycleID, ucID string) {
	if v, ok := job.Payload["cycle_id"]; ok {
		cycleID, _ = v.(string)
	}
	if v, ok := job.Payload["uc_id"]; ok {
		ucID, _ = v.(string)
	}
	return
}
