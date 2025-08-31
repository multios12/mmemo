<script lang="ts">
  import { run } from "svelte/legacy";
  import type { listType } from "../models/diaryModels.js";
  import { dom, library } from "@fortawesome/fontawesome-svg-core";
  import { faPlus, faNoteSticky } from "@fortawesome/free-solid-svg-icons";
  import { goto } from "@mateothegreat/svelte5-router";
  import { getMonthApi } from "../models/apiUrl.js";
  library.add(faPlus, faNoteSticky);
  dom.watch();

  // メニューバー表示化
  document
    .querySelector<HTMLDivElement>(".navbar")
    ?.classList.remove("is-hidden");

  interface Props {
    selectMonth?: string | null;
    route: any;
  }

  let {
    selectMonth = $bindable(
      `${new Date().getFullYear()}-${("00" + (new Date().getMonth() + 1)).slice(
        -2,
      )}`,
    ),
    route,
  }: Props = $props();
  export const showList = () => {
    let url = selectMonth !== null ? selectMonth.replace("-", "/") : "";
    let apiFetch = getMonthApi(url, route);
    apiFetch
      .then((r) => r.json())
      .then((r) => r as listType)
      .then((r) => {
        model = r;
        if (r.Lines.length == 0 && r.WritedMonths.length > 0) {
          selectMonth = r.WritedMonths[0];
        }
      });
  };

  let model: listType = $state({ WritedMonths: [], Lines: [] });

  /** 追加ボタンクリックイベント */
  const addClick = () => goto("/d/add");

  /** リストクリックイベント */
  const listClick = async (e: any, l: string) => goto("/d/" + l);

  run(() => {
    let params = route.result.path.params;

    let url = selectMonth !== null ? selectMonth.replace("-", "/") : "";
    getMonthApi(url, route)
      .then((r) => r.json())
      .then((r) => r as listType)
      .then((r) => {
        model = r;
        if (r.Lines.length == 0 && r.WritedMonths.length > 0) {
          selectMonth = r.WritedMonths[0];
        }
      });
  });
</script>

<div class="card px-10">
  <div class="card-content">
    <div class="columns">
      <div class="column">
        <div class="select">
          <select class="select" bind:value={selectMonth}>
            {#each model.WritedMonths as v}
              <option value={v}>{v}</option>
            {/each}
          </select>
        </div>
      </div>
      <div class="column">
        <button class="button is-primary" onclick={addClick} aria-label="add">
          <i class="fa-solid fa-plus"></i>
        </button>
      </div>
    </div>
    <table class="table is-hoverable is-fullwidth">
      <tbody>
        {#each model.Lines as v}
          <tr onclick={(e) => listClick(e, v.Day)}>
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
            <td>
              <div class="tags are-medium">
                {#each v.Tags as t}
                  <span class="tag">{t}</span>
                {/each}
              </div>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
