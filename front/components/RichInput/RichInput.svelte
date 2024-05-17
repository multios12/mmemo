<script lang="ts">
  import { onMount } from "svelte";
  import ToolbarPlugin from "./ToolbarPlugin.svelte";
  import { mergeRegister } from "@lexical/utils";
  import { $isRangeSelection as _isRangeSelection } from "lexical";
  import { $createParagraphNode as _createParagraphNode } from "lexical";
  import { $getSelection as _getSelection } from "lexical";
  import { ElementNode, createEditor, type LexicalEditor } from "lexical";
  import { FORMAT_TEXT_COMMAND } from "lexical";
  import { CAN_REDO_COMMAND, CAN_UNDO_COMMAND } from "lexical";
  import { $setBlocksType as _setBlocksType } from "@lexical/selection";
  import { $convertFromMarkdownString as _convertFromMarkdownString } from "@lexical/markdown";
  import { $convertToMarkdownString as _convertToMarkdownString } from "@lexical/markdown";
  import { TRANSFORMERS } from "@lexical/markdown";
  import { $createQuoteNode as _createQuoteNode } from "@lexical/rich-text";
  import { HeadingNode, QuoteNode, registerRichText } from "@lexical/rich-text";
  import { $createCodeNode as _createCodeNode, CodeNode } from "@lexical/code";
  import { ListNode, ListItemNode } from "@lexical/list";
  import { createEmptyHistoryState, registerHistory } from "@lexical/history";

  /** 入力値 */
  export let value: string;
  /** lexical Editor */
  let editor: LexicalEditor;
  /** 選択ノードの段落種別 */
  let para: string;
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
      editor.update(() => _convertFromMarkdownString(value, TRANSFORMERS));
    }

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
  /** LexicalEditorの初期化 */
  export const InitializeLixicalEditor = async (detailElement: HTMLElement) => {
    // LexicalEditorの設定
    const initialConfig = {
      namespace: "Detail",
      nodes: [CodeNode, HeadingNode, ListNode, ListItemNode, QuoteNode],
      onError: (error: Error) => {
        throw error;
      },
    };
    editor = createEditor(initialConfig);
    editor.setRootElement(detailElement);

    // プラグイン登録
    mergeRegister(
      editor.registerUpdateListener(({ editorState }) => {
        editorState.read(updateToolbar);
      }),
      registerRichText(editor),
      registerHistory(editor, createEmptyHistoryState(), 300),
    );

    // デバッグ情報出力
    // const stateRef = document.getElementById("lexical-state") as HTMLTextAreaElement;
    //    editor.registerUpdateListener(({ editorState }) => {
    //      stateRef!.value = JSON.stringify(editorState.toJSON(), undefined, 2);
    //    });

    return editor;
  };

  /** ツールバー更新イベント */
  export const updateToolbar = () => {
    const selection = _getSelection();
    if (!_isRangeSelection(selection)) {
      return;
    }
    const node = selection.getNodes()[0];
    const p = node.getParent() as ElementNode;

    if (p === null) {
      return;
    }

    if (p.getType() === "heading") {
      para = (p as HeadingNode).getTag();
    } else if (p.getType() === "listitem") {
      para = (p.getParent() as ListNode)?.getTag();
    } else if (p.getType() === "quote") {
      para = "quote";
    } else if (p.getType() === "code") {
      para = "code";
    } else {
      para = "0";
    }

    updateButton("boldButton", selection.hasFormat("bold"));
    updateButton("italicButton", selection.hasFormat("italic"));
    updateButton("underButton", selection.hasFormat("underline"));
  };

  /** ツールバーボタンの表示を更新 */
  const updateButton = (buttonId: string, isSelect: boolean) => {
    const button = <HTMLButtonElement>document.getElementById(buttonId);
    if (isSelect) {
      button.classList.add("is-light");
    } else {
      button.classList.remove("is-light");
    }
  };

  /** 選択ノードを太字に変更 */
  export const formatBold = () =>
    editor.dispatchCommand(FORMAT_TEXT_COMMAND, "bold");
  /** 選択ノードをイタリックに変更 */
  export const formatItalic = () =>
    editor.dispatchCommand(FORMAT_TEXT_COMMAND, "italic");
  /** 選択ノードを下線有に変更 */
  export const formatUnderline = () =>
    editor.dispatchCommand(FORMAT_TEXT_COMMAND, "underline");
</script>

<div class="panel is-dark">
  <ToolbarPlugin {editor} {para} {canUndo} {canRedo} />
  <div class="panel-block p-0">
    <div id="detail" class="content editor-input" contenteditable></div>
  </div>
</div>

<style>
  #detail {
    width: 100%;
  }
</style>
