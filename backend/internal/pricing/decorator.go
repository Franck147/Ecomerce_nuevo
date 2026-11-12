package pricing

import (
	"context"
	"fmt"
	"math"
)

// =========================================================================
// PATRON DE DISEÑO: DECORATOR (Estructural)
// =========================================================================
// Calculator define la interfaz componente comun.
// BaseCalculator es el componente concreto (calcula el subtotal base).
// Cada Decorator envuelve a otro Calculator y le agrega comportamiento
// (descuentos, envios, impuestos, servicios adicionales) de forma dinamica,
// sin alterar el codigo del componente base ni acoplarse entre decoradores.

type CartItem struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	PriceCents int64 `json:"price_cents"`
	Quantity   int   `json:"quantity"`
}

type QuoteRequest struct {
	Items      []CartItem `json:"items"`
	DecoratorIDs []string `json:"decorator_ids"` // IDs o codigos de decoradores a aplicar en orden
}

type Line struct {
	Concept     string `json:"concept"`
	Kind        string `json:"kind"` // base, discount, shipping, service, tax
	AmountCents int64  `json:"amount_cents"`
}

type Breakdown struct {
	Currency           string `json:"currency"`
	SubtotalCents      int64  `json:"subtotal_cents"`
	TotalDiscountCents int64  `json:"total_discount_cents"`
	TotalShippingCents int64  `json:"total_shipping_cents"`
	TotalServicesCents int64  `json:"total_services_cents"`
	TotalTaxCents      int64  `json:"total_tax_cents"`
	TotalCents         int64  `json:"total_cents"`
	Lines              []Line `json:"lines"`
}

func (b *Breakdown) Recalc() {
	b.TotalCents = b.SubtotalCents - b.TotalDiscountCents + b.TotalShippingCents + b.TotalServicesCents + b.TotalTaxCents
}

func (b *Breakdown) NetSubtotal() int64 {
	return b.SubtotalCents - b.TotalDiscountCents
}

// -------------------------------------------------------------------------
// Componente Interfaz
// -------------------------------------------------------------------------

type Calculator interface {
	Calculate(ctx context.Context, req QuoteRequest) (Breakdown, error)
}

// DecoratorFunc es una funcion que envuelve un Calculator con otro.
type DecoratorFunc func(Calculator) Calculator

// Chain encadena decoradores sobre una base. El orden de ejecucion fluye de la base hacia afuera.
func Chain(base Calculator, decorators ...DecoratorFunc) Calculator {
	curr := base
	for _, dec := range decorators {
		curr = dec(curr)
	}
	return curr
}

// -------------------------------------------------------------------------
// Componente Concreto: BaseCalculator
// -------------------------------------------------------------------------

type BaseCalculator struct {
	Currency string
}

func NewBaseCalculator(currency string) BaseCalculator {
	return BaseCalculator{Currency: currency}
}

func (b BaseCalculator) Calculate(_ context.Context, req QuoteRequest) (Breakdown, error) {
	if len(req.Items) == 0 {
		return Breakdown{}, fmt.Errorf("debe seleccionar al menos un producto para cotizar")
	}

	out := Breakdown{Currency: b.Currency}
	for _, it := range req.Items {
		if it.Quantity <= 0 {
			return Breakdown{}, fmt.Errorf("la cantidad de %q debe ser mayor a 0", it.Name)
		}
		itemTotal := it.PriceCents * int64(it.Quantity)
		out.SubtotalCents += itemTotal
		out.Lines = append(out.Lines, Line{
			Concept:     fmt.Sprintf("%s (x%d)", it.Name, it.Quantity),
			Kind:        "base",
			AmountCents: itemTotal,
		})
	}

	out.Recalc()
	return out, nil
}

// -------------------------------------------------------------------------
// Decoradores Concretos
// -------------------------------------------------------------------------

// 1. Decorador de Descuento (porcentaje o monto fijo)
type discountDecorator struct {
	inner        Calculator
	name         string
	isPercentage bool
	value        float64
}

func WithDiscount(name string, isPercentage bool, value float64) DecoratorFunc {
	return func(c Calculator) Calculator {
		return discountDecorator{inner: c, name: name, isPercentage: isPercentage, value: value}
	}
}

