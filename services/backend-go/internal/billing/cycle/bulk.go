package cycle

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// BulkActionRequest define uma ação em massa.
type BulkActionRequest struct {
	Action    string     `json:"action"`              // sync, recalculate, generate_pdf, approve
	UCCodes   []string   `json:"uc_codes,omitempty"`  // vazio = todas
	ForceAll  bool       `json:"force_all"`           // ignora checks
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
}

// BulkActionResult retorna o resultado.
type BulkActionResult struct {
	JobsCreated     int      `json:"jobs_created"`
	JobsSkipped     int      `json:"jobs_skipped"`
	SkippedReasons  []string `json:"skipped_reasons,omitempty"`
}

// Bulk executa uma ação em massa nas UCs de um ciclo.
func (s *Service) Bulk(ctx context.Context, cycleID uuid.UUID, req BulkActionRequest) (*BulkActionResult, error) {
	switch req.Action {
	case "sync":
		return s.bulkSync(ctx, cycleID, req)
	case "recalculate":
		return s.bulkRecalculate(ctx, cycleID, req)
	case "generate_pdf":
		return s.bulkGeneratePDF(ctx, cycleID, req)
	case "approve":
		return s.bulkApprove(ctx, cycleID, req)
	default:
		return nil, fmt.Errorf("ação inválida: %s", req.Action)
	}
}

func (s *Service) bulkSync(ctx context.Context, cycleID uuid.UUID, req BulkActionRequest) (*BulkActionResult, error) {
	ucs, err := s.listUCsForBulk(ctx, cycleID, req.UCCodes)
	if err != nil {
		return nil, err
	}

	result := &BulkActionResult{}
	for _, uc := range ucs {
		if !req.ForceAll && uc.Status != "pending" {
			result.JobsSkipped++
			result.SkippedReasons = append(result.SkippedReasons,
				fmt.Sprintf("UC %s already %s", uc.UCCode, uc.Status))
			continue
		}
		if err := s.enqueueJob(ctx, "sync_uc", map[string]any{
			"cycle_id": cycleID.String(),
			"uc_id":    uc.ConsumerUnitID.String(),
			"uc_code":  uc.UCCode,
		}); err != nil {
			return nil, err
		}
		result.JobsCreated++
	}
	return result, nil
}

func (s *Service) bulkRecalculate(ctx context.Context, cycleID uuid.UUID, req BulkActionRequest) (*BulkActionResult, error) {
	ucs, err := s.listUCsForBulk(ctx, cycleID, req.UCCodes)
	if err != nil {
		return nil, err
	}

	result := &BulkActionResult{}
	for _, uc := range ucs {
		if !req.ForceAll && uc.Status != "synced" {
			result.JobsSkipped++
			result.SkippedReasons = append(result.SkippedReasons,
				fmt.Sprintf("UC %s not synced (status=%s)", uc.UCCode, uc.Status))
			continue
		}
		if err := s.enqueueJob(ctx, "calculate", map[string]any{
			"cycle_id": cycleID.String(),
			"uc_id":    uc.ConsumerUnitID.String(),
		}); err != nil {
			return nil, err
		}
		result.JobsCreated++
	}
	return result, nil
}

func (s *Service) bulkGeneratePDF(ctx context.Context, cycleID uuid.UUID, req BulkActionRequest) (*BulkActionResult, error) {
	// Buscar cálculos aprovados
	rows, err := s.pool.Query(ctx, `
		SELECT bc.id, cu.uc_code
		FROM public.billing_calculation bc
		JOIN public.consumer_unit cu ON cu.id = bc.consumer_unit_id
		WHERE bc.billing_cycle_id = $1 AND bc.status = 'approved'
	`, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &BulkActionResult{}
	for rows.Next() {
		var calcID uuid.UUID
		var ucCode string
		if err := rows.Scan(&calcID, &ucCode); err != nil {
			return nil, err
		}
		// TODO: verificar se PDF já existe (skip se não force_all)
		if err := s.enqueueJob(ctx, "generate_pdf", map[string]any{
			"cycle_id":       cycleID.String(),
			"calculation_id": calcID.String(),
			"uc_code":        ucCode,
		}); err != nil {
			return nil, err
		}
		result.JobsCreated++
	}
	return result, rows.Err()
}

func (s *Service) bulkApprove(ctx context.Context, cycleID uuid.UUID, req BulkActionRequest) (*BulkActionResult, error) {
	ucs, err := s.listUCsForBulk(ctx, cycleID, req.UCCodes)
	if err != nil {
		return nil, err
	}

	result := &BulkActionResult{}
	for _, uc := range ucs {
		if !req.ForceAll && uc.Status != "calculated" {
			result.JobsSkipped++
			result.SkippedReasons = append(result.SkippedReasons,
				fmt.Sprintf("UC %s not calculated (status=%s)", uc.UCCode, uc.Status))
			continue
		}
		if err := s.enqueueJob(ctx, "approve", map[string]any{
			"cycle_id": cycleID.String(),
			"uc_id":    uc.ConsumerUnitID.String(),
		}); err != nil {
			return nil, err
		}
		result.JobsCreated++
	}
	return result, nil
}

// ucForBulk representa uma UC para processamento em massa.
type ucForBulk struct {
	ConsumerUnitID uuid.UUID
	UCCode         string
	Status         string
}

func (s *Service) listUCsForBulk(ctx context.Context, cycleID uuid.UUID, ucCodes []string) ([]ucForBulk, error) {
	query := `
		SELECT cu.id, cu.uc_code, ccu.status
		FROM public.cycle_consumer_unit ccu
		JOIN public.consumer_unit cu ON cu.id = ccu.consumer_unit_id
		WHERE ccu.billing_cycle_id = $1`
	args := []interface{}{cycleID}

	if len(ucCodes) > 0 {
		placeholders := make([]string, len(ucCodes))
		for i, code := range ucCodes {
			placeholders[i] = fmt.Sprintf("$%d", i+2)
			args = append(args, code)
		}
		query += fmt.Sprintf(" AND cu.uc_code IN (%s)", joinStrings(placeholders, ","))
	}
	query += " ORDER BY cu.uc_code"

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ucForBulk
	for rows.Next() {
		var u ucForBulk
		if err := rows.Scan(&u.ConsumerUnitID, &u.UCCode, &u.Status); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, rows.Err()
}

func (s *Service) enqueueJob(ctx context.Context, jobType string, payload map[string]any) error {
	payloadJSON, _ := json.Marshal(payload)
	idempotencyKey := fmt.Sprintf("%s:%s", jobType, payloadJSON)
	_, err := s.pool.Exec(ctx, `
		INSERT INTO public.sync_job (type, payload_json, status, idempotency_key, scheduled_for)
		VALUES ($1, $2, 'pending', $3, NOW())
		ON CONFLICT (type, idempotency_key) WHERE status IN ('pending','running','retrying','success')
		DO NOTHING
	`, jobType, payloadJSON, idempotencyKey)
	return err
}

func joinStrings(ss []string, sep string) string {
	if len(ss) == 0 {
		return ""
	}
	result := ss[0]
	for i := 1; i < len(ss); i++ {
		result += sep + ss[i]
	}
	return result
}
