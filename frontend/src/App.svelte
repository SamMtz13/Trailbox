<script lang="ts">
  import { onMount } from "svelte";
  import type { ComponentType } from "svelte";

  import Home from "./routes/Home.svelte";
  import Users from "./routes/Users.svelte";
  import Routes from "./routes/Routes.svelte";
  import Workouts from "./routes/Workouts.svelte";
  import Leaderboard from "./routes/Leaderboard.svelte";
  import Maps from "./routes/Maps.svelte";
  import NotFound from "./routes/NotFound.svelte";
  import { api } from "./lib/api";

  type Notification = {
    id: string;
    user_id: string;
    message: string;
    read: boolean;
    created_at: string;
  };

  const routeTable = new Map<string, ComponentType>([
    ["/", Home],
    ["/users", Users],
    ["/routes", Routes],
    ["/workouts", Workouts],
    ["/leaderboard", Leaderboard],
    ["/maps", Maps]
  ]);

  const navItems = [
    { path: "/", label: "Inicio" },
    { path: "/users", label: "Usuarios" },
    { path: "/routes", label: "Rutas" },
    { path: "/workouts", label: "Entrenamientos" },
    { path: "/leaderboard", label: "Leaderboard" },
    { path: "/maps", label: "Mapas" }
  ];

  let currentPath = "/";
  let Component: ComponentType = Home;
  let notifications: Notification[] = [];
  let notificationsLoading = false;
  let showNotifications = false;

  function normalizePath(path: string) {
    if (path !== "/" && path.endsWith("/")) return path.slice(0, -1);
    return path || "/";
  }

  function setRouteFromPath(path: string) {
    const normalized = normalizePath(path);
    currentPath = normalized;
    Component = routeTable.get(normalized) ?? NotFound;
  }

  function navigate(path: string) {
    const normalized = normalizePath(path);
    if (normalized === currentPath) return;

    history.pushState({}, "", normalized);
    setRouteFromPath(normalized);
  }

  async function loadNotifications() {
    notificationsLoading = true;
    try {
      const data = await api.getNotifications('11111111-1111-1111-1111-111111111111');
      notifications = data.notifications || [];
    } catch (err) {
      console.error(err);
      notifications = [];
    } finally {
      notificationsLoading = false;
    }
  }

  onMount(() => {
    setRouteFromPath(window.location.pathname);
    loadNotifications();

    const onPopState = () =>
      setRouteFromPath(window.location.pathname);

    window.addEventListener("popstate", onPopState);

    return () => {
      window.removeEventListener("popstate", onPopState);
    };
  });

  function handleGlobalClick(event: MouseEvent) {
    const target = event.target as HTMLElement | null;
    const link = target?.closest("a");

    if (link && link.getAttribute("data-nav") === "true") {
      event.preventDefault();
      navigate((link as HTMLAnchorElement).getAttribute("href") || "/");
    }
  }
</script>

<svelte:window on:click={handleGlobalClick} />

<div class="min-h-screen">
  <header class="sticky top-0 z-20 border-b border-emerald-500/40 bg-black/40 backdrop-blur">
    <div class="mx-auto flex max-w-6xl items-center gap-6 px-6 py-4">
      <h1 class="text-xl font-bold text-emerald-300">TRAILBOX</h1>
      <nav class="ml-auto hidden gap-3 md:flex">
        {#each navItems as item}
          <a
            data-nav="true"
            href={item.path}
            class={`px-3 py-2 rounded-lg text-sm font-semibold transition ${
              currentPath === item.path
                ? 'bg-emerald-500 text-black'
                : 'text-emerald-300 hover:bg-emerald-700/20'
            }`}
          >
            {item.label}
          </a>
        {/each}
      </nav>
      <div class="relative ml-4">
        <button
          class="flex h-10 w-10 items-center justify-center rounded-full border border-emerald-400/60 bg-slate-900/70 text-emerald-200 transition hover:bg-emerald-500/20"
          on:click={() => (showNotifications = !showNotifications)}
        >
          <span class="text-lg">🔔</span>
        </button>
        {#if notifications.length > 0}
          <span class="absolute -right-1 -top-1 h-4 w-4 rounded-full bg-red-500 text-[10px] font-bold text-white flex items-center justify-center">
            {notifications.length}
          </span>
        {/if}
        {#if showNotifications}
          <div class="absolute right-0 mt-2 w-72 rounded-2xl border border-white/10 bg-slate-900/90 p-3 shadow-2xl backdrop-blur">
            <div class="flex items-center justify-between text-xs text-slate-400">
              <span>Notificaciones recientes</span>
              <button class="text-emerald-300 hover:underline" on:click={() => showNotifications = false}>Cerrar</button>
            </div>
            <div class="mt-2 max-h-72 space-y-2 overflow-y-auto pr-1">
              {#if notificationsLoading}
                <p class="text-xs text-emerald-100">Cargando...</p>
              {:else if notifications.length === 0}
                <p class="text-xs text-slate-300">Sin mensajes.</p>
              {:else}
                {#each notifications as n}
                  <div class="rounded-xl border border-white/10 bg-slate-950/70 p-2">
                    <p class="text-xs text-slate-400">{new Date(n.created_at).toLocaleString()}</p>
                    <p class="text-sm text-slate-100 mt-1">{n.message}</p>
                    <p class="text-[11px] text-slate-500 mt-1">User: <span class="font-mono">{n.user_id}</span></p>
                  </div>
                {/each}
              {/if}
            </div>
          </div>
        {/if}
      </div>
    </div>
  </header>

  <main class="mx-auto max-w-6xl space-y-6 px-6 py-10">
    <svelte:component this={Component} />
  </main>
</div>
