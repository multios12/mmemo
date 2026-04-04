<script lang="ts">
  import { goto, type RouteResult } from "@mateothegreat/svelte5-router";
  import { onMount } from "svelte";
  import TagsInput from "../components/TagsInput.svelte";
  import type { entryType } from "../models/entryModels.js";
  import type { TemplateType } from "../models/settingType.js";
  import RichInput from "../components/RichInput/index.svelte";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import {
    faArrowLeft,
    faBook,
    faCloudArrowUp,
    faTags,
    faTrash,
  } from "@fortawesome/free-solid-svg-icons";
  import { settingsStore } from "../store.js";
  import { apiPath, appPath } from "../basePath.js";

  library.add(faTrash, faTags, faCloudArrowUp, faArrowLeft, faBook);
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
      templates: categorySetting?.Templates ?? [],
      showTags: categorySetting?.UseTag ?? false,
      allowMultipleEntriesPerDate:
        categorySetting?.AllowMultipleEntriesPerDate ?? true,
      dateLabel: categorySetting?.Fields?.Date || "日付",
      tagsLabel: categorySetting?.Fields?.Tags || "タグ",
      outlineLabel: categorySetting?.Fields?.Outline || "見出し",
    };
  });

  let innerHeight: number = $state(0);
  let innerWidth: number = $state(0);

  let entry = $state<entryType>(createEmptyEntry());
  let isErr = $state(false);
  let errMessage = $state("");
  let isLoading = $state(false);
  let editorValue = $state("");
  let isNew = $state(false);
  let initialized = $state(false);
  let initialSnapshot = $state("");
  let canCheckDirty = $state(false);
  let showTagsOnMobile = $state(false);
  let outlineInput = $state<HTMLInputElement | null>(null);
  let isPageLoading = $state(true);
  let selectedTemplateIndex = $state<number | null>(null);
  let isTemplateSelectorOpen = $state(false);

  document.querySelector<HTMLDivElement>(".navbar")?.classList.add("is-hidden");

  $effect(() => {
    innerHeight;
    innerWidth;
    const headRect = document.querySelector("header")?.getBoundingClientRect();
    const footRect = document.querySelector("footer")?.getBoundingClientRect();
    const barRect = document.querySelector("#toolbar")?.getBoundingClientRect();
    if (footRect && headRect && barRect) {
      const height =
        innerHeight - headRect.height - footRect.height - barRect.height - 28;
      document
        .querySelector<HTMLDivElement>("#detail")
        ?.style.setProperty("height", height + "px");
    }
  });

  const entryId = $derived(routeParams.id ? String(routeParams.id) : undefined);
  const isAddRoute = $derived(entryId === "add");
  const imageUploadPath = $derived(
    isAddRoute
      ? apiPath(`${categoryKey}/images/tmp`)
      : apiPath(`${categoryKey}/${entryId}/images`),
  );
  const saveMethod = () => (isNew ? "put" : "post");
  const finishPageLoading = () => {
    requestAnimationFrame(() => {
      requestAnimationFrame(() => {
        isPageLoading = false;
      });
    });
  };
  const applyInitialEditorValue = (value: string) => {
    requestAnimationFrame(() => {
      requestAnimationFrame(() => {
        editorValue = value;
        entry.Value = value;
        initialSnapshot = snapshotEntry(entry, value);
      });
    });
  };
  const snapshotEntry = (target: entryType, value: string) =>
    JSON.stringify({
      Outline: target.Outline ?? "",
      Date: target.Date ?? "",
      Value: value ?? "",
      Tags: target.Tags ?? [],
    });
  const isDirty = $derived(
    canCheckDirty && snapshotEntry(entry, editorValue) !== initialSnapshot,
  );
  const saveButtonLabel = $derived(isNew ? "作成" : "保存");
  const templateOptions = $derived.by(() => {
    const options = [...resolvedSettings.templates];
    if (!isNew) {
      return options;
    }

    return [
      { Name: "空白から作成", Value: "" } satisfies TemplateType,
      ...options,
    ];
  });
  const selectableTemplateCount = $derived(resolvedSettings.templates.length);
  const showTemplateSelector = $derived(isNew && isTemplateSelectorOpen);

  const entryUrl = () => {
    if (isNew) {
      return apiPath(categoryKey);
    }
    return apiPath(`${categoryKey}/${entryId}`);
  };

  const listPath = () => appPath(`/${categoryKey}/`);

  const goToList = async () => goto(listPath());
  const shouldLeave = () =>
    !isDirty || window.confirm("未保存の変更があります。戻りますか？");

  const onOk = async () => {
    entry.Value = editorValue;
    if (!resolvedSettings.allowMultipleEntriesPerDate && entry.Date) {
      const month = entry.Date.slice(0, 7);
      const response = await fetch(`${apiPath(categoryKey)}?month=${month}`);
      const monthlyEntries = (await response.json()) as entryType[];
      const hasDuplicateDate = monthlyEntries.some(
        (item) => item.Date === entry.Date && item.Id !== entry.Id,
      );
      if (hasDuplicateDate) {
        errMessage = "同じ日付のエントリは登録できません";
        isErr = true;
        return;
      }
    }

    const response = await fetch(entryUrl(), {
      method: saveMethod(),
      body: JSON.stringify(entry),
    });

    if (response.status !== 200) {
      errMessage = (await response.json()).error;
      isErr = true;
      return;
    }

    initialSnapshot = snapshotEntry(entry, editorValue);
    canCheckDirty = true;
    await goToList();
  };

  const onCancel = async () => {
    if (!shouldLeave()) {
      return;
    }
    await goToList();
  };

  const onDelete = async () => {
    if (isNew) {
      await goToList();
      return;
    }

    await fetch(entryUrl(), { method: "delete" });
    await goToList();
  };

  const onTextChange = (value: string) => {
    editorValue = value;
    entry.Value = value;
  };
  const onSelectTemplate = (template: TemplateType, index: number) => {
    selectedTemplateIndex = index;
    isTemplateSelectorOpen = false;
    applyInitialEditorValue(template.Value ?? "");
  };
  const onSaveTemplate = async () => {
    const suggestedName = entry.Outline.trim() || "新しいテンプレート";
    const name = window.prompt("テンプレート名", suggestedName)?.trim() ?? "";
    if (name === "") {
      return;
    }
    if (editorValue.trim() === "") {
      errMessage = "本文が空のためテンプレート保存できません";
      isErr = true;
      return;
    }

    const response = await fetch(apiPath(`${categoryKey}/templates`), {
      method: "post",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ Name: name, Value: editorValue }),
    });
    if (response.status !== 200) {
      errMessage = (await response.json()).error;
      isErr = true;
      return;
    }

    const nextSettings = await response.json();
    settingsStore.set(nextSettings);
    isErr = false;
    errMessage = "";
  };
  const templatePreview = (value: string) =>
    value
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter((line) => line !== "")
      .slice(0, 3);

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

    isPageLoading = true;
    isErr = false;
    errMessage = "";
    entry = createEmptyEntry();
    isNew = isAddRoute;
    showTagsOnMobile = false;
    selectedTemplateIndex = null;
    isTemplateSelectorOpen = false;
    canCheckDirty = true;

    if (isNew) {
      const nextEntry = createEmptyEntry();
      const initialTemplateValue = resolvedSettings.templates[0]?.Value ?? "";
      entry = nextEntry;
      editorValue = "";
      nextEntry.Value = "";
      initialSnapshot = snapshotEntry(nextEntry, "");
      isLoading = false;
      finishPageLoading();
      if (selectableTemplateCount > 1) {
        isTemplateSelectorOpen = true;
      } else if (initialTemplateValue !== "") {
        applyInitialEditorValue(initialTemplateValue);
      }
      return;
    }

    if (entryId == undefined || entryId === "") {
      editorValue = "";
      initialSnapshot = snapshotEntry(entry, "");
      canCheckDirty = false;
      isLoading = false;
      finishPageLoading();
      return;
    }

    isLoading = true;
    (async () => {
      try {
        const response = await fetch(apiPath(`${categoryKey}/${entryId}`));
        const nextEntry = (await response.json()) as entryType;
        nextEntry.Tags = nextEntry.Tags ?? [];
        entry = nextEntry;
        editorValue = nextEntry.Value ?? "";
        initialSnapshot = snapshotEntry(nextEntry, editorValue);
        canCheckDirty = true;
      } finally {
        isLoading = false;
        finishPageLoading();
      }
    })();
  });

  $effect(() => {
    initialized;
    outlineInput;

    if (!initialized || outlineInput === null) {
      return;
    }

    requestAnimationFrame(() => {
      outlineInput?.focus();
    });
  });
