<script lang="ts">
  import { goto, type RouteResult } from "@mateothegreat/svelte5-router";
  import { onMount } from "svelte";
  import MDInput from "../components/MDInput/index.svelte";
  import TagsInput from "../components/TagsInput.svelte";
  import TemplateSaveModal from "../components/TemplateSaveModal.svelte";
  import TemplateSelectorModal from "../components/TemplateSelectorModal.svelte";
  import type { entryType } from "../models/entryModels.js";
  import type { TemplateType } from "../models/settingType.js";
  import { settingsStore } from "../store.js";
  import { apiPath, appPath } from "../basePath.js";
  import { createEmptyEntry, loadEntry } from "../lib/entryApi.js";

  interface Props {
    route?: RouteResult;
  }

  type EntryRouteParams = {
    id?: string | number | boolean;
    category?: string | number | boolean;
  };
  type EntryRouteQuery = Record<string, string | number | boolean | undefined>;

  const carryOverMarker = "----ここまで前回内容で置換";

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
  let isTemplateSaveModalOpen = $state(false);
  let templateSaveInitialMode = $state<"new" | "overwrite">("new");
  let templateSaveInitialDraftName = $state("");
  let templateSaveInitialOverwriteName = $state("");
  let templateModalError = $state("");
  let isTemplateSaving = $state(false);
  let isTemplateDeleting = $state(false);
  let previousCardEntryValue = $state("");
  let lastInitKey = "";
  let loadSequence = 0;

  document.querySelector<HTMLDivElement>(".navbar")?.classList.add("is-hidden");

  const entryId = $derived(routeParams.id ? String(routeParams.id) : undefined);
  const isAddRoute = $derived(entryId == undefined || entryId === "");
  const routeQuery = $derived(
    (currentRoute?.result?.querystring?.params ??
      currentRoute?.result?.querystring?.original ??
      {}) as EntryRouteQuery,
  );
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
        initialSnapshot = snapshotEntry(outlineValue, dateValue, tagsValue, value);
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
  const selectableTemplateCount = $derived(resolvedSettings.templates.length);
  const showTemplateSelector = $derived(isNew && isTemplateSelectorOpen);
  const hasTemplates = $derived(resolvedSettings.templates.length > 0);
  const initKey = $derived.by(() => {
    const templateSignature = resolvedSettings.templates
      .map((template) => `${template.Name}:${template.Value}:${(template.Tags ?? []).join("#")}`)
      .join("|");
    const querySignature = JSON.stringify(routeQuery);
    return `${initialized ? "1" : "0"}:${categoryKey}:${entryId ?? ""}:${isAddRoute ? "1" : "0"}:${templateSignature}:${querySignature}`;
  });
  const selectedTemplateName = $derived.by(() => {
    if (selectedTemplateIndex === null) {
      return "";
    }
    const template = templateOptions[selectedTemplateIndex];
    if (template == undefined || template.Name === "空白から作成") {
      return "";
    }
    return template.Name;
  });
  const entryUrl = () => {
    if (isNew) {
      return apiPath(categoryKey);
    }
    return apiPath(`${categoryKey}/${entryId}`);
  };

  const listPath = () => appPath(`/${categoryKey}/`);
  const referencePath = () => appPath(`/${categoryKey}/${entryId}`);
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
    const markerIndex = lines.findIndex((line) => line.trim() === carryOverMarker);
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
        carriedPrefix = previousLines.slice(0, previousMarkerIndex + 1).join("\n");
      }
    }

    if (carriedPrefix !== "" && templateSuffix !== "") {
      return `${carriedPrefix.replace(/\n+$/, "")}\n${templateSuffix.replace(/^\n+/, "")}`;
    }
    return carriedPrefix || templateSuffix;
  };

  const goToList = async () => goto(listPath());
  const goToReference = async () => {
    if (isNew || !entryId) {
      await goToList();
      return;
    }
    await goto(referencePath());
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
    initialSnapshot = snapshotEntry(outlineValue, dateValue, tagsValue, editorValue);
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
  const openTemplateSaveModal = () => {
    const suggestedName = outlineValue.trim() || "新しいテンプレート";
    const defaultOverwriteName =
      selectedTemplateName ||
      resolvedSettings.templates[0]?.Name ||
      "";
    templateSaveInitialMode = hasTemplates ? "overwrite" : "new";
    templateSaveInitialDraftName = suggestedName;
    templateSaveInitialOverwriteName = defaultOverwriteName;
    templateModalError = "";
    isTemplateSaveModalOpen = true;
  };
  const closeTemplateSaveModal = () => {
    if (isTemplateSaving || isTemplateDeleting) {
      return;
    }
    isTemplateSaveModalOpen = false;
    templateModalError = "";
  };
  const onSaveTemplate = async (payload: {
    mode: "new" | "overwrite";
    name: string;
  }) => {
    const name = payload.name.trim();
    if (editorValue.trim() === "") {
      templateModalError = "本文が空のためテンプレート保存できません";
      return;
    }
    if (name === "") {
      templateModalError =
        payload.mode === "overwrite"
          ? "上書きするテンプレートを選択してください"
          : "テンプレート名を入力してください";
      return;
    }

    isTemplateSaving = true;
    templateModalError = "";
    const response = await fetch(apiPath(`${categoryKey}/templates`), {
      method: "post",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ Name: name, Value: editorValue, Tags: tagsValue }),
    });
    isTemplateSaving = false;
    if (response.status !== 200) {
      templateModalError = (await response.json()).error;
      return;
    }

    const nextSettings = await response.json();
    settingsStore.set(nextSettings);
    isErr = false;
    errMessage = "";
    isTemplateSaveModalOpen = false;
    templateModalError = "";
  };
  const onDeleteTemplate = async (name: string) => {
    name = name.trim();
    if (name === "") {
      templateModalError = "削除するテンプレートを選択してください";
      return;
    }
    if (!window.confirm(`テンプレート「${name}」を削除しますか？`)) {
      return;
    }

    isTemplateDeleting = true;
    templateModalError = "";
    const response = await fetch(
      apiPath(`${categoryKey}/templates/${encodeURIComponent(name)}`),
      { method: "delete" },
    );
    isTemplateDeleting = false;
    if (response.status !== 200) {
      templateModalError = (await response.json()).error;
      return;
    }

    const nextSettings = await response.json();
    settingsStore.set(nextSettings);
    isTemplateSaveModalOpen = false;
    templateModalError = "";
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
    isTemplateSaveModalOpen = false;
    templateSaveInitialMode = "new";
    templateSaveInitialDraftName = "";
    templateSaveInitialOverwriteName = "";
    templateModalError = "";
    isTemplateSaving = false;
    isTemplateDeleting = false;
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
          const previousEntry = await loadEntry(categoryKey, String(previousEntryId));
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
      const nextTags = hasPresetTag ? (presetTag ? [presetTag] : []) : [...initialTemplateTags];
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
      initialSnapshot = snapshotEntry(nextOutline, nextDate, nextTags, nextEditorValue);
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
                  bind:value={dateValue}
                />
              </div>
            {:else}
              <div class="detail-date-label">{dateValue}</div>
            {/if}
          </div>
          {#if !isNew}
            <div class="detail-action-wrap">
              <button
                class="button detail-action-button detail-action-button-danger"
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
            <TagsInput bind:items={tagsValue} />
          </div>
        </div>
      </div>
    {/if}
  </header>

  <section class="detail-body p-0">
    <TemplateSaveModal
      isOpen={isTemplateSaveModalOpen}
      initialMode={templateSaveInitialMode}
      initialDraftName={templateSaveInitialDraftName}
      initialOverwriteName={templateSaveInitialOverwriteName}
      templates={resolvedSettings.templates}
      errorMessage={templateModalError}
      isSaving={isTemplateSaving}
      isDeleting={isTemplateDeleting}
      onClose={closeTemplateSaveModal}
      onSave={onSaveTemplate}
      onDelete={onDeleteTemplate}
    />

    <TemplateSelectorModal
      isOpen={showTemplateSelector}
      templates={templateOptions}
      initialSelectedIndex={selectedTemplateIndex}
      onSelect={onSelectTemplate}
    />

    <div class="field">
      <div class="control py-2 detail-editor-control">
        <MDInput value={editorValue} {imageUploadPath} {onTextChange} />
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
        <button class="button footer-button detail-secondary-button" onclick={onCancel}>
          <span class="icon"><i class="fa-solid fa-arrow-left"></i></span>
          <span>戻る</span>
        </button>
        <button
          class="button footer-button footer-template-button detail-secondary-button"
          onclick={openTemplateSaveModal}
        >
          <span class="icon"><i class="fa-solid fa-book"></i></span>
          <span>テンプレート</span>
        </button>
        <button
          class="button footer-button"
          class:detail-primary-button={isDirty}
          class:detail-secondary-button={!isDirty}
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
    padding-bottom: 5.5rem;
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

  .detail-secondary-button {
    background: color-mix(in srgb, var(--bulma-border) 74%, var(--bulma-scheme-main));
    border: 1px solid color-mix(in srgb, var(--bulma-border) 86%, white 14%);
    color: color-mix(in srgb, var(--bulma-text) 90%, white 10%);
    box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.04);
  }

  .detail-primary-button {
    background: color-mix(in srgb, #2d8f86 62%, var(--bulma-scheme-main));
    border: 1px solid color-mix(in srgb, #2d8f86 74%, black 26%);
    color: #edf8f6;
    box-shadow: 0 10px 24px rgba(10, 31, 29, 0.18);
  }

  .detail-action-button-danger {
    background: color-mix(in srgb, #8a4f55 52%, var(--bulma-scheme-main));
    border: 1px solid color-mix(in srgb, #8a4f55 70%, black 30%);
    color: #f8ecee;
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
    border-radius: 0.9rem;
    font-weight: 600;
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

  .detail-body .field {
    height: 100%;
    margin-bottom: 0;
  }

  .detail-editor-control {
    display: flex;
    min-height: max(22rem, calc(100vh - 23rem));
  }

  .detail-editor-control :global(.md-input) {
    width: 100%;
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
      padding-bottom: 6rem;
    }

    .detail-editor-control {
      min-height: max(18rem, calc(100vh - 19rem));
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