func (d discountDecorator) Calculate(ctx context.Context, req QuoteRequest) (Breakdown, error) {
	b, err := d.inner.Calculate(ctx, req)
	if err != nil {
		return b, err
	}

	var discountAmount int64
	if d.isPercentage {
		discountAmount = int64(math.Round(float64(b.SubtotalCents) * (d.value / 100.0)))
	} else {
		discountAmount = int64(d.value)
	}

	if discountAmount > b.SubtotalCents {
		discountAmount = b.SubtotalCents
	}

	b.TotalDiscountCents += discountAmount
	b.Lines = append(b.Lines, Line{
		Concept:     fmt.Sprintf("Descuento: %s", d.name),
		Kind:        "discount",
		AmountCents: -discountAmount,
	})

	b.Recalc()
	return b, nil
}

// 2. Decorador de Envio (costo fijo con umbral opcional para envio gratuito)
type shippingDecorator struct {
	inner           Calculator
	name            string
	costCents       int64
	freeOverSubtotal int64
}

func WithShipping(name string, costCents int64, freeOverSubtotal int64) DecoratorFunc {
	return func(c Calculator) Calculator {
		return shippingDecorator{inner: c, name: name, costCents: costCents, freeOverSubtotal: freeOverSubtotal}
	}
}

func (s shippingDecorator) Calculate(ctx context.Context, req QuoteRequest) (Breakdown, error) {
	b, err := s.inner.Calculate(ctx, req)
	if err != nil {
		return b, err
	}

	if s.freeOverSubtotal > 0 && b.NetSubtotal() >= s.freeOverSubtotal {
		b.Lines = append(b.Lines, Line{
			Concept:     fmt.Sprintf("Envío: %s (¡GRATIS por compra mayor a S/ %.2f!)", s.name, float64(s.freeOverSubtotal)/100.0),
			Kind:        "shipping",
			AmountCents: 0,
		})
	} else {
		b.TotalShippingCents += s.costCents
		b.Lines = append(b.Lines, Line{
			Concept:     fmt.Sprintf("Envío: %s", s.name),
			Kind:        "shipping",
			AmountCents: s.costCents,
		})
	}

	b.Recalc()
	return b, nil
}

// 3. Decorador de Servicio Adicional (envoltorio, garantia, seguro)
type serviceDecorator struct {
	inner        Calculator
	name         string
	isPercentage bool
	value        float64
}

func WithService(name string, isPercentage bool, value float64) DecoratorFunc {
	return func(c Calculator) Calculator {
		return serviceDecorator{inner: c, name: name, isPercentage: isPercentage, value: value}
	}
}

func (s serviceDecorator) Calculate(ctx context.Context, req QuoteRequest) (Breakdown, error) {
	b, err := s.inner.Calculate(ctx, req)
	if err != nil {
		return b, err
	}

	var serviceCost int64
	if s.isPercentage {
		serviceCost = int64(math.Round(float64(b.NetSubtotal()) * (s.value / 100.0)))
	} else {
		serviceCost = int64(s.value)
	}

	b.TotalServicesCents += serviceCost
	b.Lines = append(b.Lines, Line{
		Concept:     fmt.Sprintf("Servicio: %s", s.name),
		Kind:        "service",
		AmountCents: serviceCost,
	})

	b.Recalc()
	return b, nil
}

// 4. Decorador de Impuesto / IVA
type taxDecorator struct {
	inner   Calculator
	name    string
	ratePct float64
}

func WithTax(name string, ratePct float64) DecoratorFunc {
	return func(c Calculator) Calculator {
		return taxDecorator{inner: c, name: name, ratePct: ratePct}
	}
}

func (t taxDecorator) Calculate(ctx context.Context, req QuoteRequest) (Breakdown, error) {
	b, err := t.inner.Calculate(ctx, req)
	if err != nil {
		return b, err
	}

	taxableBase := b.NetSubtotal()
	taxAmount := int64(math.Round(float64(taxableBase) * (t.ratePct / 100.0)))

	b.TotalTaxCents += taxAmount
	b.Lines = append(b.Lines, Line{
		Concept:     fmt.Sprintf("Impuesto: %s (%.1f%%)", t.name, t.ratePct),
		Kind:        "tax",
		AmountCents: taxAmount,
	})

	b.Recalc()
	return b, nil
}
