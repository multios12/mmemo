<script lang="ts" module>
  export type { ToolbarButton, ToolbarButtonContext } from "./types.js";
</script>

<script lang="ts">
  import { tick } from "svelte";
  import ArrowLeft from "lucide-svelte/icons/arrow-left";
  import Image from "lucide-svelte/icons/image";
  import Link from "lucide-svelte/icons/link";
  import Trash from "lucide-svelte/icons/trash";
  import Toolbar from "./Toolbar.svelte";
  import type { ToolbarButton, ToolbarButtonContext } from "./types.js";

  interface Props {
    value?: string;
    imageUploadPath?: string;
    onTextChange?: (value: string) => void;
    toolbarButtons?: ToolbarButton[];
  }

  let {
    value = $bindable(""),
    imageUploadPath = "",
    onTextChange,
    toolbarButtons = [],
  }: Props = $props();
  let textarea = $state<HTMLTextAreaElement | null>(null);
  let paragraph = $state("normal");
  let isBold = $state(false);
  let isItalic = $state(false);
  let isLink = $state(false);
  let isStrike = $state(false);
  let isLinkModalOpen = $state(false);
  let linkValue = $state("");
  let linkInput = $state<HTMLInputElement | null>(null);
  let pendingLinkSelection = $state({ start: 0, end: 0 });
  let isImageModalOpen = $state(false);
  let imageInput = $state<HTMLInputElement | null>(null);
  let imageError = $state("");
  let isImageUploading = $state(false);
  let pendingImageCursor = $state(0);

  const resizeTextarea = () => {
    if (textarea === null) {
      return;
    }

    textarea.style.height = "0px";
    textarea.style.height = `${textarea.scrollHeight}px`;
  };

  const handleInput = (event: Event) => {
    const target = event.currentTarget as HTMLTextAreaElement;
    onTextChange?.(target.value);
    syncToolbarState(target);
    resizeTextarea();
  };

  const handleKeydown = (event: KeyboardEvent) => {
    if (textarea === null) {
      return;
    }
    if (!event.metaKey && !event.ctrlKey) {
      return;
    }

    const key = event.key.toLowerCase();
    switch (key) {
      case "b":
        event.preventDefault();
        applyInline("bold");
        return;
      case "i":
        event.preventDefault();
        applyInline("italic");
        return;
      case "k":
        event.preventDefault();
        openLinkModal();
        return;
      case "x":
        if (event.shiftKey) {
          event.preventDefault();
          applyInline("strike");
        }
        return;
      case "1":
        if (event.altKey) {
          event.preventDefault();
          applyParagraph("h1");
        }
        return;
      case "2":
        if (event.altKey) {
          event.preventDefault();
          applyParagraph("h2");
        }
        return;
      case "3":
        if (event.altKey) {
          event.preventDefault();
          applyParagraph("h3");
        }
        return;
      case "7":
        if (event.shiftKey) {
          event.preventDefault();
          applyParagraph("ol");
        }
        return;
      case "8":
        if (event.shiftKey) {
          event.preventDefault();
          applyParagraph("ul");
        }
        return;
      case "9":
        if (event.shiftKey) {
          event.preventDefault();
          applyParagraph("quote");
        }
        return;
      case "\\":
        if (event.shiftKey) {
          event.preventDefault();
          applyParagraph("code");
        }
        return;
      case "0":
        if (event.altKey) {
          event.preventDefault();
          applyParagraph("normal");
        }
        return;
      default:
        return;
    }
  };

  const listMarkerPattern = /^(\s*)(?:\d+\.\s+|-\s+|>\s+|#{1,3}\s+)/;
  const headingPattern = /^\s{0,3}(#{1,3})\s+/;
  const orderedListPattern = /^\s*\d+\.\s+/;
  const unorderedListPattern = /^\s*-\s+/;
  const quotePattern = /^\s*>\s+/;
  const codeFencePattern = /^\s*```/;

  const normalizeLine = (line: string) =>
    line.replace(listMarkerPattern, "");

  const updateValue = (nextValue: string) => {
    value = nextValue;
    onTextChange?.(nextValue);
  };

  const applyParagraph = (kind: string) => {
    if (textarea === null) {
      return;
    }

    const rawValue = textarea.value;
    const selectionStart = textarea.selectionStart;
    const lineStart = rawValue.lastIndexOf("\n", Math.max(0, selectionStart - 1)) + 1;
    const lineEndIndex = rawValue.indexOf("\n", selectionStart);
    const lineEnd = lineEndIndex === -1 ? rawValue.length : lineEndIndex;
    const currentLine = rawValue.slice(lineStart, lineEnd);
    const nextBlock =
      kind === "code"
        ? applyCodeBlock(rawValue, currentLine, lineStart, lineEnd)
        : applyLineFormat(currentLine, kind);
    const nextValue =
      rawValue.slice(0, lineStart) + nextBlock + rawValue.slice(lineEnd);
    const cursorOffset = getNextCursorOffset(currentLine, nextBlock, selectionStart - lineStart);

    updateValue(nextValue);
    paragraph = kind;

    requestAnimationFrame(() => {
      if (textarea === null) {
        return;
      }

      textarea.focus();
      const nextCursor = lineStart + cursorOffset;
      textarea.setSelectionRange(nextCursor, nextCursor);
      syncToolbarState(textarea);
    });
  };

  const inlineMarkers = {
    bold: "**",
    italic: "*",
    strike: "~~",
  } as const;

  type InlineRange = {
    openStart: number;
    contentStart: number;
    contentEnd: number;
    closeEnd: number;
  };

  type LinkRange = InlineRange & {
    urlStart: number;
    urlEnd: number;
  };

  const applyInline = (kind: keyof typeof inlineMarkers) => {
    if (textarea === null) {
      return;
    }

    const marker = inlineMarkers[kind];
    const rawValue = textarea.value;
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const range = getInlineRange(rawValue, start, end, marker);
    let nextValue = rawValue;
    let nextStart = start;
    let nextEnd = end;

    if (range !== null) {
      nextValue =
        rawValue.slice(0, range.openStart) +
        rawValue.slice(range.contentStart, range.contentEnd) +
        rawValue.slice(range.closeEnd);
      nextStart = range.contentStart - marker.length;
      nextEnd = range.contentEnd - marker.length;
    } else if (start !== end) {
      nextValue =
        rawValue.slice(0, start) +
        marker +
        rawValue.slice(start, end) +
        marker +
        rawValue.slice(end);
      nextStart = start + marker.length;
      nextEnd = end + marker.length;
    } else {
      nextValue =
        rawValue.slice(0, start) +
        marker +
        marker +
        rawValue.slice(end);
      nextStart = start + marker.length;
      nextEnd = nextStart;
    }

    updateValue(nextValue);

    requestAnimationFrame(() => {
      if (textarea === null) {
        return;
      }

      textarea.focus();
      textarea.setSelectionRange(nextStart, nextEnd);
      syncToolbarState(textarea);
    });
  };

  const openLinkModal = () => {
    if (textarea === null) {
      return;
    }

    pendingLinkSelection = {
      start: textarea.selectionStart,
      end: textarea.selectionEnd,
    };

    const currentLink = getLinkRange(
      textarea.value,
      textarea.selectionStart,
      textarea.selectionEnd,
    );
    linkValue =
      currentLink === null
        ? ""
        : textarea.value.slice(currentLink.urlStart, currentLink.urlEnd);
    isLinkModalOpen = true;

    requestAnimationFrame(() => {
      linkInput?.focus();
      linkInput?.select();
    });
  };

  const openImageModal = () => {
    if (textarea === null) {
      return;
    }

    pendingImageCursor = textarea.selectionStart;
    imageError = "";
    isImageModalOpen = true;
  };

  const closeImageModal = () => {
    if (isImageUploading) {
      return;
    }

    isImageModalOpen = false;
    imageError = "";
    requestAnimationFrame(() => {
      if (textarea === null) {
        return;
      }

      textarea.focus();
      textarea.setSelectionRange(pendingImageCursor, pendingImageCursor);
      syncToolbarState(textarea);
    });
  };

  const openImagePicker = () => {
    imageInput?.click();
  };

  const onImageSelected = async (event: Event) => {
    const target = event.currentTarget as HTMLInputElement;
    const file = target.files?.item(0);
    target.value = "";

    if (file == null) {
      return;
    }
    if (imageUploadPath.trim() === "") {
      imageError = "画像アップロード先が設定されていません";
      return;
    }

    isImageUploading = true;
    imageError = "";

    try {
      const data = new FormData();
      data.append("file", file);
      const response = await fetch(imageUploadPath, {
        method: "post",
        body: data,
      });
      if (response.status !== 200) {
        imageError = await response.text();
        return;
      }

      const imageUrl = await response.text();
      insertImageMarkdown(imageUrl, file.name);
      isImageModalOpen = false;
    } catch {
      imageError = "画像をアップロードできませんでした";
    } finally {
      isImageUploading = false;
    }
  };

  const closeLinkModal = () => {
    isLinkModalOpen = false;
    linkValue = "";
    requestAnimationFrame(() => {
      textarea?.focus();
      if (textarea !== null) {
        textarea.setSelectionRange(
          pendingLinkSelection.start,
          pendingLinkSelection.end,
        );
        syncToolbarState(textarea);
      }
    });
  };

  const applyLink = () => {
    if (textarea === null) {
      return;
    }

    const url = linkValue.trim();
    if (url === "") {
      closeLinkModal();
      return;
    }

    const rawValue = textarea.value;
    const start = pendingLinkSelection.start;
    const end = pendingLinkSelection.end;
    const currentLink = getLinkRange(rawValue, start, end);
    let nextValue = rawValue;
    let nextStart = start;
    let nextEnd = end;

    if (currentLink !== null) {
      const text = rawValue.slice(currentLink.contentStart, currentLink.contentEnd);
      const replacement = `[${text}](${url})`;
      nextValue =
        rawValue.slice(0, currentLink.openStart) +
        replacement +
        rawValue.slice(currentLink.closeEnd);
      nextStart = currentLink.openStart;
      nextEnd = currentLink.openStart + replacement.length;
    } else {
      const text = rawValue.slice(start, end).trim() || url;
      const replacement = `[${text}](${url})`;
      nextValue = rawValue.slice(0, start) + replacement + rawValue.slice(end);
      nextStart = start;
      nextEnd = start + replacement.length;
    }

    updateValue(nextValue);
    isLinkModalOpen = false;
    linkValue = "";

    requestAnimationFrame(() => {
      if (textarea === null) {
        return;
      }

      textarea.focus();
      textarea.setSelectionRange(nextStart, nextEnd);
      syncToolbarState(textarea);
    });
  };

  const insertImageMarkdown = (imageUrl: string, fileName: string) => {
    if (textarea === null) {
      return;
    }

    const rawValue = textarea.value;
    const cursor = pendingImageCursor;
    const lineStart = rawValue.lastIndexOf("\n", Math.max(0, cursor - 1)) + 1;
    const lineEndIndex = rawValue.indexOf("\n", cursor);
    const lineEnd = lineEndIndex === -1 ? rawValue.length : lineEndIndex;
    const currentLine = rawValue.slice(lineStart, lineEnd);
    const altText = toImageAltText(fileName);
    const imageMarkdown = `![${altText}](${imageUrl})`;

    let nextValue = rawValue;
    let nextCursor = cursor;
    if (currentLine.trim() === "") {
      nextValue =
        rawValue.slice(0, lineStart) + imageMarkdown + rawValue.slice(lineEnd);
      nextCursor = lineStart + imageMarkdown.length;
    } else {
      const insertion = `\n${imageMarkdown}`;
      nextValue =
        rawValue.slice(0, lineEnd) + insertion + rawValue.slice(lineEnd);
      nextCursor = lineEnd + insertion.length;
    }

    updateValue(nextValue);
    imageError = "";

    requestAnimationFrame(() => {
      if (textarea === null) {
        return;
      }

      textarea.focus();
      textarea.setSelectionRange(nextCursor, nextCursor);
      syncToolbarState(textarea);
    });
  };

  const focusTextarea = () => {
    requestAnimationFrame(() => {
      textarea?.focus();
      if (textarea !== null) {
        syncToolbarState(textarea);
      }
    });
  };

  const getToolbarButtonContext = (): ToolbarButtonContext => ({
    value,
    textarea,
    selectionStart: textarea?.selectionStart ?? 0,
    selectionEnd: textarea?.selectionEnd ?? 0,
    updateValue,
    focus: focusTextarea,
    openImageUploadModal: openImageModal,
    syncToolbarState: () => {
      if (textarea !== null) {
        syncToolbarState(textarea);
      }
    },
  });

  const handleExtraToolbarButton = (button: ToolbarButton) => {
    button.onClick?.(getToolbarButtonContext());
  };

  const toImageAltText = (fileName: string) => {
    const baseName = fileName.replace(/\.[^.]+$/, "").trim();
    return baseName === "" ? "image" : baseName;
  };

  const removeLink = () => {
    if (textarea === null) {
      return;
    }

    const rawValue = textarea.value;
    const start = pendingLinkSelection.start;
    const end = pendingLinkSelection.end;
    const currentLink = getLinkRange(rawValue, start, end);
    if (currentLink === null) {
      closeLinkModal();
      return;
    }

    const text = rawValue.slice(currentLink.contentStart, currentLink.contentEnd);
    const nextValue =
      rawValue.slice(0, currentLink.openStart) +
      text +
      rawValue.slice(currentLink.closeEnd);
    const nextStart = currentLink.openStart;
    const nextEnd = currentLink.openStart + text.length;

    updateValue(nextValue);
    isLinkModalOpen = false;
    linkValue = "";

    requestAnimationFrame(() => {
      if (textarea === null) {
        return;
      }

      textarea.focus();
      textarea.setSelectionRange(nextStart, nextEnd);
      syncToolbarState(textarea);
    });
  };

  const applyCodeBlock = (
    rawValue: string,
    line: string,
    lineStart: number,
    lineEnd: number,
  ) => {
    const before = rawValue.slice(0, lineStart);
    const after = rawValue.slice(lineEnd);
    const previousLineStart = before.lastIndexOf("\n", Math.max(0, before.length - 2)) + 1;
    const previousLine = before.slice(previousLineStart).replace(/\n$/, "");
    const nextLineEnd = after.indexOf("\n");
    const nextLine =
      nextLineEnd === -1 ? after : after.slice(0, nextLineEnd);

    const hasWrappedFences =
      codeFencePattern.test(previousLine) && codeFencePattern.test(nextLine);

    if (hasWrappedFences) {
      return line;
    }

    return `\`\`\`\n${line}\n\`\`\``;
  };

  const applyLineFormat = (line: string, kind: string) => {
    if (line.trim() === "") {
      return line;
    }

    const baseLine = normalizeLine(line);
    switch (kind) {
      case "h1":
        return `# ${baseLine}`;
      case "h2":
        return `## ${baseLine}`;
      case "h3":
        return `### ${baseLine}`;
      case "ol":
        return `1. ${baseLine}`;
      case "ul":
        return `- ${baseLine}`;
      case "quote":
        return `> ${baseLine}`;
      case "normal":
      default:
        return baseLine;
    }
  };

  const getLineAtCursor = (rawValue: string, cursor: number) => {
    const lineStart = rawValue.lastIndexOf("\n", Math.max(0, cursor - 1)) + 1;
    const lineEndIndex = rawValue.indexOf("\n", cursor);
    const lineEnd = lineEndIndex === -1 ? rawValue.length : lineEndIndex;
    return rawValue.slice(lineStart, lineEnd);
  };

  const detectParagraph = (rawValue: string, cursor: number) => {
    const line = getLineAtCursor(rawValue, cursor);
    if (isCursorInsideCodeBlock(rawValue, cursor)) {
      return "code";
    }
    if (headingPattern.test(line)) {
      const level = line.match(headingPattern)?.[1].length ?? 0;
      return level === 1 ? "h1" : level === 2 ? "h2" : "h3";
    }
    if (orderedListPattern.test(line)) {
      return "ol";
    }
    if (unorderedListPattern.test(line)) {
      return "ul";
    }
    if (quotePattern.test(line)) {
      return "quote";
    }
    return "normal";
  };

  const isCursorInsideCodeBlock = (rawValue: string, cursor: number) => {
    const lines = rawValue.slice(0, cursor).split("\n");
    let fenceCount = 0;
    for (const line of lines) {
      if (codeFencePattern.test(line)) {
        fenceCount += 1;
      }
    }
    return fenceCount % 2 === 1;
  };

  const syncParagraphWithCursor = (target: HTMLTextAreaElement) => {
    paragraph = detectParagraph(target.value, target.selectionStart);
  };

  const getInlineRange = (
    rawValue: string,
    start: number,
    end: number,
    marker: string,
  ): InlineRange | null => {
    if (start !== end) {
      const openStart = start - marker.length;
      const closeEnd = end + marker.length;
      if (
        openStart >= 0 &&
        closeEnd <= rawValue.length &&
        rawValue.slice(openStart, start) === marker &&
        rawValue.slice(end, closeEnd) === marker
      ) {
        return {
          openStart,
          contentStart: start,
          contentEnd: end,
          closeEnd,
        };
      }
      return null;
    }

    const lineStart = rawValue.lastIndexOf("\n", Math.max(0, start - 1)) + 1;
    const lineEndIndex = rawValue.indexOf("\n", start);
    const lineEnd = lineEndIndex === -1 ? rawValue.length : lineEndIndex;
    const lineText = rawValue.slice(lineStart, lineEnd);
    const offset = start - lineStart;

    let searchIndex = 0;
    while (searchIndex <= lineText.length) {
      const openIndex = lineText.indexOf(marker, searchIndex);
      if (openIndex === -1) {
        return null;
      }
      const contentStart = openIndex + marker.length;
      const closeIndex = lineText.indexOf(marker, contentStart);
      if (closeIndex === -1) {
        return null;
      }
      if (offset >= contentStart && offset <= closeIndex) {
        return {
          openStart: lineStart + openIndex,
          contentStart: lineStart + contentStart,
          contentEnd: lineStart + closeIndex,
          closeEnd: lineStart + closeIndex + marker.length,
        };
      }
      searchIndex = closeIndex + marker.length;
    }

    return null;
  };

  const getLinkRange = (
    rawValue: string,
    start: number,
    end: number,
  ): LinkRange | null => {
    const checkStart = start === end ? start : start + 1;
    const lineStart = rawValue.lastIndexOf("\n", Math.max(0, start - 1)) + 1;
    const lineEndIndex = rawValue.indexOf("\n", end);
    const lineEnd = lineEndIndex === -1 ? rawValue.length : lineEndIndex;
    const lineText = rawValue.slice(lineStart, lineEnd);
    const linkPattern = /\[([^\]]+)\]\(([^)\s]+)\)/g;

    for (const match of lineText.matchAll(linkPattern)) {
      const matched = match[0];
      const text = match[1];
      const url = match[2];
      const localStart = match.index ?? 0;
      const openStart = lineStart + localStart;
      const contentStart = openStart + 1;
      const contentEnd = contentStart + text.length;
      const urlStart = contentEnd + 2;
      const urlEnd = urlStart + url.length;
      const closeEnd = openStart + matched.length;

      if (
        (checkStart >= contentStart && checkStart <= contentEnd) ||
        (start === end && start >= openStart && start <= closeEnd) ||
        (start !== end && start >= contentStart && end <= contentEnd)
      ) {
        return {
          openStart,
          contentStart,
          contentEnd,
          urlStart,
          urlEnd,
          closeEnd,
        };
      }
    }

    return null;
  };

  const syncInlineState = (target: HTMLTextAreaElement) => {
    const { value: currentValue, selectionStart, selectionEnd } = target;
    isBold = getInlineRange(currentValue, selectionStart, selectionEnd, inlineMarkers.bold) !== null;
    isItalic =
      getInlineRange(currentValue, selectionStart, selectionEnd, inlineMarkers.italic) !== null;
    isLink = getLinkRange(currentValue, selectionStart, selectionEnd) !== null;
    isStrike =
      getInlineRange(currentValue, selectionStart, selectionEnd, inlineMarkers.strike) !== null;
  };

  const syncToolbarState = (target: HTMLTextAreaElement) => {
    syncParagraphWithCursor(target);
    syncInlineState(target);
  };

  const getNextCursorOffset = (
    previousLine: string,
    nextBlock: string,
    previousOffset: number,
  ) => {
    if (nextBlock.includes("\n")) {
      const lines = nextBlock.split("\n");
      const firstContentLine = lines[1] ?? "";
      return Math.min(nextBlock.length, 4 + Math.min(previousOffset, firstContentLine.length));
    }

    const previousIndentStripped = normalizeLine(previousLine);
    const removedPrefix = previousLine.length - previousIndentStripped.length;
    const nextIndentStripped = normalizeLine(nextBlock);
    const addedPrefix = nextBlock.length - nextIndentStripped.length;
    const baseOffset = Math.max(0, previousOffset - removedPrefix);
    return Math.min(nextBlock.length, addedPrefix + baseOffset);
  };

  $effect(() => {
    value;
    tick().then(resizeTextarea);
  });
