<script lang="ts">
  import AppIcon from "./AppIcon.svelte";

  interface PreviewImage {
    id: string;
    src: string;
    alt: string;
    markdown?: string;
  }

  interface Props {
    images?: PreviewImage[];
    onEmbedImage?: (markdown: string) => void;
    onPreviewImage?: (image: PreviewImage) => void;
    imageUploadPath?: string;
    interactive?: boolean;
    showAddButton?: boolean;
  }

  let {
    images = [
      {
        id: "sample-1",
        src: "https://images.unsplash.com/photo-1520975916090-3105956dac38?auto=format&fit=crop&w=1200&q=80",
        alt: "sample image 1",
        markdown: "![sample image 1](https://images.unsplash.com/photo-1520975916090-3105956dac38?auto=format&fit=crop&w=1200&q=80)",
      },
      {
        id: "sample-2",
        src: "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=80",
        alt: "sample image 2",
        markdown: "![sample image 2](https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=80)",
      },
      {
        id: "sample-3",
        src: "https://images.unsplash.com/photo-1516035069371-29a1b244cc32?auto=format&fit=crop&w=1200&q=80",
        alt: "sample image 3",
        markdown: "![sample image 3](https://images.unsplash.com/photo-1516035069371-29a1b244cc32?auto=format&fit=crop&w=1200&q=80)",
      },
    ],
    onEmbedImage,
    onPreviewImage,
    imageUploadPath = "",
    interactive = true,
    showAddButton = true,
  }: Props = $props();

  let imageList = $state<PreviewImage[]>([]);
  let fileInput = $state<HTMLInputElement | null>(null);
  let isUploading = $state(false);
  let uploadError = $state("");

  $effect(() => {
    imageList = [...images];
  });

  const embedImage = (image: PreviewImage) => {
    if (onPreviewImage != null) {
      onPreviewImage(image);
      return;
    }
    if (!interactive) {
      return;
    }
    const markdown = image.markdown ?? `![${image.alt}](${image.src})`;
    onEmbedImage?.(markdown);
  };

  const removeImage = (id: string) => {
    if (!interactive) {
      return;
    }
    const image = imageList.find((item) => item.id === id);
    if (image == null) {
      return;
    }
    (async () => {
      try {
        const response = await fetch(image.src, { method: "DELETE" });
        if (!response.ok) {
          uploadError = await response.text();
          return;
        }
        imageList = imageList.filter((item) => item.id !== id);
      } catch {
        uploadError = "画像を削除できませんでした";
      }
    })();
  };

  const toPreviewImage = (src: string, fileName: string, id: string): PreviewImage => {
    const alt = fileName.replace(/\.[^.]+$/, "").trim() || "image";
    return {
      id,
      src,
      alt,
      markdown: `![${alt}](${src})`,
    };
  };

  const openImagePicker = () => {
    if (!interactive || isUploading) {
      return;
    }
    fileInput?.click();
  };

  const onImageSelected = async (event: Event) => {
    const target = event.currentTarget as HTMLInputElement;
    const file = target.files?.item(0);
    target.value = "";

    if (file == null || !interactive) {
      return;
    }
    if (imageUploadPath.trim() === "") {
      uploadError = "画像アップロード先が設定されていません";
      return;
    }

    isUploading = true;
    uploadError = "";

    try {
      const data = new FormData();
      data.append("file", file);
      const response = await fetch(imageUploadPath, {
        method: "post",
        body: data,
      });
      if (response.status !== 200) {
        uploadError = await response.text();
        return;
      }

      const imageUrl = await response.text();
      const nextImage = toPreviewImage(imageUrl, file.name, crypto.randomUUID());
      imageList = [nextImage, ...imageList];
    } catch {
      uploadError = "画像をアップロードできませんでした";
    } finally {
      isUploading = false;
    }
  };
</script>

