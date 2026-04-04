<script lang="ts">
  import { monthStart, monthEnd, weekStart } from "@formkit/tempo";
  import { addDay, diffDays, format } from "@formkit/tempo";

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

  const normalizeDate = (value: Date) => format(value, "YYYY-MM-DD", "ja");
  const activeDateSet = $derived(new Set(ActiveDates));
  const holidayDateSet = $derived(new Set(HolidayDates));

  const buildTable = (value: Date) => {
    const calendarTable: dayType[][] = [];
    let targetDate = weekStart(value, -1);
    let dayIndex = 6;
    let weekIndex = -1;
    while (
      diffDays(targetDate, monthEnd(value)) < 0 ||
      (diffDays(targetDate, monthEnd(value)) >= 0 && dayIndex < 6)
    ) {
      if (dayIndex == 6) {
        dayIndex = 0;
        weekIndex += 1;
      } else {
        dayIndex += 1;
      }
      targetDate = addDay(targetDate);
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
      } else if (dayIndex == 0 || holidayDateSet.has(normalizeDate(targetDate))) {
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
            {format(date.date, "ddd", "ja")}
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
                {format(date.date, "D", "ja")}
              </button>
            {:else}
              <span class="day-number">{format(date.date, "D", "ja")}</span>
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
    font-weight: 500;
    color: color-mix(in srgb, var(--bulma-text-weak) 82%, white 18%);
    background: color-mix(in srgb, var(--bulma-scheme-main) 88%, black 12%);
  }

  .holiday-head {
    color: #c98f8f;
  }

  .saturday-head {
    color: #8ea0df;
  }

  .calendar-day {
    padding: 0.2rem;
    height: 2.35rem;
    font-variant-numeric: tabular-nums;
    color: var(--bulma-text);
    background: color-mix(in srgb, var(--bulma-scheme-main) 93%, black 7%);
  }

  .day-number {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 1.8rem;
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
    outline: 2px solid rgba(255, 236, 186, 0.5);
    outline-offset: 2px;
  }

  .active-day .day-number {
    font-weight: 700;
    color: #f5efe2;
    background: linear-gradient(180deg, #9d7a2a 0%, #775915 100%);
    box-shadow:
      inset 0 0 0 1px rgba(255, 236, 186, 0.18),
      0 0 0 1px rgba(125, 93, 24, 0.3);
  }

  .holiday {
    color: #f3dede;
    background: linear-gradient(180deg, #5f1818 0%, #471111 100%);
  }

  .disable {
    color: color-mix(in srgb, var(--bulma-text-weak) 55%, transparent);
    background: color-mix(in srgb, var(--bulma-scheme-main) 96%, black 4%);
  }

  .saturday {
    color: #dde5ff;
    background: linear-gradient(180deg, #22306f 0%, #18234d 100%);
  }

</style>
