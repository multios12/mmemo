<script lang="ts">
  import { onMount } from "svelte";
  import { onParagraphChange } from "./ToolbarEvent";
  import { redo, undo } from "./ToolbarEvent";
  import { formatBold, formatItalic, formatUnderline } from "./ToolbarEvent";
  /** 段落セレクトボックス選択値 */
  let para: string;
  /** アンドゥボタン表示状態 */
  export let canUndo: boolean;
  /** リドゥボタン表示状態 */
  export let canRedo: boolean;

  /** マウントイベント */
  onMount(async () => {});
  /** 段落変更イベント */
  const onChange = () => (para = onParagraphChange());
  /*
  const onLink = () => {
    const element = document.getElementById("linkModal") as HTMLElement;
    if (!element.classList.contains("is-active")) {
      element.classList.add("is-active");
    } else {
      element.classList.remove("is-active");
    }
  };
  */
</script>

<div class="panel-heading p-0">
  <nav class="navbar">
    <button class="button" disabled={!canUndo} on:click={undo}>
      <i class="material-icons">undo</i>
    </button>
    <button class="button" disabled={!canRedo} on:click={redo}>
      <i class="material-icons">redo</i>
    </button>
    <div class="select">
      <select id="paragraph" on:change={() => onChange()} bind:value={para}>
        <option value="0"> normal </option>
        <option value="ol"> number list </option>
        <option value="ul"> bulleonParagraphChange> </option><option
          value="quote"
        >
          quote
        </option>
        <option value="code"> code </option>
        <option value="h1"> heading h1 </option>
        <option value="h2"> heading h2 </option>
        <option value="h3"> heading h3 </option>
      </select>
    </div>
    <button id="boldButton" class="button is-light" on:click={formatBold}>
      <i class="material-icons">format_bold</i>
    </button>
    <button id="italicButton" class="button" on:click={formatItalic}>
      <i class="material-icons">format_italic</i>
    </button>
    <button id="underButton" class="button" on:click={formatUnderline}>
      <i class="material-icons">format_underlined</i>
    </button>

    <!--
    <button class="button" data-target="linkModal" disabled>
      <i class="material-icons">add_link</i>
    </button>
    <button class="button" disabled>
      <i class="material-icons">add_photo_alternate</i>
    </button>
    <button class="button" disabled>
      <i class="material-icons">horizontal_rule</i>
    </button>
   -->
  </nav>

  <!--
  <div id="linkModal" class="modal">
    <div class="modal-background"></div>
    <div class="modal-content">
      <header class="modal-card-head">
        <p class="modal-card-title">Add link</p>
        <button class="delete" aria-label="close"></button>
      </header>
      <section class="modal-card-body">
        <input type="url" class="input" />
      </section>
      <footer class="modal-card-foot">
        <div class="buttons">
          <button class="button is-success">Add</button>
        </div>
      </footer>
    </div>
  </div>
  -->
</div>
