<script lang="ts">
  import { onMount } from "svelte";
  import { CLEAR_HISTORY_COMMAND } from "lexical";
  import { type LexicalEditor, type UpdateListener } from "lexical";
  import { $toggleLink as _toggleLink } from "@lexical/link";
  import { $convertFromMarkdownString as _convertFromMarkdownString } from "@lexical/markdown";
  import { TRANSFORMERS } from "@lexical/markdown";

  import ToolbarPlugin from "./ToolbarPlugin/ToolbarPlugin.svelte";
  import { updateToolbar } from "./ToolbarPlugin/ToolbarEvent.js";
  import { InitialEditor } from "./EditorEvent.js";
  import { IMAGE } from "./MarkdownTransformers.js";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faLink, faTrash } from "@fortawesome/free-solid-svg-icons";
  library.add(faLink, faTrash);
  dom.watch();

  interface Props {
    /** 入力値 */
    value?: string;
    imageUploadPath?: string;
    onTextChange?: (value: string) => void;
  }

  let { value = $bindable(""), imageUploadPath = "", onTextChange }: Props = $props();

  /** lexical Editor */
  let editor = $state<LexicalEditor | undefined>(undefined);
  /** 選択ノードの段落種別 */
  let para = $state("");
  /** アンドゥボタン表示状態 */
  let canUndo = $state(false);
  /** リドゥボタン表示状態 */
  let canRedo = $state(false);
  /** リンクポップアップ表示値 */
  let linkValue = $state("");
  /** 直近でエディタから通知された値 */
  let emittedValue = $state("");

  const handleTextChange = (nextValue: string) => {
    emittedValue = nextValue;
    onTextChange?.(nextValue);
  };

  /** マウントイベント */
  onMount(async () => {
    let detailRect = document.querySelector("#detail")?.getBoundingClientRect();
    let footerRect = document.querySelector("footer")?.getBoundingClientRect();
    if (detailRect !== undefined && footerRect !== undefined) {
      let height = footerRect.top - detailRect.top - 36;
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
      ((arg) => {
        arg.editorState.read(() => {
          const { para: ppara, linkValue: plinkValue } = updateToolbar();
          para = ppara;
          linkValue = <string>plinkValue;
        });
      }) satisfies UpdateListener,
      handleTextChange,
    );
  });

  $effect(() => {
    // 入力値が更新されたとき、マークダウン変換と、履歴クリア
    if (
      editor !== undefined &&
      value !== undefined &&
      value !== emittedValue
    ) {
      const trans = [IMAGE, ...TRANSFORMERS];
      editor.update(() => _convertFromMarkdownString(value ?? "", trans));
      editor.dispatchCommand(CLEAR_HISTORY_COMMAND, undefined);
    }
  });

  // リンクURLの更新
  const linkOK = () => {
    if (!editor) return;
    editor.update(() => {
      _toggleLink(linkValue);
    });
  };
  const linkDel = () => {
    if (!editor) return;
    editor.update(() => {
      _toggleLink(null);
    });
  };
</script>

<div id="rich" class="panel is-dark">
  {#if editor}
    <ToolbarPlugin {editor} {imageUploadPath} {para} {canUndo} {canRedo} />
  {/if}

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
              <button class="button is-info is-small" onclick={linkOK}>
                <i class="fa-solid fa-link"></i>
              </button>
              <!-- svelte-ignore a11y_consider_explicit_label -->
              <button class="button is-danger is-small" onclick={linkDel}>
                <i class="fa-solid fa-trash"></i>
              </button>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
  <textarea id="lexical-state" style="display: none;"> </textarea>
  <!-- デバッグテキストエリア
  -->
</div>

<style>
  #rich {
    overflow: visible;
  }

  #rich :global(.panel-block) {
    overflow: visible;
    align-items: stretch;
  }

  #detail {
    overflow: auto;
    width: 100%;
    min-height: 150px;
    margin-left: 0;
    padding: 8px 10px;
    box-sizing: border-box;
  }

  #link-menu {
    position: absolute;
    width: 400px;
    z-index: 9999;
  }
</style>
