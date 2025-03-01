<script lang="ts">
  import type { getParentElement } from "node_modules/lexical/LexicalUtils.js";
  let html: string;
  export let value: string;

  const html2md = (html: HTMLElement): string => {
    let result = "";
    for (const child of html.children) {
      if (result !== "") {
        result += "\n";
      }
      result += child.textContent;
    }
    return result;
  };

  const oninput = (e) => {
    html = e.currentTarget.innerHTML;
    value = html2md(e.currentTarget);
  };

  const isValue = (target: string, value: string): boolean => {
    return (
      target !== undefined &&
      value !== undefined &&
      target.length > value.length &&
      target.substring(0, value.length) === value
    );
  };

  const getSelectElement = () => {
    const s = window.getSelection();
    let offset = s.getRangeAt(0).startOffset;
    let soffset = s.getRangeAt(0).startOffset;
    let eoffset = s.getRangeAt(0).endOffset;
    const element = s.getRangeAt(0).commonAncestorContainer.parentElement;

    if (element.id === "main") {
      return { element: undefined, offset: 0, soffset: 0, eoffset: 0 };
    }
    return { s, element, offset, soffset, eoffset };
  };

  // 選択している行にヘッダを設定
  const changeHeader = () => {
    let s = getSelectElement();
    if (s.element === undefined) {
      return;
    }
    if (isValue(s.element.innerText, "#")) {
      if (!isValue(s.element.innerText, "###")) {
        s.element.innerHTML = "#" + s.element.innerHTML;
        s.offset += 1;
      } else {
        s.element.innerHTML = s.element.innerHTML.substring(4);
        s.offset -= 4;
      }
    } else {
      s.element.innerHTML = "# " + s.element.innerHTML;
      s.offset += 2;
    }

    const range = document.createRange();
    range.setStart(s.element.firstChild, s.offset);
    range.setEnd(s.element.firstChild, s.offset);
    selection.removeAllRanges();
    selection.addRange(range);
  };

  // 選択している行にリストを設定
  const changeList = () => {
    let s = getSelectElement();
    if (s.element === undefined) {
      return;
    }

    if (isValue(s.element.innerText, "*")) {
      s.element.innerHTML = s.element.innerHTML.substring(2);
      s.offset -= 2;
    } else {
      s.element.innerHTML = "* " + s.element.innerHTML;
      s.offset += 2;
    }

    const range = document.createRange();
    range.setStart(s.element.firstChild, s.offset);
    range.setEnd(s.element.firstChild, s.offset);
    s.s.removeAllRanges();
    s.s.addRange(range);
  };

  // 選択している文字を太字に設定
  const changeCharactor = (c) => {
    const s = window.getSelection();
    let soffset = s.getRangeAt(0).startOffset;
    let eoffset = s.getRangeAt(0).endOffset;
    const element = s.getRangeAt(0).commonAncestorContainer.parentElement;

    if (element.id === "main") {
      return;
    }

    let i = element.innerHTML;
    let t = i.substring(soffset, eoffset);
    console.log(t.length);
    // 太字「**」→インライン「`」→取消線「~~」→文字装飾なし
    let neweoffset;
    if (t.length > 4 && t.indexOf("**") > -1) {
      t = t.replaceAll("**", "");
      console.log(t);
      t = `\`${t}\``;
      neweoffset = eoffset - 2;
    } else if (t.length > 2 && t.indexOf("`") > -1) {
      t = t.replaceAll("`", "");
      t = `~~${t}~~`;
      neweoffset = eoffset + 2;
    } else if (t.length > 4 && t.indexOf("~~") > -1) {
      t = t.replaceAll("~~", "");
      c = "";
      neweoffset = eoffset - 4;
    } else if (t.length === 0) {
      return;
    } else {
      t = `**${t}**`;
      neweoffset = eoffset + 4;
    }

    element.innerHTML = `${i.substring(0, soffset)}${t}${i.substring(eoffset)}`;

    const range = document.createRange();
    range.setStart(element.firstChild, soffset);
    range.setEnd(element.firstChild, neweoffset);
    s.removeAllRanges();
    s.addRange(range);
  };

  const onkeydown = (e) => {
    if (e.ctrlKey) {
      if (event.code === "KeyH") {
        changeHeader();
        e.preventDefault();
      } else if (event.code === "KeyL") {
        changeList();
        e.preventDefault();
      } else if (event.code === "KeyB") {
        changeCharactor();
        e.preventDefault();
      } else if (event.code === "Enter") {
      }
      html = e.currentTarget.innerHTML;
      value = html2md(e.currentTarget);
    }
  };
