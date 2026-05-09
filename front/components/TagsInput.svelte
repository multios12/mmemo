<script lang="ts">
  interface Props {
    items?: string[];
    inputId?: string;
  }

  let { items = $bindable([]), inputId = "tagsInput" }: Props = $props();
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
          >{i}<button
            class="delete"
            aria-label={`delete tag ${i}`}
            data-value={i}
            onclick={deleteClick}
          ></button></span
        >{/each}
    </div>
  </div>
  <div class="control is-expanded">
    <input
      id={inputId}
      class="input"
      type="text"
      placeholder="タグ（複数指定可）"
      bind:value
      onkeydown={onKeydown}
      onblur={onBlur}
    />
  </div>
</div>

<style>
</style>
