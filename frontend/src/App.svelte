<script lang="ts">
  import Home from './routes/Home.svelte';
  import Users from './routes/Users.svelte';
  import Routes from './routes/Routes.svelte';
  import Workouts from './routes/Workouts.svelte';
  import Reviews from './routes/Reviews.svelte';
  import Leaderboard from './routes/Leaderboard.svelte';
  import Notifications from './routes/Notifications.svelte';
  import Maps from './routes/Maps.svelte';

  type Page = {
    id: string;
    label: string;
    component: typeof Home;
  };

  const pages: Page[] = [
    { id: 'home', label: 'Overview', component: Home },
    { id: 'users', label: 'Users', component: Users },
    { id: 'routes', label: 'Routes', component: Routes },
    { id: 'workouts', label: 'Workouts', component: Workouts },
    { id: 'reviews', label: 'Reviews', component: Reviews },
    { id: 'leaderboard', label: 'Leaderboard', component: Leaderboard },
    { id: 'notifications', label: 'Notifications', component: Notifications },
    { id: 'maps', label: 'Maps', component: Maps }
  ];

  let current = pages[0].id;
  $: ActiveComponent = pages.find((page) => page.id === current)?.component ?? Home;
</script>

<div class="min-h-screen bg-slate-100">
  <header class="border-b border-slate-200 bg-white">
    <div class="mx-auto flex max-w-6xl flex-col gap-3 px-6 py-6 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs font-semibold uppercase tracking-wide text-emerald-600">Trailbox</p>
        <h1 class="text-2xl font-bold text-slate-900">Outdoor Intelligence Dashboard</h1>
        <p class="text-sm text-slate-500">Gateway-powered UI connected to every Go microservice.</p>
      </div>
      <nav class="flex flex-wrap gap-2 text-sm">
        {#each pages as page}
          <button
            class={`rounded-full border px-3 py-1 font-medium transition ${
              page.id === current
                ? 'border-slate-900 bg-slate-900 text-white'
                : 'border-slate-300 bg-white text-slate-600 hover:border-slate-400'
            }`}
            on:click={() => (current = page.id)}
          >
            {page.label}
          </button>
        {/each}
      </nav>
    </div>
  </header>

  <main class="mx-auto max-w-6xl px-6 py-8">
    <svelte:component this={ActiveComponent} />
  </main>
</div>
