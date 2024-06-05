import { type TextMatchTransformer } from "@lexical/markdown"
import { $createImageNode, $isImageNode, ImageNode } from "./ImagesPlugin/ImageNode.js"

/**
 * 画像ノード変換
 * ```
 *    ![エビフライトライアングル](http://i.imgur.com/Jjwsc.jpg "サンプル")
 *    ![Altのテキスト](/path/to/img.jpg)
 *    ![Altのテキスト](/path/to / img.png "タイトル")
 * ```
 */
export const IMAGE: TextMatchTransformer = {
    dependencies: [ImageNode],
    export: (node, exportChildren, exportFormat) => {
        if (!$isImageNode(node)) {
            return null;
        }
        const title = undefined;
        const linkContent = title
            ? `![${node.getAltText()}](${node.getSrc()} "${title}")`
            : `![${node.getAltText()}](${node.getSrc()})`;

        return linkContent;
    },
    importRegExp:
        /(?:!\[([^[]+)\])(?:\((?:([^()\s]+)(?:\s"((?:[^"]*\\")*[^"]*)"\s*)?)\))/,
    regExp:
        /(?:!\[([^[]+)\])(?:\((?:([^()\s]+)(?:\s"((?:[^"]*\\")*[^"]*)"\s*)?)\))$/,
    replace: (textNode, match) => {
        const [, linkText, linkSrc, linkTitle] = match;
        const linkNode = $createImageNode({ src: linkSrc, altText: linkText });
        textNode.replace(linkNode);
    },
    trigger: ')',
    type: 'text-match',
};