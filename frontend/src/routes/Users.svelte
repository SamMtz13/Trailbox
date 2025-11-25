<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';

  type User = { id: string; name: string; email: string };

  let users: User[] = [];
  let loading = true;
  let error = '';

  async function load() {
    loading = true;
    error = '';
    try {
      const data = await api.listUsers();
      users = data.users || [];
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    load();
  });
</script>

<section class="space-y-4 rounded-3xl border border-white/10 bg-slate-900/70 p-6 shadow-2xl backdrop-blur">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <p class="badge bg-white/5 text-emerald-200">Usuarios</p>
      <h2 class="text-2xl font-semibold text-white">Catálogo</h2>
      <p class="text-sm text-emerald-200/80">GET /api/users → gRPC Users.ListUsers</p>
    </div>
    <button class="button-primary" on:click={load}>Refrescar datos</button>
  </div>

  {#if loading}
    <p class="animate-pulse text-emerald-200">Sincronizando usuarios...</p>
  {:else if error}
    <p class="text-red-400 font-semibold">{error}</p>
  {:else if users.length === 0}
    <p class="text-emerald-100">No hay usuarios registrados.</p>
  {:else}
    <div class="overflow-hidden rounded-2xl border border-white/5 bg-slate-950/40">
      <table class="min-w-full text-sm text-white/90">
        <thead class="bg-white/5 text-left text-emerald-200 uppercase tracking-wide text-xs">
          <tr>
            <th class="px-4 py-3">ID</th>
            <th class="px-4 py-3">Nombre</th>
            <th class="px-4 py-3">Email</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/80">
          {#each users as user}
            <tr class="hover:bg-emerald-500/5 transition">
              <td class="px-4 py-3 font-mono text-xs text-emerald-200">{user.id}</td>
              <td class="px-4 py-3">{user.name}</td>
              <td class="px-4 py-3 text-emerald-300">{user.email}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</section>
