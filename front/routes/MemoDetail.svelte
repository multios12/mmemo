<script lang="ts">
  import { run } from 'svelte/legacy';

  import { pop } from "svelte-spa-router";
  import type { memoType } from "../models/memoModels.js";
  import { onMount } from "svelte";
  import RichInput from "../components/RichInput/index.svelte";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faTrash } from "@fortawesome/free-solid-svg-icons";
  library.add(faTrash);
  dom.watch();

  let innerHeight: number = $state();
  let innerWidth: number = $state();

  // メニューバー非表示化
  document.querySelector<HTMLDivElement>(".navbar")?.classList.add("is-hidden");

  // テキスト部の高さ調整
  run(() => {
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
  });

  const m: memoType = {
    Id: undefined,
    Name: "",
    Date: new Date().toISOString().substring(0, 10),
    Value: "",
  };

  
  interface Props {
    // ルーティングパラメータ
    params?: { id: string | undefined; category: string | undefined };
  }

  let { params = { id: undefined, category: undefined } }: Props = $props();
  let memo = $state(m);
  let isErr = false;
  let errMessage = $state("");
  let isLoading = $state(false);
  let changedValue: string;
  let template: string;

  /*
  onMount(async () => {
    settingsStore.subscribe((s) => {
      for (let index = 0; index < s.Categories.length; index++) {
        const c = s.Categories[index];
      }
    });
  });
  */

  const onOk = () => {
    isLoading = true;
    let url = `./api/memos/${params.category}`;
    url += params.id === "add" ? "" : `/${params.id}`;
    let o = {};
    memo.Value = changedValue;
    if (params.id === "add") {
      o = { method: "put", body: JSON.stringify(memo) };
    } else {
      o = { method: "post", body: JSON.stringify(memo) };
    }

    fetch(url, o)
      .then(pop)
      .catch((res) => {
        errMessage = res.response.data.error;
        isErr = true;
      })
      .finally(() => {
        isLoading = false;
      });
  };

  /** キャンセル クリックイベント */
  const onCancel = () => pop();

  const deleteClick = async () => {
    let url = `./api/memos/${params.category}/${params.id}`;
    await fetch(url, { method: "delete" });
    pop();
  };

  onMount(async () => {
    if (params.id === "add") {
      memo.Value = template;
      return;
    }
    isLoading = true;
    const r = await fetch(`./api/memos/${params.category}/${params.id}`);
    const v = await r.json();
    memo = v[0];
    isLoading = false;
  });

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
      <input id="dateInput" type="date" class="input" bind:value={memo.Date} />
      <input
        type="text"
        placeholder="outline"
        class="input"
        bind:value={memo.Name}
      />
    </div>
    <div class="level-right">
      <div class="level-item">
        <div class="column p-0">
          <button class="button has-text-danger" onclick={deleteClick}>
            <i class="fa-solid fa-trash"></i>
          </button>
        </div>
      </div>
    </div>
  </nav>
</header>

<section class="p-0">
  <div class="field">
    <div class="control py-2">
      <RichInput bind:value={memo.Value} on:textChange={onTextChange} />
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
        onclick={onOk}
      >
        保存
      </button>
    </div>
    <div class="level-item">
      <button class="button" onclick={onCancel}> cancel </button>
    </div>
  </div>
</footer>

<style>
  /* タイトルテキストボックスを300px以下に縮小できるよう調整 */
  .title-left {
    flex-basis: 100px;
  }

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
