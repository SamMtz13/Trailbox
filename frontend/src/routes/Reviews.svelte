<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';

  type Review = { id: string; user_id: string; route_id: string; rating: number; comment: string; created_at: string };

  let routeId = '';
  let reviews: Review[] = [];
  let error = '';
  let loading = false;

  let form = {
    userId: '',
    routeId: '',
    rating: 5,
    comment: ''
  };

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
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }

  onMount(() => {
    load();
  });
</script>

<section class="card space-y-4">
  <div class="flex items-center justify-between gap-3 flex-wrap">
    <div>
      <p class="badge">Reseñas</p>
      <h2 class="text-xl font-semibold text-forest">Feedback por ruta</h2>
      <p class="text-sm text-emerald-800">GET/POST /api/reviews → gRPC Reviews</p>
    </div>
    <button class="button-ghost" on:click={load}>Refrescar</button>
  </div>

  <div class="grid md:grid-cols-2 gap-4">
    <div class="card space-y-3">
      <label class="text-sm text-emerald-800 font-semibold" for="route-filter">Filtrar por Route ID</label>
      <input id="route-filter" class="input" placeholder="route-id" bind:value={routeId} />
      <button class="button-primary" on:click={load}>Buscar reseñas</button>
    </div>
    <div class="card space-y-3">
      <label class="text-sm text-emerald-800 font-semibold" for="review-user">Crear reseña</label>
      <input id="review-user" class="input" placeholder="User ID" bind:value={form.userId} />
      <input id="review-route" class="input" placeholder="Route ID" bind:value={form.routeId} />
      <input id="review-rating" class="input" type="number" min="1" max="5" bind:value={form.rating} />
      <textarea id="review-comment" class="input" rows="3" placeholder="Comentario" bind:value={form.comment}></textarea>
      <button class="button-primary" on:click={submit}>Enviar</button>
    </div>
  </div>

  {#if loading}
    <p class="text-emerald-700">Cargando reseñas...</p>
  {:else if error}
    <p class="text-red-700">{error}</p>
  {:else if reviews.length === 0}
    <p class="text-emerald-700">No hay reseñas disponibles.</p>
  {:else}
    <div class="grid md:grid-cols-2 gap-3">
      {#each reviews as r}
        <div class="card">
          <div class="flex items-center justify-between gap-3">
            <p class="badge">{r.rating} ★</p>
            <span class="text-xs text-emerald-700">{new Date(r.created_at).toLocaleString()}</span>
          </div>
          <p class="text-sm text-emerald-900 mt-2">{r.comment}</p>
          <p class="text-xs text-emerald-700 mt-2">User: <span class="font-mono">{r.user_id}</span></p>
          <p class="text-xs text-emerald-700">Route: <span class="font-mono">{r.route_id}</span></p>
        </div>
      {/each}
    </div>
  {/if}
</section>
