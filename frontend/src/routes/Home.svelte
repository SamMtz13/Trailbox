<script lang="ts">
  import Card from '../lib/components/Card.svelte';
</script>

<div class="grid gap-6">
  <Card title="Trailbox" description="Full stack observability for runners, hikers and cyclists.">
    <p class="text-sm text-slate-600">
      This dashboard speaks to the Go microservices through the gateway. Each section on the left loads live data
      through the Kubernetes LoadBalancer entrypoint. Use it to validate routes, workouts, reviews, leaderboard entries,
      notifications and map information.
    </p>
  </Card>

  <Card title="How to use it">
    <ul class="list-disc space-y-2 pl-4 text-sm text-slate-600">
      <li>Users, routes and workouts load automatically for quick sanity checks.</li>
      <li>Reviews, notifications and leaderboard panels include small forms to exercise POST flows.</li>
      <li>The Maps panel lets you store arbitrary GeoJSON payloads mapped to an existing route.</li>
      <li>Everything runs through the same API gateway, so the UI works locally and inside the cluster.</li>
    </ul>
  </Card>
</div>
