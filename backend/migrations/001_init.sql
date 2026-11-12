-- Tablas para los 2 CRUDs: Composite (Catálogo) y Decorator (Reglas de Precio)

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================================================================
-- CRUD 1: PATRÓN COMPOSITE (Productos Simples y Combos)
-- =========================================================================
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    kind VARCHAR(16) NOT NULL CHECK (kind IN ('simple', 'bundle')),
    price_cents BIGINT NOT NULL DEFAULT 0,
    stock INT NOT NULL DEFAULT 0,
    discount_pct INT NOT NULL DEFAULT 0 CHECK (discount_pct >= 0 AND discount_pct <= 100),
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS bundle_items (
    bundle_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (bundle_id, product_id)
);

-- =========================================================================
-- CRUD 2: PATRÓN DECORATOR (Reglas y Decoradores de Precios)
-- =========================================================================
CREATE TABLE IF NOT EXISTS pricing_decorators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    kind VARCHAR(16) NOT NULL CHECK (kind IN ('discount', 'shipping', 'tax', 'service')),
    value NUMERIC(10,2) NOT NULL DEFAULT 0,
    is_percentage BOOLEAN NOT NULL DEFAULT false,
    min_subtotal_cents BIGINT NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT true,
    priority INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =========================================================================
-- Datos iniciales de prueba (Seeds en UTF-8)
-- =========================================================================
INSERT INTO products (id, sku, name, description, kind, price_cents, stock, discount_pct) VALUES
('11111111-1111-1111-1111-111111111111', 'MS-GPRO', 'Mouse Gamer Inalámbrico', 'Sensor óptico de alta precisión 25K DPI y batería recargable.', 'simple', 15000, 30, 0),
('22222222-2222-2222-2222-222222222222', 'KB-MECH', 'Teclado Mecánico RGB', 'Switches lineales, chasis de aluminio e iluminación RGB.', 'simple', 28000, 20, 0),
('33333333-3333-3333-3333-333333333333', 'HS-71', 'Auriculares Surround 7.1', 'Almohadillas viscoelásticas y micrófono con cancelación de ruido.', 'simple', 22000, 25, 0),
('44444444-4444-4444-4444-444444444444', 'KIT-STREAMER', 'Combo Streamer Pro', 'Combo con Mouse, Teclado y Auriculares con 15% de descuento.', 'bundle', 0, 0, 15)
ON CONFLICT (sku) DO NOTHING;

INSERT INTO bundle_items (bundle_id, product_id, quantity) VALUES
('44444444-4444-4444-4444-444444444444', '11111111-1111-1111-1111-111111111111', 1),
('44444444-4444-4444-4444-444444444444', '22222222-2222-2222-2222-222222222222', 1),
('44444444-4444-4444-4444-444444444444', '33333333-3333-3333-3333-333333333333', 1)
ON CONFLICT (bundle_id, product_id) DO NOTHING;

INSERT INTO pricing_decorators (code, name, description, kind, value, is_percentage, min_subtotal_cents, active, priority) VALUES
('DESC10', 'Cupón Bienvenida 10%', 'Descuento del 10% sobre el subtotal.', 'discount', 10.00, true, 0, true, 1),
('ENVIO-STD', 'Envío Estándar', 'Costo fijo de S/ 10.00 (Gratis si la compra supera S/ 100.00).', 'shipping', 1000, false, 10000, true, 2),
('REGALO', 'Envoltorio de Regalo', 'Empaque festivo y tarjeta de felicitación personalizada (+S/ 5.00).', 'service', 500, false, 0, true, 3),
('GARANTIA', 'Garantía Extendida (+1 año)', 'Protección adicional del 8% sobre el subtotal.', 'service', 8.00, true, 0, true, 4),
('IGV', 'Impuesto General a las Ventas (18%)', 'Tasa impositiva legal del 18% aplicada al subtotal neto.', 'tax', 18.00, true, 0, true, 5)
ON CONFLICT (code) DO NOTHING;
