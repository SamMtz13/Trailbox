<script lang="ts">
  import Card from '../lib/components/Card.svelte';
  import LoadingState from '../lib/components/LoadingState.svelte';
  import type { Review } from '../lib/api';
  import { createReview, getReviews } from '../lib/api';

  let routeId = '';
  let reviews: Review[] = [];
  let loading = false;
  let error = '';
  let successMessage = '';

  let form = {
    userId: '',
    routeId: '',
    rating: 5,
    comment: ''
  };

  async function loadReviews(route: string) {
    if (!route) return;
    loading = true;
    error = '';
    successMessage = '';
    try {
      reviews = await getReviews(route);
      form.routeId = form.routeId || route;
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to load reviews';
    } finally {
      loading = false;
    }
  }

  async function submitReview(event: SubmitEvent) {
    event.preventDefault();
    if (!form.userId || !form.routeId || !form.comment) {
      error = 'All fields are required';
      return;
    }
    loading = true;
    error = '';
    successMessage = '';
    try {
      await createReview(form);
      successMessage = 'Review created!';
      form.comment = '';
      await loadReviews(form.routeId);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to create review';
    } finally {
      loading = false;
    }
  }
</script>

<div class="grid gap-6">
  <Card title="Reviews" description="Read and create comments for a specific route">
    <form
      class="flex flex-col gap-3 md:flex-row"
      on:submit|preventDefault={() => loadReviews(routeId.trim())}
    >
      <input
        class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
        placeholder="Route ID"
        bind:value={routeId}
      />
      <button
        type="submit"
        class="rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-emerald-500"
      >
        Load reviews
      </button>
    </form>

    {#if loading}
      <div class="mt-4">
        <LoadingState message="Talking to the reviews service..." />
      </div>
    {:else if error}
      <p class="mt-4 text-sm text-red-600">{error}</p>
    {:else if reviews.length === 0 && routeId}
      <p class="mt-4 text-sm text-slate-500">No reviews yet for this route. Be the first one!</p>
    {:else if reviews.length > 0}
      <ul class="mt-4 space-y-3">
        {#each reviews as review}
          <li class="rounded-lg border border-slate-100 bg-slate-50/70 p-3">
            <div class="flex items-center justify-between text-sm text-slate-500">
              <span class="font-mono text-xs">{review.user_id}</span>
              <span class="rounded bg-emerald-100 px-2 py-0.5 text-xs font-semibold text-emerald-700">
                {review.rating}/5
              </span>
            </div>
            <p class="mt-2 text-sm text-slate-700">{review.comment}</p>
            <p class="mt-1 text-xs text-slate-400">{review.created_at}</p>
          </li>
        {/each}
      </ul>
    {/if}
  </Card>

  <Card title="Create review">
    <form class="space-y-3" on:submit={submitReview}>
      <input
        class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
        placeholder="User ID"
        bind:value={form.userId}
      />
      <input
        class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
        placeholder="Route ID"
        bind:value={form.routeId}
      />
      <div>
        <label class="block text-xs font-semibold uppercase text-slate-500" for="rating">Rating</label>
        <input
          id="rating"
          name="rating"
          type="number"
          min="1"
          max="5"
          class="mt-1 w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
          bind:value={form.rating}
        />
      </div>
      <textarea
        rows="3"
        class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
        placeholder="Comment"
        bind:value={form.comment}
      ></textarea>
      <button
        type="submit"
        class="w-full rounded-lg bg-slate-900 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-slate-800"
      >
        Publish
      </button>
      {#if successMessage}
        <p class="text-sm text-emerald-600">{successMessage}</p>
      {/if}
    </form>
  </Card>
</div>
