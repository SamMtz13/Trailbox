<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';

  type Route = { id: string; name: string; distance_km: number; elevation_gain: number };
  type Review = { id: string; user_id: string; route_id: string; rating: number; comment: string; created_at: string };

  let routes: Route[] = [];
  let loading = true;
  let error = '';

  let selectedRoute: Route | null = null;
  let reviews: Review[] = [];
  let reviewsLoading = false;
  let reviewsError = '';

  async function loadRoutes() {
    loading = true;
    error = '';
    try {
      const data = await api.listRoutes();
      routes = data.routes || [];
      if (!selectedRoute && routes.length > 0) {
        selectRoute(routes[0]);
      } else if (selectedRoute) {
        const updated = routes.find((r) => r.id === selectedRoute?.id);
        if (updated) {
          selectedRoute = updated;
        }
      }
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  async function loadReviews(routeId: string) {
    reviewsLoading = true;
    reviewsError = '';
    try {
      const data = await api.getReviews(routeId);
      reviews = data.reviews || [];
    } catch (err) {
      reviewsError = err instanceof Error ? err.message : String(err);
      reviews = [];
    } finally {
      reviewsLoading = false;
    }
  }

  function selectRoute(route: Route) {
    if (selectedRoute?.id === route.id) return;
    selectedRoute = route;
    loadReviews(route.id);
  }

  onMount(() => {
    loadRoutes();
  });
</script>

<section class="space-y-5 rounded-3xl border border-white/10 bg-slate-900/70 p-6 shadow-2xl backdrop-blur">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <p class="badge bg-white/5 text-cyan-200">Rutas</p>
      <h2 class="text-2xl font-semibold text-white">Catálogo y reseñas</h2>
      <p class="text-sm text-cyan-100/80">GET /api/routes + GET /api/reviews?routeId=...</p>
    </div>
    <button class="button-primary" on:click={loadRoutes}>Refrescar catálogo</button>
  </div>

  {#if loading}
    <p class="animate-pulse text-cyan-100">Sincronizando rutas...</p>
  {:else if error}
    <p class="rounded-xl border border-red-500/40 bg-red-500/10 px-3 py-2 text-red-200">{error}</p>
  {:else if routes.length === 0}
    <p class="text-slate-200">No hay rutas registradas.</p>
  {:else}
    <div class="grid gap-4 lg:grid-cols-[320px,1fr]">
      <div class="space-y-3">
        {#each routes as route}
          <button
            class={`w-full rounded-2xl border px-4 py-3 text-left transition ${
              selectedRoute?.id === route.id
                ? 'border-cyan-400 bg-cyan-500/10 text-white shadow-lg'
                : 'border-white/10 bg-slate-900/70 text-slate-200 hover:border-cyan-300/40'
            }`}
            on:click={() => selectRoute(route)}
          >
            <h3 class="text-lg font-semibold">{route.name}</h3>
            <p class="text-xs text-slate-400 font-mono truncate">{route.id}</p>
            <div class="mt-2 flex gap-2 text-xs text-cyan-200">
              <span class="badge bg-cyan-500/10 border-cyan-400/40 text-cyan-100">{route.distance_km} km</span>
              <span class="badge bg-cyan-500/10 border-cyan-400/40 text-cyan-100">{route.elevation_gain} m D+</span>
            </div>
          </button>
        {/each}
      </div>

      <div class="card border-white/10 bg-slate-900/80 text-slate-100">
        {#if selectedRoute}
          <div class="flex flex-wrap items-center justify-between gap-3">
            <div>
              <p class="text-sm text-cyan-200/80 uppercase tracking-wide">Resumen de ruta</p>
              <h3 class="text-2xl font-semibold">{selectedRoute.name}</h3>
              <p class="text-xs text-slate-400 font-mono">{selectedRoute.id}</p>
            </div>
            <button class="button-primary" on:click={() => loadReviews(selectedRoute!.id)}>Recargar reseñas</button>
          </div>
          <div class="mt-4 grid gap-3 sm:grid-cols-2">
            <div class="rounded-2xl border border-white/10 bg-slate-950/60 p-3 text-center">
              <p class="text-xs text-slate-400">Distancia</p>
              <p class="text-xl font-semibold text-white">{selectedRoute.distance_km} km</p>
            </div>
            <div class="rounded-2xl border border-white/10 bg-slate-950/60 p-3 text-center">
              <p class="text-xs text-slate-400">Desnivel</p>
              <p class="text-xl font-semibold text-white">{selectedRoute.elevation_gain} m</p>
            </div>
          </div>

          <div class="mt-6 space-y-3">
            <div class="flex items-center justify-between">
              <h4 class="text-lg font-semibold text-white">Reseñas</h4>
              {#if reviewsLoading}
                <span class="text-xs text-cyan-200/70">Actualizando…</span>
              {/if}
            </div>
            {#if reviewsLoading}
              <p class="text-cyan-100">Cargando reseñas de esta ruta...</p>
            {:else if reviewsError}
              <p class="text-red-200 text-sm">{reviewsError}</p>
            {:else if reviews.length === 0}
              <p class="text-slate-300 text-sm">No hay reseñas para esta ruta.</p>
            {:else}
              <div class="space-y-3 max-h-[320px] overflow-y-auto pr-1">
                {#each reviews as review}
                  <div class="rounded-2xl border border-white/10 bg-slate-950/70 p-3">
                    <div class="flex items-center justify-between text-xs text-slate-400">
                      <span class="badge bg-orange-500/20 text-orange-200">{review.rating} ★</span>
                      <span>{new Date(review.created_at).toLocaleString()}</span>
                    </div>
                    <p class="mt-2 text-sm text-slate-100">{review.comment}</p>
                    <p class="mt-2 text-[11px] text-slate-400">User: <span class="font-mono">{review.user_id}</span></p>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {:else}
          <p class="text-slate-200">Selecciona una ruta para ver su resumen y reseñas.</p>
        {/if}
      </div>
    </div>
  {/if}
</section>
