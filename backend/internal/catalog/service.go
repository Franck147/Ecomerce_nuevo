package catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductDTO struct {
	ID          string    `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Kind        string    `json:"kind"`
	PriceCents  int64     `json:"price_cents"`
	Stock       int       `json:"stock"`
	DiscountPct int       `json:"discount_pct"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	Available   int       `json:"available"`
	Tree        *Node     `json:"tree,omitempty"`
}

type CreateSimpleInput struct {
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int64  `json:"price_cents"`
	Stock       int    `json:"stock"`
}

type BundleItemInput struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CreateBundleInput struct {
	SKU         string            `json:"sku"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	DiscountPct int               `json:"discount_pct"`
	Items       []BundleItemInput `json:"items"`
}

type UpdateProductInput struct {
	Name        *string            `json:"name"`
	Description *string            `json:"description"`
	PriceCents  *int64             `json:"price_cents"`
	Stock       *int               `json:"stock"`
	DiscountPct *int               `json:"discount_pct"`
	Active      *bool              `json:"active"`
	Items       *[]BundleItemInput `json:"items"`
}

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// -------------------------------------------------------------------------
// CRUD: CREATE
// -------------------------------------------------------------------------

func (s *Service) CreateSimple(ctx context.Context, in CreateSimpleInput) (*ProductDTO, error) {
	if in.SKU == "" || in.Name == "" {
		return nil, fmt.Errorf("SKU y Nombre son obligatorios")
	}
	if in.PriceCents < 0 {
		return nil, fmt.Errorf("el precio no puede ser negativo")
	}

	var p ProductDTO
	err := s.pool.QueryRow(ctx, `
		INSERT INTO products (sku, name, description, kind, price_cents, stock, discount_pct)
		VALUES ($1, $2, $3, 'simple', $4, $5, 0)
		RETURNING id, sku, name, description, kind, price_cents, stock, discount_pct, active, created_at
	`, in.SKU, in.Name, in.Description, in.PriceCents, in.Stock).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Description, &p.Kind, &p.PriceCents, &p.Stock, &p.DiscountPct, &p.Active, &p.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error al crear producto simple: %w", err)
	}
	p.Available = p.Stock
	return &p, nil
}

func (s *Service) CreateBundle(ctx context.Context, in CreateBundleInput) (*ProductDTO, error) {
	if in.SKU == "" || in.Name == "" {
		return nil, fmt.Errorf("SKU y Nombre son obligatorios")
	}
	if len(in.Items) == 0 {
		return nil, fmt.Errorf("un combo debe incluir al menos un producto componente")
	}
	if in.DiscountPct < 0 || in.DiscountPct > 100 {
		return nil, fmt.Errorf("el descuento debe estar entre 0 y 100%%")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var p ProductDTO
	err = tx.QueryRow(ctx, `
		INSERT INTO products (sku, name, description, kind, price_cents, stock, discount_pct)
		VALUES ($1, $2, $3, 'bundle', 0, 0, $4)
		RETURNING id, sku, name, description, kind, price_cents, stock, discount_pct, active, created_at
	`, in.SKU, in.Name, in.Description, in.DiscountPct).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Description, &p.Kind, &p.PriceCents, &p.Stock, &p.DiscountPct, &p.Active, &p.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error al registrar combo: %w", err)
	}

	for _, it := range in.Items {
		if it.Quantity <= 0 {
			return nil, fmt.Errorf("la cantidad de cada componente debe ser mayor a 0")
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO bundle_items (bundle_id, product_id, quantity)
			VALUES ($1, $2, $3)
		`, p.ID, it.ProductID, it.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error al asociar componente %s: %w", it.ProductID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.Get(ctx, p.ID)
}

// -------------------------------------------------------------------------
// CRUD: READ
// -------------------------------------------------------------------------

