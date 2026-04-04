<script lang="ts">
  import { goto, type RouteResult } from "@mateothegreat/svelte5-router";
  import type { entryType, listType } from "../models/entryModels.js";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faPlus, faNoteSticky } from "@fortawesome/free-solid-svg-icons";
  import { settingsStore } from "../store.js";

  library.add(faPlus, faNoteSticky);
  dom.watch();

  document
    .querySelector<HTMLDivElement>(".navbar")
    ?.classList.remove("is-hidden");

  interface Props {
    route?: RouteResult;
  }

  const emptyListModel = (): listType => ({ WritedMonths: [], Lines: [] });

  let { route: currentRoute = undefined }: Props = $props();
  let selectMonth = $state(
    `${new Date().getFullYear()}-${("00" + (new Date().getMonth() + 1)).slice(-2)}`,
  );

  const categoryKey = $derived(String(currentRoute?.result?.path?.params?.category ?? ""));
  const resolvedSettings = $derived.by(() => {
    const categorySetting = $settingsStore?.Categories?.find(
      (category) => category.Key === categoryKey,
    );

    return {
      ready: categorySetting !== undefined,
      useMonthFilter: categorySetting?.UseDate ?? false,
      showTags: categorySetting?.UseTag ?? false,
    };
  });

  let model: listType = $state(emptyListModel());
  let entries: entryType[] = $state([]);

  const addPath = () => `/${categoryKey}/add`;

  const detailPath = (entry: entryType) => `/${categoryKey}/${entry.Id}`;

  const showList = async () => {
    if (!categoryKey) {
      model = emptyListModel();
      entries = [];
      return;
    }

    if (resolvedSettings.useMonthFilter) {
      const month = selectMonth ?? "";
      const [entriesResponse, monthsResponse] = await Promise.all([
        fetch(`/api/${categoryKey}?month=${month}`),
        fetch(`/api/${categoryKey}?months=1`),
      ]);

      const entries = (await entriesResponse.json()) as entryType[];
      const writedMonths = (await monthsResponse.json()) as string[];
      const lines: entryType[] = entries.map((entry) => ({
        Id: entry.Id,
        Date: entry.Date,
        Outline: entry.Outline,
        Value: entry.Value,
        Tags: entry.Tags ?? [],
        HasDetail: entry.HasDetail ?? false,
      }));

      model = {
        WritedMonths: writedMonths,
        Lines: lines,
      };

      if (model.Lines.length === 0 && model.WritedMonths.length > 0) {
        selectMonth = model.WritedMonths[0];
      }
      return;
    }

    const response = await fetch(`/api/${categoryKey}`);
    entries = (await response.json()) as entryType[];
  };

  const listClick = (entry: entryType) => goto(detailPath(entry));

  $effect(() => {
    categoryKey;
    resolvedSettings.ready;
    selectMonth;

    if (!resolvedSettings.ready) {
      return;
    }

    showList();
  });
</script>

<div class="card px-10">
  <div class="card-content">
    <div class="columns">
      {#if resolvedSettings.useMonthFilter}
        <div class="column">
          <div class="select">
            <select class="select" bind:value={selectMonth}>
              {#each model.WritedMonths as v}
                <option value={v}>{v}</option>
              {/each}
            </select>
          </div>
        </div>
      {/if}
      <div class="column">
        <button
          class="button is-primary"
          aria-label={`add ${categoryKey}`}
          onclick={() => goto(addPath())}
        >
          <i class="fa-solid fa-plus"></i>
        </button>
      </div>
    </div>

    <table
      class="table is-hoverable is-fullwidth"
      class:is-striped={!resolvedSettings.useMonthFilter}
    >
      <tbody>
        {#if resolvedSettings.useMonthFilter}
          {#each model.Lines as v}
            <tr
              onclick={() =>
                listClick({
                  Id: v.Id,
                  Outline: v.Outline,
                  Date: v.Date,
                  Value: "",
                  Tags: v.Tags,
                })}
            >
              <td>
                <button>
                  {v.Date}
                  {v.Outline}{#if v.HasDetail}
                    <i
                      class="fa-solid fa-note-sticky has-text-grey-light"
                      style="vertical-align:middle"
                    ></i>
                  {/if}
                </button>
              </td>
              {#if resolvedSettings.showTags}
                <td>
                  <div class="tags are-medium">
                    {#each v.Tags as t}
                      <span class="tag">{t}</span>
                    {/each}
                  </div>
                </td>
              {/if}
            </tr>
          {/each}
        {:else}
          {#each entries as entry}
            <tr>
              <td onclick={() => listClick(entry)} onkeypress={() => listClick(entry)}>
                <button class="is-fullwidth">
                  {entry.Date}&nbsp;{entry.Outline}
                </button>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>
