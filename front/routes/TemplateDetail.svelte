<script lang="ts">
  import { useNavigate, useRoute } from "@dvcol/svelte-simple-router/router";
  import {
    MuDangerButton,
    MuActionsFooter,
    MuPrimaryButton,
    MuSecondaryButton,
    MuTagsField,
  } from "mu-ui-lib";
  import AppIcon from "../components/AppIcon.svelte";
  import MuMdField from "../components/mu-md-field/index.svelte";
  import { apiPath, apiSettingsPath, appPath } from "../basePath.js";
  import type { settingType } from "../models/settingType.js";
  import { settingsStore } from "../store.js";

  const { location } = $derived(useRoute());
  const { push } = useNavigate();
  let outlineValue = $state("");
  let tagsValue = $state<string[]>([]);
  let bodyValue = $state("");
  let initialSnapshot = $state("");
  let templateCategoryKey = $state("");
  let originalTemplateName = $state("");
  let loading = $state(true);

  const routeParams = $derived(
    (location?.params ?? {}) as {
      categoryKey?: string | number | boolean;
      name?: string | number | boolean;
    },
  );
  const routeCategoryKey = $derived(
    decodeURIComponent(String(routeParams.categoryKey ?? "")).trim(),
  );
  const templateName = $derived(
    decodeURIComponent(String(routeParams.name ?? "")).trim(),
  );
  const isNewTemplate = $derived(templateName === "new");
  const categories = $derived.by(
    () =>
      $settingsStore.Categories?.map((category) => ({
        key: category.Key,
        name: category.Name,
      })) ?? [],
  );
  const template = $derived.by(() => {
    if (isNewTemplate) {
      return undefined;
    }
    if (routeCategoryKey !== "") {
      const category = $settingsStore.Categories?.find(
        (item) => item.Key === routeCategoryKey,
      );
      const match = category?.Templates?.find(
        (item) => item.Name === templateName,
      );
      if (category !== undefined && match !== undefined) {
        return { categoryKey: category.Key, template: match };
      }
      return undefined;
    }
    for (const category of $settingsStore.Categories ?? []) {
      const match = category.Templates?.find(
        (item) => item.Name === templateName,
      );
      if (match !== undefined) {
        return { categoryKey: category.Key, template: match };
      }
    }
    return undefined;
  });
  const blankCategoryKey = $derived(
    $settingsStore.Categories?.[0]?.Key ?? "idea",
  );

  $effect(() => {
    const currentTemplate = template?.template;
    const nextTemplateCategoryKey =
      template?.categoryKey ??
      (routeCategoryKey !== "" ? routeCategoryKey : isNewTemplate ? blankCategoryKey : "");
    const nextOriginalTemplateName = currentTemplate?.Name ?? "";
    const nextOutlineValue = currentTemplate?.Name ?? "";
    const nextTagsValue = [...(currentTemplate?.Tags ?? [])];
    const nextBodyValue = currentTemplate?.Value ?? "";

    templateCategoryKey = nextTemplateCategoryKey;
    originalTemplateName = nextOriginalTemplateName;
    outlineValue = nextOutlineValue;
    tagsValue = nextTagsValue;
    bodyValue = nextBodyValue;
    initialSnapshot = snapshotTemplate(
      nextTemplateCategoryKey,
      nextOutlineValue,
      nextTagsValue,
      nextBodyValue,
    );
    loading = false;
  });

  const snapshotTemplate = (
    categoryKey: string,
    outline: string,
    tags: string[],
    body: string,
  ) =>
    JSON.stringify({
      Category: categoryKey ?? "",
      Outline: outline ?? "",
      Tags: tags ?? [],
      Value: body ?? "",
    });

  const isDirty = $derived(
    snapshotTemplate(templateCategoryKey, outlineValue, tagsValue, bodyValue) !==
      initialSnapshot,
  );

  const backToSettings = async () => {
    await push({ path: "/settings" });
  };

  const refreshSettings = async () => {
    const response = await fetch(apiSettingsPath());
    if (!response.ok) {
      throw new Error("設定の読み込みに失敗しました");
    }
    settingsStore.set((await response.json()) as settingType);
  };

  const onDelete = async () => {
    if (templateCategoryKey === "" || originalTemplateName.trim() === "") {
      await backToSettings();
      return;
    }
    const response = await fetch(
      apiPath(
        `${templateCategoryKey}/templates/${encodeURIComponent(originalTemplateName.trim())}`,
      ),
      { method: "DELETE" },
    );
    if (!response.ok) {
      return;
    }
    await refreshSettings();
    await backToSettings();
  };

  const onSave = async () => {
    if (templateCategoryKey === "" || outlineValue.trim() === "") {
      return;
    }
    const response = await fetch(apiPath(`${templateCategoryKey}/templates`), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        Name: outlineValue.trim(),
        Value: bodyValue,
        Tags: tagsValue,
      }),
    });
    if (!response.ok) {
      return;
    }
    await refreshSettings();
    originalTemplateName = outlineValue.trim();
    initialSnapshot = snapshotTemplate(
      templateCategoryKey,
      outlineValue,
      tagsValue,
      bodyValue,
    );
  };

  const onTextChange = (value: string) => {
    bodyValue = value;
  };
</script>

