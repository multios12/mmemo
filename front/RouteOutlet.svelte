<script lang="ts">
  import { useRoute } from "@dvcol/svelte-simple-router/router";

  const currentRoute = $derived.by(() => useRoute().route);
</script>

{#if currentRoute?.component}
  {@const CurrentRoute = currentRoute.component}
  <CurrentRoute {...(currentRoute.props ?? {})} />
{/if}
