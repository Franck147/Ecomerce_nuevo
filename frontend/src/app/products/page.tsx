"use client";

import { useEffect, useState } from "react";
import {
  fetchProducts,
  createSimpleProduct,
  deleteProduct,
  formatMoney,
  Product,
} from "@/lib/api";
import {
  Plus,
  Trash2,
  RefreshCw,
  AlertCircle,
  Search,
  X,
} from "lucide-react";

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState("");

  const [showSimpleModal, setShowSimpleModal] = useState(false);

  // Formulario Producto Simple
  const [simpleForm, setSimpleForm] = useState({
    sku: "",
    name: "",
    description: "",
    price_soles: 120,
    stock: 20,
  });

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchProducts();
      // Solo productos simples
      setProducts(data.filter((p) => p.kind === "simple"));
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleCreateSimple = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createSimpleProduct({
        sku: simpleForm.sku,
        name: simpleForm.name,
        description: simpleForm.description,
        price_cents: Math.round((parseFloat(String(simpleForm.price_soles)) || 0) * 100),
        stock: simpleForm.stock,
      });
      setShowSimpleModal(false);
      setSimpleForm({ sku: "", name: "", description: "", price_soles: 120, stock: 20 });
      await loadData();
    } catch (err: any) {
      alert("Error: " + err.message);
    }
  };

  const handleDelete = async (id: string, name: string) => {
    if (!confirm(`¿Estás seguro de eliminar el producto "${name}"?`)) return;
    try {
      await deleteProduct(id);
      await loadData();
    } catch (err: any) {
      alert("Error al eliminar: " + err.message);
    }
  };

  const filteredProducts = products.filter(
    (p) =>
      p.name.toLowerCase().includes(search.toLowerCase()) ||
      p.sku.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-6">
      {/* Encabezado */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-4 border-b border-slate-200">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Productos</h1>
          <p className="text-sm text-slate-500 mt-0.5">
            Administra el catálogo de productos individuales de la tienda.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          <button
            onClick={() => setShowSimpleModal(true)}
            className="inline-flex items-center gap-2 px-3.5 py-2 rounded-xl bg-slate-900 hover:bg-slate-800 text-white text-xs font-semibold shadow-sm transition"
          >
            <Plus className="w-4 h-4" /> Nuevo Producto
          </button>
          <button
            onClick={loadData}
            title="Actualizar catálogo"
            className="p-2 rounded-xl border border-slate-200 text-slate-600 hover:bg-slate-100 transition"
          >
            <RefreshCw className="w-4 h-4" />
          </button>
        </div>
      </div>

      {/* Buscador */}
      <div className="flex items-center justify-between gap-4 bg-white p-3.5 rounded-2xl border border-slate-200/80 shadow-sm">
        <div className="relative w-full sm:w-80">
          <Search className="w-4 h-4 absolute left-3 top-2.5 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Buscar por nombre o SKU..."
            className="w-full pl-9 pr-3 py-1.5 border border-slate-200 rounded-lg text-xs bg-slate-50 focus:bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
        <span className="text-xs text-slate-500 font-mono">
          Total: {filteredProducts.length} productos
        </span>
      </div>

      {/* Tabla de Productos Full Width */}
      {loading ? (
        <div className="text-center py-16 text-slate-400 text-sm">Cargando productos...</div>
      ) : error ? (
        <div className="p-4 rounded-xl bg-red-50 border border-red-200 text-red-700 flex items-center gap-2 text-sm">
          <AlertCircle className="w-5 h-5 shrink-0" /> Error: {error}
        </div>
      ) : (
        <div className="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
          <table className="w-full text-left text-xs border-collapse">
            <thead>
              <tr className="bg-slate-50/80 border-b border-slate-200 text-slate-500 font-semibold uppercase tracking-wider">
                <th className="p-3.5">Artículo / SKU</th>
                <th className="p-3.5">Descripción</th>
                <th className="p-3.5">Precio</th>
                <th className="p-3.5">Stock Disponible</th>
                <th className="p-3.5 text-right">Acciones</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {filteredProducts.map((p) => (
                <tr key={p.id} className="hover:bg-slate-50/70 transition">
                  <td className="p-3.5">
                    <div className="font-semibold text-slate-900 text-sm">{p.name}</div>
                    <div className="text-[11px] text-slate-400 font-mono mt-0.5">SKU: {p.sku}</div>
                  </td>
                  <td className="p-3.5 text-slate-500 max-w-md">
                    {p.description || "—"}
                  </td>
                  <td className="p-3.5 font-mono font-bold text-slate-900 text-sm">
                    {formatMoney(p.price_cents)}
                  </td>
                  <td className="p-3.5">
                    <span
                      className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-medium ${
                        p.stock > 0
                          ? "bg-emerald-50 text-emerald-700 border border-emerald-200/50"
                          : "bg-red-50 text-red-700 border border-red-200/50"
                      }`}
                    >
                      <span
                        className={`w-1.5 h-1.5 rounded-full ${
                          p.stock > 0 ? "bg-emerald-500" : "bg-red-500"
                        }`}
                      ></span>
                      {p.stock > 0 ? `${p.stock} unidades` : "Agotado"}
                    </span>
                  </td>
                  <td className="p-3.5 text-right">
                    <button
                      onClick={() => handleDelete(p.id, p.name)}
                      title="Eliminar producto"
                      className="p-1.5 rounded-lg text-slate-400 hover:text-red-600 hover:bg-red-50 transition"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Modal Nuevo Producto */}
      {showSimpleModal && (
        <div className="fixed inset-0 z-50 bg-black/40 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4">
            <div className="flex items-center justify-between border-b border-slate-100 pb-3">
              <h2 className="text-base font-bold text-slate-900">Nuevo Producto</h2>
              <button onClick={() => setShowSimpleModal(false)} className="text-slate-400 hover:text-slate-600">
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleCreateSimple} className="space-y-3 text-xs">
              <div>
                <label className="block font-semibold text-slate-700 mb-1">SKU</label>
                <input
                  required
                  value={simpleForm.sku}
                  onChange={(e) => setSimpleForm({ ...simpleForm, sku: e.target.value.toUpperCase() })}
                  placeholder="EJ: PROD-01"
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg font-mono focus:ring-2 focus:ring-indigo-500"
                />
              </div>
              <div>
                <label className="block font-semibold text-slate-700 mb-1">Nombre</label>
                <input
                  required
                  value={simpleForm.name}
                  onChange={(e) => setSimpleForm({ ...simpleForm, name: e.target.value })}
                  placeholder="Nombre del artículo"
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg focus:ring-2 focus:ring-indigo-500"
                />
              </div>
              <div>
                <label className="block font-semibold text-slate-700 mb-1">Descripción</label>
                <textarea
                  value={simpleForm.description}
                  onChange={(e) => setSimpleForm({ ...simpleForm, description: e.target.value })}
                  placeholder="Descripción breve..."
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg"
                  rows={2}
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block font-semibold text-slate-700 mb-1">Precio (S/)</label>
                  <div className="relative">
                    <span className="absolute left-3 top-2 text-slate-400 font-mono text-xs">S/</span>
                    <input
                      type="number"
                      step="0.01"
                      min="0"
                      required
                      value={simpleForm.price_soles}
                      onChange={(e) => setSimpleForm({ ...simpleForm, price_soles: parseFloat(e.target.value) || 0 })}
                      className="w-full pl-8 pr-3 py-2 border border-slate-200 rounded-lg font-mono focus:ring-2 focus:ring-indigo-500"
                      placeholder="0.00"
                    />
                  </div>
                </div>
                <div>
                  <label className="block font-semibold text-slate-700 mb-1">Stock</label>
                  <input
                    type="number"
                    min="0"
                    required
                    value={simpleForm.stock}
                    onChange={(e) => setSimpleForm({ ...simpleForm, stock: parseInt(e.target.value) || 0 })}
                    className="w-full px-3 py-2 border border-slate-200 rounded-lg font-mono focus:ring-2 focus:ring-indigo-500"
                  />
                </div>
              </div>
              <div className="flex justify-end gap-2 pt-4 border-t border-slate-100">
                <button
                  type="button"
                  onClick={() => setShowSimpleModal(false)}
                  className="px-4 py-2 text-slate-600 hover:bg-slate-100 rounded-lg"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-lg shadow-sm"
                >
                  Guardar
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
