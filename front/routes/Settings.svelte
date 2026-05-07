<script lang="ts">
  import { goto } from "@mateothegreat/svelte5-router";
  import { onMount } from "svelte";
  import AppIcon from "../components/AppIcon.svelte";
  import { apiSettingsPath, appPath } from "../basePath.js";
  import type { settingType } from "../models/settingType.js";
  import { settingsStore } from "../store.js";

  let selectedCategoryKey = $state("");

  const categories = $derived.by(() => [
    { key: "", name: "すべて" },
    ...($settingsStore.Categories?.map((category) => ({
      key: category.Key,
      name: category.Name,
    })) ?? []),
  ]);

  const templates = $derived.by(() =>
    ($settingsStore.Categories ?? []).flatMap((category) =>
      (category.Templates ?? []).map((template) => ({
        name: template.Name,
        categoryKey: category.Key,
        summary: template.Value?.trim()
          ? `${template.Value.split(/\r?\n/, 1)[0]}`
          : "空テンプレート",
        tags: template.Tags ?? [],
        body: template.Value ?? "",
        notes: `${category.Name}のテンプレート`,
      })),
    ),
  );

  const filteredTemplates = $derived.by(() =>
    selectedCategoryKey === ""
      ? templates
      : templates.filter((template) => template.categoryKey === selectedCategoryKey),
  );

  onMount(() => {
    (async () => {
      const response = await fetch(apiSettingsPath());
      if (!response.ok) {
        return;
      }
      settingsStore.set((await response.json()) as settingType);
    })();
  });

  const openTemplateDetail = async (name: string) => {
    await goto(appPath(`/settings/templates/${encodeURIComponent(name)}`));
  };

  const createNewTemplate = async () => {
    await goto(appPath("/settings/templates/new"));
  };

</script>

