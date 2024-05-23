import { $createParagraphNode, $getSelection, $isRangeSelection, type LexicalEditor } from "lexical";
import { $createHeadingNode, $createQuoteNode, type HeadingTagType } from "@lexical/rich-text"
import { $createCodeNode } from "@lexical/code";
import { insertList } from "@lexical/list"
import { $setBlocksType } from "@lexical/selection";

export const paragraphs = [
  {
    key: "",
    value: "normal",
  },
  {
    key: "ol",
    value: "number list",
  },
  {
    key: "ul",
    value: "bullet list",
  },
  {
    key: "quote",
    value: "quote",
  },
  {
    key: "code",
    value: "code",
  },
  {
    key: "h1",
    value: "heading h1",
  },
  {
    key: "h2",
    value: "heading h2",
  },
  {
    key: "h3",
    value: "heading h3",
  },
];

/** 段落セレクトボックス 値変更イベント */
export const onParagraphChange = (editor: LexicalEditor, value: string) => {
  switch (value) {
    case "h1":
    case "h2":
    case "h3":
      editor.update(() => {
        const selection = $getSelection();
        if ($isRangeSelection(selection)) {
          $setBlocksType(selection, () => $createHeadingNode(<HeadingTagType>value));
        }
      });
      break;
    case "ol": // number list
      insertList(editor, "number");
      // editor.dispatchCommand(INSERT_ORDERED_LIST_COMMAND, undefined);
      break;
    case "ul": // bullet list
      insertList(editor, "bullet");
      // editor.dispatchCommand(INSERT_UNORDERED_LIST_COMMAND, undefined);
      break;
    case "quote": // quote
      formatQuote(editor);
      break;
    case "code": // code
      formatCode(editor);
      break;

    default:
      formatParagraph(editor);
      break;
  }
};

/** 選択ノードをクオートに変更 */
const formatQuote = (editor: LexicalEditor) => {
  editor.update(() => {
    const selection = $getSelection();
    if ($isRangeSelection(selection)) {
      $setBlocksType(selection, () => $createQuoteNode());
    }
  });
}

/** 選択ノードをコードに変更 */
const formatCode = (editor: LexicalEditor) => {
  editor.update(() => {
    const selection = $getSelection();
    if ($isRangeSelection(selection)) {
      $setBlocksType(selection, () => $createCodeNode());
    }
  });
}

/** 選択ノードを段落に変更 */
const formatParagraph = (editor: LexicalEditor) => {
  editor.update(() => {
    const selection = $getSelection();
    if ($isRangeSelection(selection)) {
      $setBlocksType(selection, () => $createParagraphNode());
    }
  });
};
