<script lang="ts">
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faCaretDown } from "@fortawesome/free-solid-svg-icons";
  library.add(faCaretDown);
  dom.watch();

  interface Props {
    /** 選択アイテムのリスト */
    items: { key: string; value: string }[];
    /** 表示値のキー */
    key?: string;
    tabindex: number | null | undefined;
    onchange?: (value: string) => void;
  }

  let { items, key = $bindable(""), tabindex, onchange }: Props = $props();
  /** 表示値 */
  let value = $derived(items.find((item) => item.key === key)?.value ?? "");
  /** ドロップダウン トグルイベント */
  const onToggle = () =>
    document.getElementById("dropdown")?.classList.toggle("is-active");

  /** アイテム変更イベント */
  const onChange = (nextKey: string) => {
    key = nextKey;
    onchange?.(nextKey);
    onToggle();
  };
</script>

<div class="dropdown" id="dropdown">
  <div class="dropdown-trigger">
    <button
      class="button"
      aria-haspopup="true"
      aria-controls="dropdown-menu"
      id="para-button"
      {tabindex}
      onclick={onToggle}
    >
      <span>{value}</span>
      <span class="icon is-small">
        <i class="fa-solid fa-caret-down"></i>
      </span>
    </button>
  </div>
  <div class="dropdown-menu" id="dropdown-menu" role="menu">
    <div class="dropdown-content">
      {#each items as i}
        <button class="dropdown-item" onclick={() => onChange(i.key)}>
          {i.value}
        </button>
      {/each}
    </div>
  </div>
</div>
