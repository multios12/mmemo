<script lang="ts">
  import Router, { link, push } from "svelte-spa-router";
  import { onMount } from "svelte";
  import type { memoType } from "../models/memoModels.js";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faPlus, faNoteSticky } from "@fortawesome/free-solid-svg-icons";
  library.add(faPlus, faNoteSticky);
  dom.watch();

  // メニューバー表示化
  document
    .querySelector<HTMLDivElement>(".navbar")
    ?.classList.remove("is-hidden");

  
  interface Props {
    // ルーティングパラメータ
    params?: { category: string | undefined };
  }

  let { params = { category: undefined } }: Props = $props();

  let memos: memoType[] = $state([]);
  const showEdit = (id: string | undefined) => {
    push(`/${params.category}/${id}`);
  };
  onMount(async () => {
    const r = await fetch(`./api/memos/${params.category}`);
    memos = await r.json();
  });
</script>

<div class="card px-10">
  <div class="card-content">
    <div class="columns">
      <div class="column">
        <a class="button is-primary" href={`/${params.category}/add`} use:link>
          <i class="fa-solid fa-plus"></i>
          add
        </a>
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
