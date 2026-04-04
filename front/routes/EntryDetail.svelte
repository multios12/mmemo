<script lang="ts">
  import { goto, type RouteResult } from "@mateothegreat/svelte5-router";
  import { onMount } from "svelte";
  import TagsInput from "../components/TagsInput.svelte";
  import type { memoType } from "../models/memoModels.js";
  import RichInput from "../components/RichInput/index.svelte";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faTrash } from "@fortawesome/free-solid-svg-icons";

  library.add(faTrash);
  dom.watch();

  interface Props {
    category?: string;
    route?: RouteResult;
    template?: string;
    showTags?: boolean;
    dateEditableOnCreate?: boolean;
  }

  let {
    category = "",
    route: currentRoute = undefined,
    template = "",
    showTags = undefined,
    dateEditableOnCreate = undefined,
  }: Props = $props();

  const routeParams = $derived(
    (currentRoute?.result?.path?.params ?? {}) as {
      id?: string | number | boolean;
      category?: string | number | boolean;
    },
  );
  const entryCategory = $derived(
    category || String(routeParams.category ?? ""),
  );
  const showsTags = $derived(showTags);
  const allowsDateEditOnCreate = $derived(dateEditableOnCreate);

  let innerHeight: number = $state(0);
  let innerWidth: number = $state(0);

  let memo = $state<memoType>({
    Id: undefined,
    Title: "",
    Date: new Date().toISOString().substring(0, 10),
    Value: "",
    Tags: [],
  });
  let isErr = $state(false);
  let errMessage = $state("");
  let isLoading = $state(false);
  let changedValue = $state("");
  let isNew = $state(false);
  let initialized = $state(false);

  document.querySelector<HTMLDivElement>(".navbar")?.classList.add("is-hidden");

  $effect(() => {
    innerHeight;
    innerWidth;
    const headRect = document.querySelector("header")?.getBoundingClientRect();
    const footRect = document.querySelector("footer")?.getBoundingClientRect();
    const barRect = document.querySelector("#toolbar")?.getBoundingClientRect();
    if (footRect && headRect && barRect) {
      const height =
        innerHeight - headRect.height - footRect.height - barRect.height - 30;
      document
        .querySelector<HTMLDivElement>("#detail")
        ?.style.setProperty("height", height + "px");
    }
  });

  const entryId = $derived(routeParams.id ? String(routeParams.id) : undefined);

  const isAddRoute = () => entryId === "add";
  const saveMethod = () => (isNew ? "put" : "post");

  const entryUrl = () => {
    if (isNew) {
      return `/api/${entryCategory}`;
    }
    return `/api/${entryCategory}/${entryId}`;
  };

  const listPath = () => `/${entryCategory}/`;

  const goToList = async () => goto(listPath());

  const onOk = async () => {
    memo.Value = changedValue || memo.Value;
    const response = await fetch(entryUrl(), {
      method: saveMethod(),
      body: JSON.stringify(memo),
    });

    if (response.status !== 200) {
      errMessage = (await response.json()).error;
      isErr = true;
      return;
    }

    await goToList();
  };

  const onCancel = async () => await goToList();

  const onDelete = async () => {
    if (isNew) {
      await goToList();
      return;
    }

    await fetch(entryUrl(), { method: "delete" });
    await goToList();
  };

  const onTextChange = (value: string) => {
    changedValue = value;
  };

  const resetMemo = () => ({
    Id: undefined,
    Title: "",
    Date: new Date().toISOString().substring(0, 10),
    Value: "",
    Tags: [],
  });

  onMount(() => {
    const headerRect = document
      .querySelector("header")
      ?.getBoundingClientRect();
    const barRect = document.querySelector("#toolbar")?.getBoundingClientRect();
    if (headerRect !== undefined && barRect !== undefined) {
      const bar = document.querySelector<HTMLDivElement>("#toolbar");
      bar?.style.setProperty("position", "fixed");
      bar?.style.setProperty("top", headerRect.height + 5 + "px");
      bar?.style.setProperty("width", "100%");

      const rich = document.querySelector<HTMLDivElement>("#rich");
      const top = headerRect.height + barRect.height;
      rich?.style.setProperty("margin-top", top + "px");
    }
    initialized = true;
  });

  $effect(() => {
    initialized;
    entryCategory;
    entryId;
    currentRoute;

    if (!initialized || !entryCategory) {
      return;
    }

    isErr = false;
    errMessage = "";
    memo = resetMemo();
    isNew = isAddRoute();

    if (isNew) {
      const nextMemo = { ...resetMemo(), Value: template ?? "" };
      memo = nextMemo;
      changedValue = nextMemo.Value;
      isLoading = false;
      return;
    }

    if (entryId == undefined || entryId === "") {
      changedValue = "";
      isLoading = false;
      return;
    }

    isLoading = true;
    (async () => {
      try {
        const response = await fetch(`/api/${entryCategory}/${entryId}`);
        const nextMemo = (await response.json()) as memoType;
        nextMemo.Tags = nextMemo.Tags ?? [];
        memo = nextMemo;
        changedValue = nextMemo.Value ?? "";
      } finally {
        isLoading = false;
      }
    })();
  });
</script>

<svelte:window bind:innerHeight bind:innerWidth />

{#if isErr && errMessage != ""}
  <div class="notification is-danger p-1">{errMessage}</div>
{/if}

<header>
  <nav class="level is-mobile m-0">
    <div class="level-item title-left">
      <input
        id="dateInput"
        type="date"
        class="input"
        bind:value={memo.Date}
        disabled={!allowsDateEditOnCreate || !isNew}
      />
      <input
        type="text"
        placeholder="outline"
        class="input"
        bind:value={memo.Title}
      />
    </div>
    {#if !isNew}
      <div class="level-right">
        <div class="level-item">
          <div class="column p-0">
            <button
              class="button has-text-danger"
              aria-label={`delete ${entryCategory}`}
              onclick={onDelete}
            >
              <i class="fa-solid fa-trash"></i>
            </button>
          </div>
        </div>
      </div>
    {/if}
  </nav>
  {#if showsTags}
    <div class="control">
      <TagsInput bind:items={memo.Tags} />
    </div>
  {/if}
</header>

<section class="p-0">
  <div class="field">
    <div class="control py-2">
      <RichInput bind:value={memo.Value} {onTextChange} />
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
  .title-left {
    flex-basis: 100px;
  }

  header {
    background-color: var(--bulma-border);
    width: 100%;
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