</script>

<div class="md-input">
  <Toolbar
    bind:value={paragraph}
    bold={isBold}
    italic={isItalic}
    link={isLink}
    strike={isStrike}
    onBold={() => applyInline("bold")}
    onItalic={() => applyInline("italic")}
    onLink={openLinkModal}
    onStrike={() => applyInline("strike")}
    extraButtons={toolbarButtons}
    onExtraButton={handleExtraToolbarButton}
  />

  <textarea
    class="md-input-area"
    bind:this={textarea}
    bind:value
    oninput={handleInput}
    onkeydown={handleKeydown}
    onclick={(event) => syncToolbarState(event.currentTarget as HTMLTextAreaElement)}
    onkeyup={(event) => syncToolbarState(event.currentTarget as HTMLTextAreaElement)}
    onselect={(event) => syncToolbarState(event.currentTarget as HTMLTextAreaElement)}
    placeholder={`# 見出し1

本文

- リスト
- リスト

> 引用

\`\`\`
コード
\`\`\``}
    spellcheck="false"
  ></textarea>

  {#if isLinkModalOpen}
    <div class="md-link-modal-backdrop">
      <div class="md-link-modal">
        <div class="md-link-controls">
          <div class="md-link-input-wrap">
            <input
              class="md-link-input"
              type="text"
              bind:this={linkInput}
              bind:value={linkValue}
              placeholder="https://example.com"
              onkeydown={(event) => {
                if (event.key === "Enter") {
                  event.preventDefault();
                  applyLink();
                }
                if (event.key === "Escape") {
                  event.preventDefault();
                  closeLinkModal();
                }
              }}
            />
          </div>
          <div class="md-link-action">
            <button
              class="md-button md-button-primary"
              type="button"
              aria-label="apply link"
              onclick={applyLink}
            >
              <span class="md-icon"><Link /></span>
            </button>
          </div>
          <div class="md-link-action">
            <button
              class="md-button md-button-light"
              type="button"
              aria-label="remove link"
              onclick={removeLink}
            >
              <span class="md-icon"><Trash /></span>
            </button>
          </div>
        </div>
      </div>
    </div>
  {/if}

  {#if isImageModalOpen}
    <div
      class="md-link-modal-backdrop"
      role="button"
      tabindex="0"
      onclick={(event) => {
        if (event.target === event.currentTarget) {
          closeImageModal();
        }
      }}
      onkeydown={(event) => {
        if (event.target !== event.currentTarget) {
          return;
        }
        if (event.key === "Enter" || event.key === " ") {
          event.preventDefault();
          closeImageModal();
        }
        if (event.key === "Escape") {
          event.preventDefault();
          closeImageModal();
        }
      }}
    >
      <div class="md-link-modal">
        {#if imageError !== ""}
          <div class="md-error-message">{imageError}</div>
        {/if}
        <div class="md-image-actions">
          <button
            class="md-button md-button-light"
            type="button"
            aria-label="back"
            disabled={isImageUploading}
            onclick={closeImageModal}
          >
            <span class="md-icon">
              <ArrowLeft />
            </span>
          </button>
          <button
            class="md-button md-button-primary md-image-select-button"
            type="button"
            aria-label="add image"
            disabled={isImageUploading}
            onclick={openImagePicker}
          >
            <span class="md-icon">
              <Image />
            </span>
            <span>画像選択</span>
          </button>
          {#if isImageUploading}
            <button
              class="md-button md-button-light md-loading-button"
              aria-label="uploading"
            ></button>
          {/if}
        </div>
        <input
          class="md-hidden"
          type="file"
          accept="image/*"
          bind:this={imageInput}
          onchange={onImageSelected}
        />
      </div>
    </div>
  {/if}
</div>

<style>
  .md-input {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    min-height: 100%;
  }

  .md-input-area {
    width: 100%;
    min-height: 0;
    height: auto;
    max-height: none;
    resize: none;
    overflow-y: hidden;
    box-sizing: border-box;
    align-self: stretch;
    padding: 0.9rem 1rem;
    border: 1px solid #394152;
    border-radius: 0.75rem;
    background-color: #161b28;
    background-image: linear-gradient(180deg, #1b2030, #161b28);
    color: #eef2ff;
    font-family:
      "Iosevka Custom",
      "SFMono-Regular",
      "Consolas",
      "Liberation Mono",
      monospace;
    font-size: 0.96rem;
    line-height: 1.65;
    caret-color: #8ec5ff;
    tab-size: 2;
    white-space: pre-wrap;
    color-scheme: dark;
  }

  .md-input-area:focus {
    border-color: #5ea7f0;
    box-shadow: 0 0 0 0.15rem rgba(94, 167, 240, 0.28);
    outline: none;
  }

  .md-input-area::placeholder {
    color: rgba(238, 242, 255, 0.52);
  }

  .md-input-area::selection {
    background: rgba(47, 125, 244, 0.38);
    color: #ffffff;
  }

  .md-link-modal-backdrop {
    position: fixed;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    background: color-mix(in srgb, black 16%, transparent);
    z-index: 1200;
  }

  .md-link-modal {
    width: min(34rem, 100%);
    padding: 1rem;
    border-radius: 0.75rem;
    background: #ffffff;
    box-shadow: 0 1.2rem 3rem color-mix(in srgb, black 18%, transparent);
  }

  .md-link-controls {
    display: flex;
    align-items: stretch;
    width: 100%;
  }

  .md-link-input-wrap {
    flex: 1 1 auto;
  }

  .md-link-action {
    flex: 0 0 auto;
  }

  .md-link-input {
    width: 100%;
    height: 2.5rem;
    box-sizing: border-box;
    border: 1px solid #d9dee8;
    border-right: 0;
    border-radius: 0.45rem 0 0 0.45rem;
    color: #1f2937;
    font: inherit;
    padding: 0 0.75rem;
  }

  .md-link-input:focus {
    border-color: #3e8ed0;
    outline: none;
  }

  .md-button {
    appearance: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 2.5rem;
    min-height: 2.5rem;
    border: 1px solid transparent;
    border-radius: 0.45rem;
    cursor: pointer;
    font: inherit;
    gap: 0.35rem;
    line-height: 1;
    padding: 0 0.75rem;
    user-select: none;
  }

  .md-button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  .md-button-primary {
    background: #3e8ed0;
    border-color: #3e8ed0;
    color: #ffffff;
  }

  .md-button-primary:not(:disabled):hover {
    background: #3273dc;
    border-color: #3273dc;
  }

  .md-button-light {
    background: #f5f7fb;
    border-color: #d9dee8;
    color: #1f2937;
  }

  .md-button-light:not(:disabled):hover {
    background: #e8edf5;
  }

  .md-link-action .md-button {
    border-radius: 0;
  }

  .md-link-action:last-child .md-button {
    border-radius: 0 0.45rem 0.45rem 0;
  }

  .md-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 1.25rem;
    height: 1.25rem;
  }

  .md-error-message {
    margin-bottom: 0.75rem;
    border-radius: 0.5rem;
    background: #fee2e2;
    color: #991b1b;
    padding: 0.75rem;
  }

  .md-image-actions {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .md-image-select-button {
    flex: 1 1 auto;
    justify-content: center;
  }

  .md-loading-button {
    position: relative;
    color: transparent;
  }

  .md-loading-button::after {
    position: absolute;
    width: 1rem;
    height: 1rem;
    border: 2px solid color-mix(in srgb, #1f2937 18%, transparent);
    border-top-color: #1f2937;
    border-radius: 999px;
    content: "";
    animation: md-spin 0.8s linear infinite;
  }

  .md-hidden {
    display: none;
  }

  @keyframes md-spin {
    to {
      transform: rotate(360deg);
    }
  }

  @media screen and (max-width: 768px) {
    .md-input-area {
      min-height: 0;
      height: auto;
      max-height: none;
      padding: 0.8rem 0.9rem;
      font-size: 0.92rem;
    }

    .md-link-modal {
      padding: 0.85rem;
    }

    .md-image-actions {
      gap: 0.5rem;
    }
  }
</style>
