<script lang="ts">
  import { onMount } from "svelte";
  import { InitializeLixicalEditor, editor } from "./ToolbarEvent";
  import ToolbarPlugin from "./ToolbarPlugin.svelte";
  import { mergeRegister } from "@lexical/utils";
  import { CAN_REDO_COMMAND, CAN_UNDO_COMMAND } from "lexical";
  import {
    $convertFromMarkdownString as _convertFromMarkdownString,
    $convertToMarkdownString as _convertToMarkdownString,
    TRANSFORMERS,
  } from "@lexical/markdown";

  /** 入力値 */
  export let value: string;
  /** アンドゥボタン表示状態 */
  let canUndo: boolean;
  /** リドゥボタン表示状態 */
  let canRedo: boolean;
  const LowPriority = 1;

  /** マウントイベント */
  onMount(async () => {
    const detailElement = document.getElementById("detail") as HTMLElement;
    await InitializeLixicalEditor(detailElement);

    mergeRegister(
      editor.registerCommand(
        CAN_UNDO_COMMAND,
        (payload) => {
          canUndo = payload;
          return false;
        },
        LowPriority,
      ),
      editor.registerCommand(
        CAN_REDO_COMMAND,
        (payload) => {
          canRedo = payload;
          return false;
        },
        LowPriority,
      ),
    );

    if (value !== null) {
      console.log(`初期表示:${value}`);
      editor.update(() => _convertFromMarkdownString(value, TRANSFORMERS));
    }
    // editor.update(prepopulatedRichText, { tag: "history-merge" });

    // 更新タイミングでのmarkdown出力
    editor.registerUpdateListener(({}) => {
      editor.update(() => (value = _convertToMarkdownString(TRANSFORMERS)));
    });
  });

  $: {
    if (value !== null && value !== "") {
      editor.update(() => _convertFromMarkdownString(value, TRANSFORMERS));
    }
  }
</script>

<div class="panel">
  <ToolbarPlugin {canUndo} {canRedo} />
  <div class="box px-0">
    <div class="panel-block p-0">
      <div id="detail" class="content editor-input" contenteditable></div>
    </div>
  </div>
</div>

<style>
  #detail {
    width: 100%;
  }
</style>
