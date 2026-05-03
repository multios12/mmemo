<script lang="ts">
  import AppIcon from "./AppIcon.svelte";

  type PreviewImage = {
    id: string;
    src: string;
    alt: string;
    markdown?: string;
  };

  interface Props {
    images?: PreviewImage[];
    currentIndex?: number;
    open?: boolean;
    onClose?: () => void;
    onPrev?: () => void;
    onNext?: () => void;
  }

  let {
    images = [],
    currentIndex = 0,
    open = false,
    onClose,
    onPrev,
    onNext,
  }: Props = $props();

  const close = () => {
    onClose?.();
  };

  let touchStartX = 0;
  let touchDeltaX = 0;

  const handleKeydown = (event: KeyboardEvent) => {
    if (event.key === "Escape") {
      event.preventDefault();
      close();
    }
    if (event.key === "ArrowLeft") {
      event.preventDefault();
      onPrev?.();
    }
    if (event.key === "ArrowRight") {
      event.preventDefault();
      onNext?.();
    }
  };

  const handleTouchStart = (event: TouchEvent) => {
    touchStartX = event.touches[0]?.clientX ?? 0;
    touchDeltaX = 0;
  };

  const handleTouchMove = (event: TouchEvent) => {
    touchDeltaX = (event.touches[0]?.clientX ?? 0) - touchStartX;
  };

  const handleTouchEnd = () => {
    if (Math.abs(touchDeltaX) < 50) {
      return;
    }
    if (touchDeltaX > 0) {
      onPrev?.();
    } else {
      onNext?.();
    }
  };

  const currentImage = $derived(images[currentIndex]);
</script>

{#if open}
  <div
    class="image-preview-modal"
    role="dialog"
    aria-modal="true"
    tabindex="0"
    aria-label="image preview"
    onkeydown={handleKeydown}
    ontouchstart={handleTouchStart}
    ontouchmove={handleTouchMove}
    ontouchend={handleTouchEnd}
  >
    <button
      class="image-preview-modal-close-area"
      type="button"
      aria-label="close image preview"
      onclick={close}
    ></button>
    <div class="image-preview-modal-shell">
      <button class="image-preview-close" type="button" aria-label="close" onclick={close}>
        <AppIcon name="xmark" />
      </button>
      {#if images.length > 1}
        <button
          class="image-preview-nav image-preview-nav-prev"
          type="button"
          aria-label="previous image"
          onclick={() => onPrev?.()}
        >
          <AppIcon name="chevron-left" />
        </button>
        <button
          class="image-preview-nav image-preview-nav-next"
          type="button"
          aria-label="next image"
          onclick={() => onNext?.()}
        >
          <AppIcon name="chevron-right" />
        </button>
      {/if}
      {#if currentImage != null}
        <button
          class="image-preview-click-close"
          type="button"
          aria-label="close image preview"
          onclick={close}
        >
          <img class="image-preview-image" src={currentImage.src} alt={currentImage.alt} />
        </button>
      {/if}
    </div>
  </div>
{/if}

<style>
  .image-preview-modal {
    position: fixed;
    inset: 0;
    z-index: 1300;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0.75rem;
    background: rgba(6, 10, 18, 0.9);
    backdrop-filter: blur(10px);
  }

  .image-preview-modal-close-area {
    position: absolute;
    inset: 0;
    border: none;
    background: transparent;
    padding: 0;
  }

  .image-preview-modal-shell {
    position: relative;
    width: 100%;
    height: 100%;
  }

  .image-preview-image {
    display: block;
    width: 100%;
    height: 100%;
    object-fit: contain;
    user-select: none;
  }

  .image-preview-click-close {
    position: absolute;
    inset: 0;
    border: none;
    background: transparent;
    padding: 0;
  }

  .image-preview-close {
    position: absolute;
    top: 0.35rem;
    right: 0.35rem;
    z-index: 1;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 999px;
    color: #fff;
    background: rgba(255, 255, 255, 0.12);
    box-shadow: 0 0.35rem 1rem rgba(0, 0, 0, 0.22);
  }

  .image-preview-close:hover {
    background: rgba(255, 255, 255, 0.18);
  }

  .image-preview-nav {
    position: absolute;
    top: 50%;
    z-index: 1;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 3.5rem;
    height: 3.5rem;
    margin-top: -1.75rem;
    border-radius: 999px;
    color: #fff;
    background: rgba(255, 255, 255, 0.16);
    box-shadow: 0 0.45rem 1.25rem rgba(0, 0, 0, 0.28);
  }

  .image-preview-nav:hover {
    background: rgba(255, 255, 255, 0.22);
  }

  .image-preview-nav-prev {
    left: 0.5rem;
  }

  .image-preview-nav-next {
    right: 0.5rem;
  }
</style>
