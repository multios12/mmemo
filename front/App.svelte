<script lang="ts">
  import "bulma/css/bulma.css";
  import Router, { location, link } from "svelte-spa-router";
  import { type RouteDefinition } from "svelte-spa-router";
  import DiaryList from "./routes/DiaryList.svelte";
  import DiaryEdit from "./routes/DiaryDetail.svelte";
  import HMemoList from "./routes/MemoList.svelte";
  import HMemoEdit from "./routes/MemoDetail.svelte";
  import Planner from "./routes/Planner/Planner.svelte";
  import { onMount } from "svelte";
  import type { settingType } from "./models/settingType.js";
  import { settingsStore } from "./store.js";
  let page = "";

  let routes: RouteDefinition = {
    "/": DiaryList,
    "/d/:id": DiaryEdit,
    "/d/add": DiaryEdit,
    "/planner/": Planner,
    "/:category/": HMemoList,
    "/:category/:id": HMemoEdit,
    "/:category/add": HMemoEdit,
  };
  let settings: settingType;

  onMount(async () => {
    const r = await fetch("./api/memos");
    settings = <settingType>await r.json();
    settingsStore.update((s) => settings);
  });

  $: {
    if (settings !== undefined) {
      page = "";
      for (const category of settings.Categories) {
        if ($location.indexOf("/" + category.Key) >= 0) {
          page = category.Key;
        }
      }
    }
  }

  // navbarのバーガー開閉イベント
  document.addEventListener("DOMContentLoaded", () => {
    // Get all "navbar-burger" elements
    const $navbarBurgers = Array.prototype.slice.call(
      document.querySelectorAll(".navbar-burger"),
      0,
    );

    // Add a click event on each of them
    $navbarBurgers.forEach((el) => {
      el.addEventListener("click", () => {
        // Get the target from the "data-target" attribute
        const target = el.dataset.target;
        const _target = <HTMLElement>document.getElementById(target);

        // Toggle the "is-active" class on both the "navbar-burger" and the "navbar-menu"
        el.classList.toggle("is-active");
        _target.classList.toggle("is-active");
      });
    });
  });
</script>

<nav class="navbar is-transparent is-dark">
  <div class="navbar-brand">
    <div class="navbar-item is-unselectable has-text-weight-bold">memo</div>
    <div class="navbar-burger js-burger" data-target="navbarMMemo">
      <span></span>
      <span></span>
      <span></span>
      <span></span>
    </div>
  </div>

  <div id="navbarMMemo" class="navbar-menu">
    <div class="navbar-start">
      <a
        class="navbar-item is-unselectable is-tab"
        class:is-active={page === ""}
        href="/"
        use:link
      >
        {settings?.Diary?.Name}
      </a>
      {#if settings !== undefined}
        {#each settings.Categories as category}
          <a
            class="navbar-item is-unselectable is-tab"
            class:is-active={page === category.Key}
            href={`/${category.Key}/`}
            use:link
          >
            {category.Name}
          </a>
        {/each}
      {/if}
      <!--
      <a
        class="navbar-item is-unselectable is-tab"
        class:is-active={page === "/planner/"}
        href="/planner"
        use:link
        >planner
      </a>
      -->
    </div>
  </div>
</nav>
<main>
  <Router {routes} />
</main>

<style>
</style>
