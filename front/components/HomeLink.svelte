<script lang="ts">
  import House from "lucide-svelte/icons/house";

  let {
    canLeave,
    onNavigate,
  }: {
    canLeave?: () => boolean;
    onNavigate?: () => void;
  } = $props();

  const handleClick = (event: MouseEvent) => {
    if (canLeave && !canLeave()) {
      event.preventDefault();
      return;
    }
    onNavigate?.();
  };
</script>

<a
  class="home-link"
  href="/"
  aria-label="トップに戻る"
  title="トップに戻る"
  onclick={handleClick}
>
  <House size={20} aria-hidden="true" />
</a>

<style>
  .home-link {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    width: 2.75rem;
    height: 2.75rem;
    border-radius: 0.4rem;
    color: inherit;
  }

  .home-link:hover,
  .home-link:focus-visible {
    background: color-mix(in srgb, currentColor 12%, transparent);
  }
</style>
