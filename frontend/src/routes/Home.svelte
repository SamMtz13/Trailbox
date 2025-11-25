<section class="relative overflow-hidden rounded-3xl border border-white/10 bg-gradient-to-br from-emerald-500/20 to-slate-900/60 p-6 shadow-[0_25px_80px_rgba(6,182,212,0.25)]">
  <div class="flex flex-col gap-6 lg:flex-row lg:items-center">
    <div class="flex-1">
      <p class="badge bg-white/10 border-white/30 text-white/80">Gateway REST → gRPC</p>
      <h1 class="mt-3 text-4xl font-bold text-white drop-shadow-md">Trailbox Dashboard</h1>
      <p class="mt-4 text-slate-200/90 sm:text-lg">
        Observa en tiempo real el flujo HTTP→gRPC. Cada sección consume el gateway y éste traduce las peticiones
        hacia los microservicios usuarios, rutas, workouts, reseñas, leaderboard, notificaciones y mapas.
      </p>
      <div class="mt-6 flex flex-wrap gap-4">
        <div class="flex items-center gap-2 rounded-xl border border-white/20 px-4 py-2 text-sm text-white/70 bg-slate-900/30">
          <span class="text-xl">🛰️</span>
          <span>HTTP → gRPC Gateway</span>
        </div>
        <div class="flex items-center gap-2 rounded-xl border border-white/20 px-4 py-2 text-sm text-white/70 bg-slate-900/30">
          <span class="text-xl">📡</span>
          <span>PostgreSQL bootstrap</span>
        </div>
        <div class="flex items-center gap-2 rounded-xl border border-white/20 px-4 py-2 text-sm text-white/70 bg-slate-900/30">
          <span class="text-xl">🧭</span>
          <span>Kind / Kubernetes</span>
        </div>
      </div>
    </div>
    <div class="flex-1 space-y-4">
      <div class="card bg-slate-900/80 border-white/5 text-white">
        <h3 class="text-xl font-semibold">Micros conectados</h3>
        <ul class="mt-3 grid grid-cols-2 gap-2 text-sm text-slate-200/80">
          <li class="flex items-center gap-2">✅ Usuarios</li>
          <li class="flex items-center gap-2">✅ Rutas</li>
          <li class="flex items-center gap-2">✅ Workouts</li>
          <li class="flex items-center gap-2">✅ Reseñas</li>
          <li class="flex items-center gap-2">✅ Leaderboard</li>
          <li class="flex items-center gap-2">✅ Notificaciones</li>
          <li class="flex items-center gap-2">✅ Mapas</li>
        </ul>
      </div>
      <div class="card bg-slate-900/80 border-white/5 text-white">
        <h3 class="text-xl font-semibold">Cómo usar</h3>
        <ol class="mt-3 list-decimal list-inside space-y-1 text-sm text-slate-200/80">
          <li>Selecciona una pestaña del menú superior.</li>
          <li>Dispara peticiones de ejemplo (GET/POST).</li>
          <li>Observa cómo el gateway enruta hacia gRPC.</li>
        </ol>
      </div>
    </div>
  </div>
</section>
