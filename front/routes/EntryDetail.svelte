<script lang="ts">
  import { useNavigate, useRoute } from "@dvcol/svelte-simple-router/router";
  import {
    MuDangerButton,
    MuDateField,
    MuActionsFooter,
    MuPrimaryButton,
    MuSecondaryButton,
    MuTagsInput,
  } from "mu-ui-lib";
  import { onMount } from "svelte";
  import AppIcon from "../components/AppIcon.svelte";
  import ImageSelectionCard from "../components/ImageSelectionCard.svelte";
  import MDInput from "../components/MDInput/index.svelte";
  import TemplateSelectorModal from "../components/TemplateSelectorModal.svelte";
  import type { entryType } from "../models/entryModels.js";
  import type { TemplateType } from "../models/settingType.js";
  import { settingsStore } from "../store.js";
  import { apiPath } from "../basePath.js";
  import { createEmptyEntry, loadEntry } from "../lib/entryApi.js";

  type EntryRouteParams = {
    id?: string | number | boolean;
    category?: string | number | boolean;
  };
  type EntryRouteQuery = Record<string, string | number | boolean | undefined>;

  const carryOverMarker = "----ここまで前回内容で置換";

  const { location } = $derived(useRoute());
  const { push } = useNavigate();

  const routeParams = $derived((location?.params ?? {}) as EntryRouteParams);
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

  let entry = $state<entryType>(createEmptyEntry());
  let isErr = $state(false);
  let errMessage = $state("");
  let isLoading = $state(false);
  let editorValue = $state("");
  let outlineValue = $state("");
  let dateValue = $state("");
  let tagsValue = $state<string[]>([]);
  let isNew = $state(false);
  let initialized = $state(false);
  let initialSnapshot = $state("");
  let canCheckDirty = $state(false);
  let showTagsOnMobile = $state(false);
  let outlineInput = $state<HTMLInputElement | null>(null);
  let isPageLoading = $state(true);
  let selectedTemplateIndex = $state<number | null>(null);
  let isTemplateSelectorOpen = $state(false);
  let isImagePickerOpen = $state(false);
  let previousCardEntryValue = $state("");
  let lastInitKey = "";
  let loadSequence = 0;

  document.querySelector<HTMLDivElement>(".navbar")?.classList.add("is-hidden");

  const entryId = $derived(routeParams.id ? String(routeParams.id) : undefined);
  const isAddRoute = $derived(entryId == undefined || entryId === "");
  const routeQuery = $derived(
    (location?.query ?? {}) as EntryRouteQuery,
  );
  const imageUploadPath = $derived(
    isAddRoute
      ? apiPath(`${categoryKey}/images/tmp`)
      : apiPath(`${categoryKey}/${entryId}/images`),
  );
  const previewImages = $derived.by(() =>
    (entry.Images ?? []).map((image) => ({
      id: image.Id,
      src: image.Src,
      alt: image.Alt,
      markdown: image.Markdown,
    })),
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
        initialSnapshot = snapshotEntry(
          outlineValue,
          dateValue,
          tagsValue,
          value,
        );
      });
    });
  };
  const snapshotEntry = (
    outline: string,
    date: string,
    tags: string[],
    value: string,
  ) =>
    JSON.stringify({
      Outline: outline ?? "",
      Date: date ?? "",
      Value: value ?? "",
      Tags: tags ?? [],
    });
  const isDirty = $derived(
    canCheckDirty &&
      snapshotEntry(outlineValue, dateValue, tagsValue, editorValue) !==
        initialSnapshot,
  );
  const saveButtonLabel = $derived(isNew ? "作成" : "保存");
  const templateOptions = $derived.by(() => {
    const options = [...resolvedSettings.templates];
    if (!isNew) {
      return options;
    }

    return [
      { Name: "空白から作成", Value: "", Tags: [] } satisfies TemplateType,
      ...options,
    ];
  });
  const showTemplateSelector = $derived(isNew && isTemplateSelectorOpen);
  const initKey = $derived.by(() => {
    const templateSignature = resolvedSettings.templates
      .map(
        (template) =>
          `${template.Name}:${template.Value}:${(template.Tags ?? []).join("#")}`,
      )
      .join("|");
    const querySignature = JSON.stringify(routeQuery);
    return `${initialized ? "1" : "0"}:${categoryKey}:${entryId ?? ""}:${isAddRoute ? "1" : "0"}:${templateSignature}:${querySignature}`;
  });
  const entryUrl = () => {
    if (isNew) {
      return apiPath(categoryKey);
    }
    return apiPath(`${categoryKey}/${entryId}`);
  };

  const listPath = () => `/${categoryKey}/`;
  const referencePath = () => `/${categoryKey}/${entryId}`;
  const queryValue = (key: string) => {
    const value = routeQuery[key];
    if (typeof value !== "string") {
      return undefined;
    }

    try {
      return decodeURIComponent(value);
    } catch {
      return value;
    }
  };
  const hasQueryValue = (key: string) =>
    Object.prototype.hasOwnProperty.call(routeQuery, key);
  const parsePreviousEntryId = () => {
    const value = queryValue("previousEntryId");
    if (!value) {
      return undefined;
    }

    const parsed = Number(value);
    return Number.isInteger(parsed) && parsed > 0 ? parsed : undefined;
  };
  const applyTemplateValue = (templateValue: string) => {
    const lines = templateValue.split(/\r?\n/);
    const markerIndex = lines.findIndex(
      (line) => line.trim() === carryOverMarker,
    );
    if (markerIndex < 0) {
      return templateValue;
    }

    const templatePrefix = lines.slice(0, markerIndex + 1).join("\n");
    const templateSuffix = lines.slice(markerIndex + 1).join("\n");

    let carriedPrefix = templatePrefix;
    if (previousCardEntryValue.trim() !== "") {
      const previousLines = previousCardEntryValue.split(/\r?\n/);
      const previousMarkerIndex = previousLines.findIndex(
        (line) => line.trim() === carryOverMarker,
      );
      if (previousMarkerIndex >= 0) {
        carriedPrefix = previousLines
          .slice(0, previousMarkerIndex + 1)
          .join("\n");
      }
    }

    if (carriedPrefix !== "" && templateSuffix !== "") {
      return `${carriedPrefix.replace(/\n+$/, "")}\n${templateSuffix.replace(/^\n+/, "")}`;
    }
    return carriedPrefix || templateSuffix;
  };

  const goToList = async () => push({ path: listPath() });
  const goToReference = async () => {
    if (isNew || !entryId) {
      await goToList();
      return;
    }
    await push({ path: referencePath() });
  };
  const shouldLeave = () =>
    !isDirty || window.confirm("未保存の変更があります。戻りますか？");

  const onOk = async () => {
    const nextEntry = {
      ...entry,
      Outline: outlineValue,
      Date: dateValue,
      Tags: [...tagsValue],
      Value: editorValue,
    } satisfies entryType;

    if (!resolvedSettings.allowMultipleEntriesPerDate && nextEntry.Date) {
      const month = nextEntry.Date.slice(0, 7);
      const response = await fetch(`${apiPath(categoryKey)}?month=${month}`);
      const monthlyEntries = (await response.json()) as entryType[];
      const hasDuplicateDate = monthlyEntries.some(
        (item) => item.Date === nextEntry.Date && item.Id !== nextEntry.Id,
      );
      if (hasDuplicateDate) {
        errMessage = "同じ日付のエントリは登録できません";
        isErr = true;
        return;
      }
    }

    const response = await fetch(entryUrl(), {
      method: saveMethod(),
      body: JSON.stringify(nextEntry),
    });

    if (response.status !== 200) {
      errMessage = (await response.json()).error;
      isErr = true;
      return;
    }

    entry = nextEntry;
    initialSnapshot = snapshotEntry(
      outlineValue,
      dateValue,
      tagsValue,
      editorValue,
    );
    canCheckDirty = true;
    await goToReference();
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
  };
  const onEmbedImage = (markdown: string) => {
    const nextMarkdown = markdown.trim();
    if (nextMarkdown === "") {
      return;
    }

    const separator = editorValue.trim() === "" ? "" : "\n";
    editorValue = `${editorValue}${separator}${nextMarkdown}`;
  };
  const openImagePicker = () => {
    isImagePickerOpen = true;
  };
  const closeImagePicker = () => {
    isImagePickerOpen = false;
  };
  const onSelectTemplate = (payload: {
    template: TemplateType;
    index: number;
  }) => {
    const { template, index } = payload;
    selectedTemplateIndex = index;
    isTemplateSelectorOpen = false;
    tagsValue = [...(template.Tags ?? [])];
    applyInitialEditorValue(applyTemplateValue(template.Value ?? ""));
  };
  onMount(() => {
    initialized = true;
  });

  const resetDetailState = (nextIsNew: boolean) => {
    isPageLoading = true;
    isErr = false;
    errMessage = "";
    entry = createEmptyEntry();
    isNew = nextIsNew;
    outlineValue = "";
    dateValue = createEmptyEntry().Date;
    tagsValue = [];
    showTagsOnMobile = false;
    selectedTemplateIndex = null;
    isTemplateSelectorOpen = false;
    previousCardEntryValue = "";
    canCheckDirty = true;
  };

  const initializeDetail = async () => {
    if (!initialized || !categoryKey) {
      return;
    }

    const nextIsNew = isAddRoute;
    const currentLoadSequence = ++loadSequence;
    const initialTemplateValue = resolvedSettings.templates[0]?.Value ?? "";
    const initialTemplateTags = resolvedSettings.templates[0]?.Tags ?? [];
    const shouldAutoApplyTemplate =
      resolvedSettings.templates.length <= 1 && initialTemplateValue !== "";
    const presetOutline = queryValue("presetOutline") ?? "";
    const presetTag = queryValue("presetTag") ?? "";
    const hasPresetOutline = hasQueryValue("presetOutline");
    const hasPresetTag = hasQueryValue("presetTag");
    const previousEntryId = parsePreviousEntryId();

    resetDetailState(nextIsNew);

    if (nextIsNew) {
      if (previousEntryId != undefined) {
        try {
          const previousEntry = await loadEntry(
            categoryKey,
            String(previousEntryId),
          );
          if (currentLoadSequence === loadSequence) {
            previousCardEntryValue = previousEntry.Value ?? "";
          }
        } catch {
          if (currentLoadSequence === loadSequence) {
            previousCardEntryValue = "";
          }
        }
      }

      const nextEntry = createEmptyEntry();
      const nextOutline = hasPresetOutline ? presetOutline : nextEntry.Outline;
      const nextDate = nextEntry.Date;
      const nextTags = hasPresetTag
        ? presetTag
          ? [presetTag]
          : []
        : [...initialTemplateTags];
      entry = nextEntry;
      outlineValue = nextOutline;
      dateValue = nextDate;
      editorValue = "";
      tagsValue = nextTags;
      initialSnapshot = snapshotEntry(nextOutline, nextDate, nextTags, "");
      isLoading = false;
      finishPageLoading();
      if (resolvedSettings.templates.length > 1) {
        isTemplateSelectorOpen = true;
      } else if (shouldAutoApplyTemplate) {
        applyInitialEditorValue(applyTemplateValue(initialTemplateValue));
      }
      return;
    }

    if (entryId == undefined || entryId === "") {
      const nextOutline = entry.Outline ?? "";
      const nextDate = entry.Date ?? "";
      const nextTags = [...(entry.Tags ?? [])];
      outlineValue = nextOutline;
      dateValue = nextDate;
      tagsValue = nextTags;
      editorValue = "";
      initialSnapshot = snapshotEntry(nextOutline, nextDate, nextTags, "");
      canCheckDirty = false;
      isLoading = false;
      finishPageLoading();
      return;
    }

    isLoading = true;
    try {
      const nextEntry = await loadEntry(categoryKey, entryId);
      if (currentLoadSequence !== loadSequence) {
        return;
      }
      const nextOutline = nextEntry.Outline ?? "";
      const nextDate = nextEntry.Date ?? "";
      const nextTags = [...(nextEntry.Tags ?? [])];
      const nextEditorValue = nextEntry.Value ?? "";
      entry = nextEntry;
      outlineValue = nextOutline;
      dateValue = nextDate;
      tagsValue = nextTags;
      editorValue = nextEditorValue;
      initialSnapshot = snapshotEntry(
        nextOutline,
        nextDate,
        nextTags,
        nextEditorValue,
      );
      canCheckDirty = true;
    } finally {
      if (currentLoadSequence === loadSequence) {
        isLoading = false;
        finishPageLoading();
      }
    }
  };

  $effect(() => {
    const key = initKey;
    if (!initialized || !categoryKey || key === lastInitKey) {
      return;
    }
    lastInitKey = key;
    void initializeDetail();
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
            <MuDateField bind:value={dateValue} editable={isNew} />
          </div>
          <div class="field detail-outline-field">
            <div class="field has-addons mobile-outline-row">
              <div class="control is-expanded">
                <input
                  id="outlineInput"
                  type="text"
                  placeholder={resolvedSettings.outlineLabel}
                  class="input is-medium"
                  bind:this={outlineInput}
                  bind:value={outlineValue}
                />
              </div>
              {#if resolvedSettings.showTags}
                <div class="control is-hidden-tablet">
                  <button
                    class="button mobile-tag-button detail-secondary-button"
                    class:is-link={showTagsOnMobile}
                    type="button"
                    aria-label="toggle tags"
                    onclick={() => (showTagsOnMobile = !showTagsOnMobile)}
                  >
                    <span class="icon"><AppIcon name="tags" /></span>
                  </button>
                </div>
              {/if}
            </div>
          </div>
          {#if !isNew}
            <div class="detail-action-wrap">
              <MuDangerButton
                ariaLabel={`delete ${categoryKey}`}
                onclick={onDelete}
              >
                <span class="icon"><AppIcon name="trash" /></span>
                <span>削除</span>
              </MuDangerButton>
            </div>
          {/if}
        </div>
      </div>
    </div>
    {#if resolvedSettings.showTags}
      <div class="detail-tags">
        <div class:mobile-hidden-tags={!showTagsOnMobile}>
          <div class="control">
            <MuTagsInput bind:items={tagsValue} />
          </div>
        </div>
      </div>
    {/if}
  </header>

  <section class="detail-body p-0">
    <TemplateSelectorModal
      isOpen={showTemplateSelector}
      templates={templateOptions}
      initialSelectedIndex={selectedTemplateIndex}
      onSelect={onSelectTemplate}
    />

    <div class="field detail-editor-field">
      <div class="control detail-editor-control">
        <MDInput value={editorValue} {imageUploadPath} {onTextChange} />
      </div>
    </div>

  </section>

  {#if isImagePickerOpen}
    <div
      class="image-picker-modal-backdrop"
      role="button"
      tabindex="0"
      onclick={(event) => {
        if (event.target === event.currentTarget) {
          closeImagePicker();
        }
      }}
      onkeydown={(event) => {
        if (event.target !== event.currentTarget) {
          return;
        }
        if (event.key === "Escape" || event.key === "Enter" || event.key === " ") {
          event.preventDefault();
          closeImagePicker();
        }
      }}
    >
      <div class="image-picker-modal">
        <div class="image-picker-modal-header">
          <div>
            <div class="image-picker-modal-title">画像選択</div>
            <div class="image-picker-modal-subtitle">画像をクリックして本文に挿入できます</div>
          </div>
          <MuSecondaryButton
            ariaLabel="close image picker"
            onclick={closeImagePicker}
          >
            <span class="icon"><AppIcon name="xmark" /></span>
          </MuSecondaryButton>
        </div>
        <ImageSelectionCard
          onEmbedImage={(markdown) => {
            onEmbedImage(markdown);
            closeImagePicker();
          }}
          {imageUploadPath}
          images={previewImages}
        />
      </div>
    </div>
  {/if}

  <MuActionsFooter hasUnsavedChanges={isDirty}>
    <MuSecondaryButton onclick={onCancel}>
      <span class="icon"><AppIcon name="arrow-left" /></span>
      <span>戻る</span>
    </MuSecondaryButton>
    <MuSecondaryButton onclick={openImagePicker}>
      <span class="icon"><AppIcon name="image" /></span>
      <span>画像</span>
    </MuSecondaryButton>
    {#if isDirty}
      <MuPrimaryButton disabled={isLoading} onclick={onOk}>
        <span class="icon"><AppIcon name="cloud-arrow-up" /></span>
        <span>{saveButtonLabel}</span>
      </MuPrimaryButton>
    {:else}
      <MuSecondaryButton disabled={true} onclick={onOk}>
        <span class="icon"><AppIcon name="cloud-arrow-up" /></span>
        <span>{saveButtonLabel}</span>
      </MuSecondaryButton>
    {/if}
  </MuActionsFooter>
</div>

<style>
  .detail-page {
    position: relative;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }

  .detail-body {
    position: relative;
    flex: 1 1 auto;
    min-height: 0;
    margin-bottom: 4.25rem;
    display: flex;
    flex-direction: column;
  }

  .detail-editor-field {
    margin-bottom: 0;
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
    padding: 0.7rem 0.8rem;
  }

  .detail-header {
    display: flex;
    justify-content: space-between;
    gap: 0.25rem;
    align-items: flex-start;
  }

  .detail-header-main {
    flex: 1;
  }

  .detail-header-top {
    display: flex;
    justify-content: flex-start;
    align-items: center;
    gap: 0.25rem;
    margin-bottom: 0.15rem;
    flex-wrap: nowrap;
  }

  .detail-date-field {
    margin-bottom: 0;
    flex: 0 0 auto;
  }

  .detail-outline-field {
    flex: 1 1 auto;
    min-width: 0;
    margin-bottom: 0;
  }

  .detail-action-wrap {
    display: flex;
    align-items: center;
    margin-bottom: 0;
  }

  .detail-tags {
    margin-top: 0.15rem;
  }

  .mobile-outline-row {
    margin-bottom: 0;
    width: 100%;
  }

  .mobile-tag-button {
    height: 2.2rem;
  }

  .detail-secondary-button {
    background: color-mix(
      in srgb,
      var(--bulma-border) 74%,
      var(--bulma-scheme-main)
    );
    border: 1px solid color-mix(in srgb, var(--bulma-border) 86%, white 14%);
    color: color-mix(in srgb, var(--bulma-text) 90%, white 10%);
    box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.04);
  }

  .detail-body .field {
    height: 100%;
    margin-bottom: 0;
  }

  .detail-editor-control {
    display: flex;
    flex: 1 1 auto;
    min-height: 0;
    align-items: flex-start;
  }

  .detail-editor-control :global(.md-input) {
    width: 100%;
    height: 100%;
  }

  .image-picker-modal-backdrop {
    position: fixed;
    inset: 0;
    z-index: 1280;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0.75rem;
    background: rgba(6, 10, 18, 0.82);
    backdrop-filter: blur(8px);
  }

  .image-picker-modal {
    width: min(72rem, 100%);
    max-height: min(90vh, 56rem);
    overflow: hidden;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    padding: 1rem;
    border-radius: 1rem;
    background: var(--bulma-scheme-main);
    box-shadow: 0 1.5rem 3.5rem rgba(0, 0, 0, 0.28);
  }

  .image-picker-modal-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 0.75rem;
  }

  .image-picker-modal-title {
    color: var(--bulma-text);
    font-size: 1rem;
    font-weight: 700;
    line-height: 1.2;
  }

  .image-picker-modal-subtitle {
    margin-top: 0.15rem;
    color: var(--bulma-text-weak);
    font-size: 0.82rem;
    line-height: 1.4;
  }

  .image-picker-modal :global(.image-card) {
    margin-top: 0;
    border-radius: 0.9rem;
  }

  @media screen and (max-width: 768px) {
    header {
      padding: 0.4rem 0.55rem;
    }

    .detail-header-top {
      display: flex;
      align-items: stretch;
      gap: 0.2rem;
      margin-bottom: 0.15rem;
      flex-wrap: wrap;
    }

    .detail-header .field {
      margin-bottom: 0.05rem;
    }

    .detail-action-wrap {
      margin-bottom: 0;
    }

    .detail-tags {
      margin-top: 0.1rem;
    }

    .mobile-hidden-tags {
      display: none;
    }

    .detail-header .input.is-medium {
      font-size: 1rem;
      padding-top: 0.35rem;
      padding-bottom: 0.35rem;
    }

    .mobile-outline-row {
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto;
      gap: 0.2rem;
      align-items: stretch;
    }

    .mobile-outline-row .control {
      min-width: 0;
      display: flex;
      align-items: stretch;
    }

    .mobile-tag-button {
      min-width: 2.5rem;
      padding-left: 0.5rem;
      padding-right: 0.5rem;
      height: 100%;
    }

    .detail-header-top :global(.button) {
      white-space: nowrap;
      height: 2.35rem;
    }

    .detail-date-field,
    .detail-action-wrap,
    .mobile-outline-row .control {
      display: flex;
      align-items: stretch;
      height: 2.35rem;
    }

    .detail-outline-field {
      width: 100%;
    }

    .detail-action-wrap {
      margin-bottom: 0;
    }

    .mobile-outline-row :global(.input),
    .mobile-tag-button {
      height: 100%;
    }

    .detail-body {
      flex: 1 1 auto;
      min-height: 0;
      margin-bottom: 4.25rem;
    }

    .detail-editor-control {
      flex: 1 1 auto;
      min-height: 0;
    }

    .image-picker-modal {
      width: calc(100vw - 1rem);
      max-height: calc(100vh - 1rem);
      padding: 0.75rem;
      border-radius: 0.85rem;
    }

    .image-picker-modal-header {
      gap: 0.5rem;
    }

    .image-picker-modal :global(.image-card-strip) {
      display: flex;
      flex-direction: column;
      grid-auto-flow: initial;
      grid-auto-columns: unset;
      overflow-x: hidden;
      overflow-y: auto;
      max-height: calc(100vh - 9rem);
      padding-right: 0.15rem;
    }

    .image-picker-modal :global(.image-card-item),
    .image-picker-modal :global(.image-card-add) {
      min-height: 7rem;
    }

    .image-picker-modal :global(.image-card-item) {
      padding: 0;
      overflow: hidden;
    }

    .image-picker-modal :global(.image-card-preview) {
      border-radius: 0;
    }

    .image-picker-modal :global(.image-card-preview img) {
      width: 100%;
      height: auto;
      aspect-ratio: 4 / 3;
    }

    .image-picker-modal :global(.image-card-item-footer) {
      padding: 0.65rem 0.75rem 0.75rem;
    }

  }
</style>
