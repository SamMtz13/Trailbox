<script lang="ts">
  import { api } from '../lib/api';

  type MapData = { route_id: string; geo_json: string; created_at: string };

  let routeId = '';
  let mapData: MapData | null = null;
  let error = '';
  let form = { routeId: '', geoJson: '' };

  async function load() {
    error = '';
    try {
      if (!routeId) {
        mapData = null;
        return;
      }
      mapData = await api.getMap(routeId);
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }

  async function submit() {
    error = '';
    try {
      await api.setMap(form);
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }
</script>

<section class="card space-y-4">
  <div class="flex items-center justify-between gap-3 flex-wrap">
    <div>
      <p class="badge">Mapas</p>
      <h2 class="text-xl font-semibold text-forest">GeoJSON por ruta</h2>
      <p class="text-sm text-emerald-800">GET/POST /api/maps → gRPC Map</p>
    </div>
    <button class="button-ghost" on:click={load}>Refrescar</button>
  </div>

  <div class="grid md:grid-cols-2 gap-4">
    <div class="card space-y-3">
      <label class="text-sm text-emerald-800 font-semibold" for="map-route-search">Consultar ruta</label>
      <input id="map-route-search" class="input" placeholder="Route ID" bind:value={routeId} />
      <button class="button-primary" on:click={load}>Buscar</button>
      {#if mapData}
        <p class="text-xs text-emerald-700">Última actualización: {new Date(mapData.created_at).toLocaleString()}</p>
      {/if}
    </div>
    <div class="card space-y-3">
      <label class="text-sm text-emerald-800 font-semibold" for="map-route-save">Guardar GeoJSON</label>
      <input id="map-route-save" class="input" placeholder="Route ID" bind:value={form.routeId} />
      <textarea
        id="map-geojson"
        class="input"
        rows="4"
        placeholder="GeoJSON, por ejemplo: type=Feature, geometry=..."
        bind:value={form.geoJson}
      ></textarea>
      <button class="button-primary" on:click={submit}>Guardar</button>
    </div>
  </div>

  {#if error}
    <p class="text-red-700">{error}</p>
  {/if}

  {#if mapData}
    <pre class="bg-mist border border-emerald-100 rounded-xl p-3 text-xs text-forest overflow-auto">{JSON.stringify(mapData, null, 2)}</pre>
  {/if}
</section>
