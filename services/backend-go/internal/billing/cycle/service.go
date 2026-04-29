package cycle

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service orquestra ciclos de faturamento.
type Service struct {
	pool *pgxpool.Pool
}

// NewService cria um novo Service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// -------------------------------------------------------------------
// CYCLE
// -------------------------------------------------------------------

// Cycle representa um ciclo de faturamento.
type Cycle struct {
	ID            uuid.UUID  `json:"id"`
	Year          int16      `json:"year"`
	Month         int16      `json:"month"`
	ReferenceDate time.Time  `json:"reference_date"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	ClosedAt      *time.Time `json:"closed_at,omitempty"`
	ClosedBy      *uuid.UUID `json:"closed_by,omitempty"`
	// Aggregated counts (populated by List/Get)
	TotalUCs        int `json:"total_ucs,omitempty"`
	SyncedCount     int `json:"synced_count,omitempty"`
	CalculatedCount int `json:"calculated_count,omitempty"`
	ApprovedCount   int `json:"approved_count,omitempty"`
}

// CreateRequest cria um novo ciclo.
type CreateRequest struct {
	Year             int16       `json:"year"`
	Month            int16       `json:"month"`
	IncludeAllActive bool        `json:"include_all_active"`
	UCIDs            []uuid.UUID `json:"uc_ids,omitempty"`
	CreatedBy        *uuid.UUID  `json:"created_by,omitempty"`
}

// Create cria um billing_cycle e associa consumer units.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*Cycle, error) {
	if req.Year < 2020 || req.Year > 2100 || req.Month < 1 || req.Month > 12 {
		return nil, fmt.Errorf("ano/mês inválido")
	}

	refDate := time.Date(int(req.Year), time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var cycle Cycle
	err = tx.QueryRow(ctx, `
		INSERT INTO public.billing_cycle (year, month, reference_date, status, created_by)
		VALUES ($1, $2, $3, 'open', $4)
		RETURNING id, year, month, reference_date, status, created_at, created_by
	`, req.Year, req.Month, refDate, req.CreatedBy).Scan(
		&cycle.ID, &cycle.Year, &cycle.Month, &cycle.ReferenceDate,
		&cycle.Status, &cycle.CreatedAt, &cycle.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("insert cycle: %w", err)
	}

	// Associar consumer units
	var ucIDs []uuid.UUID
	if req.IncludeAllActive {
		rows, err := tx.Query(ctx, `
			SELECT id FROM public.consumer_unit WHERE ativa = true
		`)
		if err != nil {
			return nil, fmt.Errorf("list active UCs: %w", err)
		}
		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return nil, err
			}
			ucIDs = append(ucIDs, id)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return nil, err
		}
	} else {
		ucIDs = req.UCIDs
	}

	for _, ucID := range ucIDs {
		_, err := tx.Exec(ctx, `
			INSERT INTO public.cycle_consumer_unit (billing_cycle_id, consumer_unit_id, status)
			VALUES ($1, $2, 'pending')
			ON CONFLICT DO NOTHING
		`, cycle.ID, ucID)
		if err != nil {
			return nil, fmt.Errorf("insert cycle_consumer_unit: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return &cycle, nil
}

// Get retorna um ciclo por ID com contagens agregadas.
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Cycle, error) {
	var c Cycle
	err := s.pool.QueryRow(ctx, `
		SELECT c.id, c.year, c.month, c.reference_date, c.status,
		       c.created_at, c.created_by, c.closed_at, c.closed_by,
		       COUNT(ccu.consumer_unit_id) FILTER (WHERE ccu.consumer_unit_id IS NOT NULL),
		       COUNT(ccu.consumer_unit_id) FILTER (WHERE ccu.status = 'synced'),
		       COUNT(ccu.consumer_unit_id) FILTER (WHERE ccu.status = 'calculated'),
		       COUNT(ccu.consumer_unit_id) FILTER (WHERE ccu.status = 'approved')
		FROM public.billing_cycle c
		LEFT JOIN public.cycle_consumer_unit ccu ON ccu.billing_cycle_id = c.id
		WHERE c.id = $1
		GROUP BY c.id
	`, id).Scan(
		&c.ID, &c.Year, &c.Month, &c.ReferenceDate, &c.Status,
		&c.CreatedAt, &c.CreatedBy, &c.ClosedAt, &c.ClosedBy,
		&c.TotalUCs, &c.SyncedCount, &c.CalculatedCount, &c.ApprovedCount,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("not_found")
		}
		return nil, err
	}
	return &c, nil
}

// ListCyclesRequest filtros para listagem.
type ListCyclesRequest struct {
	Year   int16
	Status string
	Limit  int
	Offset int
}

// List retorna ciclos com filtros.
func (s *Service) List(ctx context.Context, req ListCyclesRequest) ([]Cycle, error) {
	if req.Limit <= 0 {
		req.Limit = 50
	}
	query := `
		SELECT c.id, c.year, c.month, c.reference_date, c.status,
		       c.created_at, c.created_by, c.closed_at, c.closed_by,
		       COUNT(ccu.consumer_unit_id) FILTER (WHERE ccu.consumer_unit_id IS NOT NULL),
		       COUNT(ccu.consumer_unit_id) FILTER (WHERE ccu.status = 'synced'),
		       COUNT(ccu.consumer_unit_id) FILTER (WHERE ccu.status = 'calculated'),
		       COUNT(ccu.consumer_unit_id) FILTER (WHERE ccu.status = 'approved')
		FROM public.billing_cycle c
		LEFT JOIN public.cycle_consumer_unit ccu ON ccu.billing_cycle_id = c.id
		WHERE 1=1`
	args := []interface{}{}
	argNum := 1

	if req.Year > 0 {
		query += fmt.Sprintf(" AND c.year = $%d", argNum)
		args = append(args, req.Year)
		argNum++
	}
	if req.Status != "" {
		query += fmt.Sprintf(" AND c.status = $%d", argNum)
		args = append(args, req.Status)
		argNum++
	}

	query += fmt.Sprintf(`
		GROUP BY c.id
		ORDER BY c.year DESC, c.month DESC
		LIMIT $%d OFFSET $%d`, argNum, argNum+1)
	args = append(args, req.Limit, req.Offset)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cycles []Cycle
	for rows.Next() {
		var c Cycle
		err := rows.Scan(
			&c.ID, &c.Year, &c.Month, &c.ReferenceDate, &c.Status,
			&c.CreatedAt, &c.CreatedBy, &c.ClosedAt, &c.ClosedBy,
			&c.TotalUCs, &c.SyncedCount, &c.CalculatedCount, &c.ApprovedCount,
		)
		if err != nil {
			return nil, err
		}
		cycles = append(cycles, c)
	}
	return cycles, rows.Err()
}

// CloseRequest fecha um ciclo.
type CloseRequest struct {
	ClosedBy *uuid.UUID `json:"closed_by,omitempty"`
}

// Close fecha um ciclo. Requer que todos os cálculos estejam approved.
func (s *Service) Close(ctx context.Context, id uuid.UUID, req CloseRequest) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Verificar se há cálculos não-approved
	var nonApproved int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*) FROM public.billing_calculation
		WHERE billing_cycle_id = $1 AND status NOT IN ('approved', 'superseded')
	`, id).Scan(&nonApproved)
	if err != nil {
		return fmt.Errorf("check calculations: %w", err)
	}
	if nonApproved > 0 {
		return fmt.Errorf("ciclo possui %d cálculo(s) não aprovado(s)", nonApproved)
	}

	_, err = tx.Exec(ctx, `
		UPDATE public.billing_cycle
		SET status = 'closed', closed_at = NOW(), closed_by = $1
		WHERE id = $2 AND status != 'closed'
	`, req.ClosedBy, id)
	if err != nil {
		return fmt.Errorf("update cycle: %w", err)
	}

	return tx.Commit(ctx)
}
