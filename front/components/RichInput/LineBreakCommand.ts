import { $createLineBreakNode, $createParagraphNode, $getSelection, $isRangeSelection } from "lexical"
import type { LexicalEditor } from "lexical"
import { COMMAND_PRIORITY_LOW, KEY_ENTER_COMMAND } from "lexical"
import { } from "@lexical/rich-text"
import { $isListItemNode } from "@lexical/list"

/** Enterキー押下時の動作を変更するコマンド */
export const LineBreakPlugin = (editor: LexicalEditor) => editor.registerCommand(
    KEY_ENTER_COMMAND,
    (payload) => {
        if (!payload) {
            return true
        }
        const event: KeyboardEvent = payload;
        event.preventDefault();
        // 選択範囲の検出
        const selection = $getSelection();
        if (!selection) {
            return false
        }

        if (!$isRangeSelection(selection)) {
            return false
        }

        let node = selection.anchor.getNode()
        if ($isListItemNode(node)) {
            // アンカーがListItemを指している場合、ListItemが空と判断して、次段落を作成
            selection.insertNodes([$createParagraphNode()])
        } else if ($isListItemNode(node.getParent())) {
            // アンカーの親が、ListItemの場合、ListItemの中での入力があるとみなして、何もせず次のコマンドを実行
            return false
        } else {
            if (!event.shiftKey) {
                // シフトキーが押されていない場合、改行
                selection.insertNodes([$createLineBreakNode()])
            } else {
                // シフトキーが押されている場合、次段落を作成
                selection.insertNodes([$createParagraphNode()])
            }
        }

        return true
    },
    COMMAND_PRIORITY_LOW
);
