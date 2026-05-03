<script lang="ts">
  import { goto, type RouteResult } from "@mateothegreat/svelte5-router";
  import { onMount } from "svelte";
  import type { entryType } from "../models/entryModels.js";
  import { settingsStore } from "../store.js";
  import { appPath } from "../basePath.js";
  import { createEmptyEntry, loadEntry } from "../lib/entryApi.js";

  interface Props {
    route?: RouteResult;
  }

  type EntryRouteParams = {
    id?: string | number | boolean;
    category?: string | number | boolean;
  };

  let { route: currentRoute = undefined }: Props = $props();

  const routeParams = $derived(
    (currentRoute?.result?.path?.params ?? {}) as EntryRouteParams,
  );
  const categoryKey = $derived(String(routeParams.category ?? ""));
  const entryId = $derived(routeParams.id ? String(routeParams.id) : "");
  const resolvedSettings = $derived.by(() => {
    const categorySetting = $settingsStore?.Categories?.find(
      (category) => category.Key === categoryKey,
    );
    return {
      showTags: categorySetting?.UseTag ?? false,
    };
  });

  let entry = $state<entryType>(createEmptyEntry());
  let initialized = $state(false);
  let isLoading = $state(true);
  let isPageLoading = $state(true);
  let isErr = $state(false);
  let errMessage = $state("");

  document.querySelector<HTMLDivElement>(".navbar")?.classList.add("is-hidden");

  const listPath = () => appPath(`/${categoryKey}/`);
  const editPath = () => appPath(`/${categoryKey}/${entryId}/edit`);
  const goToList = async () => goto(listPath());
  const goToEdit = async () => {
    if (!entryId) {
      return;
    }
    await goto(editPath());
  };
  const finishPageLoading = () => {
    requestAnimationFrame(() => {
      requestAnimationFrame(() => {
        isPageLoading = false;
      });
    });
  };
  const renderedHTML = $derived.by(() => {
    if (entry.HTML && entry.HTML.trim() !== "") {
      return entry.HTML;
    }
    return "";
  });

  onMount(() => {
    initialized = true;
  });

  $effect(() => {
    initialized;
    categoryKey;
    entryId;
    currentRoute;

    if (!initialized || !categoryKey || !entryId) {
      return;
    }

    isLoading = true;
    isPageLoading = true;
    isErr = false;
    errMessage = "";
    entry = createEmptyEntry();

    (async () => {
      try {
        entry = await loadEntry(categoryKey, entryId);
      } catch (error) {
        isErr = true;
        errMessage =
          error instanceof Error ? error.message : "エントリを読み込めませんでした";
      } finally {
        isLoading = false;
        finishPageLoading();
      }
    })();
  });
</script>

