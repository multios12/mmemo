<script lang="ts">
  import { paragraphs, onParagraphChange } from "./ToolbarEvent";
  import { REDO_COMMAND, UNDO_COMMAND, type LexicalEditor } from "lexical";
  import { FORMAT_TEXT_COMMAND } from "lexical";
  import { onMount } from "svelte";
  import Dropdown from "../Dropdown.svelte";

  /** lexical Editor */
  export let editor: LexicalEditor;
  export let para: string;
  /** アンドゥボタン表示状態 */
  export let canUndo: boolean;
  /** リドゥボタン表示状態 */
  export let canRedo: boolean;

  /** 選択ノードを太字に変更 */
  const formatBold = () => editor.dispatchCommand(FORMAT_TEXT_COMMAND, "bold");
  /** 選択ノードをイタリックに変更 */
  const formatItalic = () =>
    editor.dispatchCommand(FORMAT_TEXT_COMMAND, "italic");
  /** 選択ノードを下線有に変更 */
  const formatUnderline = () =>
    editor.dispatchCommand(FORMAT_TEXT_COMMAND, "underline");
  /** アンドゥコマンド実行 */
  export const undo = () => editor.dispatchCommand(UNDO_COMMAND, undefined);
  /** リドゥコマンド実行 */
  export const redo = () => editor.dispatchCommand(REDO_COMMAND, undefined);

  /** マウントイベント */
  onMount(async () => {
    const el = <HTMLButtonElement>document.getElementById("para-button");
  });

  /** 段落ドロップダウン 値変更イベント */
  const onChange = (e: CustomEvent<any>) => onParagraphChange(editor, e.detail);
</script>

<div class="panel-heading py-1">
  <div class="level is-mobile">
    <div class="level-left">
      <div class="level-item">
        <button class="button is-ghost" disabled={!canUndo} on:click={undo}>
          <i class="material-icons">undo</i>
        </button>
      </div>
      <div class="level-item">
        <button class="button is-ghost" disabled={!canRedo} on:click={redo}>
          <i class="material-icons">redo</i>
        </button>
      </div>

      <div class="level-item">
        <Dropdown items={paragraphs} key={para} on:change={onChange} />
      </div>

      <div class="level-item">
        <button id="boldButton" class="button is-ghost" on:click={formatBold}>
          <i class="material-icons">format_bold</i>
        </button>
      </div>
      <div class="level-item">
        <button
          id="italicButton"
          class="button is-ghost"
          on:click={formatItalic}
        >
          <i class="material-icons">format_italic</i>
        </button>
      </div>
      <div class="level-item">
        <button
          id="underButton"
          class="button is-ghost"
          on:click={formatUnderline}
        >
          <i class="material-icons">format_underlined</i>
        </button>
      </div>
    </div>
  </div>
</div>
