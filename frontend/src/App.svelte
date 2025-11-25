<script lang="ts">
  import { onMount } from "svelte";
  import type { ComponentType } from "svelte";

  import Home from "./routes/Home.svelte";
  import Users from "./routes/Users.svelte";
  import RoutesPage from "./routes/Routes.svelte";
  import Workouts from "./routes/Workouts.svelte";
  import Reviews from "./routes/Reviews.svelte";
  import Leaderboard from "./routes/Leaderboard.svelte";
  import Notifications from "./routes/Notifications.svelte";
  import Maps from "./routes/Maps.svelte";
  import NotFound from "./routes/NotFound.svelte";

  const routeTable = new Map<string, ComponentType>([
    ["/", Home],
    ["/users", Users],
    ["/routes", RoutesPage],
    ["/workouts", Workouts],
    ["/reviews", Reviews],
    ["/leaderboard", Leaderboard],
    ["/notifications", Notifications],
    ["/maps", Maps]
  ]);

  const navItems = [
    { path: "/", label: "Inicio" },
    { path: "/users", label: "Usuarios" },
    { path: "/routes", label: "Rutas" },
    { path: "/workouts", label: "Entrenamientos" },
    { path: "/reviews", label: "Reseñas" },
    { path: "/leaderboard", label: "Leaderboard" },
    { path: "/notifications", label: "Notificaciones" },
    { path: "/maps", label: "Mapas" }
  ];

  let currentPath = "/";
  let Component: ComponentType = Home;

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

  onMount(() => {
    setRouteFromPath(window.location.pathname);

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
  <header class="bg-black/40 backdrop-blur border-b border-emerald-700 sticky top-0 z-20">
    <div class="max-w-6xl mx-auto px-6 py-4 flex items-center gap-6">
      <h1 class="text-emerald-300 font-bold text-xl">TRAILBOX</h1>
      <nav class="ml-auto hidden md:flex gap-3">
        {#each navItems as item}
          <a
            data-nav="true"
            href={item.path}
            class={`px-3 py-2 rounded-lg text-sm font-semibold transition ${
              currentPath === item.path
                ? "bg-emerald-500 text-black"
                : "text-emerald-300 hover:bg-emerald-700/20"
            }`}
          >
            {item.label}
          </a>
        {/each}
      </nav>
    </div>
  </header>

  <main class="max-w-6xl mx-auto px-6 py-10 space-y-6">
    <svelte:component this={Component} />
  </main>
</div>
