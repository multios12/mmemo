import { $createParagraphNode, $getSelection, $isRangeSelection, ElementNode, createEditor, type LexicalEditor } from "lexical";
import { FORMAT_TEXT_COMMAND, REDO_COMMAND, UNDO_COMMAND } from "lexical";
import { $createHeadingNode, $createQuoteNode, HeadingNode, QuoteNode, registerRichText, type HeadingTagType } from "@lexical/rich-text"
import { $createCodeNode, CodeNode } from "@lexical/code";
import { ListNode, ListItemNode, insertList } from "@lexical/list"
import { $setBlocksType } from "@lexical/selection";
import { createEmptyHistoryState, registerHistory } from "@lexical/history";
import { mergeRegister } from "@lexical/utils";

/** lexical Editor */
export let editor: LexicalEditor;
/** 選択ノードの段落種別 */
export let para: string;

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
}

/** ツールバー更新イベント */
export const updateToolbar = () => {
  const selection = $getSelection();
  if (!$isRangeSelection(selection)) {
    return
  }
  const node = selection.getNodes()[0];
  const p = node.getParent() as ElementNode;

  if (p === null) {
    return
  }

  if (p.getType() === "heading") {
    para = (p as HeadingNode).getTag()
  } else if (p.getType() === "listitem") {
    para = (p.getParent() as ListNode)?.getTag()
  } else if (p.getType() === "quote") {
    para = "quote"
  } else if (p.getType() === "code") {
    para = "code"
  } else {
    para = "0"
  }
  const s = document.getElementById("paragraph") as HTMLSelectElement
  s.value = para

  updateButton("boldButton", selection.hasFormat('bold'))
  updateButton("italicButton", selection.hasFormat('italic'))
  updateButton("underButton", selection.hasFormat('underline'))
};

/** ツールバーボタンの表示を更新 */
const updateButton = (buttonId: string, isSelect: boolean) => {
  const button = document.getElementById(buttonId) as HTMLButtonElement;
  if (isSelect) {
    button.classList.add("is-light")
  } else {
    button.classList.remove("is-light")
  }
}

/** 段落セレクトボックス 値変更イベント */
export const onParagraphChange = (): string => {
  let p = document.getElementById("paragraph") as HTMLSelectElement;
  switch (p.value) {
    case "h1":
    case "h2":
    case "h3":
      editor.update(() => {
        const selection = $getSelection();
        if ($isRangeSelection(selection)) {
          let tagType = p.value as HeadingTagType;
          $setBlocksType(selection, () => $createHeadingNode(tagType));
        }
      });
      break;
    case "ol": // number list
      insertList(editor, "number");
      // editor.dispatchCommand(INSERT_ORDERED_LIST_COMMAND, undefined);
      break;
    case "ul": //bullet list
      insertList(editor, "bullet");
      // editor.dispatchCommand(INSERT_UNORDERED_LIST_COMMAND, undefined);
      break;
    case "quote": //        <option value="4"> quote </option>
      formatQuote();
      break;
    case "code": //        <option value="5"> code </option>
      formatCode();
      break;

    default:
      formatParagraph()
      break;
  }
  return para;
};

/** 選択ノードをクオートに変更 */
const formatQuote = () => {
  editor.update(() => {
    const selection = $getSelection();
    if ($isRangeSelection(selection)) {
      $setBlocksType(selection, () => $createQuoteNode());
    }
  });
}

/** 選択ノードをコードに変更 */
const formatCode = () => {
  editor.update(() => {
    const selection = $getSelection();
    if ($isRangeSelection(selection)) {
      $setBlocksType(selection, () => $createCodeNode());
    }
  });
}

/** 選択ノードを段落に変更 */
const formatParagraph = () => {
  editor.update(() => {
    const selection = $getSelection();
    if ($isRangeSelection(selection)) {
      $setBlocksType(selection, () => $createParagraphNode());
    }
  });
};

/** 選択ノードを太字に変更 */
export const formatBold = () => editor.dispatchCommand(FORMAT_TEXT_COMMAND, "bold")
/** 選択ノードをイタリックに変更 */
export const formatItalic = () => editor.dispatchCommand(FORMAT_TEXT_COMMAND, "italic")
/** 選択ノードを下線有に変更 */
export const formatUnderline = () => editor.dispatchCommand(FORMAT_TEXT_COMMAND, "underline")
/** アンドゥコマンド実行 */
export const undo = () => editor.dispatchCommand(UNDO_COMMAND, undefined)
/** リドゥコマンド実行 */
export const redo = () => editor.dispatchCommand(REDO_COMMAND, undefined)
