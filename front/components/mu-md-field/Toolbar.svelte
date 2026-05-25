<script lang="ts">
  import Bold from "lucide-svelte/icons/bold";
  import ChevronDown from "lucide-svelte/icons/chevron-down";
  import Code from "lucide-svelte/icons/code";
  import Grip from "lucide-svelte/icons/grip";
  import Italic from "lucide-svelte/icons/italic";
  import LinkIcon from "lucide-svelte/icons/link";
  import List from "lucide-svelte/icons/list";
  import ListOrdered from "lucide-svelte/icons/list-ordered";
  import Quote from "lucide-svelte/icons/quote";
  import Strikethrough from "lucide-svelte/icons/strikethrough";
  import type { ToolbarButton, ToolbarIcon } from "./types.js";

  interface ParagraphOption {
    key: string;
    value: string;
    icon?: ToolbarIcon;
    badge?: string;
  }

  export const paragraphs: ParagraphOption[] = [
    { key: "normal", value: "本文", icon: Grip },
    { key: "h1", value: "見出し1", badge: "H1" },
    { key: "h2", value: "見出し2", badge: "H2" },
    { key: "h3", value: "見出し3", badge: "H3" },
    { key: "ol", value: "番号リスト", icon: ListOrdered },
    { key: "ul", value: "段落リスト", icon: List },
    { key: "code", value: "コード", icon: Code },
    { key: "quote", value: "引用", icon: Quote },
  ];

  interface Props {
    value?: string;
    bold?: boolean;
    italic?: boolean;
    link?: boolean;
    strike?: boolean;
    onChange?: (value: string) => void;
    onBold?: () => void;
    onItalic?: () => void;
    onLink?: () => void;
    onStrike?: () => void;
    extraButtons?: ToolbarButton[];
    onExtraButton?: (button: ToolbarButton) => void;
  }

  let {
    value = $bindable("normal"),
    bold = false,
    italic = false,
    link = false,
    strike = false,
    onChange,
    onBold,
    onItalic,
    onLink,
    onStrike,
    extraButtons = [],
    onExtraButton,
  }: Props = $props();
  let isParagraphMenuOpen = $state(false);

  const currentParagraph = $derived(
    paragraphs.find((item) => item.key === value) ?? paragraphs[0],
  );
  const CurrentParagraphIcon = $derived(currentParagraph.icon ?? Grip);

  const handleChange = (nextValue: string) => {
    value = nextValue;
    isParagraphMenuOpen = false;
    onChange?.(nextValue);
  };

  const closeParagraphMenu = () => {
    isParagraphMenuOpen = false;
  };
</script>

