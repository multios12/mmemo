<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faCaretDown } from "@fortawesome/free-solid-svg-icons";
  library.add(faCaretDown);
  dom.watch();

  const dispatch = createEventDispatcher();

  /** 選択アイテムのリスト */
  export let items: { key: string; value: string }[];
  /** 表示値のキー */
  export let key: string;
  /** 表示値 */
  let value: string;

  export let tabindex: number | null | undefined;
  /** ドロップダウン トグルイベント */
  const onToggle = () =>
    document.getElementById("dropdown")?.classList.toggle("is-active");

  /** アイテム変更イベント */
  const onChange = (value: string) => {
    dispatch("change", value);
    onToggle();
  };

  $: {
    for (let i = 0; i < items.length; i++) {
      const p = items[i];
      if (p.key === key) {
        value = p.value;
      }
    }
  }
</script>

<div class="dropdown" id="dropdown">
  <div class="dropdown-trigger">
    <button
      class="button"
      aria-haspopup="true"
      aria-controls="dropdown-menu"
      id="para-button"
      {tabindex}
      on:click={onToggle}
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
        <button class="dropdown-item" on:click={() => onChange(i.key)}>
          {i.value}
        </button>
      {/each}
    </div>
  </div>
</div>