{#if isErr && errMessage !== ""}
  <div class="notification is-danger m-3">{errMessage}</div>
{/if}

<div class="reference-page" class:is-page-loading={isPageLoading}>
  {#if isLoading}
    <div class="reference-loading">
      <button class="button is-dark is-loading" aria-label="loading"></button>
    </div>
  {/if}

  <header class="reference-header">
    <div class="reference-date-value">{entry.Date}</div>

    <div class="reference-outline-value">{entry.Outline || "-"}</div>

    {#if resolvedSettings.showTags && entry.Tags.length > 0}
      <div class="reference-tags">
        <span class="reference-tags-icon" aria-hidden="true">
          <i class="fa-solid fa-tags"></i>
        </span>
        {#each entry.Tags as tag}
          <span class="reference-tag">{tag}</span>
        {/each}
      </div>
    {/if}
  </header>

  <button
    class="reference-body"
    type="button"
    aria-label="本文を編集"
    onclick={goToEdit}
  >
    <div class="reference-body-inner">
      {#if renderedHTML !== ""}
        <div class="reference-markdown" tabindex="-1">
          {@html renderedHTML}
        </div>
      {:else}
        <div class="reference-empty-body">本文がありません</div>
      {/if}
    </div>
  </button>

  <footer class="reference-footer">
    <div class="reference-footer-actions">
      <button class="button reference-footer-button reference-footer-button-secondary" onclick={goToList}>
        <span class="icon"><i class="fa-solid fa-arrow-left"></i></span>
        <span>戻る</span>
      </button>
      <button class="button reference-footer-button reference-footer-button-primary" onclick={goToEdit}>
        <span class="icon"><i class="fa-solid fa-pen"></i></span>
        <span>編集</span>
      </button>
    </div>
  </footer>
</div>

<style>
  .reference-page {
    min-height: 100vh;
    background: var(--bulma-scheme-main);
  }

  .reference-page.is-page-loading {
    visibility: hidden;
  }

  .reference-loading {
    position: fixed;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: color-mix(in srgb, var(--bulma-scheme-main) 86%, transparent);
    z-index: 1200;
    visibility: visible;
  }

  .reference-header {
    padding: 0.9rem 1rem 0.8rem;
    background: var(--bulma-border);
    position: sticky;
    top: 0;
    z-index: 20;
    box-shadow: 0 0.2rem 0.8rem rgba(15, 23, 42, 0.14);
  }

  .reference-date-value {
    font-size: 0.98rem;
    font-weight: 600;
    line-height: 1.2;
    font-variant-numeric: tabular-nums;
    color: var(--bulma-text-weak);
  }

  .reference-outline-value {
    margin-top: 0.3rem;
    font-size: clamp(1.5rem, 4.8vw, 2.15rem);
    font-weight: 700;
    line-height: 1.22;
    letter-spacing: 0.01em;
    white-space: pre-wrap;
    word-break: break-word;
  }

  .reference-tags {
    display: flex;
    align-items: center;
    gap: 0.45rem;
    flex-wrap: wrap;
    margin-top: 0.55rem;
    color: var(--bulma-text-weak);
    font-size: 0.94rem;
  }

  .reference-tags-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    color: var(--bulma-text);
  }

  .reference-tag {
    color: var(--bulma-text);
  }

  .reference-body {
    width: 100%;
    min-height: calc(100vh - 16rem);
    padding: 0;
    border: none;
    background: transparent;
    text-align: left;
    cursor: pointer;
  }

  .reference-body-inner {
    min-height: calc(100vh - 16rem);
    padding: 1rem 1rem 6rem;
  }

  .reference-markdown,
  .reference-empty-body {
    padding: 1.1rem 1rem;
    border-radius: 1rem;
    background: color-mix(in srgb, var(--bulma-scheme-main) 92%, black 8%);
    box-shadow: 0 0.5rem 1.4rem rgba(15, 23, 42, 0.06);
  }

  .reference-empty-body {
    color: var(--bulma-text-weak);
  }

  .reference-markdown :global(*) {
    max-width: 100%;
  }

  .reference-markdown :global(h1),
  .reference-markdown :global(h2),
  .reference-markdown :global(h3),
  .reference-markdown :global(h4),
  .reference-markdown :global(h5),
  .reference-markdown :global(h6) {
    margin-top: 1.4rem;
    margin-bottom: 0.65rem;
    line-height: 1.3;
  }

  .reference-markdown :global(h1),
  .reference-markdown :global(h2),
  .reference-markdown :global(h3) {
    color: color-mix(in srgb, var(--bulma-link) 58%, white 42%);
    font-weight: 700;
    letter-spacing: 0.01em;
  }

  .reference-markdown :global(h1) {
    margin-top: 1.8rem;
    padding-bottom: 0.45rem;
    border-bottom: 1px solid color-mix(in srgb, var(--bulma-link) 35%, transparent);
    font-size: 1.45rem;
  }

  .reference-markdown :global(h2) {
    margin-top: 1.55rem;
    padding-left: 0.7rem;
    border-left: 0.25rem solid color-mix(in srgb, var(--bulma-link) 55%, transparent);
    font-size: 1.2rem;
  }

  .reference-markdown :global(h3) {
    margin-top: 1.3rem;
    font-size: 1.05rem;
  }

  .reference-markdown :global(h1:first-child),
  .reference-markdown :global(h2:first-child),
  .reference-markdown :global(h3:first-child),
  .reference-markdown :global(h4:first-child),
  .reference-markdown :global(h5:first-child),
  .reference-markdown :global(h6:first-child),
  .reference-markdown :global(p:first-child) {
    margin-top: 0;
  }

  .reference-markdown :global(p),
  .reference-markdown :global(ul),
  .reference-markdown :global(ol),
  .reference-markdown :global(blockquote),
  .reference-markdown :global(pre),
  .reference-markdown :global(table) {
    margin-top: 0;
    margin-bottom: 1rem;
  }

  .reference-markdown :global(ul) {
    padding-left: 0;
    list-style: none;
  }

  .reference-markdown :global(ul li) {
    position: relative;
    padding-left: 1.1rem;
    margin-bottom: 0.35rem;
  }

  .reference-markdown :global(ul li::before) {
    content: "・";
    position: absolute;
    left: 0;
    color: color-mix(in srgb, var(--bulma-link) 45%, white 55%);
  }

  .reference-markdown :global(ol) {
    padding-left: 1.4rem;
  }

  .reference-markdown :global(table) {
    width: 100%;
    border-collapse: collapse;
    overflow: hidden;
    border-radius: 0.75rem;
    border: 1px solid color-mix(in srgb, var(--bulma-border) 70%, white 10%);
    background: color-mix(in srgb, var(--bulma-scheme-main) 94%, black 6%);
  }

  .reference-markdown :global(th),
  .reference-markdown :global(td) {
    padding: 0.7rem 0.8rem;
    border: 1px solid color-mix(in srgb, var(--bulma-border) 74%, white 8%);
    text-align: left;
    vertical-align: top;
  }

  .reference-markdown :global(th) {
    background: color-mix(in srgb, var(--bulma-border) 86%, black 14%);
    color: color-mix(in srgb, var(--bulma-text) 82%, white 18%);
    font-weight: 600;
  }

  .reference-markdown :global(img) {
    display: block;
    height: auto;
    margin: 0.75rem 0;
    border-radius: 0.75rem;
  }

  .reference-markdown :global(pre) {
    overflow-x: auto;
    padding: 0.85rem 1rem;
    border-radius: 0.75rem;
    background: color-mix(in srgb, var(--bulma-border) 70%, black 30%);
  }

  .reference-markdown :global(blockquote) {
    padding-left: 1rem;
    border-left: 0.25rem solid var(--bulma-link);
    color: var(--bulma-text-strong);
  }

  .reference-footer {
    position: fixed;
    left: 0;
    right: 0;
    bottom: 0;
    background: var(--bulma-border);
  }

  .reference-footer-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
  }

  .reference-footer-button {
    min-width: 7rem;
    border-radius: 0.9rem;
    border: 1px solid transparent;
    font-weight: 600;
  }

  .reference-footer-button-secondary {
    background: color-mix(in srgb, var(--bulma-border) 74%, var(--bulma-scheme-main));
    border-color: color-mix(in srgb, var(--bulma-border) 88%, white 12%);
    color: color-mix(in srgb, var(--bulma-text) 90%, white 10%);
  }

  .reference-footer-button-primary {
    background: color-mix(in srgb, #2d8f86 62%, var(--bulma-scheme-main));
    border-color: color-mix(in srgb, #2d8f86 74%, black 26%);
    color: #edf8f6;
  }

  @media screen and (max-width: 768px) {
    .reference-header {
      padding: 0.7rem 0.75rem 0.65rem;
    }

    .reference-body-inner {
      padding: 0.75rem 0.75rem 5.5rem;
    }

    .reference-markdown,
    .reference-empty-body {
      padding: 0.95rem 0.85rem;
      border-radius: 0.85rem;
    }

    .reference-date-value {
      font-size: 0.92rem;
    }

    .reference-footer-actions {
      gap: 0.5rem;
      padding: 0.65rem 0.75rem;
    }

    .reference-footer-button {
      flex: 1;
      min-width: 0;
    }
  }
</style>