<div id="toolbar" class="md-toolbar">
  <div class="md-paragraph-select">
    {#if isParagraphMenuOpen}
      <button
        class="md-paragraph-backdrop"
        type="button"
        aria-label="close paragraph menu"
        onclick={closeParagraphMenu}
      ></button>
    {/if}

    <button
      class="md-button md-paragraph-trigger"
      type="button"
      aria-haspopup="menu"
      aria-expanded={isParagraphMenuOpen}
      onclick={() => (isParagraphMenuOpen = !isParagraphMenuOpen)}
    >
      {#key currentParagraph.key}
        {#if currentParagraph.badge}
          <span class="md-icon md-paragraph-icon-slot">
            <span class="md-paragraph-badge">{currentParagraph.badge}</span>
          </span>
        {:else}
          <span class="md-icon md-paragraph-icon-slot">
            <CurrentParagraphIcon />
          </span>
        {/if}
      {/key}
      <span class="md-icon md-icon-small">
        <ChevronDown size={16} />
      </span>
    </button>

    {#if isParagraphMenuOpen}
      <div class="md-paragraph-menu" role="menu">
        {#each paragraphs as item}
          {@const ParagraphIcon = item.icon ?? Grip}
          <button
            class="md-button md-paragraph-option"
            class:md-active={item.key === value}
            type="button"
            role="menuitemradio"
            aria-checked={item.key === value}
            onclick={() => handleChange(item.key)}
          >
            {#if item.badge}
              <span class="md-icon md-paragraph-icon-slot">
                <span class="md-paragraph-badge">{item.badge}</span>
              </span>
            {:else}
              <span class="md-icon md-paragraph-icon-slot">
                <ParagraphIcon />
              </span>
            {/if}
            <span>{item.value}</span>
          </button>
        {/each}
      </div>
    {/if}
  </div>
  <button
    class="md-button md-toolbar-button"
    class:md-active={bold}
    type="button"
    aria-label="bold"
    onclick={onBold}
  >
    <Bold />
  </button>
  <button
    class="md-button md-toolbar-button"
    class:md-active={italic}
    type="button"
    aria-label="italic"
    onclick={onItalic}
  >
    <Italic />
  </button>
  <button
    class="md-button md-toolbar-button"
    class:md-active={strike}
    type="button"
    aria-label="strike"
    onclick={onStrike}
  >
    <Strikethrough />
  </button>
  <button
    class="md-button md-toolbar-button"
    class:md-active={link}
    type="button"
    aria-label="link"
    onclick={onLink}
  >
    <LinkIcon />
  </button>
  {#each extraButtons as button (button.key)}
    <button
      class="md-button md-toolbar-button"
      class:md-active={button.active}
      type="button"
      aria-label={button.ariaLabel}
      title={button.title}
      disabled={button.disabled}
      onclick={() => onExtraButton?.(button)}
    >
      {#if button.icon}
        {@const ButtonIcon = button.icon}
        <ButtonIcon />
      {:else if button.label}
        <span class="md-toolbar-button-label">{button.label}</span>
      {/if}
    </button>
  {/each}
</div>

<style>
  .md-toolbar {
    display: flex;
    align-items: center;
    flex-wrap: nowrap;
    gap: 0.25rem;
    overflow: visible;
    margin: 0 0.35rem 0;
  }

  .md-paragraph-select {
    position: relative;
  }

  .md-button {
    appearance: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border: 1px solid color-mix(in srgb, white 12%, transparent);
    border-radius: 0.45rem;
    background: transparent;
    color: inherit;
    cursor: pointer;
    font: inherit;
    gap: 0.35rem;
    line-height: 1;
    min-height: 2rem;
    min-width: 2rem;
    padding: 0 0.55rem;
    user-select: none;
  }

  .md-button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  .md-button:hover:not(:disabled),
  .md-button:focus-visible {
    background: color-mix(in srgb, white 8%, transparent);
    outline: none;
  }

  .md-button.md-active {
    background: color-mix(in srgb, #3e8ed0 22%, transparent);
    border-color: color-mix(in srgb, #3e8ed0 55%, white 12%);
    color: #8ec5ff;
  }

  .md-paragraph-trigger {
    min-width: 3.3rem;
    padding: 0 0.45rem;
    justify-content: space-between;
    gap: 0.2rem;
    color: #f6f7fb;
    border-color: color-mix(in srgb, white 14%, transparent);
    background: #1b2030;
    font-size: 0.84rem;
  }

  .md-paragraph-menu {
    position: absolute;
    top: calc(100% + 0.25rem);
    left: 0;
    z-index: 10;
    display: flex;
    flex-direction: column;
    min-width: 13rem;
    padding: 0.25rem;
    border: 1px solid color-mix(in srgb, white 12%, transparent);
    border-radius: 0.75rem;
    background: #161b28;
    box-shadow: 0 0.85rem 2rem color-mix(in srgb, black 28%, transparent);
  }

  .md-paragraph-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 1.1rem;
    height: 1.1rem;
    color: #f6f7fb;
    font-size: 0.66rem;
    font-weight: 800;
    letter-spacing: 0.04em;
    line-height: 1;
  }

  .md-paragraph-icon-slot {
    flex: 0 0 1.25rem;
    width: 1.25rem;
    justify-content: center;
    margin-right: 0;
  }

  .md-paragraph-option {
    justify-content: flex-start;
    gap: 0.35rem;
    color: #eef2ff;
    padding: 0.4rem 0.55rem;
    font-size: 0.88rem;
  }

  .md-paragraph-option.md-active {
    color: #8ec5ff;
    background: color-mix(in srgb, #2f7df4 18%, transparent);
  }

  .md-toolbar-button {
    flex: 0 0 auto;
    min-width: 2rem;
    padding: 0 0.5rem;
  }

  .md-toolbar-button-label {
    font-size: 0.8rem;
    font-weight: 700;
    letter-spacing: 0.01em;
  }

  .md-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 1rem;
    height: 1rem;
  }

  .md-icon-small {
    width: 1rem;
    height: 1rem;
  }

  .md-paragraph-backdrop {
    position: fixed;
    inset: 0;
    z-index: 9;
    border: none;
    background: transparent;
    padding: 0;
  }
</style>
