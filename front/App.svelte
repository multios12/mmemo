<script lang="ts">
  import "bulma/css/bulma.css";
  import { Router, route } from "@mateothegreat/svelte5-router";
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
  let settings: settingType;

  onMount(async () => {
    const r = await fetch("/settings");
    settings = <settingType>await r.json();
    settingsStore.update((s) => settings);
  });

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
      {#if settings !== undefined}
        {#each settings.Categories as category}
          <a
            class="navbar-item is-unselectable is-tab"
            href={`/${category.Key}/`}
            use:route={{ active: { class: "is-active", absolute: false } }}
          >
            {category.Name}
          </a>
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
