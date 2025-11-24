<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';

  type Workout = { id: string; user_id: string; route_id: string; date: string; duration: number; calories: number };

  let workouts: Workout[] = [];
  let loading = true;
  let error = '';

  async function load() {
    loading = true;
    error = '';
    try {
      const data = await api.listWorkouts();
      workouts = data.workouts || [];
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
      <p class="badge">Entrenamientos</p>
      <h2 class="text-xl font-semibold text-forest">Historial</h2>
      <p class="text-sm text-emerald-800">GET /api/workouts → gRPC Workouts.ListWorkouts</p>
    </div>
    <button class="button-ghost" on:click={load}>Refrescar</button>
  </div>

  {#if loading}
    <p class="text-emerald-700">Cargando entrenamientos...</p>
  {:else if error}
    <p class="text-red-700">{error}</p>
  {:else if workouts.length === 0}
    <p class="text-emerald-700">No hay entrenamientos aún.</p>
  {:else}
    <div class="space-y-3">
      {#each workouts as w}
        <div class="card">
          <div class="flex items-center justify-between gap-3 flex-wrap">
            <div>
              <p class="text-sm text-emerald-800">Usuario: <span class="font-mono text-xs">{w.user_id}</span></p>
              <p class="text-sm text-emerald-800">Ruta: <span class="font-mono text-xs">{w.route_id}</span></p>
            </div>
            <span class="badge">{new Date(w.date).toLocaleString()}</span>
          </div>
          <div class="flex gap-4 text-sm text-emerald-900 mt-2">
            <span class="badge">{w.duration} min</span>
            <span class="badge">{w.calories} cal</span>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</section>
