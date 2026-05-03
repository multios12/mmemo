<script lang="ts">
  import AppIcon from "./AppIcon.svelte";
  import type { TemplateType } from "../models/settingType.js";

  interface Props {
    isOpen: boolean;
    initialMode: "new" | "overwrite";
    initialDraftName: string;
    initialOverwriteName: string;
    templates: TemplateType[];
    errorMessage: string;
    isSaving: boolean;
    isDeleting: boolean;
    onClose: () => void;
    onSave: (payload: { mode: "new" | "overwrite"; name: string }) => void;
    onDelete: (name: string) => void;
  }

  let {
    isOpen,
    initialMode,
    initialDraftName,
    initialOverwriteName,
    templates,
    errorMessage,
    isSaving,
    isDeleting,
    onClose,
    onSave,
    onDelete,
  }: Props = $props();

  let mode = $state<"new" | "overwrite">("new");
  let draftName = $state("");
  let overwriteName = $state("");
  let wasOpen = $state(false);

  const hasTemplates = $derived(templates.length > 0);

  $effect(() => {
    if (isOpen && !wasOpen) {
      mode = initialMode;
      draftName = initialDraftName;
      overwriteName =
        initialOverwriteName || templates[0]?.Name || "";
    }
    if (isOpen && mode === "overwrite" && !templates.some((template) => template.Name === overwriteName)) {
      overwriteName = templates[0]?.Name || "";
    }
    wasOpen = isOpen;
  });

  const onSubmit = () => {
    onSave({
      mode,
      name: mode === "overwrite" ? overwriteName.trim() : draftName.trim(),
    });
  };

  const onDeleteClick = () => {
    onDelete(overwriteName.trim());
  };
</script>

{#if isOpen}
  <div class="template-selector-modal">
    <div class="template-selector template-save-modal">
      <div class="template-selector-header">
        <h2 class="title is-6 mb-2">テンプレート保存</h2>
      </div>

      <div
        class="template-save-mode"
        role="tablist"
        aria-label="テンプレート保存モード"
      >
        <button
          class="template-save-mode-option"
          class:is-active={mode === "new"}
          type="button"
          role="tab"
          aria-selected={mode === "new"}
          onclick={() => (mode = "new")}
        >
          <span class="icon"><AppIcon name="plus" /></span>
          <span>新規作成</span>
        </button>
        <button
          class="template-save-mode-option"
          class:is-active={mode === "overwrite"}
          type="button"
          role="tab"
          aria-selected={mode === "overwrite"}
          disabled={!hasTemplates}
          onclick={() => (mode = "overwrite")}
        >
          <span class="icon"><AppIcon name="pen-to-square" /></span>
          <span>上書き</span>
        </button>
      </div>

      {#if mode === "new"}
        <div class="field">
          <label class="label" for="templateNameInput">テンプレート名</label>
          <div class="control">
            <input
              id="templateNameInput"
              class="input"
              type="text"
              bind:value={draftName}
              placeholder="テンプレート名"
            />
          </div>
        </div>
      {:else}
        <div class="field">
          <label class="label" for="templateOverwriteSelect">上書き先テンプレート</label>
          <div class="control">
            <div class="select is-fullwidth template-save-select">
              <select
                id="templateOverwriteSelect"
                bind:value={overwriteName}
              >
                {#each templates as template}
                  <option value={template.Name}>{template.Name}</option>
                {/each}
              </select>
            </div>
          </div>
        </div>
      {/if}

      {#if errorMessage !== ""}
        <div class="notification is-danger is-light p-3">
          {errorMessage}
        </div>
      {/if}

      <div class="template-save-actions">
        <button
          class="button modal-secondary-button"
          type="button"
          disabled={isSaving || isDeleting}
          onclick={onClose}
        >
          <span class="icon"><AppIcon name="xmark" /></span>
          <span>キャンセル</span>
        </button>
        {#if mode === "overwrite" && hasTemplates}
          <button
            class="button modal-danger-button"
            type="button"
            class:is-loading={isDeleting}
            disabled={isSaving || isDeleting}
            onclick={onDeleteClick}
          >
            <span class="icon"><AppIcon name="trash" /></span>
            <span>削除</span>
          </button>
        {/if}
        <button
          class="button modal-primary-button"
          type="button"
          class:is-loading={isSaving}
          disabled={isSaving || isDeleting}
          onclick={onSubmit}
        >
          <span class="icon"><AppIcon name="floppy-disk" /></span>
          <span>保存</span>
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
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

  .template-save-modal {
    width: min(100%, 30rem);
  }

  .template-selector-header {
    margin-bottom: 0.75rem;
  }

  .template-save-mode {
    display: flex;
    gap: 0.4rem;
    padding: 0.35rem;
    margin-bottom: 1rem;
    border: 1px solid color-mix(in srgb, var(--bulma-link) 16%, var(--bulma-border));
    border-radius: 999px;
    background: color-mix(in srgb, var(--bulma-link) 5%, var(--bulma-scheme-main));
  }

  .template-save-mode-option {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.45rem;
    flex: 1 1 0;
    min-height: 2.5rem;
    padding: 0.55rem 0.9rem;
    border: none;
    border-radius: 999px;
    background: transparent;
    color: var(--bulma-text-weak);
    font-size: 0.95rem;
    font-weight: 600;
    transition:
      background-color 0.15s ease,
      color 0.15s ease,
      box-shadow 0.15s ease;
  }

  .template-save-mode-option.is-active {
    background: color-mix(in srgb, var(--bulma-link) 18%, var(--bulma-scheme-main));
    color: var(--bulma-link-text);
    box-shadow: 0 0 0 1px color-mix(in srgb, var(--bulma-link) 18%, transparent);
  }

  .template-save-mode-option:disabled {
    opacity: 0.45;
  }

  .template-save-mode-option:not(:disabled):hover {
    color: var(--bulma-text);
    background: color-mix(in srgb, var(--bulma-link) 10%, var(--bulma-scheme-main));
  }

  .modal-secondary-button,
  .modal-primary-button,
  .modal-danger-button {
    border-radius: 0.85rem;
    font-weight: 600;
  }

  .modal-secondary-button {
    background: color-mix(in srgb, var(--bulma-border) 74%, var(--bulma-scheme-main));
    border: 1px solid color-mix(in srgb, var(--bulma-border) 86%, white 14%);
    color: color-mix(in srgb, var(--bulma-text) 90%, white 10%);
  }

  .modal-primary-button {
    background: color-mix(in srgb, #2d8f86 62%, var(--bulma-scheme-main));
    border: 1px solid color-mix(in srgb, #2d8f86 74%, black 26%);
    color: #edf8f6;
  }

  .modal-danger-button {
    background: color-mix(in srgb, #8a4f55 52%, var(--bulma-scheme-main));
    border: 1px solid color-mix(in srgb, #8a4f55 70%, black 30%);
    color: #f8ecee;
  }

  .template-save-select,
  .template-save-select select {
    width: 100%;
  }

  .template-save-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
  }

  @media screen and (max-width: 768px) {
    .template-selector {
      padding: 0.75rem;
    }

    .template-save-actions {
      justify-content: stretch;
      flex-wrap: wrap;
    }

    .template-save-actions .button {
      flex: 1 1 0;
    }
  }
</style>