<section class="image-card">
  <div class="image-card-body">
    {#if uploadError !== ""}
      <div class="notification is-danger mb-3">{uploadError}</div>
    {/if}
    <div class="image-card-strip" aria-label="selected images">
      {#if showAddButton}
        <button
          class="image-card-add"
          type="button"
          aria-label="add images"
          disabled={!interactive || isUploading}
          onclick={openImagePicker}
        >
          <span class="image-card-add-icon">
            <AppIcon name="plus" />
          </span>
          <span class="image-card-add-text">画像を追加</span>
        </button>
      {/if}

      {#each imageList as image}
        <div
          class="image-card-item"
          role="button"
          tabindex={0}
          aria-disabled={false}
          onclick={() => embedImage(image)}
          onkeydown={(event) => {
            if (event.key === "Enter" || event.key === " ") {
              event.preventDefault();
              embedImage(image);
            }
          }}
        >
          <div class="image-card-preview">
            <img src={image.src} alt={image.alt} loading="lazy" />
          </div>
          <div class="image-card-item-footer">
            <span class="image-card-item-name">{image.alt}</span>
            {#if interactive}
              <button
                class="button is-ghost image-card-delete-button"
                type="button"
                aria-label={`delete ${image.alt}`}
                onclick={(event) => {
                  event.stopPropagation();
                  removeImage(image.id);
                }}
              >
                <AppIcon name="trash" />
              </button>
            {/if}
          </div>
        </div>
      {/each}
    </div>
    <input
      class="is-hidden"
      type="file"
      accept="image/*"
      bind:this={fileInput}
      onchange={onImageSelected}
    />
  </div>
</section>

<style>
  .image-card {
    margin-top: 0.9rem;
    padding: 1rem;
    border: 1px solid color-mix(in srgb, var(--bulma-border) 78%, white 22%);
    border-radius: 1rem;
    background:
      linear-gradient(180deg, color-mix(in srgb, var(--bulma-scheme-main) 96%, white 4%), var(--bulma-scheme-main));
    box-shadow: 0 0.75rem 2rem rgba(12, 18, 28, 0.08);
  }

  .image-card-strip {
    display: grid;
    grid-auto-flow: column;
    grid-auto-columns: minmax(16rem, 18rem);
    gap: 0.85rem;
    overflow-x: auto;
    padding-bottom: 0.2rem;
    scroll-snap-type: x proximity;
    scrollbar-gutter: stable;
    scrollbar-width: thin;
    scrollbar-color: color-mix(in srgb, #2d8f86 58%, white 42%) color-mix(in srgb, var(--bulma-scheme-main) 82%, black 18%);
  }

  .image-card-strip::-webkit-scrollbar {
    height: 0.7rem;
  }

  .image-card-strip::-webkit-scrollbar-track {
    background: color-mix(in srgb, var(--bulma-scheme-main) 82%, black 18%);
    border-radius: 999px;
  }

  .image-card-strip::-webkit-scrollbar-thumb {
    border: 0.18rem solid color-mix(in srgb, var(--bulma-scheme-main) 82%, black 18%);
    border-radius: 999px;
    background: linear-gradient(
      90deg,
      color-mix(in srgb, #2d8f86 78%, white 22%),
      color-mix(in srgb, #2d8f86 52%, black 48%)
    );
  }

  .image-card-strip::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(
      90deg,
      color-mix(in srgb, #2d8f86 88%, white 12%),
      color-mix(in srgb, #2d8f86 62%, black 38%)
    );
  }

  .image-card-item,
  .image-card-add {
    scroll-snap-align: start;
  }

  .image-card-item {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    padding: 0.75rem;
    border: 1px solid color-mix(in srgb, var(--bulma-border) 82%, white 18%);
    border-radius: 0.95rem;
    background: var(--bulma-scheme-main);
    cursor: pointer;
    text-align: left;
  }

  .image-card-item:disabled,
  .image-card-add:disabled {
    cursor: default;
    opacity: 0.72;
  }

  .image-card-preview {
    overflow: hidden;
    aspect-ratio: 4 / 3;
    border-radius: 0.8rem;
    background: linear-gradient(135deg, #dfe7ef, #f4f7fa);
  }

  .image-card-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }

  .image-card-item-footer {
    display: flex;
    align-items: center;
    gap: 0.55rem;
  }

  .image-card-item-name {
    flex: 1 1 auto;
    color: var(--bulma-text);
    font-size: 0.84rem;
    font-weight: 600;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .image-card-delete-button {
    flex: 0 0 auto;
    min-width: 2rem;
    min-height: 2rem;
    padding: 0;
    color: #d64545;
    background: color-mix(in srgb, #d64545 10%, transparent);
    box-shadow: inset 0 0 0 1px color-mix(in srgb, #d64545 18%, transparent);
  }

  .image-card-delete-button:hover {
    color: #bf2f2f;
    background: color-mix(in srgb, #d64545 16%, transparent);
  }

  .image-card-add {
    display: inline-flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.65rem;
    min-height: 100%;
    padding: 1rem;
    border: 1px dashed color-mix(in srgb, var(--bulma-border) 82%, white 18%);
    border-radius: 0.95rem;
    background: color-mix(in srgb, var(--bulma-scheme-main) 88%, #edf4f2);
    color: var(--bulma-text-weak);
  }

  .image-card-add-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 999px;
    background: color-mix(in srgb, #2d8f86 12%, var(--bulma-scheme-main));
    color: color-mix(in srgb, #2d8f86 84%, black 16%);
  }

  .image-card-add-text {
    font-size: 0.9rem;
    font-weight: 600;
  }

  @media screen and (max-width: 768px) {
    .image-card {
      padding: 0.85rem;
    }

    .image-card-strip {
      grid-auto-columns: minmax(13.5rem, 16rem);
      gap: 0.7rem;
    }
  }
</style>
