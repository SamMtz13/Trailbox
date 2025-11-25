<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';

  type Entry = { id: string; user_id: string; score: number; position: number };

  let entries: Entry[] = [];
  let error = '';
  let limit = 5;
  let form = { userId: '', score: 0 };

  async function load() {
    error = '';
    try {
      const data = await api.getLeaderboard(limit);
      entries = data.entries || [];
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }

  async function submit() {
    try {
      await api.upsertScore({ userId: form.userId, score: Number(form.score) });
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
      <p class="badge">Leaderboard</p>
      <h2 class="text-xl font-semibold text-forest">Top atletas</h2>
      <p class="text-sm text-emerald-800">GET/POST /api/leaderboard → gRPC Leaderboard</p>
    </div>
    <button class="button-ghost" on:click={load}>Refrescar</button>
  </div>

  <div class="grid md:grid-cols-2 gap-4">
    <div class="card space-y-3">
      <label class="text-sm text-emerald-800 font-semibold" for="leaderboard-limit">Límite</label>
      <input id="leaderboard-limit" class="input" type="number" min="1" max="50" bind:value={limit} />
      <button class="button-primary" on:click={load}>Actualizar</button>
    </div>
    <div class="card space-y-3">
      <label class="text-sm text-emerald-800 font-semibold" for="leaderboard-user">Registrar puntaje</label>
      <input id="leaderboard-user" class="input" placeholder="User ID" bind:value={form.userId} />
      <input id="leaderboard-score" class="input" type="number" placeholder="Score" bind:value={form.score} />
      <button class="button-primary" on:click={submit}>Enviar</button>
    </div>
  </div>

  {#if error}
    <p class="text-red-700">{error}</p>
  {/if}

  <div class="space-y-3">
    {#each entries as e}
      <div class="card flex items-center justify-between gap-3">
        <div>
          <p class="text-sm text-emerald-800">Usuario <span class="font-mono text-xs">{e.user_id}</span></p>
          <p class="text-sm text-emerald-900 font-semibold">Score: {e.score}</p>
        </div>
        <div class="flex items-center gap-2 text-forest font-bold text-lg">
          <span class="badge">#{e.position}</span>
        </div>
      </div>
    {/each}
  </div>
</section>
