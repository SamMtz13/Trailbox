<script lang="ts">
  export let title: string;
  export let description = '';
</script>

<section class="rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
  <header class="mb-4">
    <h2 class="text-xl font-semibold text-slate-900">{title}</h2>
    {#if description}
      <p class="text-sm text-slate-500">{description}</p>
    {/if}
  </header>
  <slot />
</section>
