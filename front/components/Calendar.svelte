<script lang="ts">
  import { onMount } from "svelte";
  import { monthStart, monthEnd, weekStart } from "@formkit/tempo";
  import { addDay, diffDays, format } from "@formkit/tempo";

  export let Value: Date;

  let table: dayType[][] = [];

  type dayType = {
    date: Date;
    isHoliday: boolean;
    isSaturday: boolean;
    isDisable: boolean;
  };

  onMount(async () => {
    let targetDate = weekStart(Value, -1);
    let dayIndex = 6;
    let weekIndex = -1;
    while (
      diffDays(targetDate, monthEnd(Value)) < 0 ||
      (diffDays(targetDate, monthEnd(Value)) >= 0 && dayIndex < 6)
    ) {
      if (dayIndex == 6) {
        dayIndex = 0;
        weekIndex += 1;
      } else {
        dayIndex += 1;
      }
      targetDate = addDay(targetDate);
      if (table[weekIndex] == undefined) {
        table[weekIndex] = [];
      }

      table[weekIndex][dayIndex] = { date: targetDate };
      if (
        diffDays(targetDate, monthStart(Value)) < 0 ||
        diffDays(targetDate, monthEnd(Value)) > 0
      ) {
        table[weekIndex][dayIndex] = { date: targetDate, isDisable: true };
      } else if (dayIndex == 0) {
        table[weekIndex][dayIndex] = { date: targetDate, isHoliday: true };
      } else if (dayIndex == 6) {
        table[weekIndex][dayIndex] = { date: targetDate, isSaturday: true };
      }
    }
    console.log(dayIndex);
  });
</script>

<table class="table is-bordered">
  <caption>{format(Value, "YYYY/MM/DD", "ja")}</caption>
  {#if table != undefined && table.length >= 1}
    <thead class="p-0">
      <tr>
        {#each table[0] as date}
          <td class="p-1">{format(date.date, "ddd", "ja")}</td>
        {/each}
      </tr>
    </thead>
  {/if}
  {#each table as w}
    <tr>
      {#each w as date}
        <td
          class="p-0"
          class:holiday={date.isHoliday}
          class:saturday={date.isSaturday}
          class:disable={date.isDisable}
        >
          {format(date.date, "D", "ja")}
        </td>
      {/each}
    </tr>
  {/each}
</table>

<style>
  .holiday {
    background-color: rgb(128, 0, 0);
  }
  .disable {
    color: gray;
  }
  .saturday {
    background-color: navy;
  }
</style>
