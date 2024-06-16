<script lang="ts">
  import { pop } from "svelte-spa-router";
  import type { memoType } from "../models/memoModels.js";
  import { onMount } from "svelte";
  import RichInput from "../components/RichInput/index.svelte";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faTrash } from "@fortawesome/free-solid-svg-icons";
  library.add(faTrash);
  dom.watch();

  const m: memoType = {
    Id: undefined,
    Name: "",
    Date: new Date().toISOString().substring(0, 10),
    Value: "",
  };

  // ルーティングパラメータ
  export let params: { id: string | undefined; category: string | undefined } =
    { id: undefined, category: undefined };
  let memo = m;
  let isErr = false;
  let errMessage = "";
  let isLoading = false;
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

<div class="card">
  <header class="card-header sp-panel-heading">
    <input
      class="input"
      type="text"
      name="name"
      bind:value={memo.Name}
      placeholder="name"
    />
    <input
      id="dateInput"
      class="input"
      type="text"
      name="date"
      bind:value={memo.Date}
    />
    <button
      class="button is-inverted is-small has-text-danger sp-right"
      on:click={deleteClick}
    >
      <i class="fa-solid fa-trash"></i>
    </button>
  </header>

  <section class="card-content p-0">
    {#if errMessage != ""}
      <div class="notification is-danger">{errMessage}</div>
    {/if}
    <div class="columns m-0">
      <div class="field column">
        <div class="control"></div>
      </div>
    </div>
    <div class="field">
      <RichInput bind:value={memo.Value} on:textChange={onTextChange} />
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
