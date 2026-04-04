<script lang="ts">
  import { goto } from "@mateothegreat/svelte5-router";
  import { settingsStore } from "../store.js";
  import { appPath } from "../basePath.js";

  $: {
    const firstCategory = $settingsStore?.Categories?.[0];

    if (firstCategory?.Key) {
      goto(appPath(`/${firstCategory.Key}/`));
    }
  }
</script>
