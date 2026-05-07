<script lang="ts">
  import AppIcon from "../AppIcon.svelte";

  export const paragraphs = [
    { key: "normal", value: "本文", icon: "grip-lines" },
    { key: "h1", value: "見出し1", badge: "H1" },
    { key: "h2", value: "見出し2", badge: "H2" },
    { key: "h3", value: "見出し3", badge: "H3" },
    { key: "ol", value: "番号リスト", icon: "list-ol" },
    { key: "ul", value: "段落リスト", icon: "list-ul" },
    { key: "code", value: "コード", icon: "code" },
    { key: "quote", value: "引用", icon: "quote-left" },
  ];

  interface Props {
    value?: string;
    bold?: boolean;
    italic?: boolean;
    link?: boolean;
    strike?: boolean;
    onChange?: (value: string) => void;
    onBold?: () => void;
    onCarryOver?: () => void;
    onImage?: () => void;
    onItalic?: () => void;
    onLink?: () => void;
    onStrike?: () => void;
  }

  let {
    value = $bindable("normal"),
    bold = false,
    italic = false,
    link = false,
    strike = false,
    onChange,
    onBold,
    onCarryOver,
    onImage,
    onItalic,
    onLink,
    onStrike,
  }: Props = $props();
  let isParagraphMenuOpen = $state(false);

  const currentParagraph = $derived(
    paragraphs.find((item) => item.key === value) ?? paragraphs[0],
  );

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
      class="button md-paragraph-trigger"
      type="button"
      aria-haspopup="menu"
      aria-expanded={isParagraphMenuOpen}
      onclick={() => (isParagraphMenuOpen = !isParagraphMenuOpen)}
    >
      {#key currentParagraph.key}
        {#if currentParagraph.badge}
          <span class="icon md-paragraph-icon-slot">
            <span class="md-paragraph-badge">{currentParagraph.badge}</span>
          </span>
        {:else}
          <span class="icon md-paragraph-icon-slot">
            <AppIcon name={currentParagraph.icon ?? "grip-lines"} />
          </span>
        {/if}
      {/key}
      <span class="icon is-small">
        <AppIcon name="chevron-down" size={16} />
      </span>
    </button>

    {#if isParagraphMenuOpen}
      <div class="md-paragraph-menu" role="menu">
        {#each paragraphs as item}
          <button
            class="button is-ghost md-paragraph-option"
            class:is-active={item.key === value}
            type="button"
            role="menuitemradio"
            aria-checked={item.key === value}
            onclick={() => handleChange(item.key)}
          >
            {#if item.badge}
              <span class="icon md-paragraph-icon-slot">
                <span class="md-paragraph-badge">{item.badge}</span>
              </span>
            {:else}
              <span class="icon md-paragraph-icon-slot">
                <AppIcon name={item.icon ?? "grip-lines"} />
              </span>
            {/if}
            <span>{item.value}</span>
          </button>
        {/each}
      </div>
    {/if}
  </div>
  <button
    class="button is-ghost md-toolbar-button"
    class:is-active={bold}
    type="button"
    aria-label="bold"
    onclick={onBold}
  >
    <AppIcon name="bold" />
  </button>
  <button
    class="button is-ghost md-toolbar-button"
    class:is-active={italic}
    type="button"
    aria-label="italic"
    onclick={onItalic}
  >
    <AppIcon name="italic" />
  </button>
  <button
    class="button is-ghost md-toolbar-button"
    class:is-active={strike}
    type="button"
    aria-label="strike"
    onclick={onStrike}
  >
    <AppIcon name="strikethrough" />
  </button>
  <button
    class="button is-ghost md-toolbar-button"
    type="button"
    aria-label="carry over marker"
    title="次回グループ追加へ引き継ぐ位置を挿入"
    onclick={onCarryOver}
  >
    <AppIcon name="arrows-rotate" />
  </button>
  <button
    class="button is-ghost md-toolbar-button"
    type="button"
    aria-label="image"
    onclick={onImage}
  >
    <AppIcon name="image" />
  </button>
  <button
    class="button is-ghost md-toolbar-button"
    class:is-active={link}
    type="button"
    aria-label="link"
    onclick={onLink}
  >
    <AppIcon name="link" />
  </button>
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

  .md-paragraph-backdrop {
    position: fixed;
    inset: 0;
    z-index: 9;
    border: none;
    background: transparent;
    padding: 0;
  }

  .md-paragraph-trigger {
    min-width: 3.3rem;
    min-height: 2rem;
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
    border-radius: var(--bulma-radius-large);
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

  .md-paragraph-option.is-active {
    color: #8ec5ff;
    background: color-mix(in srgb, #2f7df4 18%, transparent);
  }

  .md-toolbar-button {
    flex: 0 0 auto;
    min-width: 2rem;
    min-height: 2rem;
    padding: 0 0.2rem;
  }

  .md-toolbar-button.is-active {
    color: var(--bulma-link);
    background: color-mix(in srgb, var(--bulma-link) 12%, transparent);
    box-shadow: inset 0 0 0 1px
      color-mix(in srgb, var(--bulma-link) 26%, transparent);
  }

  @media screen and (max-width: 768px) {
    .md-toolbar {
      gap: 0.2rem;
    }

    .md-paragraph-trigger {
      min-width: 3rem;
    }
  }
</style>
