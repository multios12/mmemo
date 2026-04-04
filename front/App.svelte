<script lang="ts">
  import "bulma/css/bulma.css";
  import { Router, goto, route } from "@mateothegreat/svelte5-router";
  import EntryList from "./routes/EntryList.svelte";
  import EntryDetail from "./routes/EntryDetail.svelte";
  import Home from "./routes/Home.svelte";
  import Planner from "./routes/Planner/Planner.svelte";
  import { onMount } from "svelte";
  import type { settingType } from "./models/settingType.js";
  import { settingsStore } from "./store.js";

  const routes = [
    { path: "/", component: Home },
    { path: "/planner/", component: Planner },
    { path: /^\/(?<category>[^/]+)\/$/, component: EntryList },
    { path: /^\/(?<category>[^/]+)\/add$/, component: EntryDetail },
    { path: /^\/(?<category>[^/]+)\/(?<id>[^/]+)$/, component: EntryDetail },
  ];
  let settings = $state<settingType | undefined>(undefined);
  let currentPath = $state("/");
  const currentCategoryName = $derived.by(() => {
    if (settings === undefined) {
      return "";
    }

    const categoryKey = currentPath.split("/").filter(Boolean)[0] ?? "";
    return settings.Categories.find((category) => category.Key === categoryKey)?.Name ?? "";
  });
  const isDetailRoute = $derived.by(() => {
    const parts = currentPath.split("/").filter(Boolean);
    return parts.length === 2;
  });
  const mobileNavItems = $derived.by(() => {
    return (
      settings?.Categories?.map((category) => ({
        key: category.Key,
        value: category.Name,
      })) ?? []
    );
  });
  const closeNavbarMenu = () => {
    document.querySelector(".navbar-burger")?.classList.remove("is-active");
    document.querySelector(".navbar-menu")?.classList.remove("is-active");
  };
  const changeMobileCategory = async (value: string) => {
    if (value) {
      await goto(`/${value}/`);
    }
  };

  onMount(() => {
    const syncCurrentPath = () => {
      currentPath = window.location.pathname;
    };
    const dispatchLocationChange = () => {
      window.dispatchEvent(new Event("locationchange"));
    };

    const originalPushState = history.pushState.bind(history);
    const originalReplaceState = history.replaceState.bind(history);
    history.pushState = (...args) => {
      originalPushState(...args);
      dispatchLocationChange();
    };
    history.replaceState = (...args) => {
      originalReplaceState(...args);
      dispatchLocationChange();
    };

    syncCurrentPath();

    window.addEventListener("popstate", syncCurrentPath);
    window.addEventListener("locationchange", syncCurrentPath);

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
      const r = await fetch("/settings");
      settings = (await r.json()) as settingType;
      settingsStore.set(settings);
    })();

    return () => {
      window.removeEventListener("popstate", syncCurrentPath);
      window.removeEventListener("locationchange", syncCurrentPath);
      history.pushState = originalPushState;
      history.replaceState = originalReplaceState;
    };
  });

  const mobileNavKey = $derived.by(() => {
    const parts = currentPath.split("/").filter(Boolean);
    return parts[0] ?? "";
  });
</script>

<nav class="navbar is-transparent is-dark" class:is-hidden={isDetailRoute}>
  <div class="navbar-brand">
    <div class="navbar-item is-unselectable has-text-weight-bold is-hidden-mobile">memo</div>
    <div class="navbar-item navbar-mobile-context is-hidden-tablet">
      <div class="mobile-nav-segments" role="tablist" aria-label="category navigation">
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
    <div class="navbar-burger js-burger is-hidden-mobile" data-target="navbarMMemo">
      <span></span>
      <span></span>
      <span></span>
      <span></span>
    </div>
  </div>

  <div id="navbarMMemo" class="navbar-menu">
    <div class="navbar-start">
      {#if settings !== undefined}
        {#each settings.Categories as category}
          <a
            class="navbar-item is-unselectable is-tab"
            href={`/${category.Key}/`}
            use:route={{ active: { class: "is-active", absolute: false } }}
            onclick={closeNavbarMenu}
          >
            {category.Name}
          </a>
        {/each}
      {/if}
    </div>
  </div>
</nav>
<main class="app-main">
  <Router {routes} />
</main>

<style>
  .app-main {
    overflow-x: hidden;
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
    color: var(--bulma-text-weak-invert);
    font-size: 0.95rem;
    font-weight: 500;
    white-space: nowrap;
  }

  .mobile-nav-segment.is-active {
    color: var(--bulma-text-invert);
    box-shadow: inset 0 -2px 0 var(--bulma-link);
  }
</style>
