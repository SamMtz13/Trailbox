<script lang="ts">
  import { api } from '../lib/api';

  type Notification = { id: string; user_id: string; message: string; read: boolean; created_at: string };

  let userId = '';
  let notifications: Notification[] = [];
  let error = '';
  let form = { userId: '', message: '' };

  async function load() {
    error = '';
    try {
      if (!userId) {
        notifications = [];
        return;
      }
      const data = await api.getNotifications(userId);
      notifications = data.notifications || [];
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }

  async function submit() {
    error = '';
    try {
      await api.sendNotification(form);
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }
</script>

<section class="card space-y-4">
  <div class="flex items-center justify-between gap-3 flex-wrap">
    <div>
      <p class="badge">Notificaciones</p>
      <h2 class="text-xl font-semibold text-forest">Mensajes por usuario</h2>
      <p class="text-sm text-emerald-800">GET/POST /api/notifications → gRPC Notifications</p>
    </div>
    <button class="button-ghost" on:click={load}>Refrescar</button>
  </div>

  <div class="grid md:grid-cols-2 gap-4">
    <div class="card space-y-3">
      <label class="text-sm text-emerald-800 font-semibold" for="notif-user-filter">Consultar por usuario</label>
      <input id="notif-user-filter" class="input" placeholder="User ID" bind:value={userId} />
      <button class="button-primary" on:click={load}>Buscar</button>
    </div>
    <div class="card space-y-3">
      <label class="text-sm text-emerald-800 font-semibold" for="notif-user">Enviar notificación</label>
      <input id="notif-user" class="input" placeholder="User ID" bind:value={form.userId} />
      <textarea id="notif-message" class="input" rows="3" placeholder="Mensaje" bind:value={form.message}></textarea>
      <button class="button-primary" on:click={submit}>Enviar</button>
    </div>
  </div>

  {#if error}
    <p class="text-red-700">{error}</p>
  {/if}

  <div class="space-y-3">
    {#each notifications as n}
      <div class="card">
        <div class="flex items-center justify-between gap-3">
          <span class="badge">{n.read ? 'LEÍDO' : 'NUEVO'}</span>
          <span class="text-xs text-emerald-700">{new Date(n.created_at).toLocaleString()}</span>
        </div>
        <p class="text-sm text-emerald-900 mt-2">{n.message}</p>
        <p class="text-xs text-emerald-700 mt-2">User: <span class="font-mono">{n.user_id}</span></p>
      </div>
    {/each}
  </div>
</section>
