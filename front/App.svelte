<script lang="ts">
  import "bulma/css/bulma.css";
  import { goto } from "@mateothegreat/svelte5-router";
  import { Router, type RouteConfig } from "@mateothegreat/svelte5-router";
  import DiaryList from "./routes/DiaryList.svelte";
  import DiaryEdit from "./routes/DiaryDetail.svelte";
  import HMemoList from "./routes/MemoList.svelte";
  import HMemoEdit from "./routes/MemoDetail.svelte";
  import { onMount } from "svelte";
  import type { settingType } from "./models/settingType.js";
  import { settingsStore } from "./store.js";
  import { getMemosSettingApi } from "./models/apiUrl.js";

  /** 選択中ページ */
  let page = $state("");

  let routes: RouteConfig[] = $state([
    { component: DiaryList },
    { path: "/d/(?<id>.?)", component: DiaryEdit },
    { path: "/(?<category>.+)", component: HMemoList },
    { path: "/d/add", component: DiaryEdit },
    { path: "/(?<category>.+)/(?<id>.+)", component: HMemoEdit },
    { path: "/(?<category>.+)/add", component: HMemoEdit },
  ]);

  let settings: settingType;

  onMount(async () => {
    const r = await getMemosSettingApi();
    settings = <settingType>await r.json();

    settingsStore.update((s) => settings);
  });

  $effect(() => {
    const location = window.location.href;
    if (settings !== undefined) {
      page = "";
      for (const category of settings.Categories) {
        if (location.indexOf("/" + category.Key) >= 0) {
          page = category.Key;
        }
      }
      routes = routes;
    }
  });

  // navbarのバーガー開閉イベント
  document.addEventListener("DOMContentLoaded", () => {
    // Get all "navbar-burger" elements
    const navbarBurgers = Array.prototype.slice.call(
      document.querySelectorAll(".navbar-burger"),
      0,
    );

    // Add a click event on each of them
    navbarBurgers.forEach((el) => {
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

  const onClick = (value: string) => {
    page = value;
    value = value === "" ? "" : `/${value}/`;
    goto(value);
  };
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
      <button
        class="navbar-item is-unselectable is-tab"
        class:is-active={page === ""}
        onclick={() => onClick("")}
      >
        {settings?.Diary?.Name}
      </button>
      {#if settings !== undefined}
        {#each settings.Categories as category}
          <button
            class="navbar-item is-unselectable is-tab"
            class:is-active={page === category.Key}
            onclick={() => onClick(category.Key)}
          >
            {category.Name}
          </button>
        {/each}
      {/if}
    </div>
  </div>
</nav>
<main>
  <Router {routes} />
</main>

<style>
</style>
