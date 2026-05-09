<script lang="ts">
  import { monthStart, monthEnd, weekStart } from "../lib/date.js";
  import { addDay, diffDays, format } from "../lib/date.js";

  interface Props {
    Value: Date;
    ActiveDates?: string[];
    HolidayDates?: string[];
    OnSelectDate?: (date: string) => void;
  }

  let {
    Value,
    ActiveDates = [],
    HolidayDates = [],
    OnSelectDate = undefined,
  }: Props = $props();

  type dayType = {
    date: Date;
    isHoliday?: boolean;
    isSaturday?: boolean;
    isDisable?: boolean;
    isActive?: boolean;
  };

  const normalizeDate = (value: Date) => format(value, "YYYY-MM-DD");
  const activeDateSet = $derived(new Set(ActiveDates));
  const holidayDateSet = $derived(new Set(HolidayDates));

  const buildTable = (value: Date) => {
    const calendarTable: dayType[][] = [];
    let targetDate = weekStart(value);
    let dayIndex = 0;
    let weekIndex = 0;
    const lastDate = monthEnd(value);
    while (
      diffDays(targetDate, lastDate) <= 0 ||
      dayIndex !== 0
    ) {
      if (calendarTable[weekIndex] == undefined) {
        calendarTable[weekIndex] = [];
      }

      calendarTable[weekIndex][dayIndex] = {
        date: targetDate,
        isActive: activeDateSet.has(normalizeDate(targetDate)),
      };
      if (
        diffDays(targetDate, monthStart(value)) < 0 ||
        diffDays(targetDate, monthEnd(value)) > 0
      ) {
        calendarTable[weekIndex][dayIndex] = {
          date: targetDate,
          isDisable: true,
          isActive: activeDateSet.has(normalizeDate(targetDate)),
        };
      } else if (
        dayIndex == 0 ||
        holidayDateSet.has(normalizeDate(targetDate))
      ) {
        calendarTable[weekIndex][dayIndex] = {
          date: targetDate,
          isHoliday: true,
          isActive: activeDateSet.has(normalizeDate(targetDate)),
        };
      } else if (dayIndex == 6) {
        calendarTable[weekIndex][dayIndex] = {
          date: targetDate,
          isSaturday: true,
          isActive: activeDateSet.has(normalizeDate(targetDate)),
        };
      }

      targetDate = addDay(targetDate);
      dayIndex += 1;
      if (dayIndex === 7) {
        dayIndex = 0;
        weekIndex += 1;
      }
    }

    return calendarTable;
  };

  const table = $derived(buildTable(Value));
</script>

<table class="table is-bordered calendar">
  {#if table != undefined && table.length >= 1}
    <thead>
      <tr>
        {#each table[0] as date, index}
          <th
            scope="col"
            class:holiday-head={index === 0}
            class:saturday-head={index === 6}
          >
            {format(date.date, "ddd")}
          </th>
        {/each}
      </tr>
    </thead>
  {/if}
  <tbody>
    {#each table as w}
      <tr>
        {#each w as date}
          <td
            class="calendar-day"
            class:holiday={date.isHoliday}
            class:saturday={date.isSaturday}
            class:disable={date.isDisable}
            class:active-day={date.isActive && !date.isDisable}
          >
            {#if date.isActive && !date.isDisable && OnSelectDate}
              <button
                class="day-number day-button"
                type="button"
                aria-label={`${normalizeDate(date.date)} を開く`}
                onclick={() => OnSelectDate?.(normalizeDate(date.date))}
              >
                {format(date.date, "D")}
              </button>
            {:else}
              <span class="day-number">{format(date.date, "D")}</span>
            {/if}
          </td>
        {/each}
      </tr>
    {/each}
  </tbody>
</table>

<style>
  .calendar {
    width: 100%;
    table-layout: fixed;
    border-collapse: collapse;
    background: color-mix(in srgb, var(--bulma-scheme-main) 94%, black 6%);
    border: 1px solid color-mix(in srgb, var(--bulma-border) 82%, white 10%);
  }

  .calendar th,
  .calendar td {
    width: calc(100% / 7);
    min-width: 0;
    text-align: center;
    vertical-align: middle;
  }

  .calendar th {
    padding: 0.5rem 0.25rem;
    font-weight: 600;
    color: color-mix(in srgb, var(--bulma-text-weak) 84%, white 16%);
    background: color-mix(in srgb, var(--bulma-border) 82%, black 18%);
    border-color: color-mix(in srgb, var(--bulma-border) 82%, white 8%);
  }

  .holiday-head {
    color: #c79aa0;
  }

  .saturday-head {
    color: #98a8d1;
  }

  .calendar-day {
    padding: 0.2rem;
    height: 2.35rem;
    font-variant-numeric: tabular-nums;
    color: var(--bulma-text);
    background: color-mix(in srgb, var(--bulma-scheme-main) 93%, black 7%);
    border-color: color-mix(in srgb, var(--bulma-border) 82%, white 8%);
  }

  .day-number {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 2.2rem;
    height: 1.8rem;
    border-radius: 999px;
    transition:
      background-color 0.15s ease,
      box-shadow 0.15s ease,
      transform 0.15s ease;
  }

  .day-button {
    border: none;
    cursor: pointer;
    font: inherit;
  }

  .day-button:hover {
    transform: translateY(-1px);
  }

  .day-button:focus-visible {
    outline: 2px solid rgba(110, 193, 183, 0.55);
    outline-offset: 2px;
  }

  .active-day .day-number {
    font-weight: 700;
    color: #edf8f6;
    background: color-mix(in srgb, #2d8f86 62%, var(--bulma-scheme-main));
    box-shadow:
      inset 0 0 0 1px rgba(170, 230, 223, 0.14),
      0 0 0 1px rgba(24, 79, 74, 0.24);
  }

  .holiday {
    color: #f1e3e5;
    background: color-mix(in srgb, #6d3d45 42%, var(--bulma-scheme-main));
  }

  .disable {
    color: color-mix(in srgb, var(--bulma-text-weak) 55%, transparent);
    background: color-mix(in srgb, var(--bulma-scheme-main) 96%, black 4%);
  }

  .saturday {
    color: #e4e9f8;
    background: color-mix(in srgb, #40527f 40%, var(--bulma-scheme-main));
  }
</style>
