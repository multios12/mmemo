<script lang="ts">
  import { pop } from "svelte-spa-router";
  import type { memoType } from "../models/memoModels.js";
  import { onMount } from "svelte";
  import RichInput from "../components/RichInput/index.svelte";

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

  const sendClick = () => {
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
  const cancelClick = () => pop();

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
    <p class="card-header-title" />
    <button
      class="button is-inverted is-small has-text-danger sp-right"
      on:click={deleteClick}
    >
      <i class="material-icons">delete</i>
    </button>
  </header>

  <section class="card-content p-0">
    {#if errMessage != ""}
      <div class="notification is-danger">{errMessage}</div>
    {/if}
    <div class="columns m-0">
      <div class="field column m-0">
        <label class="label" for="name">name</label>
        <div class="control">
          <input class="input" type="text" name="name" bind:value={memo.Name} />
        </div>
      </div>
      <div class="field column">
        <label class="label" for="date">date</label>
        <div class="control">
          <input class="input" type="text" name="date" bind:value={memo.Date} />
        </div>
      </div>
    </div>
    <div class="field">
      <label class="label" for="value"> value </label>
      <RichInput bind:value={memo.Value} on:textChange={onTextChange} />
    </div>
    <!--
    <div class="columns m-0">
      <div class="field column">
        <label class="label" for="shop">shop</label>
        <div class="control">
          <input class="input" type="text" name="shop" bind:value={memo.Shop} />
        </div>
      </div>
      <div class="field column">
        <label class="label" for="page">home page</label>
        <div class="control">
          <input class="input" type="text" name="page" bind:value={memo.Page} />
        </div>
      </div>
    </div>
    <div class="field">
      <label class="label" for="play">play</label>
      <div class="control">
        <textarea class="textarea" name="play" bind:value={memo.Play} />
      </div>
    </div>
    <div class="field">
      <label class="label" for="talk">talk</label>
      <div class="control">
        <textarea class="textarea" name="talk" bind:value={memo.Talk} />
      </div>
    </div>
    -->
  </section>

  <footer class="card-footer">
    <div class="card-footer-item buttons">
      <button
        class="button is-primary"
        disabled={isLoading}
        class:is-loading={isLoading}
        on:click={sendClick}
      >
        ok
      </button>
      <button class="button" on:click={cancelClick}> cancel </button>
    </div>
  </footer>
</div>
