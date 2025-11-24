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

  onMount(load);
</script>

<section class="card space-y-4">
  <div class="flex items-center justify-between gap-3 flex-wrap">
    <div>
      <p class="badge">Usuarios</p>
      <h2 class="text-xl font-semibold text-forest">Listado</h2>
      <p class="text-sm text-emerald-800">GET /api/users → gRPC Users.ListUsers</p>
    </div>
    <button class="button-ghost" on:click={load}>Refrescar</button>
  </div>

  {#if loading}
    <p class="text-emerald-700">Cargando usuarios...</p>
  {:else if error}
    <p class="text-red-700">{error}</p>
  {:else if users.length === 0}
    <p class="text-emerald-700">No hay usuarios registrados.</p>
  {:else}
    <div class="overflow-x-auto">
      <table class="min-w-full text-sm">
        <thead class="text-left text-emerald-700">
          <tr>
            <th class="py-2 pr-4">ID</th>
            <th class="py-2 pr-4">Nombre</th>
            <th class="py-2 pr-4">Email</th>
          </tr>
        </thead>
        <tbody class="text-forest divide-y divide-emerald-50">
          {#each users as user}
            <tr>
              <td class="py-2 pr-4 font-mono text-xs">{user.id}</td>
              <td class="py-2 pr-4">{user.name}</td>
              <td class="py-2 pr-4">{user.email}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</section>
