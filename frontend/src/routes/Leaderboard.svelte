<script lang="ts">
  import { onMount } from 'svelte';
  import Card from '../lib/components/Card.svelte';
  import LoadingState from '../lib/components/LoadingState.svelte';
  import type { LeaderboardEntry } from '../lib/api';
  import { getLeaderboard, upsertLeaderboard } from '../lib/api';

  let entries: LeaderboardEntry[] = [];
  let loading = false;
  let error = '';
  let form = { userId: '', score: 100 };
  let message = '';

  async function loadLeaderboard() {
    loading = true;
    error = '';
    message = '';
    try {
      entries = await getLeaderboard();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to load leaderboard';
    } finally {
      loading = false;
    }
  }

  async function submitScore(event: SubmitEvent) {
    event.preventDefault();
    if (!form.userId) {
      error = 'User ID is required';
      return;
    }
    loading = true;
    error = '';
    message = '';
    try {
      await upsertLeaderboard(form);
      message = 'Score updated!';
      await loadLeaderboard();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to update score';
    } finally {
      loading = false;
    }
  }

  onMount(loadLeaderboard);
</script>

<div class="grid gap-6 md:grid-cols-2">
  <Card title="Leaderboard" description="Live ranking served by the leaderboard service">
    {#if loading}
      <LoadingState message="Fetching ranking..." />
    {:else if error}
      <p class="text-sm text-red-600">{error}</p>
    {:else if entries.length === 0}
      <p class="text-sm text-slate-500">No scores yet.</p>
    {:else}
      <ol class="space-y-2">
        {#each entries as entry}
          <li class="flex items-center justify-between rounded-lg border border-slate-100 bg-slate-50/70 p-3">
            <div>
              <p class="text-sm font-semibold text-slate-900">#{entry.position} — {entry.score} pts</p>
              <p class="font-mono text-xs text-slate-500">{entry.user_id}</p>
            </div>
            <span class="text-xs uppercase text-slate-400">ID {entry.id}</span>
          </li>
        {/each}
      </ol>
    {/if}
  </Card>

  <Card title="Manual update" description="Push new scores through the gateway">
    <form class="space-y-3" on:submit={submitScore}>
      <input
        class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
        placeholder="User ID"
        bind:value={form.userId}
      />
      <input
        type="number"
        min="0"
        class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
        placeholder="Score"
        bind:value={form.score}
      />
      <button
        type="submit"
        class="w-full rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-emerald-500"
      >
        Upsert
      </button>
      {#if message}
        <p class="text-sm text-emerald-600">{message}</p>
      {/if}
    </form>
  </Card>
</div>
