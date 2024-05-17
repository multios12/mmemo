<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import { location, pop, push } from "svelte-spa-router";
  import TagsInput from "../components/TagsInput.svelte";
  import type { detailType } from "../models/diaryModels.js";
  import RichInput from "../components/RichInput/RichInput.svelte";

  export let params: { id: string | undefined } = { id: undefined };
  let editDay: string | null;
  const dispatch = createEventDispatcher();
  let Outline = "";
  let Detail = "";
  let Tags: string[] = [];
  let isDayEdit: boolean;

  onMount(async () => {
    if ($location.indexOf("/add") > -1) {
      isDayEdit = true;
      const dt = new Date();
      editDay = `${dt.getFullYear()}-`;
      editDay += ("00" + (dt.getMonth() + 1)).slice(-2);
      editDay += `-${("00" + dt.getDate()).slice(-2)}`;
    } else {
      if (params == undefined) {
        return;
      }
      isDayEdit = false;
      if (params.id == undefined) {
        return;
      }
      let url = params.id.replaceAll("-", "/");
      url = `./api/diary/${url}`;
      fetch(url, { method: "get" })
        .then((r) => r.json())
        .then((r) => {
          const s = r as detailType;
          editDay = <string | null>params.id;
          Outline = s.Outline;
          Detail = s.Detail;
          Tags = s.Tags;
        });
    }
  });

  /** 送信 クリックイベント */
  const sendClick = async () => {
    if (editDay == null) {
      return;
    }

    let url = editDay.replaceAll("-", "/");
    url = `./api/diary/${url}`;
    const init = {
      method: "post",
      body: JSON.stringify({ Tags, Outline, Detail }),
    };
    await fetch(url, init);
    editDay = null;
    push("/d/");
  };

  /** キャンセル クリックイベント */
  const cancelClick = () => pop();

  /** 削除 クリックイベント */
  const deleteClick = async () => {
    if (editDay == null) {
      return;
    }
    let url = editDay.replaceAll("-", "/");
    url = `./api/diary/${url}`;
    await fetch(url, { method: "delete" });
    pop();
  };
</script>

<div class="card" class:is-active={editDay != null}>
  <header class="card-header sp-panel-heading">
    <input
      type="date"
      class="input"
      bind:value={editDay}
      disabled={!isDayEdit}
    />
    {#if !isDayEdit}
      <button
        class="button is-inverted is-small has-text-danger"
        on:click={deleteClick}
      >
        <i class="material-icons">delete</i>
      </button>
    {/if}
  </header>
  <section class="card-content p-0">
    <!-- Card Content -->
    <input
      type="text"
      placeholder="outline"
      class="input"
      bind:value={Outline}
    />
    <div class="field">
      <div class="control">
        <TagsInput bind:items={Tags} />
      </div>
      <div class="control py-2">
        <RichInput bind:value={Detail} />
      </div>
    </div>
  </section>
  <footer class="card-footer">
    <button class="button is-primary" on:click={sendClick}>send</button>
    <button class="button" on:click={cancelClick}>cancel</button>
  </footer>
</div>
