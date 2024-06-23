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
  let innerHeight: number;
  let innerWidth: number;
  // メニューバー非表示化
  document.querySelector<HTMLDivElement>(".navbar")?.classList.add("is-hidden");

  // テキスト部の高さ調整
  $: {
    let a = innerHeight + innerWidth + 1;
    let headRect = document.querySelector("header")?.getBoundingClientRect();
    let footRect = document.querySelector("footer")?.getBoundingClientRect();
    let barRect = document.querySelector("#toolbar")?.getBoundingClientRect();
    if (footRect && headRect && barRect) {
      let height =
        innerHeight - headRect.height - footRect.height - barRect.height - 30;
      document
        .querySelector<HTMLDivElement>("#detail")
        ?.style.setProperty("height", height + "px");
    }
  }

  /** マウントイベント */
  onMount(async () => {
    // コントロールの位置調整
    let headerRect = document.querySelector("header")?.getBoundingClientRect();
    let barRect = document.querySelector("#toolbar")?.getBoundingClientRect();
    if (headerRect !== undefined && barRect !== undefined) {
      let bar = document.querySelector<HTMLDivElement>("#toolbar");
      bar?.style.setProperty("position", "fixed");
      bar?.style.setProperty("top", headerRect.height + 5 + "px");
      bar?.style.setProperty("width", "100%");

      let rich = document.querySelector<HTMLDivElement>("#rich");
      let top = headerRect.height + barRect.height;
      rich?.style.setProperty("margin-top", top + "px");
    }

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

<svelte:window bind:innerHeight bind:innerWidth />

{#if errMessage != ""}
  <div class="notification is-danger p-1">{errMessage}</div>
{/if}

<header>
  <nav class="level is-mobile m-0">
    <div class="level-item title-left">
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
    </div>
    {#if !isDayEdit}
      <div class="level-right">
        <div class="level-item">
          <div class="column p-0">
            <button class="button has-text-danger" on:click={onDelete}>
              <i class="fa-solid fa-trash" />
            </button>
          </div>
        </div>
      </div>
    {/if}
  </nav>
  <div class="control">
    <TagsInput bind:items={Tags} />
  </div>
</header>

<section class="p-0">
  <div class="field">
    <div class="control py-2">
      <RichInput bind:value={Detail} on:textChange={onTextChange} />
    </div>
  </div>
</section>
<footer class="is-dark m-0">
  <div class="level is-mobile m-0">
    <div class="level-item p-0">
      <button
        class="button is-primary"
        disabled={isLoading}
        class:is-loading={isLoading}
        on:click={onOk}
      >
        保存
      </button>
    </div>
    <div class="level-item">
      <button class="button" on:click={onCancel}> cancel </button>
    </div>
  </div>
</footer>

<style>
  /* タイトルテキストボックスを300px以下に縮小できるよう調整 */
  .title-left {
    flex-basis: 100px;
  }
  header {
    background-color: var(--bulma-border);
    width: 100%;
    position: fixed;
    top: 0;
    left: 0;
    margin: 0;
    z-index: 999;
  }
  footer {
    background-color: var(--bulma-border);
    left: 0;
    bottom: 0;
    width: 100%;
    position: fixed;
  }
  #dateInput {
    width: 150px;
  }
</style>
