<script lang="ts">
  import { replace } from "@mateothegreat/svelte5-router";
  import { settingsStore } from "../store.js";

  $: {
    const firstCategory = $settingsStore?.Categories?.[0];

    if (firstCategory?.Key) {
      replace(`/${firstCategory.Key}/`);
    }
  }
</script>