</script>

<svelte:window bind:innerHeight bind:innerWidth />

{#if isErr && errMessage != ""}
  <div class="notification is-danger m-3">{errMessage}</div>
{/if}

<div class="detail-page" class:is-page-loading={isPageLoading}>
  {#if isPageLoading}
    <div class="detail-loading">
      <button class="button is-dark is-loading" aria-label="loading"></button>
    </div>
  {/if}

  <header>
    <div class="detail-header">
      <div class="detail-header-main">
        <div class="detail-header-top">
          <div class="field detail-date-field">
            {#if isNew}
              <label class="label" for="dateInput"
                >{resolvedSettings.dateLabel}</label
              >
            {:else}
              <div class="label">{resolvedSettings.dateLabel}</div>
            {/if}
            {#if isNew}
              <div class="control">
                <input
                  id="dateInput"
                  type="date"
                  class="input"
                  class:is-fullwidth={isNew}
                  bind:value={entry.Date}
                />
              </div>
            {:else}
              <div class="detail-date-label">{entry.Date}</div>
            {/if}
          </div>
          {#if !isNew}
            <div class="detail-action-wrap">
              <button
                class="button has-text-danger is-light detail-action-button"
                aria-label={`delete ${categoryKey}`}
                onclick={onDelete}
              >
                <span class="icon"><i class="fa-solid fa-trash"></i></span>
                <span>削除</span>
              </button>
            </div>
          {/if}
        </div>

        <div class="field">
          <label class="label" for="outlineInput"
            >{resolvedSettings.outlineLabel}</label
          >
          <div class="field has-addons mobile-outline-row">
            <div class="control is-expanded">
              <input
                id="outlineInput"
                type="text"
                placeholder={resolvedSettings.outlineLabel}
                class="input is-medium"
                bind:this={outlineInput}
                bind:value={entry.Outline}
              />
            </div>
            {#if resolvedSettings.showTags}
              <div class="control is-hidden-tablet">
                <button
                  class="button is-light mobile-tag-button"
                  class:is-link={showTagsOnMobile}
                  type="button"
                  aria-label="toggle tags"
                  onclick={() => (showTagsOnMobile = !showTagsOnMobile)}
                >
                  <span class="icon"><i class="fa-solid fa-tags"></i></span>
                </button>
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>
    {#if resolvedSettings.showTags}
      <div class="detail-tags">
        <div class:mobile-hidden-tags={!showTagsOnMobile}>
          <div class="label">{resolvedSettings.tagsLabel}</div>
          <div class="control">
            <TagsInput bind:items={entry.Tags} />
          </div>
        </div>
      </div>
    {/if}
  </header>

  <section class="detail-body p-0">
    {#if showTemplateSelector}
      <div class="template-selector-modal">
        <div class="template-selector">
          <div class="template-selector-header">
            <h2 class="title is-6 mb-2">テンプレート選択</h2>
          </div>
          <div class="template-grid">
            {#each templateOptions as template, index}
              <button
                class="template-card"
                class:is-selected={selectedTemplateIndex === index}
                type="button"
                onclick={() => onSelectTemplate(template, index)}
              >
                <span class="template-card-title">{template.Name}</span>
                <span class="template-card-preview">
                  {#if templatePreview(template.Value).length > 0}
                    {templatePreview(template.Value).join(" / ")}
                  {:else}
                    空の本文で開始します
                  {/if}
                </span>
              </button>
            {/each}
          </div>
        </div>
      </div>
    {/if}

    <div class="field">
      <div class="control py-2">
        <RichInput value={editorValue} {imageUploadPath} {onTextChange} />
      </div>
    </div>
  </section>
  <footer class="is-dark m-0">
    <div class="footer-actions">
      <div class="footer-status" class:is-visible={isDirty}>
        {#if isDirty}
          <span class="tag unsaved-tag">未保存の変更あり</span>
        {/if}
      </div>
      <div class="footer-buttons">
        <button class="button is-light footer-button" onclick={onCancel}>
          <span class="icon"><i class="fa-solid fa-arrow-left"></i></span>
          <span>戻る</span>
        </button>
        <button
          class="button is-light footer-button footer-template-button"
          onclick={onSaveTemplate}
        >
          <span class="icon"><i class="fa-solid fa-book"></i></span>
          <span>テンプレート</span>
        </button>
        <button
          class="button footer-button"
          class:is-primary={isDirty}
          class:is-light={!isDirty}
          class:is-disabled-look={!isDirty}
          disabled={isLoading || !isDirty}
          class:is-loading={isLoading}
          onclick={onOk}
        >
          <span class="icon"><i class="fa-solid fa-cloud-arrow-up"></i></span>
          <span>{saveButtonLabel}</span>
        </button>
      </div>
    </div>
  </footer>
</div>

<style>
  .detail-page {
    position: relative;
  }

  .detail-body {
    position: relative;
    min-height: calc(100vh - 15rem);
  }

  .detail-page.is-page-loading {
    visibility: hidden;
  }

  .detail-loading {
    position: fixed;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: color-mix(in srgb, var(--bulma-scheme-main) 86%, transparent);
    z-index: 1200;
    visibility: visible;
  }

  header {
    background-color: var(--bulma-border);
    width: 100%;
    top: 0;
    left: 0;
    margin: 0;
    z-index: 999;
    padding: 1rem;
  }

  footer {
    background-color: var(--bulma-border);
    left: 0;
    bottom: 0;
    width: 100%;
    position: fixed;
  }

  .detail-header {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    align-items: flex-start;
  }

  .detail-header-main {
    flex: 1;
  }

  .detail-header-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    margin-bottom: 0.75rem;
    flex-wrap: wrap;
  }

  .detail-date-field {
    margin-bottom: 0;
  }

  .detail-date-field .control,
  .detail-action-wrap {
    display: flex;
    align-items: center;
    margin-bottom: 0.75rem;
  }

  .detail-date-field .control {
    width: 100%;
  }

  .detail-date-label {
    display: inline-flex;
    align-items: center;
    min-height: 2.65rem;
    color: var(--bulma-text);
    font-variant-numeric: tabular-nums;
    font-size: 1.45rem;
    font-weight: 700;
    letter-spacing: 0.02em;
    line-height: 1;
  }

  .detail-tags {
    margin-top: 0.75rem;
  }

  .mobile-outline-row {
    margin-bottom: 0;
  }

  .mobile-tag-button,
  .detail-action-button {
    height: 2.65rem;
  }

  .detail-action-button {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: 2.65rem;
  }

  .detail-tags .label,
  .detail-date-field .label,
  .detail-header .label {
    margin-bottom: 0.35rem;
  }

  .footer-actions {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
  }

  .template-selector {
    width: min(100%, 34rem);
    max-height: min(70vh, 38rem);
    overflow-y: auto;
    padding: 1rem;
    border: 1px solid
      color-mix(in srgb, var(--bulma-link) 18%, var(--bulma-border));
    border-radius: 1.1rem;
    background: var(--bulma-scheme-main);
    box-shadow:
      0 1.2rem 3rem rgba(15, 23, 42, 0.18),
      0 0 0 1px color-mix(in srgb, var(--bulma-link) 12%, transparent);
  }

  .template-selector-modal {
    position: absolute;
    inset: 0;
    z-index: 1100;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    background: color-mix(
      in srgb,
      var(--bulma-scheme-main) 55%,
      rgba(15, 23, 42, 0.45)
    );
    backdrop-filter: blur(6px);
  }

  .template-selector-header {
    margin-bottom: 0.75rem;
  }

  .template-grid {
    display: grid;
    gap: 0.75rem;
  }

  .template-card {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.45rem;
    width: 100%;
    padding: 1rem;
    border: 1px solid var(--bulma-border);
    border-radius: 0.9rem;
    background: color-mix(
      in srgb,
      var(--bulma-link) 10%,
      var(--bulma-scheme-main)
    );
    color: var(--bulma-text);
    text-align: left;
    transition:
      border-color 0.15s ease,
      transform 0.15s ease,
      box-shadow 0.15s ease;
  }

  .template-card:hover {
    border-color: var(--bulma-link);
    transform: translateY(-1px);
  }

  .template-card.is-selected {
    border-color: var(--bulma-link);
    background: color-mix(
      in srgb,
      var(--bulma-link) 18%,
      var(--bulma-scheme-main)
    );
    box-shadow: 0 0 0 1px color-mix(in srgb, var(--bulma-link) 45%, transparent);
  }

  .template-card-title {
    font-size: 1rem;
    font-weight: 700;
  }

  .template-card-preview {
    color: var(--bulma-text-weak);
    line-height: 1.5;
  }

  .footer-status {
    display: flex;
    align-items: center;
  }

  .footer-status.is-visible {
    min-height: 2rem;
  }

  .unsaved-tag {
    background: transparent;
    border: none;
    box-shadow: none;
    color: #b15b12;
    padding-left: 0;
    padding-right: 0;
    user-select: none;
  }

  .footer-buttons {
    display: flex;
    gap: 0.75rem;
    margin-left: auto;
  }

  .footer-button {
    min-width: 7rem;
  }

  .footer-template-button {
    font-size: 0.8rem;
  }

  .footer-button.is-disabled-look .icon {
    color: color-mix(in srgb, var(--bulma-text-weak) 68%, black 32%);
  }

  #dateInput {
    width: 150px;
    height: 2.65rem;
  }

  @media screen and (max-width: 768px) {
    header {
      padding: 0.65rem 0.75rem;
    }

    .detail-tags .label,
    .detail-date-field .label,
    .detail-header .label {
      display: none;
    }

    .detail-header-top {
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto;
      align-items: stretch;
      gap: 0.5rem;
      margin-bottom: 0.5rem;
    }

    .detail-header .field {
      margin-bottom: 0.35rem;
    }

    .detail-date-field .control,
    .detail-action-wrap {
      margin-bottom: 0.5rem;
    }

    .detail-tags {
      margin-top: 0.35rem;
    }

    .mobile-hidden-tags {
      display: none;
    }

    .detail-header .input.is-medium {
      font-size: 1rem;
      padding-top: 0.55rem;
      padding-bottom: 0.55rem;
    }

    .mobile-outline-row {
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto;
      gap: 0.5rem;
      align-items: stretch;
    }

    .mobile-outline-row .control {
      min-width: 0;
      display: flex;
      align-items: stretch;
    }

    .mobile-tag-button {
      min-width: 2.75rem;
      padding-left: 0.65rem;
      padding-right: 0.65rem;
      height: 100%;
    }

    #dateInput {
      width: 100%;
      min-width: 9.5rem;
    }

    .detail-header-top :global(.button) {
      white-space: nowrap;
      height: 2.65rem;
    }

    .detail-date-field,
    .detail-action-wrap,
    .mobile-outline-row .control {
      display: flex;
      align-items: stretch;
      height: 2.65rem;
    }

    .detail-date-label {
      width: 100%;
      min-width: 9.5rem;
      font-size: 1.2rem;
    }

    .detail-date-field .control,
    .detail-action-wrap {
      margin-bottom: 0;
    }

    #dateInput,
    .detail-action-button,
    .mobile-outline-row :global(.input),
    .mobile-tag-button {
      height: 100%;
    }

    .footer-actions {
      flex-wrap: wrap;
      gap: 0.4rem;
      padding: 0.45rem 0.75rem 0.6rem;
    }

    .detail-body {
      min-height: calc(100vh - 13rem);
    }

    .template-selector {
      padding: 0.75rem;
    }

    .footer-status,
    .footer-buttons {
      width: 100%;
    }

    .footer-status {
      min-height: 0;
    }

    .footer-buttons {
      justify-content: space-between;
      gap: 0.5rem;
    }

    .footer-button {
      flex: 1 1 0;
      min-width: 0;
    }

    .footer-status.is-visible {
      min-height: 1.2rem;
    }

    .unsaved-tag {
      font-size: 0.8rem;
      line-height: 1.1;
    }
  }
</style>
