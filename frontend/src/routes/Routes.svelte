<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';

  type Route = { id: string; name: string; distance_km: number; elevation_gain: number };

  let routes: Route[] = [];
  let loading = true;
  let error = '';

  async function load() {
    loading = true;
    error = '';
    try {
      const data = await api.listRoutes();
      routes = data.routes || [];
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  onMount(load);
</script>

<section class="card space-y-4">
  <div class="flex items-center justify-between gap-3 flex-wrap">
    <div>
      <p class="badge">Rutas</p>
      <h2 class="text-xl font-semibold text-forest">Catálogo</h2>
      <p class="text-sm text-emerald-800">GET /api/routes → gRPC Routes.ListRoutes</p>
    </div>
    <button class="button-ghost" on:click={load}>Refrescar</button>
  </div>

  {#if loading}
    <p class="text-emerald-700">Cargando rutas...</p>
  {:else if error}
    <p class="text-red-700">{error}</p>
  {:else if routes.length === 0}
    <p class="text-emerald-700">Sin rutas registradas.</p>
  {:else}
    <div class="grid md:grid-cols-2 gap-4">
      {#each routes as route}
        <div class="card">
          <h3 class="text-lg font-semibold text-forest">{route.name}</h3>
          <p class="text-sm text-emerald-800">ID: <span class="font-mono text-xs">{route.id}</span></p>
          <div class="flex gap-3 text-sm text-emerald-900 mt-2">
            <span class="badge">{route.distance_km} km</span>
            <span class="badge">{route.elevation_gain} m D+</span>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</section>
