<script lang="ts">
  import "bulma/css/bulma.css";
  import Router, { location, link } from "svelte-spa-router";
  import DiaryList from "./routes/DiaryList.svelte";
  import DiaryEdit from "./routes/DiaryDetail.svelte";
  import HMemoList from "./routes/HMemoList.svelte";
  import HMemoEdit from "./routes/HMemoDetail.svelte";
  let page = "";

  const routes = {
    "/":  DiaryList,
    "/d/": DiaryList,
    "/d/:id": DiaryEdit,
    "/d/add": DiaryEdit,
    "/h/": HMemoList,
    "/h/:id": HMemoEdit,
    "/h/add": HMemoEdit,
  };

  $: {
    if ($location.indexOf("/d") >= 0) {
      page = "d";
    } else if ($location.indexOf("/h") >= 0) {
      page = "h";
    } else {
      page = "d";
    }
  }

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
    <a class="navbar-item is-unselectable" href="/" use:link>memo</a>
    <div class="navbar-burger js-burger" data-target="navbarMMemo">
      <span></span>
      <span></span>
      <span></span>
      <span></span>
    </div>
  </div>

  <div id="navbarMMemo" class="navbar-menu">
    <div class="navbar-start">
      <a class="navbar-item is-unselectable" href="/d/" use:link> diary </a>
      <a class="navbar-item is-unselectable" href="/h/" use:link> memo </a>
    </div>
  </div>
</nav>
<main>
  <div class="box">
    <Router {routes} />
  </div>
</main>

<style>
</style>
