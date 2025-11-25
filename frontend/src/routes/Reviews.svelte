<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';

  type Review = { id: string; user_id: string; route_id: string; rating: number; comment: string; created_at: string };
  type Route = { id: string; name: string; distance_km: number; elevation_gain: number };

  let routeId = '';
  let reviews: Review[] = [];
  let routes: Route[] = [];
  let error = '';
  let loading = false;

  let form = {
    userId: '',
    routeId: '',
    rating: 5,
    comment: ''
  };

  let routeLookup: Record<string, Route> = {};

  async function loadRoutes() {
    try {
      const data = await api.listRoutes();
      routes = data.routes || [];
      routeLookup = routes.reduce((acc, r) => {
        acc[r.id] = r;
        return acc;
      }, {} as Record<string, Route>);
    } catch (err) {
      console.error(err);
    }
  }

  async function load() {
    loading = true;
    error = '';
    try {
      const data = await api.getReviews(routeId || undefined);
      reviews = data.reviews || [];
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  async function submit() {
    error = '';
    try {
      await api.createReview({ ...form, rating: Number(form.rating) });
      form.comment = '';
      form.routeId = routeId || '';
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }

  onMount(() => {
    loadRoutes();
    load();
  });
</script>

<section class="space-y-5 rounded-3xl border border-white/10 bg-slate-900/70 p-6 shadow-2xl backdrop-blur">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <p class="badge bg-white/5 text-orange-200">Reseñas</p>
      <h2 class="text-2xl font-semibold text-white">Feedback por ruta</h2>
      <p class="text-sm text-orange-100/80">GET/POST /api/reviews → gRPC Reviews</p>
    </div>
    <button class="button-primary" on:click={load}>Refrescar</button>
  </div>

  <div class="grid gap-4 md:grid-cols-3">
    <div class="card space-y-3 border-white/10 bg-slate-900/80 md:col-span-1">
      <h3 class="text-sm font-semibold text-slate-200">1. Selecciona una ruta</h3>
      <select class="input" bind:value={routeId}>
        <option value="">Todas las rutas</option>
        {#each routes as route}
          <option value={route.id}>{route.name} ({route.distance_km} km)</option>
        {/each}
      </select>
      <button class="button-primary w-full justify-center" on:click={load}>Ver reseñas</button>
    </div>

    <div class="card space-y-3 border-white/10 bg-slate-900/80 md:col-span-2">
      <h3 class="text-sm font-semibold text-slate-200">2. Escribe tu reseña</h3>
      <div class="grid gap-3 md:grid-cols-2">
        <input class="input" placeholder="User ID" bind:value={form.userId} />
        <input class="input" placeholder="Route ID" bind:value={form.routeId} />
      </div>
      <input class="input" type="number" min="1" max="5" placeholder="Rating (1-5)" bind:value={form.rating} />
      <textarea class="input" rows="3" placeholder="Comentario" bind:value={form.comment}></textarea>
      <button class="button-primary w-full justify-center" on:click={submit}>Publicar</button>
    </div>
  </div>

  {#if error}
    <p class="rounded-xl border border-red-500/40 bg-red-500/10 px-3 py-2 text-red-200">{error}</p>
  {/if}

  {#if loading}
    <p class="animate-pulse text-orange-200">Cargando reseñas...</p>
  {:else if reviews.length === 0}
    <p class="text-slate-200">No hay reseñas para esta ruta todavía.</p>
  {:else}
    <div class="grid gap-4 md:grid-cols-2">
      {#each reviews as r}
        <div class="card border-white/10 bg-gradient-to-br from-slate-800/80 to-slate-900/80 text-slate-100">
          <div class="flex items-center justify-between">
            <p class="badge bg-orange-500/20 text-orange-200">{r.rating} ★</p>
            <span class="text-xs text-slate-400">{new Date(r.created_at).toLocaleString()}</span>
          </div>
          <p class="mt-3 text-sm">{r.comment}</p>
          <p class="mt-3 text-xs text-slate-400">User: <span class="font-mono">{r.user_id}</span></p>
          <p class="text-xs text-slate-400">
            Route:
            <span class="font-mono">
              {routeLookup[r.route_id]?.name || r.route_id}
            </span>
          </p>
        </div>
      {/each}
    </div>
  {/if}
</section>
