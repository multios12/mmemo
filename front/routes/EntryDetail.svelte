<script lang="ts">
  import { goto, type RouteResult } from "@mateothegreat/svelte5-router";
  import { onMount } from "svelte";
  import TagsInput from "../components/TagsInput.svelte";
  import type { entryType } from "../models/entryModels.js";
  import RichInput from "../components/RichInput/index.svelte";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faTrash } from "@fortawesome/free-solid-svg-icons";
  import { settingsStore } from "../store.js";

  library.add(faTrash);
  dom.watch();

  interface Props {
    route?: RouteResult;
  }

  type EntryRouteParams = {
    id?: string | number | boolean;
    category?: string | number | boolean;
  };

  const createEmptyEntry = (): entryType => ({
    Id: undefined,
    Outline: "",
    Date: new Date().toISOString().substring(0, 10),
    Value: "",
    Tags: [],
  });

  let { route: currentRoute = undefined }: Props = $props();

  const routeParams = $derived(
    (currentRoute?.result?.path?.params ?? {}) as EntryRouteParams,
  );
  const categoryKey = $derived(String(routeParams.category ?? ""));
  const resolvedSettings = $derived.by(() => {
    const categorySetting = $settingsStore?.Categories?.find(
      (category) => category.Key === categoryKey,
    );

    return {
      template: categorySetting?.Template || "",
      showTags: categorySetting?.UseTag ?? false,
      allowsDateEditOnCreate: categorySetting?.UseDate ?? false,
    };
  });

  let innerHeight: number = $state(0);
  let innerWidth: number = $state(0);

  let entry = $state<entryType>(createEmptyEntry());
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
  const isAddRoute = $derived(entryId === "add");
  const imageUploadPath = $derived(
    isAddRoute
      ? `/api/${categoryKey}/images/tmp`
      : `/api/${categoryKey}/${entryId}/images`,
  );
  const saveMethod = () => (isNew ? "put" : "post");

  const entryUrl = () => {
    if (isNew) {
      return `/api/${categoryKey}`;
    }
    return `/api/${categoryKey}/${entryId}`;
  };

  const listPath = () => `/${categoryKey}/`;

  const goToList = async () => goto(listPath());

  const onOk = async () => {
    entry.Value = changedValue || entry.Value;
    const response = await fetch(entryUrl(), {
      method: saveMethod(),
      body: JSON.stringify(entry),
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
    categoryKey;
    entryId;
    currentRoute;

    if (!initialized || !categoryKey) {
      return;
    }

    isErr = false;
    errMessage = "";
    entry = createEmptyEntry();
    isNew = isAddRoute;

    if (isNew) {
      const nextEntry = { ...createEmptyEntry(), Value: resolvedSettings.template };
      entry = nextEntry;
      changedValue = nextEntry.Value;
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
        const response = await fetch(`/api/${categoryKey}/${entryId}`);
        const nextEntry = (await response.json()) as entryType;
        nextEntry.Tags = nextEntry.Tags ?? [];
        entry = nextEntry;
        changedValue = nextEntry.Value ?? "";
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
        bind:value={entry.Date}
        disabled={!resolvedSettings.allowsDateEditOnCreate || !isNew}
      />
      <input
        type="text"
        placeholder="outline"
        class="input"
        bind:value={entry.Outline}
      />
    </div>
    {#if !isNew}
      <div class="level-right">
        <div class="level-item">
          <div class="column p-0">
            <button
              class="button has-text-danger"
              aria-label={`delete ${categoryKey}`}
              onclick={onDelete}
            >
              <i class="fa-solid fa-trash"></i>
            </button>
          </div>
        </div>
      </div>
    {/if}
  </nav>
  {#if resolvedSettings.showTags}
    <div class="control">
      <TagsInput bind:items={entry.Tags} />
    </div>
  {/if}
</header>

<section class="p-0">
  <div class="field">
    <div class="control py-2">
      <RichInput bind:value={entry.Value} {imageUploadPath} {onTextChange} />
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
