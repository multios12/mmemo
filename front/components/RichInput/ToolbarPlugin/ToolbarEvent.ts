import { $createParagraphNode, $getSelection, $isRangeSelection, ElementNode, type LexicalEditor } from "lexical";
import { $createHeadingNode, $createQuoteNode, HeadingNode, type HeadingTagType } from "@lexical/rich-text"
import { $createCodeNode } from "@lexical/code";
import { insertList, ListNode } from "@lexical/list"
import { LinkNode, $isLinkNode } from "@lexical/link";
import { $setBlocksType } from "@lexical/selection";

/** 段落選択値 */
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

/** ツールバー更新イベント */
export const updateToolbar = (): { para: string, linkValue: string | undefined } => {
  const selection = $getSelection();
  if (!$isRangeSelection(selection)) {
    return { para: "", linkValue: undefined };
  }
  const node = selection.getNodes()[0];
  const p = node.getParent() as ElementNode;

  if (p === null) {
    return { para: "", linkValue: "" };
  }

  let para: string
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
  const isLink = $isLinkNode(p);
  updateButton("linkButton", isLink);
  const menuElement = document.getElementById("link-menu") as HTMLDivElement;
  let linkValue: string | undefined
  if (isLink) {
    menuElement.classList.remove("is-hidden");
    const linkNode = <LinkNode>p;
    linkValue = linkNode.getURL();
  } else {
    menuElement.classList.add("is-hidden");
  }
  return { para, linkValue }
};

/** ツールバーボタンの表示を更新 */
const updateButton = (buttonId: string, isSelect: boolean) => {
  const button = <HTMLButtonElement>document.getElementById(buttonId);
  if (isSelect) {
    button.classList.add("is-underlined");
    button.classList.add("has-text-warning");
  } else {
    button.classList.remove("is-underlined");
    button.classList.remove("has-text-warning");
  }
};

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
