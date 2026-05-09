<script lang="ts">
  import { useNavigate } from "@dvcol/svelte-simple-router/router";
  import { onMount } from "svelte";
  import { settingsStore } from "../store.js";

  const { push } = useNavigate();

  let redirected = false;

  $effect(() => {
    const firstCategory = $settingsStore?.Categories?.[0];

    if (!redirected && firstCategory?.Key) {
      redirected = true;
      void push({ path: `/${firstCategory.Key}/` });
    }
  });

  onMount(() => {
    redirected = false;
  });
</script>