<section class="settings-page">
  <div class="container is-fluid settings-page-inner">
    <div class="settings-card">
      <p class="settings-eyebrow">Settings</p>
      <h1 class="title is-3 settings-title">設定</h1>
      <p class="settings-description">
        ここに表示設定や入力ルールなどの項目を追加していきます。
      </p>

      <div class="settings-placeholder box">
        <p class="has-text-weight-semibold">設定画面のひな形</p>
        <p class="has-text-grey">必要な設定項目を順番に追加できます。</p>
      </div>

      <section class="template-admin">
        <div class="template-admin-header">
          <div>
            <p class="template-admin-eyebrow">Templates</p>
            <h2 class="title is-4 template-admin-title">テンプレート管理</h2>
            <p class="template-admin-description">
              カテゴリで絞り込みながら、API から読み込んだテンプレートを管理します。
            </p>
          </div>
          <div class="template-admin-actions">
            <button class="button add-button add-button-primary" type="button" onclick={createNewTemplate}>
              <span class="icon"><AppIcon name="plus" /></span>
              <span>新規</span>
            </button>
          </div>
        </div>

        <div
          class="template-category-tabs"
          role="tablist"
          aria-label="テンプレートカテゴリ"
        >
          {#each categories as category}
            <button
              class="template-category-tab"
              class:is-active={selectedCategoryKey === category.key}
              type="button"
              role="tab"
              aria-selected={selectedCategoryKey === category.key}
              onclick={() => {
                selectedCategoryKey = category.key;
              }}
            >
              {category.name}
            </button>
          {/each}
        </div>

        <div class="template-admin-layout">
          <aside class="template-list-panel">
            <div class="template-list-panel-header">
              <p class="has-text-weight-semibold">テンプレート一覧</p>
              <p class="has-text-grey">{filteredTemplates.length}件</p>
            </div>

            <div class="template-list">
              {#each filteredTemplates as template}
                <div
                  class="template-list-item"
                  role="button"
                  tabindex="0"
                  onclick={() => openTemplateDetail(template.name)}
                  onkeydown={(event) => {
                    if (event.key === "Enter" || event.key === " ") {
                      event.preventDefault();
                      void openTemplateDetail(template.name);
                    }
                  }}
                >
                  <span class="template-list-item-main">
                    <span class="template-list-item-name">{template.name}</span>
                    <span class="template-list-item-meta">{template.summary}</span>
                    {#if template.tags.length > 0}
                      <span class="template-list-item-tags">{template.tags.join(" / ")}</span>
                    {/if}
                  </span>
                </div>
              {/each}
            </div>
          </aside>
        </div>
      </section>
    </div>
  </div>

</section>

<style>
  .settings-page {
    min-height: 100vh;
    padding: 1rem 1rem 3rem;
    background: var(--bulma-scheme-main);
  }

  .settings-page-inner {
    max-width: 56rem;
  }

  .settings-card {
    padding: 1.25rem 1rem;
    border-radius: 1rem;
    background: color-mix(in srgb, var(--bulma-scheme-main) 94%, black 6%);
  }

  .settings-eyebrow,
  .template-admin-eyebrow {
    margin-bottom: 0.35rem;
    color: var(--bulma-text-weak);
    font-size: 0.8rem;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  .settings-title {
    margin-bottom: 0.65rem;
  }

  .settings-description {
    margin-bottom: 1rem;
    color: var(--bulma-text-weak);
  }

  .settings-placeholder {
    margin-bottom: 0;
  }

  .template-admin {
    margin-top: 1.5rem;
    padding-top: 1.25rem;
    border-top: 1px solid
      color-mix(in srgb, var(--bulma-border) 78%, transparent);
  }

  .template-admin-header {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    align-items: flex-start;
    margin-bottom: 1rem;
  }

  .template-admin-title {
    margin-bottom: 0.35rem;
  }

  .template-admin-description {
    color: var(--bulma-text-weak);
  }

  .template-admin-actions {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .add-button {
    min-width: 7rem;
  }

  .add-button-primary {
    background: color-mix(in srgb, #2d8f86 62%, var(--bulma-scheme-main));
    border: 1px solid color-mix(in srgb, #2d8f86 74%, black 26%);
    color: #edf8f6;
    box-shadow: 0 10px 24px rgba(10, 31, 29, 0.16);
  }

  .template-category-tabs {
    display: flex;
    gap: 0.4rem;
    flex-wrap: wrap;
    margin-bottom: 1rem;
  }

  .template-category-tab {
    min-height: 2.25rem;
    padding: 0.45rem 0.8rem;
    border: 1px solid var(--bulma-border);
    border-radius: 999px;
    background: var(--bulma-scheme-main);
    color: var(--bulma-text-weak);
    font-weight: 600;
  }

  .template-category-tab.is-active {
    border-color: color-mix(
      in srgb,
      var(--bulma-link) 45%,
      var(--bulma-border)
    );
    background: color-mix(
      in srgb,
      var(--bulma-link) 10%,
      var(--bulma-scheme-main)
    );
    color: var(--bulma-link-text);
  }

  .template-admin-layout {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    gap: 1rem;
  }

  .template-list-panel {
    border: 1px solid var(--bulma-border);
    border-radius: 1rem;
    background: color-mix(in srgb, var(--bulma-scheme-main) 98%, black 2%);
    padding: 1rem;
  }

  .template-list-panel-header {
    display: flex;
    justify-content: space-between;
    gap: 0.75rem;
    align-items: flex-start;
    margin-bottom: 0.85rem;
  }

  .template-list {
    display: grid;
    gap: 0.6rem;
  }

  .template-list-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    width: 100%;
    padding: 0.8rem 0.85rem;
    border: 1px solid var(--bulma-border);
    border-radius: 0.85rem;
    background: transparent;
    text-align: left;
  }

  .template-list-item-main {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.15rem;
    min-width: 0;
  }

  .template-list-item-name {
    font-weight: 700;
  }

  .template-list-item-meta,
  .template-list-item-tags {
    color: var(--bulma-text-weak);
    font-size: 0.9rem;
  }

  @media screen and (max-width: 768px) {
    .template-admin-header {
      flex-direction: column;
    }

    .template-admin-actions {
      justify-content: flex-start;
    }
  }
</style>
