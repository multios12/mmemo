import { mergeRegister } from "@lexical/utils";
import { $isRangeSelection as _isRangeSelection } from "lexical";
import { $createParagraphNode as _createParagraphNode } from "lexical";
import { $getSelection as _getSelection } from "lexical";
import { createEditor } from "lexical";
import { CAN_REDO_COMMAND, CAN_UNDO_COMMAND } from "lexical";
import { $isLinkNode as _isLinkNode, LinkNode } from "@lexical/link";
import { $toggleLink as _toggleLink } from "@lexical/link";
import { TOGGLE_LINK_COMMAND } from "@lexical/link";
import { $setBlocksType as _setBlocksType } from "@lexical/selection";
import { $convertFromMarkdownString as _convertFromMarkdownString } from "@lexical/markdown";
import { $convertToMarkdownString as _convertToMarkdownString } from "@lexical/markdown";
import { $createQuoteNode as _createQuoteNode } from "@lexical/rich-text";
import { HeadingNode, QuoteNode, registerRichText } from "@lexical/rich-text";
import { $createCodeNode as _createCodeNode, CodeNode } from "@lexical/code";
import { ListNode, ListItemNode } from "@lexical/list";
import { createEmptyHistoryState, registerHistory } from "@lexical/history";
import { TRANSFORMERS } from "@lexical/markdown";
import type { UpdateListener } from "node_modules/lexical/LexicalEditor.js";
import type { EventDispatcher } from "svelte";
import { ImageNode } from "./ImagesPlugin/ImageNode.js"
import { registerInsertImageCommand } from "./ImagesPlugin/index.js";
import { IMAGE } from "./MarkdownTransformers.js"
import { LineBreakPlugin } from "./LineBreakCommand.js";
const LowPriority = 1;

/** LexicalEditorの初期化 */
export const InitialEditor = async (
  detailElement: HTMLElement,
  richElement: HTMLDivElement,
  menuElement: HTMLDivElement,
  doCanUndo: (p: boolean) => boolean,
  doCanRedo: (p: boolean) => boolean,
  toolBarListener: UpdateListener,
  dispatch: EventDispatcher<any>
) => {
  // LexicalEditorの設定
  const initialConfig = {
    namespace: "Detail",
    nodes: [
      CodeNode,
      HeadingNode,
      ListNode,
      ListItemNode,
      QuoteNode,
      LinkNode,
      ImageNode,
    ],
    onError: (error: Error) => {
      throw error;
    },
  };
  let editor = createEditor(initialConfig);
  editor.setRootElement(detailElement);

  // プラグイン登録
  mergeRegister(
    // リッチテキスト登録
    registerRichText(editor),
    // ヒストリ登録
    registerHistory(editor, createEmptyHistoryState(), 300),
    // リンクノード選択時の状態変更リスナ登録
    editor.registerMutationListener(LinkNode, (mutations) => {
      const registeredElements: WeakSet<HTMLElement> = new WeakSet();
      editor.getEditorState().read(() => {
        for (const [key, mutation] of mutations) {
          const element: null | HTMLElement = editor.getElementByKey(key);
          if (
            (mutation === "created" || mutation === "updated") &&
            element !== null &&
            !registeredElements.has(element)
          ) {
            registeredElements.add(element);

            menuElement.classList.remove("is-hidden");

            element.addEventListener("click", () => {
              const pos = getAbsolutePosition(element);
              const richPos = getAbsolutePosition(richElement);
              const posY = pos.height + pos.top - richPos.top + 10;
              menuElement.style.left = pos.left + "px";
              menuElement.style.top = posY + "px";
            });
          } else {
            menuElement.classList.add("is-hidden");
          }
        }
      });
    }),
    // 更新時のマークダウン出力リスナ登録
    editor.registerUpdateListener(({ }) => {
      editor.update(() => {
        const trans = TRANSFORMERS
        trans.unshift(IMAGE);
        const markdown = _convertToMarkdownString(trans);
        dispatch("textChange", markdown);
      });
    }),
    // ツールバー状態更新リスナ登録
    editor.registerUpdateListener(toolBarListener),
  );

  // コマンド登録
  const linkCommand = (p: string | null) => {
    _toggleLink(<string | null>p);
    return true;
  }
  mergeRegister(
    editor.registerCommand(CAN_UNDO_COMMAND, (p) => doCanUndo(p), LowPriority),
    editor.registerCommand(CAN_REDO_COMMAND, (p) => doCanRedo(p), LowPriority),
    editor.registerCommand(TOGGLE_LINK_COMMAND, linkCommand, LowPriority),
    registerInsertImageCommand(editor),
    LineBreakPlugin(editor),
  );

  // デバッグ情報出力リスナ登録
  const stateElement = <HTMLTextAreaElement>(
    document.getElementById("lexical-state")
  );
  if (stateElement !== null) {
    editor.registerUpdateListener(({ editorState }) => {
      stateElement!.value = JSON.stringify(
        editorState.toJSON(),
        undefined,
        2,
      );
    });
  }

  return editor
}

function getAbsolutePosition(elm: HTMLElement): {
  left: number;
  top: number;
  height: number;
  width: number;
} {
  const { left, top, height, width } = elm.getBoundingClientRect();
  const { left: bleft, top: btop } = document.body.getBoundingClientRect();
  return { left: left - bleft, top: top - btop, height, width };
}
