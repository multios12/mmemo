<script lang="ts">
  import { run } from 'svelte/legacy';

  import { createEventDispatcher } from "svelte";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faCaretDown } from "@fortawesome/free-solid-svg-icons";
  library.add(faCaretDown);
  dom.watch();

  const dispatch = createEventDispatcher();

  
  
  /** 表示値 */
  let value: string = $state();

  interface Props {
    /** 選択アイテムのリスト */
    items: { key: string; value: string }[];
    /** 表示値のキー */
    key: string;
    tabindex: number | null | undefined;
  }

  let { items, key, tabindex }: Props = $props();
  /** ドロップダウン トグルイベント */
  const onToggle = () =>
    document.getElementById("dropdown")?.classList.toggle("is-active");

  /** アイテム変更イベント */
  const onChange = (value: string) => {
    dispatch("change", value);
    onToggle();
  };

  run(() => {
    for (let i = 0; i < items.length; i++) {
      const p = items[i];
      if (p.key === key) {
        value = p.value;
      }
    }
  });
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
