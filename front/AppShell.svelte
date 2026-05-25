<script lang="ts">
  import { onMount } from "svelte";
  import Settings from "lucide-svelte/icons/settings";
  import { useNavigate, useRoute } from "@dvcol/svelte-simple-router/router";
  import type { settingType } from "./models/settingType.js";
  import { settingsStore } from "./store.js";
  import { apiSettingsPath, appPath } from "./basePath.js";

  const { push } = useNavigate();
  const { location } = $derived(useRoute());

  const currentPath = $derived(location?.path ?? "/");
  const isDetailRoute = $derived.by(() => {
    const parts = currentPath.split("/").filter(Boolean);
    return parts.length === 2 || (parts.length === 3 && parts[1] === "edit");
  });
  const isSettingsDetailRoute = $derived.by(() =>
    /^\/settings\/templates\/[^/]+$/.test(currentPath),
  );
  const mobileNavItems = $derived.by(() => {
    return (
      $settingsStore?.Categories?.map((category) => ({
        key: category.Key,
        value: category.Name,
      })) ?? []
    );
  });
  const mobileNavKey = $derived.by(() => {
    const parts = currentPath.split("/").filter(Boolean);
    return parts[0] ?? "";
  });

  const closeNavbarMenu = () => {
    document.querySelector(".navbar-burger")?.classList.remove("is-active");
    document.querySelector(".navbar-menu")?.classList.remove("is-active");
  };

  const changeMobileCategory = async (value: string) => {
    if (value) {
      await push({ path: `/${value}/` });
    }
  };

  const openSettings = async () => {
    await push({ path: "/settings" });
  };

  onMount(() => {
    const navbarBurger = document.querySelector<HTMLElement>(".navbar-burger");
    navbarBurger?.addEventListener("click", () => {
      const target = navbarBurger.dataset.target;
      const navbarMenu = target
        ? document.getElementById(target)
        : document.querySelector<HTMLElement>(".navbar-menu");

      navbarBurger.classList.toggle("is-active");
      navbarMenu?.classList.toggle("is-active");
    });

    (async () => {
      const r = await fetch(apiSettingsPath());
      if (!r.ok) {
        return;
      }
      const settings = (await r.json()) as settingType;
      settingsStore.set(settings);
    })();
  });
</script>

<nav
  class="navbar is-transparent is-dark"
  class:is-hidden={isDetailRoute || isSettingsDetailRoute}
>
  <div class="navbar-brand">
    <div
      class="navbar-item navbar-brand-title is-unselectable is-hidden-mobile"
    >
      memo
    </div>
    <div class="navbar-item navbar-mobile-context is-hidden-tablet">
      <div
        class="mobile-nav-segments"
        role="tablist"
        aria-label="category navigation"
      >
        {#each mobileNavItems as item}
          <button
            class="mobile-nav-segment"
            class:is-active={mobileNavKey === item.key}
            type="button"
            role="tab"
            aria-selected={mobileNavKey === item.key}
            onclick={() => changeMobileCategory(item.key)}
          >
            {item.value}
          </button>
        {/each}
      </div>
    </div>
    <div
      class="navbar-burger js-burger is-hidden-mobile"
      data-target="navbarMMemo"
    >
      <span></span>
      <span></span>
      <span></span>
      <span></span>
    </div>
  </div>

  <div id="navbarMMemo" class="navbar-menu">
    <div class="navbar-start">
      {#if $settingsStore !== undefined}
        {#each $settingsStore.Categories as category}
          <a
            class="navbar-item is-unselectable is-tab"
            href={appPath(`/${category.Key}/`)}
            onclick={closeNavbarMenu}
          >
            {category.Name}
          </a>
        {/each}
      {/if}
    </div>
    <div class="navbar-end">
      <div class="navbar-item">
        <button
          class="button is-ghost navbar-settings-button"
          class:is-active={currentPath === "/settings"}
          type="button"
          aria-label="settings"
          aria-pressed={currentPath === "/settings"}
          onclick={openSettings}
        >
          <Settings size={18} strokeWidth={2.25} aria-hidden="true" />
        </button>
      </div>
    </div>
  </div>
</nav>

<style>
  .navbar-settings-button {
    color: var(--app-text-weak-invert);
    border-color: transparent;
    min-width: 2.5rem;
    height: 2.5rem;
    padding: 0;
  }

  .navbar-brand-title {
    font-weight: 700;
  }

  .navbar-settings-button:hover,
  .navbar-settings-button:focus-visible,
  .navbar-settings-button.is-active {
    color: var(--app-text-invert);
    background-color: color-mix(in srgb, var(--app-text) 16%, transparent);
  }

  .navbar-mobile-context {
    margin-left: 0;
    margin-right: 0;
    padding-left: 0.75rem;
    padding-right: 0.75rem;
    width: 100%;
  }

  .mobile-nav-segments {
    display: flex;
    align-items: center;
    gap: 0;
    width: 100%;
    overflow-x: auto;
    scrollbar-width: none;
    justify-content: center;
  }

  .mobile-nav-segments::-webkit-scrollbar {
    display: none;
  }

  .mobile-nav-segment {
    flex: 0 0 auto;
    min-height: 2.5rem;
    padding: 0.55rem 0.9rem;
    border: none;
    border-radius: 0;
    background: transparent;
    color: var(--app-text-weak-invert);
    font-size: 0.95rem;
    font-weight: 500;
    white-space: nowrap;
  }

  .mobile-nav-segment.is-active {
    color: var(--app-text-invert);
    box-shadow: inset 0 -2px 0 var(--app-link);
  }
</style>
