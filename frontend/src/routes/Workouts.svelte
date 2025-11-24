<script lang="ts">
  import { onMount } from 'svelte';
  import Card from '../lib/components/Card.svelte';
  import LoadingState from '../lib/components/LoadingState.svelte';
  import type { Workout } from '../lib/api';
  import { listWorkouts } from '../lib/api';

  let workouts: Workout[] = [];
  let loading = false;
  let error = '';

  function formatDate(value: string) {
    const date = new Date(value);
    return isNaN(date.getTime()) ? value : date.toLocaleString();
  }

  async function loadWorkouts() {
    loading = true;
    error = '';
    try {
      workouts = await listWorkouts();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to load workouts';
    } finally {
      loading = false;
    }
  }

  onMount(loadWorkouts);
</script>

<Card title="Workouts" description="Latest sessions computed by the workouts service">
  {#if loading}
    <LoadingState message="Fetching workouts..." />
  {:else if error}
    <p class="text-sm text-red-600">{error}</p>
  {:else if workouts.length === 0}
    <p class="text-sm text-slate-500">No workouts recorded yet.</p>
  {:else}
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-slate-200 text-sm">
        <thead class="bg-slate-50 text-left text-xs font-semibold uppercase text-slate-500">
          <tr>
            <th class="px-3 py-2">Workout</th>
            <th class="px-3 py-2">User</th>
            <th class="px-3 py-2">Route</th>
            <th class="px-3 py-2">Duration (min)</th>
            <th class="px-3 py-2">Calories</th>
            <th class="px-3 py-2">Date</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          {#each workouts as workout}
            <tr class="hover:bg-slate-50">
              <td class="px-3 py-2 font-mono text-xs text-slate-400">{workout.id}</td>
              <td class="px-3 py-2 font-mono text-xs text-slate-400">{workout.user_id}</td>
              <td class="px-3 py-2 font-mono text-xs text-slate-400">{workout.route_id}</td>
              <td class="px-3 py-2">{workout.duration}</td>
              <td class="px-3 py-2">{workout.calories}</td>
              <td class="px-3 py-2">{formatDate(workout.date)}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>
