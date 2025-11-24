<script lang="ts">
  import Card from '../lib/components/Card.svelte';
  import LoadingState from '../lib/components/LoadingState.svelte';
  import type { MapResponse } from '../lib/api';
  import { getRouteMap, saveRouteMap } from '../lib/api';

  let lookupRouteId = '';
  let currentMap: MapResponse | null = null;
  let loading = false;
  let error = '';
  let form = { routeId: '', geoJson: '{\n  "type": "LineString",\n  "coordinates": []\n}' };
  let message = '';

  async function fetchMap(routeId: string) {
    if (!routeId) return;
    loading = true;
    error = '';
    message = '';
    try {
      currentMap = await getRouteMap(routeId);
      form.routeId = form.routeId || routeId;
    } catch (err) {
      error = err instanceof Error ? err.message : 'Map not found';
      currentMap = null;
    } finally {
      loading = false;
    }
  }

  async function submitMap(event: SubmitEvent) {
    event.preventDefault();
    if (!form.routeId || !form.geoJson) {
      error = 'Route ID and GeoJSON are required';
      return;
    }
    loading = true;
    error = '';
    message = '';
    try {
      await saveRouteMap(form);
      message = 'Map stored!';
      await fetchMap(form.routeId);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to store map';
    } finally {
      loading = false;
    }
  }
</script>

<div class="grid gap-6 md:grid-cols-2">
  <Card title="Route map lookup" description="Reads GeoJSON blobs from the map service">
    <form
      class="flex flex-col gap-3 md:flex-row"
      on:submit|preventDefault={() => fetchMap(lookupRouteId.trim())}
    >
      <input
        class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
        placeholder="Route ID"
        bind:value={lookupRouteId}
      />
      <button
        type="submit"
        class="rounded-lg bg-slate-900 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-slate-800"
      >
        Load map
      </button>
    </form>

    {#if loading}
      <div class="mt-4">
        <LoadingState message="Fetching GeoJSON..." />
      </div>
    {:else if error}
      <p class="mt-4 text-sm text-red-600">{error}</p>
    {:else if currentMap}
      <div class="mt-4">
        <p class="text-xs font-semibold uppercase text-slate-500">Route ID</p>
        <p class="font-mono text-xs text-slate-400">{currentMap.route_id}</p>
        <p class="mt-2 text-xs uppercase text-slate-500">GeoJSON</p>
        <pre class="mt-1 overflow-auto rounded-lg bg-slate-900 p-3 text-xs text-emerald-200">
{currentMap.geo_json}
        </pre>
      </div>
    {/if}
  </Card>

  <Card title="Store GeoJSON">
    <form class="space-y-3" on:submit={submitMap}>
      <input
        class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
        placeholder="Route ID"
        bind:value={form.routeId}
      />
      <textarea
        rows="6"
        class="w-full rounded-lg border border-slate-200 font-mono text-xs text-slate-700"
        placeholder={'{"type": "LineString", "coordinates": []}'}
        bind:value={form.geoJson}
      ></textarea>
      <button
        type="submit"
        class="w-full rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-emerald-500"
      >
        Save
      </button>
      {#if message}
        <p class="text-sm text-emerald-600">{message}</p>
      {/if}
    </form>
  </Card>
</div>
