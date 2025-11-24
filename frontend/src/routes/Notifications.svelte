<script lang="ts">
  import Card from '../lib/components/Card.svelte';
  import LoadingState from '../lib/components/LoadingState.svelte';
  import type { Notification } from '../lib/api';
  import { getNotifications, sendNotification } from '../lib/api';

  let userId = '';
  let notifications: Notification[] = [];
  let loading = false;
  let error = '';
  let form = { userId: '', message: '' };
  let feedback = '';

  async function loadNotifications(target: string) {
    if (!target) return;
    loading = true;
    error = '';
    feedback = '';
    try {
      notifications = await getNotifications(target);
      form.userId = form.userId || target;
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to load notifications';
    } finally {
      loading = false;
    }
  }

  async function submitNotification(event: SubmitEvent) {
    event.preventDefault();
    if (!form.userId || !form.message) {
      error = 'Both fields are required';
      return;
    }
    loading = true;
    error = '';
    feedback = '';
    try {
      await sendNotification(form);
      feedback = 'Notification sent!';
      form.message = '';
      await loadNotifications(form.userId);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to send notification';
    } finally {
      loading = false;
    }
  }
</script>

<div class="grid gap-6 md:grid-cols-2">
  <Card title="Notifications" description="Pull user inbox from the notifications service">
    <form
      class="flex flex-col gap-3 md:flex-row"
      on:submit|preventDefault={() => loadNotifications(userId.trim())}
    >
      <input
        class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
        placeholder="User ID"
        bind:value={userId}
      />
      <button
        type="submit"
        class="rounded-lg bg-slate-900 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-slate-800"
      >
        Load inbox
      </button>
    </form>

    {#if loading}
      <div class="mt-4">
        <LoadingState message="Fetching notifications..." />
      </div>
    {:else if error}
      <p class="mt-4 text-sm text-red-600">{error}</p>
    {:else if notifications.length === 0 && userId}
      <p class="mt-4 text-sm text-slate-500">Inbox is empty.</p>
    {:else if notifications.length > 0}
      <ul class="mt-4 space-y-2">
        {#each notifications as notification}
          <li class="rounded-lg border border-slate-100 bg-slate-50/70 p-3">
            <p class="text-sm text-slate-700">{notification.message}</p>
            <p class="text-xs text-slate-400">{notification.created_at}</p>
          </li>
        {/each}
      </ul>
    {/if}
  </Card>

  <Card title="Send notification">
    <form class="space-y-3" on:submit={submitNotification}>
      <input
        class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
        placeholder="User ID"
        bind:value={form.userId}
      />
      <textarea
        rows="3"
        class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
        placeholder="Message"
        bind:value={form.message}
      ></textarea>
      <button
        type="submit"
        class="w-full rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-emerald-500"
      >
        Send
      </button>
      {#if feedback}
        <p class="text-sm text-emerald-600">{feedback}</p>
      {/if}
    </form>
  </Card>
</div>