func (s *Service) List(ctx context.Context) ([]ProductDTO, error) {
	stockMap, err := s.loadStockMap(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := s.pool.Query(ctx, `
		SELECT id, sku, name, description, kind, price_cents, stock, discount_pct, active, created_at
		FROM products
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ProductDTO
	for rows.Next() {
		var p ProductDTO
		if err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Description, &p.Kind, &p.PriceCents, &p.Stock, &p.DiscountPct, &p.Active, &p.CreatedAt,
		); err != nil {
			return nil, err
		}

		if p.Kind == KindSimple {
			p.Available = p.Stock
		} else {
			comp, err := s.buildCompositeTree(ctx, p.ID, make(map[string]bool))
			if err == nil && comp != nil {
				p.PriceCents = comp.Price()
				p.Available = comp.Available(stockMap)
			}
		}

		result = append(result, p)
	}

	return result, nil
}

func (s *Service) Get(ctx context.Context, id string) (*ProductDTO, error) {
	stockMap, err := s.loadStockMap(ctx)
	if err != nil {
		return nil, err
	}

	var p ProductDTO
	err = s.pool.QueryRow(ctx, `
		SELECT id, sku, name, description, kind, price_cents, stock, discount_pct, active, created_at
		FROM products
		WHERE id = $1
	`, id).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Description, &p.Kind, &p.PriceCents, &p.Stock, &p.DiscountPct, &p.Active, &p.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("producto no encontrado: %w", err)
	}

	if p.Kind == KindSimple {
		p.Available = p.Stock
	} else {
		comp, err := s.buildCompositeTree(ctx, p.ID, make(map[string]bool))
		if err != nil {
			return nil, err
		}
		p.PriceCents = comp.Price()
		p.Available = comp.Available(stockMap)
		tree := Describe(comp, 1)
		p.Tree = &tree
	}

	return &p, nil
}

// -------------------------------------------------------------------------
// CRUD: UPDATE
// -------------------------------------------------------------------------

func (s *Service) Update(ctx context.Context, id string, in UpdateProductInput) (*ProductDTO, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var kind string
	err = tx.QueryRow(ctx, `SELECT kind FROM products WHERE id = $1`, id).Scan(&kind)
	if err != nil {
		return nil, fmt.Errorf("producto no encontrado: %w", err)
	}

	if in.Name != nil {
		_, _ = tx.Exec(ctx, `UPDATE products SET name = $1 WHERE id = $2`, *in.Name, id)
	}
	if in.Description != nil {
		_, _ = tx.Exec(ctx, `UPDATE products SET description = $1 WHERE id = $2`, *in.Description, id)
	}
	if in.Active != nil {
		_, _ = tx.Exec(ctx, `UPDATE products SET active = $1 WHERE id = $2`, *in.Active, id)
	}

	if kind == KindSimple {
		if in.PriceCents != nil {
			_, _ = tx.Exec(ctx, `UPDATE products SET price_cents = $1 WHERE id = $2`, *in.PriceCents, id)
		}
		if in.Stock != nil {
			_, _ = tx.Exec(ctx, `UPDATE products SET stock = $1 WHERE id = $2`, *in.Stock, id)
		}
	} else if kind == KindBundle {
		if in.DiscountPct != nil {
			_, _ = tx.Exec(ctx, `UPDATE products SET discount_pct = $1 WHERE id = $2`, *in.DiscountPct, id)
		}
		if in.Items != nil {
			_, _ = tx.Exec(ctx, `DELETE FROM bundle_items WHERE bundle_id = $1`, id)
			for _, it := range *in.Items {
				if it.Quantity <= 0 {
					return nil, fmt.Errorf("la cantidad debe ser mayor a 0")
				}
				_, err = tx.Exec(ctx, `
					INSERT INTO bundle_items (bundle_id, product_id, quantity)
					VALUES ($1, $2, $3)
				`, id, it.ProductID, it.Quantity)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.Get(ctx, id)
}

// -------------------------------------------------------------------------
// CRUD: DELETE
// -------------------------------------------------------------------------

func (s *Service) Delete(ctx context.Context, id string) error {
	res, err := s.pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error al eliminar: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("producto no encontrado")
	}
	return nil
}

// -------------------------------------------------------------------------
// HELPERS DEL PATRON COMPOSITE
// -------------------------------------------------------------------------

func (s *Service) loadStockMap(ctx context.Context) (StockView, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, stock FROM products WHERE kind = 'simple'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stock := make(StockView)
	for rows.Next() {
		var id string
		var st int
		if err := rows.Scan(&id, &st); err != nil {
			return nil, err
		}
		stock[id] = st
	}
	return stock, nil
}

func (s *Service) buildCompositeTree(ctx context.Context, id string, visited map[string]bool) (Component, error) {
	if visited[id] {
		return nil, fmt.Errorf("referencia circular detectada en combo %s", id)
	}
	visited[id] = true

	var sku, name, kind string
	var priceCents int64
	var discountPct int

	err := s.pool.QueryRow(ctx, `
		SELECT sku, name, kind, price_cents, discount_pct
		FROM products
		WHERE id = $1
	`, id).Scan(&sku, &name, &kind, &priceCents, &discountPct)
	if err != nil {
		return nil, err
	}

	if kind == KindSimple {
		return NewSimple(id, sku, name, priceCents), nil
	}

	bundle, err := NewBundle(id, sku, name, discountPct)
	if err != nil {
		return nil, err
	}

	rows, err := s.pool.Query(ctx, `
		SELECT product_id, quantity
		FROM bundle_items
		WHERE bundle_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type itemEntry struct {
		prodID string
		qty    int
	}
	var entries []itemEntry
	for rows.Next() {
		var e itemEntry
		if err := rows.Scan(&e.prodID, &e.qty); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	for _, entry := range entries {
		childVisited := make(map[string]bool)
		for k, v := range visited {
			childVisited[k] = v
		}
		child, err := s.buildCompositeTree(ctx, entry.prodID, childVisited)
		if err != nil {
			return nil, err
		}
		if err := bundle.Add(child, entry.qty); err != nil {
			return nil, err
		}
	}

	return bundle, nil
}
