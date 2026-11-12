const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";

export interface Node {
  id: string;
  sku: string;
  name: string;
  kind: string;
  quantity: number;
  price_cents: number;
  children?: Node[];
}

export interface Product {
  id: string;
  sku: string;
  name: string;
  description: string;
  kind: "simple" | "bundle";
  price_cents: number;
  stock: number;
  discount_pct: number;
  active: boolean;
  created_at: string;
  available: number;
  tree?: Node;
}

export interface Decorator {
  id: string;
  code: string;
  name: string;
  description: string;
  kind: "discount" | "shipping" | "service" | "tax";
  value: number;
  is_percentage: boolean;
  min_subtotal_cents: number;
  active: boolean;
  priority: number;
  created_at: string;
}

export interface Line {
  concept: string;
  kind: string;
  amount_cents: number;
}

export interface QuoteBreakdown {
  currency: string;
  subtotal_cents: number;
  total_discount_cents: number;
  total_shipping_cents: number;
  total_services_cents: number;
  total_tax_cents: number;
  total_cents: number;
  lines: Line[];
}

export function formatMoney(cents: number): string {
  return "S/ " + (cents / 100).toFixed(2);
}

// ------------------- CRUD 1: COMPOSITE -------------------

export async function fetchProducts(): Promise<Product[]> {
  const res = await fetch(`${API_BASE}/api/v1/products`);
  if (!res.ok) throw new Error("Error al obtener productos");
  return res.json();
}

export async function fetchProduct(id: string): Promise<Product> {
  const res = await fetch(`${API_BASE}/api/v1/products/${id}`);
  if (!res.ok) throw new Error("Error al obtener detalle del producto");
  return res.json();
}

export async function createSimpleProduct(data: {
  sku: string;
  name: string;
  description: string;
  price_cents: number;
  stock: number;
}): Promise<Product> {
  const res = await fetch(`${API_BASE}/api/v1/products/simple`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "Error al crear producto");
  }
  return res.json();
}

export async function createBundleProduct(data: {
  sku: string;
  name: string;
  description: string;
  discount_pct: number;
  items: { product_id: string; quantity: number }[];
}): Promise<Product> {
  const res = await fetch(`${API_BASE}/api/v1/products/bundle`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "Error al crear combo");
  }
  return res.json();
}

export async function updateProduct(id: string, data: any): Promise<Product> {
  const res = await fetch(`${API_BASE}/api/v1/products/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error("Error al actualizar producto");
  return res.json();
}

export async function deleteProduct(id: string): Promise<void> {
  const res = await fetch(`${API_BASE}/api/v1/products/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) throw new Error("Error al eliminar producto");
}

// ------------------- CRUD 2: DECORATOR -------------------

export async function fetchDecorators(): Promise<Decorator[]> {
  const res = await fetch(`${API_BASE}/api/v1/decorators`);
  if (!res.ok) throw new Error("Error al obtener decoradores");
  return res.json();
}

export async function createDecorator(data: {
  code: string;
  name: string;
  description: string;
  kind: string;
  value: number;
  is_percentage: boolean;
  min_subtotal_cents: number;
  priority: number;
}): Promise<Decorator> {
  const res = await fetch(`${API_BASE}/api/v1/decorators`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "Error al registrar decorador");
  }
  return res.json();
}

export async function updateDecorator(id: string, data: any): Promise<Decorator> {
  const res = await fetch(`${API_BASE}/api/v1/decorators/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error("Error al actualizar decorador");
  return res.json();
}

export async function deleteDecorator(id: string): Promise<void> {
  const res = await fetch(`${API_BASE}/api/v1/decorators/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) throw new Error("Error al eliminar decorador");
}

export async function calculateQuote(req: {
  items: { product_id: string; name: string; price_cents: number; quantity: number }[];
  decorator_ids: string[];
}): Promise<QuoteBreakdown> {
  const res = await fetch(`${API_BASE}/api/v1/pricing/quote`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "Error al cotizar");
  }
  return res.json();
}
