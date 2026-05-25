<script lang="ts">
  import AppIcon from "./AppIcon.svelte";

  type SegmentedItem = {
    value: string;
    label: string;
    icon?: string;
  };

  interface Props {
    items: SegmentedItem[];
    selected: string;
    showLabel?: boolean;
    onChange: (value: string) => void;
  }

  let {
    items,
    selected,
    showLabel = true,
    onChange,
  }: Props = $props();
</script>

<div class="seg" role="group" aria-label="表示モード">
  {#if showLabel}
    <span class="seg-label">表示モード</span>
  {/if}
  <div class="seg-items">
    {#each items as item}
      <button
        class="seg-item"
        class:is-active={selected === item.value}
        type="button"
        aria-pressed={selected === item.value}
        onclick={() => onChange(item.value)}
      >
        {#if item.icon}
          <AppIcon class="icon" name={item.icon} />
        {/if}
        <span>{item.label}</span>
      </button>
    {/each}
  </div>
</div>

<style>
  .seg {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    flex-wrap: wrap;
  }

  .seg-label {
    color: var(--app-text-weak);
    font-size: 0.86rem;
    font-weight: 600;
    white-space: nowrap;
  }

  .seg-items {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.22rem;
    border: 1px solid var(--app-border);
    border-radius: 999px;
    background: color-mix(in srgb, var(--app-scheme-main) 94%, black 6%);
  }

  .seg-item {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    min-width: 0;
    padding: 0.45rem 0.8rem;
    border: none;
    border-radius: 999px;
    background: transparent;
    color: var(--app-text-weak);
    font-size: 0.9rem;
    cursor: pointer;
  }

  .seg-item.is-active {
    background: color-mix(in srgb, #2d8f86 22%, var(--app-scheme-main));
    color: color-mix(in srgb, var(--app-text) 92%, white 8%);
    box-shadow: inset 0 0 0 1px color-mix(in srgb, #2d8f86 26%, transparent);
  }

  @media screen and (max-width: 768px) {
    .seg {
      flex: 1 1 0;
      min-width: 0;
    }

    .seg-items {
      width: 100%;
      justify-content: stretch;
    }

    .seg-item {
      flex: 1 1 0;
      justify-content: center;
      font-size: 0.74rem;
      padding: 0.38rem 0.3rem;
      gap: 0.16rem;
    }

    .seg-item :global(.icon) {
      display: none;
    }
  }
</style>
