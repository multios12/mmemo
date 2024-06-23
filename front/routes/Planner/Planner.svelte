<script lang="ts">
  import { onMount } from "svelte";
  import { diffDays, format, monthEnd } from "@formkit/tempo";
  import type { dayType } from "./dayType.js";
  import Dropdown from "../../components/Dropdown.svelte";
  export let year: number = 2024;

  /** 表示情報テーブル */
  let table: dayType[][] = [];
  /** 祝日リスト */
  let holidays: dayType[];
  let viewList: { key: string; value: string }[] = [
    { key: "day", value: "日付揃え" },
    { key: "week", value: "曜日揃え" },
  ];
  let viewMode: string = "week";

  type planType = {
    /** プラン名 */
    name: string;
    /** 日付 */
    dates: {
      /** 開始日 */
      fromt: Date;
      /** 終了日 */
      to: Date;
      row: number;
    }[];
  };

  onMount(async () => {
    holidays = [
      { date: new Date(2024, 0, 1), tip: "元旦", isHoliday: true },
      { date: new Date(2024, 6, 15), tip: "海の日", isHoliday: true },
    ];

    viewList = [
      { key: "day", value: "日付揃え" },
      { key: "week", value: "曜日揃え" },
    ];

    for (let monthIndex = 0; monthIndex < 12; monthIndex++) {
      table[monthIndex] = [];

      for (let j = 1; j <= 31; j++) {
        let targetDate = new Date(year, monthIndex, j);
        let d = g(targetDate);
        table[monthIndex][j - 1] = d;
      }
    }
  });

  const g = (targetDate: Date): dayType => {
    var endDate = monthEnd(targetDate);
    for (let index = 0; index < holidays.length; index++) {
      const h = holidays[index];
      if (diffDays(targetDate, h.date) == 0) {
        return { date: targetDate, isHoliday: true, tip: h.tip };
      }
    }
    if (diffDays(targetDate, endDate) > 0) {
      return { date: targetDate, isDisable: true };
    } else if (targetDate.getDay() === 6) {
      return { date: targetDate, isSaturday: true };
    } else if (targetDate.getDay() === 0) {
      return { date: targetDate, isHoliday: true };
    } else {
      return { date: targetDate };
    }
  };
</script>

年間計画<Dropdown items={viewList} bind:key={viewMode} tabindex={undefined} />
<table class="planner table is-bordered">
  <thead class="p-0">
    <tr>
      <td></td>
      {#if table !== null && table.length > 0}
        {#each table[0] as day}
          <td class="is-small p-0"> {format(day.date, "D", "ja")} </td>
        {/each}
      {/if}
    </tr>
  </thead>

  {#each table as d}
    <tr>
      <td class="p-0"> {format(d[0].date, "M月", "ja")} </td>
      {#each d as day}
        <td
          class="day p-0"
          class:holiday={day.isHoliday}
          class:disable={day.isDisable}
          class:saturday={day.isSaturday}
        >
        </td>
      {/each}
    </tr>
  {/each}
</table>

<style>
  .planner .day {
    height: 40px;
    width: 40px;
  }

  .holiday {
    background-color: rgb(64, 0, 0);
  }
  .disable {
    background-color: gray;
  }
  .saturday {
    background-color: rgb(0, 0, 64);
  }
</style>
