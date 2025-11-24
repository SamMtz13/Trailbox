<script lang="ts">
  import { onMount } from 'svelte';
  import Home from './routes/Home.svelte';
  import Users from './routes/Users.svelte';
  import Routes from './routes/Routes.svelte';
  import Workouts from './routes/Workouts.svelte';
  import Reviews from './routes/Reviews.svelte';
  import Leaderboard from './routes/Leaderboard.svelte';
  import Notifications from './routes/Notifications.svelte';
  import Maps from './routes/Maps.svelte';
  import NotFound from './routes/NotFound.svelte';

  const routeTable = new Map([
    ['/', Home],
    ['/users', Users],
    ['/routes', Routes],
    ['/workouts', Workouts],
    ['/reviews', Reviews],
    ['/leaderboard', Leaderboard],
    ['/notifications', Notifications],
    ['/maps', Maps]
  ]);

  const navItems = [
    { path: '/', label: 'Inicio' },
    { path: '/users', label: 'Usuarios' },
    { path: '/routes', label: 'Rutas' },
    { path: '/workouts', label: 'Entrenamientos' },
    { path: '/reviews', label: 'Reseñas' },
    { path: '/leaderboard', label: 'Leaderboard' },
    { path: '/notifications', label: 'Notificaciones' },
    { path: '/maps', label: 'Mapas' }
  ];

  let currentPath = normalizePath(window.location.pathname);
  let Component = routeTable.get(currentPath) ?? NotFound;

  function normalizePath(path: string) {
    if (path !== '/' && path.endsWith('/')) return path.slice(0, -1);
    return path || '/';
  }

  function navigate(path: string) {
    const normalized = normalizePath(path);
    if (normalized === currentPath) return;
    history.pushState({}, '', normalized);
    currentPath = normalized;
    Component = routeTable.get(normalized) ?? NotFound;
  }

  onMount(() => {
    const onPopState = () => {
      const nextPath = normalizePath(window.location.pathname);
      currentPath = nextPath;
      Component = routeTable.get(nextPath) ?? NotFound;
    };
    window.addEventListener('popstate', onPopState);
    return () => window.removeEventListener('popstate', onPopState);
  });
</script>

<svelte:window on:click={(event) => {
  const target = event.target as HTMLElement;
  if (target?.closest('a')?.getAttribute('data-nav') === 'true') {
    event.preventDefault();
    navigate((target.closest('a') as HTMLAnchorElement).getAttribute('href') || '/');
  }
}} />

<div class="min-h-screen">
  <header class="bg-white/80 backdrop-blur border-b border-emerald-100 sticky top-0 z-10">
    <div class="max-w-6xl mx-auto px-6 py-4 flex items-center gap-6">
      <div class="h-10 w-10 rounded-xl bg-forest text-white font-bold flex items-center justify-center">TB</div>
      <div>
        <p class="text-xs uppercase tracking-[0.2em] text-emerald-700">Trailbox</p>
        <p class="text-sm text-forest font-semibold">Proyecto Final · Cómputo Distribuido</p>
      </div>
      <nav class="ml-auto hidden md:flex gap-3">
        {#each navItems as item}
          <a
            data-nav="true"
            href={item.path}
            class={`px-3 py-2 rounded-lg text-sm font-semibold transition ${currentPath === item.path ? 'bg-forest text-white' : 'text-forest hover:bg-emerald-50'}`}
            >{item.label}</a
          >
        {/each}
      </nav>
      <div class="md:hidden ml-auto">
        <select class="input" bind:value={currentPath} on:change={(e) => navigate(e.currentTarget.value)}>
          {#each navItems as item}
            <option value={item.path}>{item.label}</option>
          {/each}
        </select>
      </div>
    </div>
  </header>

  <main class="max-w-6xl mx-auto px-6 py-10 space-y-6">
    <svelte:component this={Component} on:navigate={(e) => navigate(e.detail)} />
  </main>
</div>
