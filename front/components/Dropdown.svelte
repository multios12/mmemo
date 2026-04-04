<script lang="ts">
  interface Props {
    /** 選択アイテムのリスト */
    items: { key: string; value: string }[];
    /** 表示値のキー */
    key?: string;
    tabindex: number | null | undefined;
    onchange?: (value: string) => void;
  }

  let { items, key = $bindable(""), tabindex, onchange }: Props = $props();
  $effect(() => {
    if (items.length > 0 && !items.some((item) => item.key === key)) {
      key = items[0].key;
    }
  });

  const onChange = (event: Event) => {
    const nextKey = (event.currentTarget as HTMLSelectElement).value;
    key = nextKey;
    onchange?.(nextKey);
  };
</script>

<div class="select">
  <select bind:value={key} {tabindex} onchange={onChange}>
    {#each items as item}
      <option value={item.key}>{item.value}</option>
    {/each}
  </select>
</div>

<style>
  .select {
    display: flex;
    align-items: center;
  }

  .select select {
    min-height: 2.35rem;
    line-height: 1.2;
    padding-top: 0;
    padding-bottom: 0;
  }
</style>