<section class="template-detail-page">
  {#if loading}
    <div class="template-detail-loading">読み込み中...</div>
  {:else}
    <header class="template-detail-header">
      <div class="template-detail-header-main">
        <div class="template-detail-header-top">
          {#if isNewTemplate}
            <div class="template-detail-category-select-wrap">
              <label class="template-detail-category-label" for="templateCategorySelect">
                カテゴリ
              </label>
              <select
                id="templateCategorySelect"
                class="input is-medium template-detail-category-select"
                bind:value={templateCategoryKey}
              >
                {#each categories as category}
                  <option value={category.key}>{category.name}</option>
                {/each}
              </select>
            </div>
          {:else}
            <div class="template-detail-category">
              {categories.find((category) => category.key === templateCategoryKey)
                ?.name ?? "不明"}
            </div>
          {/if}
          <MuDangerButton
            ariaLabel={`delete ${templateCategoryKey}`}
            onclick={onDelete}
          >
            <AppIcon class="icon" name="trash" />
            <span>削除</span>
          </MuDangerButton>
        </div>
        <div class="template-detail-header-fields">
          <div class="template-detail-form-field">
            <label class="label" for="templateOutlineInput">見出し</label>
            <input
              id="templateOutlineInput"
              class="input is-medium"
              type="text"
              bind:value={outlineValue}
            />
          </div>

          <div class="template-detail-form-field">
            <label class="label" for="templateTagsInput">タグ</label>
            <MuTagsField inputId="templateTagsInput" bind:items={tagsValue} />
          </div>
        </div>
      </div>
    </header>

    <div class="template-detail-body">
      <div class="template-detail-inner">
        <div class="template-detail-card">
          <div class="template-detail-grid">
            <div class="template-detail-body-field">
              <MuMdField value={bodyValue} {onTextChange} />
            </div>
          </div>
        </div>
      </div>
    </div>

    <MuActionsFooter hasUnsavedChanges={isDirty}>
      <MuSecondaryButton href={appPath("/settings")}>
        <AppIcon class="icon" name="arrow-left" />
        <span>戻る</span>
      </MuSecondaryButton>
      {#if isDirty}
        <MuPrimaryButton onclick={onSave}>
          <AppIcon class="icon" name="cloud-arrow-up" />
          <span>保存</span>
        </MuPrimaryButton>
      {:else}
        <MuSecondaryButton disabled={true} onclick={onSave}>
          <AppIcon class="icon" name="cloud-arrow-up" />
          <span>保存</span>
        </MuSecondaryButton>
      {/if}
    </MuActionsFooter>
  {/if}
</section>

<style>
  .template-detail-page {
    position: relative;
    min-height: 100vh;
  }

  .template-detail-header {
    background-color: var(--app-border);
    width: 100%;
    top: 0;
    left: 0;
    margin: 0;
    z-index: 999;
    padding: 0.6rem 1rem 0.55rem;
  }

  .template-detail-header-main {
    max-width: 100%;
    margin: 0 auto;
  }

  .template-detail-header-top {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    align-items: center;
    margin-bottom: 0.35rem;
  }

  .template-detail-header-fields {
    display: grid;
    gap: 0.45rem;
    margin-top: 0.45rem;
  }

  .template-detail-category {
    color: var(--app-text);
    font-size: 1.45rem;
    font-weight: 700;
    letter-spacing: 0.02em;
    line-height: 1;
  }

  .template-detail-category-select-wrap {
    display: grid;
    gap: 0.25rem;
    min-width: min(24rem, 100%);
  }

  .template-detail-category-label {
    color: var(--app-text);
    font-size: 0.85rem;
    font-weight: 700;
    letter-spacing: 0.04em;
  }

  .template-detail-category-select {
    width: 100%;
  }

  .template-detail-body {
    position: relative;
    min-height: calc(100vh - 13.5rem);
    display: flex;
    flex: 1 1 auto;
    margin-bottom: 4.25rem;
  }

  .template-detail-inner {
    display: flex;
    flex: 1 1 auto;
    min-height: 0;
    width: 100%;
    max-width: none;
    padding-left: 0.25rem;
    padding-right: 0.25rem;
  }

  .template-detail-card {
    display: flex;
    flex: 1 1 auto;
    min-height: 0;
    width: 100%;
    max-width: none;
    padding-top: 0.45rem;
  }

  .template-detail-grid {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    gap: 0.75rem;
    align-items: start;
    width: 100%;
    min-height: 0;
  }

  .template-detail-body-field {
    display: flex;
    min-height: 0;
    width: 100%;
    grid-column: 1 / 2;
  }

  .template-detail-form-field {
    margin-bottom: 0.35rem;
  }

  .template-detail-header-fields :global(.input.is-medium) {
    padding-top: 0.55rem;
    padding-bottom: 0.55rem;
    font-size: 1rem;
  }

  .template-detail-body-field :global(.md-input) {
    display: flex;
    flex: 1 1 auto;
    min-height: 0;
    width: 100%;
  }

  .template-detail-body-field :global(.md-input-area) {
    width: 100%;
  }

  .template-detail-body-field :global(.md-input-area) {
    min-height: calc(100vh - 20rem);
  }

  @media screen and (max-width: 768px) {
    .template-detail-header {
      padding: 0.65rem 0.75rem;
    }

    .template-detail-body {
      min-height: calc(100vh - 13rem);
      margin-bottom: 4.25rem;
    }

    .template-detail-grid {
      grid-template-columns: minmax(0, 1fr);
    }

    .template-detail-header-top {
      align-items: stretch;
    }

    .template-detail-body-field {
      grid-column: auto;
    }

    .template-detail-body-field :global(.md-input-area) {
      min-height: calc(100vh - 19rem);
    }
  }
</style>
