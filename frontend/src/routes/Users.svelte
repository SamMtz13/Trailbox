<script lang="ts">
  import { onMount } from 'svelte';
  import Card from '../lib/components/Card.svelte';
  import LoadingState from '../lib/components/LoadingState.svelte';
  import type { User } from '../lib/api';
  import { listUsers } from '../lib/api';

  let users: User[] = [];
  let loading = false;
  let error = '';

  async function loadUsers() {
    loading = true;
    error = '';
    try {
      users = await listUsers();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to load users';
    } finally {
      loading = false;
    }
  }

  onMount(loadUsers);
</script>

<Card title="Users" description="Data coming straight from the users microservice">
  {#if loading}
    <LoadingState message="Fetching users..." />
  {:else if error}
    <p class="text-sm text-red-600">{error}</p>
  {:else if users.length === 0}
    <p class="text-sm text-slate-500">No users registered yet.</p>
  {:else}
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-slate-200 text-sm">
        <thead class="bg-slate-50 text-left text-xs font-semibold uppercase text-slate-500">
          <tr>
            <th class="px-3 py-2">ID</th>
            <th class="px-3 py-2">Name</th>
            <th class="px-3 py-2">Email</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          {#each users as user}
            <tr class="hover:bg-slate-50">
              <td class="px-3 py-2 font-mono text-xs text-slate-400">{user.id}</td>
              <td class="px-3 py-2">{user.name}</td>
              <td class="px-3 py-2">{user.email}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>
