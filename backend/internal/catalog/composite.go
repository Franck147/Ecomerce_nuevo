package catalog

import (
	"fmt"
	"math"
)

// =========================================================================
// PATRON DE DISEÑO: COMPOSITE (Estructural)
// =========================================================================
// Component define la interfaz comun tanto para productos individuales (hojas)
// como para combos/paquetes (compuestos).
// El cliente interactua con Component de manera uniforme sin preocuparse si
// se trata de un articulo individual o de un combo anidado complejo.

const (
	KindSimple = "simple"
	KindBundle = "bundle"
)

// StockView representa una fotografia del stock disponible por ID de producto simple.
type StockView map[string]int

// Component es la interfaz comun del patron Composite.
type Component interface {
	ID() string
	SKU() string
	Name() string
	Kind() string
	Price() int64
	Available(stock StockView) int
	Requirements(qty int) map[string]int
	Children() []Entry
}

// Entry representa un sub-componente dentro de un Bundle con su cantidad.
type Entry struct {
	Component Component
	Quantity  int
}

// -------------------------------------------------------------------------
// HOJA (LEAF): Simple
// -------------------------------------------------------------------------
// Simple representa un producto individual que no tiene hijos.
type Simple struct {
	id, sku, name string
	priceCents    int64
}

func NewSimple(id, sku, name string, priceCents int64) *Simple {
	return &Simple{id: id, sku: sku, name: name, priceCents: priceCents}
}

func (s *Simple) ID() string                          { return s.id }
func (s *Simple) SKU() string                         { return s.sku }
func (s *Simple) Name() string                        { return s.name }
func (s *Simple) Kind() string                        { return KindSimple }
func (s *Simple) Price() int64                        { return s.priceCents }
func (s *Simple) Available(stock StockView) int       { return stock[s.id] }
func (s *Simple) Requirements(qty int) map[string]int { return map[string]int{s.id: qty} }
func (s *Simple) Children() []Entry                   { return nil }

// -------------------------------------------------------------------------
// COMPUESTO (COMPOSITE): Bundle
// -------------------------------------------------------------------------
// Bundle agrupa multiples Component (que pueden ser hojas u otros bundles)
// y aplica un porcentaje de descuento global sobre la suma de sus componentes.
type Bundle struct {
	id, sku, name string
	discountPct   int
	entries       []Entry
}

func NewBundle(id, sku, name string, discountPct int) (*Bundle, error) {
	if discountPct < 0 || discountPct > 100 {
		return nil, fmt.Errorf("el porcentaje de descuento debe estar entre 0 y 100")
	}
	return &Bundle{id: id, sku: sku, name: name, discountPct: discountPct}, nil
}

func (b *Bundle) Add(c Component, qty int) error {
	if c == nil {
		return fmt.Errorf("el componente hijo no puede ser nulo")
	}
	if qty <= 0 {
		return fmt.Errorf("la cantidad del componente %q debe ser mayor a 0", c.Name())
	}
	b.entries = append(b.entries, Entry{Component: c, Quantity: qty})
	return nil
}

func (b *Bundle) ID() string        { return b.id }
func (b *Bundle) SKU() string       { return b.sku }
func (b *Bundle) Name() string      { return b.name }
func (b *Bundle) Kind() string      { return KindBundle }
func (b *Bundle) Children() []Entry { return b.entries }

// Price calcula recursivamente el precio de todos sus hijos y aplica el descuento.
func (b *Bundle) Price() int64 {
	var total int64
	for _, entry := range b.entries {
		total += entry.Component.Price() * int64(entry.Quantity)
	}
	return total * int64(100-b.discountPct) / 100
}

// Available calcula el numero maximo de combos completos disponibles segun el stock de sus hijos.
func (b *Bundle) Available(stock StockView) int {
	if len(b.entries) == 0 {
		return 0
	}
	min := math.MaxInt
	for _, entry := range b.entries {
		available := entry.Component.Available(stock) / entry.Quantity
		if available < min {
			min = available
		}
	}
	if min == math.MaxInt {
		return 0
	}
	return min
}

// Requirements expande recursivamente todos los productos simples requeridos para n unidades.
func (b *Bundle) Requirements(qty int) map[string]int {
	reqs := make(map[string]int)
	for _, entry := range b.entries {
		childReqs := entry.Component.Requirements(qty * entry.Quantity)
		for prodID, count := range childReqs {
			reqs[prodID] += count
		}
	}
	return reqs
}

// -------------------------------------------------------------------------
// SERIALIZACION DEL ARBOL (Para visualizacion en Frontend)
// -------------------------------------------------------------------------
type Node struct {
	ID         string `json:"id"`
	SKU        string `json:"sku"`
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	Quantity   int    `json:"quantity"`
	PriceCents int64  `json:"price_cents"`
	Children   []Node `json:"children,omitempty"`
}

func Describe(c Component, qty int) Node {
	node := Node{
		ID:         c.ID(),
		SKU:        c.SKU(),
		Name:       c.Name(),
		Kind:       c.Kind(),
		Quantity:   qty,
		PriceCents: c.Price(),
	}
	for _, e := range c.Children() {
		node.Children = append(node.Children, Describe(e.Component, e.Quantity))
	}
	return node
}
