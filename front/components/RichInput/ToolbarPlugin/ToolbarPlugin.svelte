<script lang="ts">
  import { paragraphs, onParagraphChange } from "./ToolbarEvent.js";
  import { REDO_COMMAND, UNDO_COMMAND, type LexicalEditor } from "lexical";
  import { FORMAT_TEXT_COMMAND } from "lexical";
  import { TOGGLE_LINK_COMMAND } from "@lexical/link";
  import { onMount } from "svelte";
  import Dropdown from "../../Dropdown.svelte";
  import ImageButton from "./ImageButton.svelte";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faRotateLeft } from "@fortawesome/free-solid-svg-icons";
  import { faRotateRight } from "@fortawesome/free-solid-svg-icons";
  import {
    faBold,
    faItalic,
    faUnderline,
    faLink,
  } from "@fortawesome/free-solid-svg-icons";
  library.add(faRotateLeft, faRotateRight);
  library.add(faBold, faItalic, faUnderline, faLink);
  dom.watch();

  /** lexical Editor */
  export let editor: LexicalEditor;
  /** アンドゥボタン表示状態 */
  export let canUndo: boolean;
  /** リドゥボタン表示状態 */
  export let canRedo: boolean;
  /** 段落種別 */
  export let para: string;

  /** 選択ノードを太字に変更 */
  const formatBold = () => editor.dispatchCommand(FORMAT_TEXT_COMMAND, "bold");
  /** 選択ノードをイタリックに変更 */
  const formatItalic = () =>
    editor.dispatchCommand(FORMAT_TEXT_COMMAND, "italic");
  /** 選択ノードを下線有に変更 */
  const formatUnderline = () =>
    editor.dispatchCommand(FORMAT_TEXT_COMMAND, "underline");
  /** リンクを更新 */
  const toggleLink = () =>
    editor.dispatchCommand(TOGGLE_LINK_COMMAND, "http://");
  /** アンドゥコマンド実行 */
  const undo = () => editor.dispatchCommand(UNDO_COMMAND, undefined);
  /** リドゥコマンド実行 */
  const redo = () => editor.dispatchCommand(REDO_COMMAND, undefined);

  /** マウントイベント */
  onMount(async () => {
    const el = <HTMLButtonElement>document.getElementById("para-button");
  });

  /** 段落ドロップダウン 値変更イベント */
  const onChange = (e: CustomEvent<any>) => onParagraphChange(editor, e.detail);
</script>

<div id="toolbar" class="panel-heading py-1">
  <div class="level is-mobile">
    <div class="level-left">
      <div class="level-item">
        <!-- svelte-ignore a11y_consider_explicit_label -->
        <button
          class="button is-ghost p-0"
          disabled={!canUndo}
          tabindex="-1"
          on:click={undo}
        >
          <i class="fa-solid fa-rotate-left"></i>
        </button>
      </div>
      <div class="level-item">
        <!-- svelte-ignore a11y_consider_explicit_label -->
        <button
          class="button is-ghost p-0"
          disabled={!canRedo}
          tabindex="-1"
          on:click={redo}
        >
          <i class="fa-solid fa-rotate-right"></i>
        </button>
      </div>

      <div class="level-item">
        <Dropdown
          items={paragraphs}
          key={para}
          tabindex={-1}
          on:change={onChange}
        />
      </div>

      <div class="level-item">
        <!-- svelte-ignore a11y_consider_explicit_label -->
        <button
          id="boldButton"
          class="button is-ghost p-0"
          tabindex="-1"
          on:click={formatBold}
        >
          <i class="fa-solid fa-bold"></i>
        </button>
      </div>
      <div class="level-item is-hidden">
        <!-- svelte-ignore a11y_consider_explicit_label -->
        <button
          id="italicButton"
          class="button is-ghost p-0"
          tabindex="-1"
          on:click={formatItalic}
        >
          <i class="fa-solid fa-italic"></i>
        </button>
      </div>
      <div class="level-item is-hidden">
        <!-- svelte-ignore a11y_consider_explicit_label -->
        <button
          id="underButton"
          class="button is-ghost p-0"
          tabindex="-1"
          on:click={formatUnderline}
        >
          <i class="fa-solid fa-underline"></i>
        </button>
      </div>
      <div class="level-item">
        <!-- svelte-ignore a11y_consider_explicit_label -->
        <button
          id="linkButton"
          class="button is-ghost p-0"
          tabindex="-1"
          on:click={toggleLink}
        >
          <i class="fa-solid fa-link"></i>
        </button>
      </div>
      <div class="level-item">
        <ImageButton {editor} />
      </div>
    </div>
  </div>
</div>

<style>
  .level-left {
    flex-direction: row;
  }
</style>
