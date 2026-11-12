package pricing

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DecoratorDTO struct {
	ID               string    `json:"id"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Kind             string    `json:"kind"` // discount, shipping, service, tax
	Value            float64   `json:"value"`
	IsPercentage     bool      `json:"is_percentage"`
	MinSubtotalCents int64     `json:"min_subtotal_cents"`
	Active           bool      `json:"active"`
	Priority         int       `json:"priority"`
	CreatedAt        time.Time `json:"created_at"`
}

type CreateDecoratorInput struct {
	Code             string  `json:"code"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Kind             string  `json:"kind"`
	Value            float64 `json:"value"`
	IsPercentage     bool    `json:"is_percentage"`
	MinSubtotalCents int64   `json:"min_subtotal_cents"`
	Priority         int     `json:"priority"`
}

type UpdateDecoratorInput struct {
	Name             *string  `json:"name"`
	Description      *string  `json:"description"`
	Value            *float64 `json:"value"`
	IsPercentage     *bool    `json:"is_percentage"`
	MinSubtotalCents *int64   `json:"min_subtotal_cents"`
	Active           *bool    `json:"active"`
	Priority         *int     `json:"priority"`
}

type Service struct {
	pool     *pgxpool.Pool
	currency string
}

func NewService(pool *pgxpool.Pool, currency string) *Service {
	return &Service{pool: pool, currency: currency}
}

// -------------------------------------------------------------------------
// CRUD: CREATE
// -------------------------------------------------------------------------

func (s *Service) Create(ctx context.Context, in CreateDecoratorInput) (*DecoratorDTO, error) {
	if in.Code == "" || in.Name == "" || in.Kind == "" {
		return nil, fmt.Errorf("código, nombre y tipo son obligatorios")
	}

	if in.Kind != "discount" && in.Kind != "shipping" && in.Kind != "service" && in.Kind != "tax" {
		return nil, fmt.Errorf("tipo de decorador no válido: %s", in.Kind)
	}

	var d DecoratorDTO
	err := s.pool.QueryRow(ctx, `
		INSERT INTO pricing_decorators (code, name, description, kind, value, is_percentage, min_subtotal_cents, priority)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, code, name, description, kind, value, is_percentage, min_subtotal_cents, active, priority, created_at
	`, in.Code, in.Name, in.Description, in.Kind, in.Value, in.IsPercentage, in.MinSubtotalCents, in.Priority).Scan(
		&d.ID, &d.Code, &d.Name, &d.Description, &d.Kind, &d.Value, &d.IsPercentage, &d.MinSubtotalCents, &d.Active, &d.Priority, &d.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error al registrar decorador: %w", err)
	}

	return &d, nil
}

// -------------------------------------------------------------------------
// CRUD: READ
// -------------------------------------------------------------------------

func (s *Service) List(ctx context.Context) ([]DecoratorDTO, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, code, name, description, kind, value, is_percentage, min_subtotal_cents, active, priority, created_at
		FROM pricing_decorators
		ORDER BY priority ASC, created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []DecoratorDTO
	for rows.Next() {
		var d DecoratorDTO
		if err := rows.Scan(
			&d.ID, &d.Code, &d.Name, &d.Description, &d.Kind, &d.Value, &d.IsPercentage, &d.MinSubtotalCents, &d.Active, &d.Priority, &d.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}

	return list, nil
}

func (s *Service) Get(ctx context.Context, idOrCode string) (*DecoratorDTO, error) {
	var d DecoratorDTO
	err := s.pool.QueryRow(ctx, `
		SELECT id, code, name, description, kind, value, is_percentage, min_subtotal_cents, active, priority, created_at
		FROM pricing_decorators
		WHERE id::text = $1 OR code = $1
	`, idOrCode).Scan(
		&d.ID, &d.Code, &d.Name, &d.Description, &d.Kind, &d.Value, &d.IsPercentage, &d.MinSubtotalCents, &d.Active, &d.Priority, &d.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("decorador no encontrado: %w", err)
	}
	return &d, nil
}

// -------------------------------------------------------------------------
// CRUD: UPDATE
// -------------------------------------------------------------------------

func (s *Service) Update(ctx context.Context, id string, in UpdateDecoratorInput) (*DecoratorDTO, error) {
	if in.Name != nil {
		_, _ = s.pool.Exec(ctx, `UPDATE pricing_decorators SET name = $1 WHERE id = $2`, *in.Name, id)
	}
	if in.Description != nil {
		_, _ = s.pool.Exec(ctx, `UPDATE pricing_decorators SET description = $1 WHERE id = $2`, *in.Description, id)
	}
	if in.Value != nil {
		_, _ = s.pool.Exec(ctx, `UPDATE pricing_decorators SET value = $1 WHERE id = $2`, *in.Value, id)
	}
	if in.IsPercentage != nil {
		_, _ = s.pool.Exec(ctx, `UPDATE pricing_decorators SET is_percentage = $1 WHERE id = $2`, *in.IsPercentage, id)
	}
	if in.MinSubtotalCents != nil {
		_, _ = s.pool.Exec(ctx, `UPDATE pricing_decorators SET min_subtotal_cents = $1 WHERE id = $2`, *in.MinSubtotalCents, id)
	}
	if in.Active != nil {
		_, _ = s.pool.Exec(ctx, `UPDATE pricing_decorators SET active = $1 WHERE id = $2`, *in.Active, id)
	}
	if in.Priority != nil {
		_, _ = s.pool.Exec(ctx, `UPDATE pricing_decorators SET priority = $1 WHERE id = $2`, *in.Priority, id)
	}

	return s.Get(ctx, id)
}

// -------------------------------------------------------------------------
// CRUD: DELETE
// -------------------------------------------------------------------------

func (s *Service) Delete(ctx context.Context, id string) error {
	res, err := s.pool.Exec(ctx, `DELETE FROM pricing_decorators WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error al eliminar decorador: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("decorador no encontrado")
	}
	return nil
}

// -------------------------------------------------------------------------
// EJECUCION DE LA CADENA DECORATOR
// -------------------------------------------------------------------------

func (s *Service) CalculateQuote(ctx context.Context, req QuoteRequest) (Breakdown, error) {
	base := NewBaseCalculator(s.currency)

	var decoratorChain []DecoratorFunc

	if len(req.DecoratorIDs) > 0 {
		for _, decID := range req.DecoratorIDs {
			dec, err := s.Get(ctx, decID)
			if err != nil {
				continue
			}
			if !dec.Active {
				continue
			}

			switch dec.Kind {
			case "discount":
				decoratorChain = append(decoratorChain, WithDiscount(dec.Name, dec.IsPercentage, dec.Value))
			case "shipping":
				costCents := int64(dec.Value)
				decoratorChain = append(decoratorChain, WithShipping(dec.Name, costCents, dec.MinSubtotalCents))
			case "service":
				decoratorChain = append(decoratorChain, WithService(dec.Name, dec.IsPercentage, dec.Value))
			case "tax":
				decoratorChain = append(decoratorChain, WithTax(dec.Name, dec.Value))
			}
		}
	}

	calculator := Chain(base, decoratorChain...)
	return calculator.Calculate(ctx, req)
}
