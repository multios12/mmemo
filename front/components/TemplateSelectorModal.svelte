<script lang="ts">
  import type { TemplateType } from "../models/settingType.js";

  interface Props {
    isOpen: boolean;
    templates: TemplateType[];
    initialSelectedIndex: number | null;
    onSelect: (payload: { template: TemplateType; index: number }) => void;
  }

  let { isOpen, templates, initialSelectedIndex, onSelect }: Props = $props();
  let selectedIndex = $state<number | null>(null);
  let wasOpen = $state(false);

  const templatePreview = (value: string) =>
    value
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter((line) => line !== "")
      .slice(0, 3);

  $effect(() => {
    if (isOpen && !wasOpen) {
      selectedIndex = initialSelectedIndex;
    }
    wasOpen = isOpen;
  });

  const onChoose = (template: TemplateType, index: number) => {
    selectedIndex = index;
    onSelect({ template, index });
  };
</script>

{#if isOpen}
  <div class="template-selector-modal">
    <div class="template-selector">
      <div class="template-selector-header">
        <h2 class="template-selector-title">テンプレート選択</h2>
      </div>
      <div class="template-grid">
        {#each templates as template, index}
          <button
            class="template-card"
            class:is-selected={selectedIndex === index}
            type="button"
            onclick={() => onChoose(template, index)}
          >
            <span class="template-card-title">{template.Name}</span>
            <span class="template-card-preview">
              {#if templatePreview(template.Value).length > 0}
                {templatePreview(template.Value).join(" / ")}
              {:else}
                空の本文で開始します
              {/if}
            </span>
            {#if (template.Tags?.length ?? 0) > 0}
              <span class="template-card-tags">
                {template.Tags?.join(" / ")}
              </span>
            {/if}
          </button>
        {/each}
      </div>
    </div>
  </div>
{/if}

<style>
  .template-selector {
    width: min(100%, 34rem);
    max-height: min(70vh, 38rem);
    overflow-y: auto;
    padding: 1rem;
    border: 1px solid
      color-mix(in srgb, var(--app-link) 18%, var(--app-border));
    border-radius: 1.1rem;
    background: var(--app-scheme-main);
    box-shadow:
      0 1.2rem 3rem rgba(15, 23, 42, 0.18),
      0 0 0 1px color-mix(in srgb, var(--app-link) 12%, transparent);
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
      var(--app-scheme-main) 55%,
      rgba(15, 23, 42, 0.45)
    );
    backdrop-filter: blur(6px);
  }

  .template-selector-header {
    margin-bottom: 0.75rem;
  }

  .template-selector-title {
    margin: 0 0 0.5rem;
    color: var(--app-text-strong);
    font-size: 1rem;
    font-weight: 700;
    line-height: 1.15;
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
    border: 1px solid var(--app-border);
    border-radius: 0.9rem;
    background: color-mix(
      in srgb,
      var(--app-link) 10%,
      var(--app-scheme-main)
    );
    color: var(--app-text);
    text-align: left;
    transition:
      border-color 0.15s ease,
      transform 0.15s ease,
      box-shadow 0.15s ease;
  }

  .template-card:hover {
    border-color: var(--app-link);
    transform: translateY(-1px);
  }

  .template-card.is-selected {
    border-color: var(--app-link);
    background: color-mix(
      in srgb,
      var(--app-link) 18%,
      var(--app-scheme-main)
    );
    box-shadow: 0 0 0 1px color-mix(in srgb, var(--app-link) 45%, transparent);
  }

  .template-card-title {
    font-size: 1rem;
    font-weight: 700;
  }

  .template-card-preview {
    color: var(--app-text-weak);
    line-height: 1.5;
  }

  .template-card-tags {
    color: var(--app-link-text);
    font-size: 0.8rem;
    font-weight: 600;
    line-height: 1.4;
  }

  @media screen and (max-width: 768px) {
    .template-selector {
      padding: 0.75rem;
    }
  }
</style>
