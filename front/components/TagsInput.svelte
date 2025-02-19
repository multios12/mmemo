<script lang="ts">
  interface Props {
    items?: string[];
  }

  let { items = $bindable([]) }: Props = $props();
  let value = $state("");

  const onKeydown = (e: KeyboardEvent) => {
    if (e.code == "Enter" && value != "" && e.isComposing == false) {
      if (!items.includes(value)) {
        var i = items;
        i.push(value);
        items = i;
        value = "";
      }
    }
  };
  const onBlur = () => {
    if (value == "" || items.includes(value)) {
      return;
    }
    var i = items;
    i.push(value);
    items = i;
    value = "";
  };
  const deleteClick = (e: MouseEvent) => {
    const t = e.target as HTMLButtonElement;
    items = items.filter((value) => value != t.dataset.value);
  };
</script>

<div class="field has-addons">
  <div class="control">
    <div class="tags are-medium">
      {#each items as i}<span class="tag is-rounded"
          >{i}<button class="delete" data-value={i} onclick={deleteClick}
          ></button></span
        >{/each}
    </div>
  </div>
  <div class="control is-expanded">
    <input
      class="input"
      type="text"
      placeholder="Choose Tags"
      bind:value
      onkeydown={onKeydown}
      onblur={onBlur}
    />
  </div>
</div>

<style>
</style>
