<script lang="ts">
  import { onMount } from "svelte";
  import type { memoType } from "../models/memoModels.js";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faPlus, faNoteSticky } from "@fortawesome/free-solid-svg-icons";
  import { goto } from "@mateothegreat/svelte5-router";
  import { getMemosApi } from "../models/apiUrl.js";
  library.add(faPlus, faNoteSticky);
  dom.watch();

  // メニューバー表示化
  document
    .querySelector<HTMLDivElement>(".navbar")
    ?.classList.remove("is-hidden");
  let { route } = $props();
  let params = { category: "" };

  let memos: memoType[] = $state([]);
  const showEdit = (id: string | undefined) => {
    goto(`/${params.category}/${id}`);
  };
  onMount(async () => {
    params = route.result.path.params;
    const r = await getMemosApi(params.category);
    memos = await r.json();
  });
</script>

<div class="card px-10">
  <div class="card-content">
    <div class="columns">
      <div class="column">
        <button
          class="button is-primary"
          onclick={() => goto(`/${params.category}/add`)}
        >
          <i class="fa-solid fa-plus"></i>
          add
        </button>
      </div>
    </div>

    <table class="table is-striped is-hoverable is-fullwidth">
      <tbody>
        {#each memos as m}
          <tr>
            <td
              onclick={() => showEdit(m.Id)}
              onkeypress={() => showEdit(m.Id)}
            >
              <button class="is-fullwidth">
                {m.Date}&nbsp;{m.Name}
              </button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
