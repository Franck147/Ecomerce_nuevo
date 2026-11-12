# Ecommerce — Arquitectura Monolítica Modular (Patrones GoF)

Este proyecto implementa una plataforma de **Ecommerce** construida con arquitectura **Monolítica Modular** en Go (Backend) y Next.js 15 (Frontend), aplicando 2 CRUDs principales basados en **Patrones de Diseño GoF**:

---

## 📁 Ubicación del Proyecto

- **Ruta Raíz:** `D:\ECOMERCE-TRABAJO-GRUPAL\ecommerce-patrones`
- **Backend (Go REST API :8081):** `D:\ECOMERCE-TRABAJO-GRUPAL\ecommerce-patrones\backend`
- **Frontend (Next.js :3001):** `D:\ECOMERCE-TRABAJO-GRUPAL\ecommerce-patrones\frontend`
- **Script de Inicio:** `D:\ECOMERCE-TRABAJO-GRUPAL\ecommerce-patrones\start.ps1`

---

## 🏛️ Estructura Limpia — Monolito Modular

```text
ecommerce-patrones/
├── backend/
│   ├── cmd/api/main.go               # Orquestador del Monolito Modular
│   ├── internal/
│   │   ├── catalog/                  # MÓDULO 1: CRUD Catálogo (Patrón Composite)
│   │   │   ├── composite.go          # Implementación Component, Simple y Bundle
│   │   │   └── service.go            # Lógica de negocio y persistencia SQL
│   │   ├── pricing/                  # MÓDULO 2: CRUD Precios y Reglas (Patrón Decorator)
│   │   │   ├── decorator.go          # Implementación Calculator, Base y Decoradores
│   │   │   └── service.go            # Lógica de reglas, cálculo y persistencia SQL
│   │   ├── singleton/                # MÓDULO TRANSVERSAL: Patrón Singleton
│   │   │   ├── db.go                 # Pool de conexiones PostgreSQL thread-safe (sync.Once)
│   │   │   └── config.go             # Carga única de variables de entorno
│   │   └── httpapi/                  # Capa de transporte HTTP (Rutas y Handlers)
│   │       ├── api.go                # Router y configuración CORS
│   │       ├── handlers_composite.go # Endpoints REST de Productos
│   │       └── handlers_decorator.go # Endpoints REST de Precios
│   └── migrations/
│       ├── 001_init.sql              # Esquema PostgreSQL y datos iniciales
│       └── embed.go                  # Migraciones embebidas en el binario
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── layout.tsx            # Sidebar profesional (Ecommerce, Productos, Precios)
│   │   │   ├── page.tsx              # Redirección directa a /products
│   │   │   ├── products/page.tsx     # CRUD de Productos (precios en S/ y stock)
│   │   │   └── pricing/page.tsx      # CRUD de Reglas de Precio (descuentos, envíos, IGV)
│   │   └── lib/
│   │       └── api.ts                # Cliente API tipado conectado a :8081
├── start.ps1                         # Script de inicio rápido (PostgreSQL + Backend + Frontend)
└── README.md                         # Documentación técnica
```

---

## 🔍 ¿Cómo saber y comprobar que tiene Composite y Decorator?

### 1. Patrón Composite (Módulo `catalog`)
- **Archivo:** `backend/internal/catalog/composite.go`
- **Interfaz Componente:**
  ```go
  type Component interface {
      GetID() string
      GetName() string
      GetPrice() int64
      GetStock() int
      IsComposite() bool
      GetChildren() []ChildItem
  }
  ```
- **Hoja (`Simple`):** Representa un producto individual con su precio y stock físico directo.
- **Compuesto (`Bundle`):** Representa un producto compuesto o paquete que contiene una lista de `Component`. Su precio se calcula dinámicamente sumando recursivamente los precios de sus hijos (con descuento opcional), y su disponibilidad de stock es el mínimo del stock de sus componentes.
- **Verificación en el código:** `Simple` y `Bundle` implementan exactamente la misma interfaz `Component`, permitiendo que el catálogo procese productos simples o compuestos de manera uniforme y transparente.

---

### 2. Patrón Decorator (Módulo `pricing`)
- **Archivo:** `backend/internal/pricing/decorator.go`
- **Interfaz Componente Base:**
  ```go
  type Calculator interface {
      Calculate(ctx context.Context, req QuoteRequest) (*QuoteResult, error)
  }
  ```
- **Componente Concreto:** `BaseCalculator` calcula el subtotal directo (`precio × cantidad`).
- **Decoradores Concretos:**
  - `discountDecorator`: Envuelve al calculador anterior y aplica descuentos en porcentaje o monto fijo.
  - `shippingDecorator`: Envuelve y añade costo de envío según la regla de pedido mínimo.
  - `serviceDecorator`: Envuelve y añade recargos por servicios opcionales (garantía, empaque).
  - `taxDecorator`: Envuelve y calcula el impuesto (IGV 18%) sobre la base imponible acumulada.
- **Función de Encadenamiento:** `Chain(base Calculator, decorators ...DecoratorFunc) Calculator` envuelve dinámicamente los decoradores en tiempo de ejecución según las reglas activas configuradas en el CRUD.

---

### 3. Patrón Singleton (Módulo transversal `singleton`)
- **Archivo:** `backend/internal/singleton/db.go` y `backend/internal/singleton/config.go`
- Garantiza mediante `sync.Once` que el pool de conexiones a la base de datos PostgreSQL (`*pgxpool.Pool`) se inicialice una sola vez y sea reutilizado por todos los módulos del monolito sin crear conexiones redundantes.

---

## 🚀 Ejecución

Para iniciar el sistema completo:
```powershell
.\start.ps1
```

- **Frontend:** [http://localhost:3001/products](http://localhost:3001/products)
- **Precios:** [http://localhost:3001/pricing](http://localhost:3001/pricing)
- **Backend API:** [http://localhost:8081/health](http://localhost:8081/health)
"# Ecomerce_nuevo" 
