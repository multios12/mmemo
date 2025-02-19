<!-- @migration-task Error while migrating Svelte code: Unexpected token `}`. Did you mean `&rbrace;` or `{"}"}`?
https://svelte.dev/e/js_parse_error -->
<!-- @migration-task Error while migrating Svelte code: Unexpected token `}`. Did you mean `&rbrace;` or `{"}"}`?
https://svelte.dev/e/js_parse_error -->
<script lang="ts">
  import { onMount } from "svelte";
  import { CLEAR_HISTORY_COMMAND } from "lexical";
  import { type LexicalEditor } from "lexical";
  import { $toggleLink as _toggleLink } from "@lexical/link";
  import { $convertFromMarkdownString as _convertFromMarkdownString } from "@lexical/markdown";
  import { TRANSFORMERS } from "@lexical/markdown";

  import ToolbarPlugin from "./ToolbarPlugin/ToolbarPlugin.svelte";
  import { updateToolbar } from "./ToolbarPlugin/ToolbarEvent.js";
  import { InitialEditor } from "./EditorEvent.js";
  import { createEventDispatcher } from "svelte";
  import { IMAGE } from "./MarkdownTransformers.js";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faLink, faTrash } from "@fortawesome/free-solid-svg-icons";
  library.add(faLink, faTrash);
  dom.watch();

  /** 入力値 */
  export let value: string;
  const dispatch = createEventDispatcher();

  /** lexical Editor */
  let editor: LexicalEditor;
  /** 選択ノードの段落種別 */
  let para: string = "";
  /** アンドゥボタン表示状態 */
  let canUndo: boolean;
  /** リドゥボタン表示状態 */
  let canRedo: boolean;
  /** リンクポップアップ表示値 */
  let linkValue: string;

  /** マウントイベント */
  onMount(async () => {
    let detailRect = document.querySelector("#detail")?.getBoundingClientRect();
    let footerRect = document.querySelector("footer")?.getBoundingClientRect();
    if (detailRect !== undefined && footerRect !== undefined) {
      let height = footerRect.top - detailRect.top - 100;
      let a = height + "px";
      document
        .querySelector<HTMLDivElement>("#detail")
        ?.style.setProperty("height", a);
    }

    // LexicalEditorの初期化
    editor = await InitialEditor(
      document.getElementById("detail") as HTMLElement,
      document.getElementById("rich") as HTMLDivElement,
      document.getElementById("link-menu") as HTMLDivElement,
      (p: boolean) => (canUndo = p),
      (p: boolean) => (canRedo = p),
      (arg) => {
        arg.editorState.read(() => {
          const { para: ppara, linkValue: plinkValue } = updateToolbar();
          para = ppara;
          linkValue = <string>plinkValue;
        });
      },
      dispatch,
    );
  });

  $: {
    // 入力値が更新されたとき、マークダウン変換と、履歴クリア
    if (value !== undefined && value !== "") {
      const trans = TRANSFORMERS;
      trans.unshift(IMAGE);
      editor.update(() => _convertFromMarkdownString(value, trans));
      editor.dispatchCommand(CLEAR_HISTORY_COMMAND, undefined);
    }
  }

  // リンクURLの更新
  const linkOK = () => {
    editor.update(() => {
      _toggleLink(linkValue);
    });
  };
  const linkDel = () => {
    editor.update(() => {
      _toggleLink(null);
    });
  };
</script>

<div id="rich" class="panel is-dark">
  <ToolbarPlugin {editor} {para} {canUndo} {canRedo} />

  <div class="panel-block p-0 is-fullwidth">
    <div id="detail" class="content editor-input" contenteditable></div>
    <div id="link-menu" class="card has-background-dark is-hidden p-0">
      <div class="card-content p-3">
        <div class="content">
          <div class="field is-grouped">
            <p class="control is-expanded m-0">
              <input
                class="input is-small"
                type="text"
                placeholder="URL"
                bind:value={linkValue}
              />
            </p>
            <p class="buttons m-0">
              <!-- svelte-ignore a11y_consider_explicit_label -->
              <button class="button is-info is-small" on:click={linkOK}>
                <i class="fa-solid fa-link"></i>
              </button>
              <!-- svelte-ignore a11y_consider_explicit_label -->
              <button class="button is-danger is-small" on:click={linkDel}>
                <i class="fa-solid fa-trash"></i>
              </button>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
  <textarea id="lexical-state"> </textarea>
  <!-- デバッグテキストエリア
  -->
</div>

<style>
  #detail {
    overflow: auto;
    width: 99%;
    min-height: 150px;
    margin-left: 2px;
    padding: 5px;
  }

  #link-menu {
    position: absolute;
    width: 400px;
    z-index: 9999;
  }
</style>
