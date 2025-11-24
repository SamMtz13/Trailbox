<section class="card">
  <p class="badge">Gateway REST → gRPC</p>
  <h1 class="text-3xl md:text-4xl font-bold text-forest mt-3">Trailbox</h1>
  <p class="text-emerald-800 mt-2 max-w-2xl">
    UI ligera en Svelte 5 sobre el gateway HTTP. Cada sección consulta a un microservicio gRPC a través del
    adaptador REST /api/*.
  </p>
  <div class="grid md:grid-cols-2 gap-4 mt-6">
    <div class="card">
      <h3 class="text-xl font-semibold text-forest">Servicios</h3>
      <ul class="mt-3 space-y-2 text-emerald-800 text-sm">
        <li>Usuarios, Rutas, Entrenamientos, Reseñas</li>
        <li>Leaderboard, Notificaciones, Mapas</li>
        <li>Gateway público (LoadBalancer) y PostgreSQL</li>
      </ul>
    </div>
    <div class="card">
      <h3 class="text-xl font-semibold text-forest">Cómo usar</h3>
      <ol class="mt-3 list-decimal list-inside text-emerald-800 text-sm space-y-1">
        <li>Selecciona la sección en el menú.</li>
        <li>Consulta datos y envía peticiones de ejemplo.</li>
        <li>El gateway las redirige a los microservicios gRPC.</li>
      </ol>
    </div>
  </div>
</section>
