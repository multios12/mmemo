<script lang="ts">
  import type { LexicalEditor } from "lexical";
  import { INSERT_IMAGE_COMMAND } from "../ImagesPlugin/index.js";
  import type { ImagePayload } from "../ImagesPlugin/ImageNode.js";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faImage, faUpload } from "@fortawesome/free-solid-svg-icons";
  library.add(faImage, faUpload);
  dom.watch();

  interface Props {
    editor: LexicalEditor;
  }

  let { editor }: Props = $props();
  /** メッセージ */
  let message = $state("");

  /** ファイル選択イベント */
  const fileChange = async () => {
    const e = document.querySelector("#fileInput") as HTMLInputElement;
    const file = e.files?.item(0);
    document.querySelector("#progress")?.classList.remove("is-hidden");
    if (!file) {
      message = "select file is null";
      return;
    }
    const data = new FormData();
    data.append("file", file);
    try {
      const init = { method: "post", body: data };
      const r = await fetch(`./api/images`, init);
      if (r.status === 200) {
        document.querySelector("#dialog")?.classList.remove("is-active");
        const f = await r.text();
        const payload: ImagePayload = {
          src: f,
          altText: "イメージ",
        };
        editor.dispatchCommand(INSERT_IMAGE_COMMAND, payload);
        //pop();
      } else {
        message = await r.text();
      }
    } catch (e) {
      message = "ファイルを保存できませんでした";
    } finally {
      document.querySelector("#progress")?.classList.add("is-hidden");
    }
  };

  /** モーダルトグルイベント */
  const toggleClick = () => {
    document.querySelector("#dialog")?.classList.toggle("is-active");
  };
</script>

<!-- svelte-ignore a11y_consider_explicit_label -->
<button
  id="linkButton"
  class="button is-ghost p-0"
  tabindex="-1"
  onclick={toggleClick}
>
  <i class="fa-solid fa-image"></i>
</button>

<div class="modal" id="dialog">
  <div class="modal-background"></div>
  <div class="modal-card">
    <header class="modal-card-head">
      <p class="modal-card-title">Image Upload</p>
      <button class="delete" aria-label="close" onclick={toggleClick}></button>
    </header>
    <section class="modal-card-body">
      {#if message != ""}
        <div class="notification is-danger">{message}</div>
      {/if}
      <div class="file has-name is-boxed is-fullwidth">
        <label class="file-label">
          <input
            id="fileInput"
            class="file-input"
            type="file"
            name="resume"
            onchange={fileChange}
          />
          <span class="file-cta">
            <span class="file-icon">
              <i class="fa-solid fa-upload"> file_upload </i>
            </span>
            <span class="file-label"> Choose a file… </span>
            <progress
              id="progress"
              class="progress is-small is-primary is-hidden"
              max="100"
            ></progress>
          </span>
        </label>
      </div>
    </section>
  </div>
</div>
