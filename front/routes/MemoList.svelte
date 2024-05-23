<script lang="ts">
  import Router, { link, push } from "svelte-spa-router";
  import { onMount } from "svelte";
  import type { memoType } from "../models/memoModels.js";

  // ルーティングパラメータ
  export let params: { category: string | undefined } = { category: undefined };

  let memos: memoType[] = [];
  const showEdit = (id: string | undefined) => {
    push(`/${params.category}/${id}`);
  };
  onMount(async () => {
    console.log(params.category);
    const r = await fetch(`./api/memos/${params.category}`);
    memos = await r.json();
  });
</script>

<div class="card px-10">
  <div class="card-content">
    <div class="columns">
      <div class="column">
        <a class="button is-primary" href={`/${params.category}/add`} use:link
          ><i class="material-icons">add</i>add</a
        >
      </div>
    </div>

    <table class="table is-striped is-hoverable is-fullwidth">
      <tbody>
        {#each memos as m}
          <tr>
            <td
              on:click={() => showEdit(m.Id)}
              on:keypress={() => showEdit(m.Id)}
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
