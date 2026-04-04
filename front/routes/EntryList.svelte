<script lang="ts">
  import { goto, type RouteResult, route } from "@mateothegreat/svelte5-router";
  import type { listType } from "../models/diaryModels.js";
  import type { memoType } from "../models/memoModels.js";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faPlus, faNoteSticky } from "@fortawesome/free-solid-svg-icons";
  import { settingsStore } from "../store.js";

  library.add(faPlus, faNoteSticky);
  dom.watch();

  document
    .querySelector<HTMLDivElement>(".navbar")
    ?.classList.remove("is-hidden");

  interface Props {
    category?: string;
    route?: RouteResult;
    selectMonth?: string | null;
    useMonthFilter?: boolean;
    showTags?: boolean;
  }

  let {
    category = "",
    route: currentRoute = undefined,
    selectMonth = $bindable(
      `${new Date().getFullYear()}-${("00" + (new Date().getMonth() + 1)).slice(-2)}`,
    ),
    useMonthFilter = undefined,
    showTags = undefined,
  }: Props = $props();

  const entryCategory = $derived(
    category || String(currentRoute?.result?.path?.params?.category ?? ""),
  );
  const categorySetting = $derived(
    $settingsStore?.Categories?.find(
      (category) => category.Key === entryCategory,
    ),
  );
  const usesMonthFilter = $derived(
    useMonthFilter ?? categorySetting?.UseDate ?? false,
  );
  const showsTags = $derived(showTags ?? categorySetting?.UseTag ?? false);
  const settingsResolved = $derived(
    useMonthFilter !== undefined ||
      showTags !== undefined ||
      categorySetting !== undefined,
  );

  let model: listType = $state({ WritedMonths: [], Lines: [] });
  let memos: memoType[] = $state([]);

  const addPath = () => `/${entryCategory}/add`;

  const detailPath = (entry: memoType) => `/${entryCategory}/${entry.Id}`;

  const showList = async () => {
    if (!entryCategory) {
      model = { WritedMonths: [], Lines: [] };
      memos = [];
      return;
    }

    if (usesMonthFilter) {
      const month = selectMonth ?? "";
      const [entriesResponse, monthsResponse] = await Promise.all([
        fetch(`/api/${entryCategory}?month=${month}`),
        fetch(`/api/${entryCategory}?months=1`),
      ]);

      const entries = (await entriesResponse.json()) as memoType[];
      const writedMonths = (await monthsResponse.json()) as string[];

      model = {
        WritedMonths: writedMonths,
        Lines: entries.map((entry) => ({
          Id: entry.Id,
          Day: entry.Date,
          Outline: entry.Title,
          Tags: entry.Tags ?? [],
          IsDetail: entry.HasDetail ?? false,
          HCount: 0,
        })),
      };

      if (model.Lines.length === 0 && model.WritedMonths.length > 0) {
        selectMonth = model.WritedMonths[0];
      }
      return;
    }

    const response = await fetch(`/api/${entryCategory}`);
    memos = (await response.json()) as memoType[];
  };

  const listClick = (entry: memoType) => goto(detailPath(entry));

  $effect(() => {
    entryCategory;
    settingsResolved;
    usesMonthFilter;
    selectMonth;

    if (!settingsResolved) {
      return;
    }

    showList();
  });
</script>

<div class="card px-10">
  <div class="card-content">
    <div class="columns">
      {#if usesMonthFilter}
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
          aria-label={`add ${entryCategory}`}
          onclick={() => goto(addPath())}
        >
          <i class="fa-solid fa-plus"></i>
        </button>
      </div>
    </div>

    <table
      class="table is-hoverable is-fullwidth"
      class:is-striped={!usesMonthFilter}
    >
      <tbody>
        {#if usesMonthFilter}
          {#each model.Lines as v}
            <tr
              onclick={() =>
                listClick({
                  Id: v.Id,
                  Title: v.Outline,
                  Date: v.Day,
                  Value: "",
                  Tags: v.Tags,
                })}
            >
              <td>
                <button>
                  {v.Day}
                  {v.Outline}{#if v.IsDetail}
                    <i
                      class="fa-solid fa-note-sticky has-text-grey-light"
                      style="vertical-align:middle"
                    ></i>
                  {/if}
                </button>
              </td>
              {#if showsTags}
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
          {#each memos as m}
            <tr>
              <td onclick={() => listClick(m)} onkeypress={() => listClick(m)}>
                <button class="is-fullwidth">
                  {m.Date}&nbsp;{m.Title}
                </button>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>