</script>

<div id="main">
  <!-- svelte-ignore a11y_consider_explicit_label -->
  <div id="header">
    <select>
      <option value="no">normal</option>
      <option value="ol">number list</option>
      <option value="ul">bullet list</option>
      <option value="qu">quote</option>
      <option value="co">code</option>
      <option value="h1">heading h1</option>
      <option value="h2">heading h2</option>
      <option value="h3">heading h3</option>
    </select>
    <!-- svelte-ignore a11y_consider_explicit_label -->
    <button class="button-icon">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 -960 960 960"
        fill="#e3e3e3"
      >
        <path
          d="M266-192v-576h227.95q67.05 0 123.55 41.32Q674-685.35 674-612q0 51-22.5 79.5T609-490.96Q635-479 665-448t30 91q0 91-67.03 128t-125.81 37H266Zm127-118h104.68Q546-310 556-334.5t10-35.5q0-11-10.5-35.5T494-430H393v120Zm0-232h93q33 0 48.5-17.5T550-597q0-24-17.11-39t-44.28-15H393v109Z"
        />
      </svg>
    </button>
    <!-- svelte-ignore a11y_consider_explicit_label -->
    <button class="button-icon">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 -960 960 960"
        fill="#e3e3e3"
      >
        <path
          d="M440-280H280q-83 0-141.5-58.5T80-480q0-83 58.5-141.5T280-680h160v80H280q-50 0-85 35t-35 85q0 50 35 85t85 35h160v80ZM320-440v-80h320v80H320Zm200 160v-80h160q50 0 85-35t35-85q0-50-35-85t-85-35H520v-80h160q83 0 141.5 58.5T880-480q0 83-58.5 141.5T680-280H520Z"
        />
      </svg>
    </button>
    <button class="button-icon">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 -960 960 960"
        fill="#e3e3e3"
        ><path
          d="M200-120q-33 0-56.5-23.5T120-200v-560q0-33 23.5-56.5T200-840h560q33 0 56.5 23.5T840-760v560q0 33-23.5 56.5T760-120H200Zm0-80h560v-560H200v560Zm40-80h480L570-480 450-320l-90-120-120 160Zm-40 80v-560 560Z"
        /></svg
      >
    </button>
  </div>

  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div id="content" contenteditable="true" {oninput} {onkeydown}>{value}</div>
</div>

<style>
  :root {
    --color-text: #ddddff;
    --color-bg: #222233;
    --color-bg-content: #334;
    --spacing: 0.25em;
    --height-header: 2em;
  }

  #main {
    color: var(--color-text);
    background: var(--color-bg);

    height: 200px;
    width: calc(100%- var(--spacing));

    margin: var(--spacing);
    padding: var(--spacing);

    border-radius: 10px 10px 0px 0px;

    /* ヘッダ部 */
    #header {
      height: var(--height-header);

      .button-icon {
        height: 24px;
        width: 24px;
        border: none;
        background: transparent;
        padding: 0;
        vertical-align: middle;
      }
      .button-icon:hover {
        background: #88a;
      }

      .button-icon:active {
        background: #444466;
      }

      select {
        color: var(--color-text);
        background: transparent;
        border: solid 1px;
        margin: 0 0.5em;
      }

      option {
        appearance: none;
        background: #222233;
        padding: 0;
        margin: 0;
      }
    }

    /* テキスト入力部 */
    #content {
      background: var(--color-bg-content);

      height: calc(100% - (var(--spacing) * 2 + var(--height-header) * 2));
      width: calc(100% - 0.7em);

      border: solid 0.1em;
      padding: var(--spacing);

      border-radius: 10px 10px 0px 0px;
    }
  }
</style>
