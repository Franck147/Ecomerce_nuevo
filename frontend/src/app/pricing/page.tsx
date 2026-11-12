"use client";

import { useEffect, useState } from "react";
import {
  fetchDecorators,
  createDecorator,
  deleteDecorator,
  updateDecorator,
  formatMoney,
  Decorator,
} from "@/lib/api";
import {
  Plus,
  Trash2,
  RefreshCw,
  AlertCircle,
  Tag,
  Truck,
  Shield,
  Receipt,
  X,
  Search,
} from "lucide-react";

export default function PricingPage() {
  const [decorators, setDecorators] = useState<Decorator[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState("");

  const [showModal, setShowModal] = useState(false);
  const [form, setForm] = useState({
    code: "",
    name: "",
    description: "",
    kind: "discount",
    value: 10,
    is_percentage: true,
    min_subtotal_soles: 0,
    priority: 1,
  });

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);
      const decs = await fetchDecorators();
      setDecorators(decs);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createDecorator({
        code: form.code,
        name: form.name,
        description: form.description,
        kind: form.kind,
        value: form.value,
        is_percentage: form.is_percentage,
        min_subtotal_cents: Math.round((parseFloat(String(form.min_subtotal_soles)) || 0) * 100),
        priority: form.priority,
      });
      setShowModal(false);
      setForm({
        code: "",
        name: "",
        description: "",
        kind: "discount",
        value: 10,
        is_percentage: true,
        min_subtotal_soles: 0,
        priority: 1,
      });
      await loadData();
    } catch (err: any) {
      alert("Error: " + err.message);
    }
  };

  const handleDelete = async (id: string, name: string) => {
    if (!confirm(`¿Eliminar la regla "${name}"?`)) return;
    try {
      await deleteDecorator(id);
      await loadData();
    } catch (err: any) {
      alert("Error al eliminar: " + err.message);
    }
  };

  const handleToggleActive = async (dec: Decorator) => {
    try {
      await updateDecorator(dec.id, { active: !dec.active });
      await loadData();
    } catch (err: any) {
      alert("Error al actualizar: " + err.message);
    }
  };

  const filteredDecorators = decorators.filter(
    (d) =>
      d.name.toLowerCase().includes(search.toLowerCase()) ||
      d.code.toLowerCase().includes(search.toLowerCase())
  );

  const kindBadge = (kind: string) => {
    switch (kind) {
      case "discount":
        return (
          <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-rose-50 text-rose-700 border border-rose-200/60">
            <Tag className="w-3 h-3" /> Descuento
          </span>
        );
      case "shipping":
        return (
          <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-blue-50 text-blue-700 border border-blue-200/60">
            <Truck className="w-3 h-3" /> Envío
          </span>
        );
      case "service":
        return (
          <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-amber-50 text-amber-700 border border-amber-200/60">
            <Shield className="w-3 h-3" /> Servicio
          </span>
        );
      case "tax":
        return (
          <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-purple-50 text-purple-700 border border-purple-200/60">
            <Receipt className="w-3 h-3" /> Impuesto
          </span>
        );
      default:
        return <span className="px-2 py-0.5 rounded text-xs font-semibold bg-slate-100">{kind}</span>;
    }
  };

  return (
    <div className="space-y-6">
      {/* Encabezado */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-4 border-b border-slate-200">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Precios</h1>
          <p className="text-sm text-slate-500 mt-0.5">
            Administra cupones, tarifas de envío, servicios adicionales e impuestos de la tienda.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          <button
            onClick={() => setShowModal(true)}
            className="inline-flex items-center gap-2 px-3.5 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold shadow-sm transition"
          >
            <Plus className="w-4 h-4" /> Nueva Regla
          </button>
          <button
            onClick={loadData}
            title="Refrescar datos"
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
            placeholder="Buscar por nombre o código..."
            className="w-full pl-9 pr-3 py-1.5 border border-slate-200 rounded-lg text-xs bg-slate-50 focus:bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
          />
        </div>
        <span className="text-xs text-slate-500 font-mono">
          Total: {filteredDecorators.length} reglas
        </span>
      </div>

      {/* Tabla Completa de Reglas de Precio */}
      {loading ? (
        <div className="text-center py-16 text-slate-400 text-sm">Cargando reglas de precio...</div>
      ) : error ? (
        <div className="p-4 rounded-xl bg-red-50 border border-red-200 text-red-700 flex items-center gap-2 text-sm">
          <AlertCircle className="w-5 h-5 shrink-0" /> Error: {error}
        </div>
      ) : (
        <div className="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
          <table className="w-full text-left text-xs border-collapse">
            <thead>
              <tr className="bg-slate-50/80 border-b border-slate-200 text-slate-500 font-semibold uppercase tracking-wider">
                <th className="p-3.5">Tipo</th>
                <th className="p-3.5">Nombre / Descripción</th>
                <th className="p-3.5">Código</th>
                <th className="p-3.5">Valor Aplicado</th>
                <th className="p-3.5">Condición</th>
                <th className="p-3.5">Estado</th>
                <th className="p-3.5 text-right">Acción</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {filteredDecorators.map((d) => (
                <tr key={d.id} className="hover:bg-slate-50/70 transition">
                  <td className="p-3.5">{kindBadge(d.kind)}</td>
                  <td className="p-3.5">
                    <div className="font-semibold text-slate-900 text-sm">{d.name}</div>
                    <div className="text-[11px] text-slate-500 mt-0.5">{d.description || "—"}</div>
                  </td>
                  <td className="p-3.5 font-mono text-xs font-semibold text-slate-700">
                    <span className="px-2 py-0.5 rounded bg-slate-100 border border-slate-200">
                      {d.code}
                    </span>
                  </td>
                  <td className="p-3.5 font-mono font-bold text-slate-900 text-sm">
                    {d.is_percentage ? `${d.value}%` : formatMoney(d.value)}
                  </td>
                  <td className="p-3.5 text-slate-600">
                    {d.min_subtotal_cents > 0 ? (
                      <span className="font-mono text-[11px]">
                        Mínimo: {formatMoney(d.min_subtotal_cents)}
                      </span>
                    ) : (
                      <span className="text-slate-400">Sin mínimo</span>
                    )}
                  </td>
                  <td className="p-3.5">
                    <button
                      onClick={() => handleToggleActive(d)}
                      className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold transition cursor-pointer ${
                        d.active
                          ? "bg-emerald-50 text-emerald-700 border border-emerald-200/60"
                          : "bg-slate-100 text-slate-500"
                      }`}
                    >
                      <span
                        className={`w-1.5 h-1.5 rounded-full ${
                          d.active ? "bg-emerald-500" : "bg-slate-400"
                        }`}
                      ></span>
                      {d.active ? "Activo" : "Inactivo"}
                    </button>
                  </td>
                  <td className="p-3.5 text-right">
                    <button
                      onClick={() => handleDelete(d.id, d.name)}
                      title="Eliminar regla"
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

      {/* Modal Nueva Regla */}
      {showModal && (
        <div className="fixed inset-0 z-50 bg-black/40 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4">
            <div className="flex items-center justify-between border-b border-slate-100 pb-3">
              <h2 className="text-base font-bold text-slate-900">Nueva Regla de Precio</h2>
              <button onClick={() => setShowModal(false)} className="text-slate-400 hover:text-slate-600">
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleCreate} className="space-y-3 text-xs">
              <div>
                <label className="block font-semibold text-slate-700 mb-1">Código</label>
                <input
                  required
                  value={form.code}
                  onChange={(e) => setForm({ ...form, code: e.target.value.toUpperCase() })}
                  placeholder="EJ: VERANO20"
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg font-mono focus:ring-2 focus:ring-emerald-500"
                />
              </div>
              <div>
                <label className="block font-semibold text-slate-700 mb-1">Nombre</label>
                <input
                  required
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="Descuento de Temporada 20%"
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg focus:ring-2 focus:ring-emerald-500"
                />
              </div>
              <div>
                <label className="block font-semibold text-slate-700 mb-1">Descripción</label>
                <textarea
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  placeholder="Condiciones o notas de la regla..."
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg"
                  rows={2}
                />
              </div>
              <div>
                <label className="block font-semibold text-slate-700 mb-1">Tipo</label>
                <select
                  value={form.kind}
                  onChange={(e) => setForm({ ...form, kind: e.target.value })}
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg bg-white"
                >
                  <option value="discount">Descuento (Resta del subtotal)</option>
                  <option value="shipping">Envío (Costo de transporte)</option>
                  <option value="service">Servicio Adicional (Envoltorio, Garantía)</option>
                  <option value="tax">Impuesto (Aplicado al neto)</option>
                </select>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block font-semibold text-slate-700 mb-1">Valor</label>
                  <input
                    type="number"
                    step="0.01"
                    required
                    value={form.value}
                    onChange={(e) => setForm({ ...form, value: parseFloat(e.target.value) || 0 })}
                    className="w-full px-3 py-2 border border-slate-200 rounded-lg font-mono"
                  />
                </div>
                <div className="flex items-center pt-5">
                  <label className="flex items-center gap-2 cursor-pointer font-semibold text-slate-700">
                    <input
                      type="checkbox"
                      checked={form.is_percentage}
                      onChange={(e) => setForm({ ...form, is_percentage: e.target.checked })}
                      className="rounded text-emerald-600"
                    />
                    ¿Es porcentaje (%)?
                  </label>
                </div>
              </div>
              <div>
                <label className="block font-semibold text-slate-700 mb-1">
                  Mínimo de Compra (Soles, opcional)
                </label>
                <div className="relative">
                  <span className="absolute left-3 top-2 text-slate-400 font-mono text-xs">S/</span>
                  <input
                    type="number"
                    step="0.01"
                    min="0"
                    value={form.min_subtotal_soles}
                    onChange={(e) => setForm({ ...form, min_subtotal_soles: parseFloat(e.target.value) || 0 })}
                    className="w-full pl-8 pr-3 py-2 border border-slate-200 rounded-lg font-mono focus:ring-2 focus:ring-emerald-500"
                    placeholder="0.00"
                  />
                </div>
              </div>

              <div className="flex justify-end gap-2 pt-4 border-t border-slate-100">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="px-4 py-2 text-slate-600 hover:bg-slate-100 rounded-lg"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold rounded-lg shadow-sm"
                >
                  Guardar Regla
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
