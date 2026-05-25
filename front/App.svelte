<script lang="ts">
  import "mu-ui-lib/mu-theme.css";
  import "./app.css";
  import { RouterContext } from "@dvcol/svelte-simple-router/components";
  import EntryList from "./routes/EntryList.svelte";
  import EntryDetail from "./routes/EntryDetail.svelte";
  import EntryReference from "./routes/EntryReference.svelte";
  import Home from "./routes/Home.svelte";
  import Planner from "./routes/Planner/Planner.svelte";
  import SettingsPage from "./routes/Settings.svelte";
  import SettingsTemplatePage from "./routes/TemplateDetail.svelte";
  import AppShell from "./AppShell.svelte";
  import RouteOutlet from "./RouteOutlet.svelte";
  import {
    routerBasePath,
  } from "./basePath.js";

  const routes = [
    { path: "/", component: Home },
    { path: "/planner", component: Planner },
    { path: "/settings", component: SettingsPage },
    {
      path: "/settings/templates/:name",
      component: SettingsTemplatePage,
    },
    { path: "/:category/add", component: EntryDetail },
    { path: "/:category/:id/edit", component: EntryDetail },
    { path: "/:category/:id", component: EntryReference },
    { path: "/:category", component: EntryList },
  ];
  const options = {
    routes,
    base: routerBasePath || "/",
    hash: true,
    strict: false,
  };
</script>

<RouterContext {options}>
  <AppShell />
  <RouteOutlet />
</RouterContext>
