<script lang="ts">
  import { onMount } from "svelte";
  import { location, pop } from "svelte-spa-router";
  import TagsInput from "../components/TagsInput.svelte";
  import type { detailType } from "../models/diaryModels.js";
  import RichInput from "../components/RichInput/index.svelte";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faTrash } from "@fortawesome/free-solid-svg-icons";
  library.add(faTrash);
  dom.watch();

  export let params: { id: string | undefined } = { id: undefined };
  export let Template: string;
  let isErr = false;
  let errMessage = "";

  let editDay: string | null;
  let Outline = "";
  let Detail = "";
  let Tags: string[] = [];
  let isDayEdit: boolean;
  let isLoading: boolean;
  let changedValue: string;

  onMount(async () => {
    isLoading = true;
    if ($location.indexOf("/add") > -1) {
      isDayEdit = true;
      const dt = new Date();
      editDay = `${dt.getFullYear()}-`;
      editDay += ("00" + (dt.getMonth() + 1)).slice(-2);
      editDay += `-${("00" + dt.getDate()).slice(-2)}`;
      Detail = Template;
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

    isLoading = false;
  });

  /** 送信 クリックイベント */
  const onOk = async () => {
    if (editDay == null) {
      return;
    }

    let url = editDay.replaceAll("-", "/");
    url = `./api/diary/${url}`;
    const init = {
      method: "post",
      body: JSON.stringify({ Tags, Outline, Detail: changedValue }),
    };

    let r = await fetch(url, init);
    if (r.status != 200) {
      errMessage = (await r.json()).error;
      isErr = true;
    } else {
      editDay = null;
      pop();
    }
  };

  /** キャンセル クリックイベント */
  const onCancel = () => pop();

  /** 削除 クリックイベント */
  const onDelete = async () => {
    if (editDay == null) {
      return;
    }
    let url = editDay.replaceAll("-", "/");
    url = `./api/diary/${url}`;
    await fetch(url, { method: "delete" });
    pop();
  };

  const onTextChange = (e: CustomEvent<string>) => {
    changedValue = e.detail;
  };
</script>

{#if errMessage != ""}
  <div class="notification is-danger p-1">{errMessage}</div>
{/if}
<div class="card m-0" class:is-active={editDay != null}>
  <header class="card-header sp-panel-heading">
    <input
      id="dateInput"
      type="date"
      class="input"
      bind:value={editDay}
      disabled={!isDayEdit}
    />
    <input
      type="text"
      placeholder="outline"
      class="input"
      bind:value={Outline}
    />
    {#if !isDayEdit}
      <button
        class="button is-inverted is-small has-text-danger"
        on:click={onDelete}
      >
        <i class="fa-solid fa-trash"></i>
      </button>
    {/if}
  </header>

  <section class="card-content p-0">
    <!-- Card Content -->
    <div class="field">
      <div class="control">
        <TagsInput bind:items={Tags} />
      </div>
      <div class="control py-2">
        <RichInput bind:value={Detail} on:textChange={onTextChange} />
      </div>
    </div>
  </section>
</div>
<footer class="columns is-dark is-mobile">
  <button
    class="column button is-primary"
    disabled={isLoading}
    class:is-loading={isLoading}
    on:click={onOk}
  >
    ok
  </button>
  <button class="column button" on:click={onCancel}> cancel </button>
</footer>

<style>
  footer {
    left: 0;
    bottom: 0;
    width: 100%;
    position: fixed;
  }
  #dateInput {
    width: 150px;
  }
</style>
