<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { api } from '../lib/api';
  import L, { Map as LeafletMap, LayerGroup, GeoJSON } from 'leaflet';
  import type { Feature } from 'geojson';

  type Route = {
    id: string;
    name: string;
    distance_km: number;
    elevation_gain: number;
  };

  type RouteMap = {
    route: Route;
    feature: Feature;
    color: string;
    created_at: string;
  };

  const colors = ['#22c55e', '#0ea5e9', '#f97316'];

  let mapContainer: HTMLDivElement;

  let map: LeafletMap | null = null;
  let layerGroup: LayerGroup | null = null;
  let layers: Record<string, GeoJSON> = {};

  let routeMaps: RouteMap[] = [];
  let loading = false;
  let error = '';
  let selectedRouteId = '';

  // Ruta seleccionada (tarjeta inferior)
  $: activeRoute = routeMaps.find((r) => r.route.id === selectedRouteId) ?? routeMaps[0];

  async function loadRoutesWithMaps() {
    loading = true;
    error = '';
    try {
      const routesResp = await api.listRoutes();
      const routes: Route[] = routesResp.routes || [];
      const aggregated: RouteMap[] = [];

      let idx = 0;
      for (const route of routes) {
        const mapResp = await api.getMap(route.id);

        if (mapResp?.geo_json) {
          const feature: Feature =
            typeof mapResp.geo_json === 'string'
              ? JSON.parse(mapResp.geo_json)
              : mapResp.geo_json;

          aggregated.push({
            route,
            feature,
            color: colors[idx % colors.length],
            created_at: mapResp.created_at ?? new Date().toISOString()
          });
          idx++;
        }
      }

      routeMaps = aggregated;

      if (routeMaps.length > 0) {
        if (!selectedRouteId || !routeMaps.find((r) => r.route.id === selectedRouteId)) {
          selectedRouteId = routeMaps[0].route.id;
        }
        renderLayers();
      } else {
        // si no hay rutas, limpiamos capas
        if (layerGroup) {
          layerGroup.clearLayers();
          layers = {};
        }
      }
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  function renderLayers() {
    if (!map || !layerGroup) return;

    layerGroup.clearLayers();
    layers = {};

    routeMaps.forEach((item) => {
      const layer = L.geoJSON(item.feature, {
        style: () => ({
          color: item.color,
          weight: selectedRouteId === item.route.id ? 6 : 4,
          opacity: 0.9
        })
      });

      layer.bindPopup(
        `<strong>${item.route.name}</strong><br/>${item.route.distance_km} km · ${item.route.elevation_gain} m D+`
      );

      layer.addTo(layerGroup);
      layers[item.route.id] = layer;
    });

    // 🔹 Calcular bounds sin usar layerGroup.getBounds()
    const allLayers = Object.values(layers);
    if (allLayers.length > 0) {
      const group = L.featureGroup(allLayers as any);
      const bounds = group.getBounds();
      if (bounds.isValid()) {
        map.fitBounds(bounds, { padding: [40, 40] });
      }
    }
  }

  function focusRoute(routeId: string) {
    selectedRouteId = routeId;
    renderLayers();

    const layer = layers[routeId];
    if (layer) {
      const bounds = layer.getBounds();
      if (bounds.isValid()) {
        map!.fitBounds(bounds, { padding: [60, 60] });
      }
      layer.openPopup();
    }
  }

  onMount(() => {
    if (!mapContainer) return;

    console.log('📐 Map container size:', mapContainer.clientWidth, mapContainer.clientHeight);

    map = L.map(mapContainer, {
      zoomControl: true,
      scrollWheelZoom: false,
      attributionControl: true,
      center: [19.4, -99.15],
      zoom: 11
    });

    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);

    layerGroup = L.layerGroup().addTo(map);

    setTimeout(() => {
      if (map) {
        map.invalidateSize();
        console.log('🗺 Leaflet map size:', map.getSize());
      }
    }, 0);

    loadRoutesWithMaps();
  });

  onDestroy(() => {
    if (map) map.remove();
    map = null;
    layerGroup = null;
    layers = {};
  });
</script>

<section class="space-y-5 rounded-3xl border border-white/10 bg-slate-900/70 p-6 shadow-2xl backdrop-blur">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <p class="badge bg-emerald-500/10 text-emerald-200">Mapas</p>
      <h2 class="text-2xl font-semibold text-white">Tracks GeoJSON</h2>
      <p class="text-sm text-emerald-100/80">Rutas reales almacenadas en PostgreSQL</p>
    </div>
    <button class="button-primary" on:click={loadRoutesWithMaps}>Recargar</button>
  </div>

  {#if error}
    <p class="rounded-xl border border-red-500/40 bg-red-500/10 px-3 py-2 text-red-200">
      {error}
    </p>
  {/if}

  {#if loading}
    <p class="text-emerald-100 text-sm">Cargando rutas y GeoJSON...</p>
  {/if}

  <div class="grid gap-4 lg:grid-cols-[320px,1fr] mt-2">
    <!-- Panel lateral -->
    <div class="space-y-3">
      {#if routeMaps.length === 0 && !loading}
        <p class="text-slate-200 text-sm">No hay tracks registrados todavía.</p>
      {/if}

      {#each routeMaps as item}
        <button
          class={`w-full rounded-2xl border px-4 py-3 text-left transition ${
            selectedRouteId === item.route.id
              ? 'border-emerald-400 bg-emerald-500/10 text-white shadow-lg'
              : 'border-white/10 bg-slate-900/70 text-slate-200 hover:border-emerald-300/40'
          }`}
          on:click={() => focusRoute(item.route.id)}
        >
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">{item.route.name}</h3>
            <span class="h-3 w-3 rounded-full" style={`background:${item.color}`}></span>
          </div>
          <p class="text-xs text-slate-400 font-mono truncate">{item.route.id}</p>
          <p class="mt-2 text-xs text-slate-300">
            {item.route.distance_km} km · {item.route.elevation_gain} m D+
          </p>
        </button>
      {/each}
    </div>

    <!-- MAPA + tarjeta -->
    <div class="space-y-3">
      <div class="relative h-[420px] overflow-hidden rounded-3xl border border-white/10 bg-slate-950/60">
        <div bind:this={mapContainer} class="absolute inset-0"></div>
      </div>

      {#if activeRoute}
        <div class="rounded-2xl border border-white/10 bg-slate-950/70 p-4 text-slate-100">
          <div class="flex flex-wrap items-center justify-between gap-3">
            <div>
              <p class="text-xs uppercase tracking-wide text-emerald-200/70">
                Ruta seleccionada
              </p>
              <h3 class="text-xl font-semibold">{activeRoute.route.name}</h3>
            </div>
            <span class="badge bg-white/5 text-slate-100">
              Actualizada: {new Date(activeRoute.created_at).toLocaleDateString()}
            </span>
          </div>
          <p class="mt-3 text-xs text-slate-400 font-mono">{activeRoute.route.id}</p>
          <div class="mt-4 grid gap-3 sm:grid-cols-2">
            <div class="rounded-2xl border border-white/10 bg-slate-900/70 p-3 text-center">
              <p class="text-xs text-slate-400">Distancia</p>
              <p class="text-xl font-semibold text-white">
                {activeRoute.route.distance_km} km
              </p>
            </div>
            <div class="rounded-2xl border border-white/10 bg-slate-900/70 p-3 text-center">
              <p class="text-xs text-slate-400">Desnivel</p>
              <p class="text-xl font-semibold text-white">
                {activeRoute.route.elevation_gain} m
              </p>
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>
</section>
