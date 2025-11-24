<script lang="ts">
  import { onMount } from 'svelte';
  import Card from '../lib/components/Card.svelte';
  import LoadingState from '../lib/components/LoadingState.svelte';
  import type { Route } from '../lib/api';
  import { listRoutes } from '../lib/api';

  let routes: Route[] = [];
  let loading = false;
  let error = '';

  async function loadRoutes() {
    loading = true;
    error = '';
    try {
      routes = await listRoutes();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to load routes';
    } finally {
      loading = false;
    }
  }

  onMount(loadRoutes);
</script>

<Card title="Routes" description="Registered segments provided by the routes service">
  {#if loading}
    <LoadingState message="Fetching routes..." />
  {:else if error}
    <p class="text-sm text-red-600">{error}</p>
  {:else if routes.length === 0}
    <p class="text-sm text-slate-500">There are no routes yet.</p>
  {:else}
    <div class="grid gap-3">
      {#each routes as route}
        <article class="rounded-lg border border-slate-100 bg-slate-50/70 p-4">
          <p class="text-xs font-mono uppercase text-slate-400">ID {route.id}</p>
          <h3 class="text-lg font-semibold text-slate-800">{route.name}</h3>
          <div class="mt-2 flex flex-wrap gap-4 text-sm text-slate-600">
            <span>
              Distance:
              <strong>{route.distance_km ?? 0}</strong>
              km
            </span>
            <span>
              Elevation gain:
              <strong>{route.elevation_gain ?? 0}</strong>
              m
            </span>
          </div>
        </article>
      {/each}
    </div>
  {/if}
</Card>
