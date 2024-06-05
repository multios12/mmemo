import { $insertNodes, COMMAND_PRIORITY_EDITOR, createCommand, type LexicalCommand, type LexicalEditor } from "lexical";
import {
    $createImageNode,
    type ImagePayload,
} from './ImageNode.js';

export type InsertImagePayload = Readonly<ImagePayload>;

export const registerInsertImageCommand = (editor: LexicalEditor) => {
    return editor.registerCommand<InsertImagePayload>(
        INSERT_IMAGE_COMMAND,
        (payload) => {
            const imageNode = $createImageNode(payload);
            $insertNodes([imageNode]);
            //if ($isRootOrShadowRoot(imageNode.getParentOrThrow())) {
            //    $wrapNodeInElement(imageNode, $createParagraphNode).selectEnd();
            //}

            return true;
        },
        COMMAND_PRIORITY_EDITOR,
    )
}


export const INSERT_IMAGE_COMMAND: LexicalCommand<InsertImagePayload> =
    createCommand('INSERT_IMAGE_COMMAND');
