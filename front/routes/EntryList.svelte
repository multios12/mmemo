<script lang="ts">
  import { useNavigate, useRoute } from "@dvcol/svelte-simple-router/router";
  import type { entryType, listType } from "../models/entryModels.js";
  import { format } from "@formkit/tempo";
  import { settingsStore } from "../store.js";
  import AppIcon from "../components/AppIcon.svelte";
  import Calendar from "../components/Calendar.svelte";
  import { apiPath } from "../basePath.js";

  const emptyListModel = (): listType => ({ WritedMonths: [], Lines: [] });
  const noOutlineLabel = "(no outline)";
  const noTagLabel = "タグなし";

  type viewModeType = "default" | "tag" | "outline";
  type groupedEntriesByTagType = {
    label: string;
    entries: entryType[];
  };
  type groupedEntriesByOutlineType = {
    outline: string;
    entries: entryType[];
  };

  const { location } = $derived(useRoute());
  const { push } = useNavigate();
  let selectMonth = $state(
    `${new Date().getFullYear()}-${("00" + (new Date().getMonth() + 1)).slice(-2)}`,
  );
  let sortOrder = $state<"desc" | "asc">("desc");
  let viewMode = $state<viewModeType>("default");

  const categoryKey = $derived(String(location?.params?.category ?? ""));
  const sortOrderStorageKey = $derived(
    `entry-list-sort-order:${categoryKey || "default"}`,
  );
  const viewModeStorageKey = $derived(
    `entry-list-view-mode:${categoryKey || "default"}`,
  );
  const resolvedSettings = $derived.by(() => {
    const categorySetting = $settingsStore?.Categories?.find(
      (category) => category.Key === categoryKey,
    );

    return {
      ready: categorySetting !== undefined,
      useMonthFilter: categorySetting?.UseDate ?? false,
      showTags: categorySetting?.UseTag ?? false,
      tagsLabel: categorySetting?.Fields?.Tags || "タグ",
      outlineLabel: categorySetting?.Fields?.Outline || "見出し",
      outlineIcon: categorySetting?.Fields?.OutlineIcon || "folder-tree",
    };
  });

  let model: listType = $state(emptyListModel());
  let entries: entryType[] = $state([]);
  let isLoading = $state(false);
  let hasLoadedInitialList = $state(false);
  let isCalendarExpanded = $state(false);

  const addPath = () => `/${categoryKey}/add`;
  const addPathWithPreset = (preset: {
    outline?: string;
    tag?: string;
    previousEntryId?: number;
  } = {}) => {
    const searchParams = new URLSearchParams();
    if (preset.outline !== undefined) {
      searchParams.set("presetOutline", preset.outline);
    }
    if (preset.tag !== undefined) {
      searchParams.set("presetTag", preset.tag);
    }
    if (preset.previousEntryId !== undefined) {
      searchParams.set("previousEntryId", String(preset.previousEntryId));
    }

    const query = searchParams.toString();
    return query ? `${addPath()}?${query}` : addPath();
  };
  const detailPath = (entry: entryType) => `/${categoryKey}/${entry.Id}`;
  const isPreviewLine = (line: string) =>
    line !== "" &&
    !/^[-*_]{3,}$/.test(line) &&
    !/^([-*+]\s+|\d+\.\s+)/.test(line);
  const carryOverMarker = "----ここまで前回内容で置換";
  const extractPreviewLine = (value: string) => {
    const lines = value
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter((line) => line !== "");

    for (let index = 0; index < lines.length; index += 1) {
      if (/^#{1,3}\s+/.test(lines[index])) {
        const nextLine = lines[index + 1] ?? "";
        return isPreviewLine(nextLine) ? nextLine : "";
      }
    }

    return "";
  };
  const splitPreviewSections = (value: string) => {
    const markerPattern = new RegExp(`^${carryOverMarker.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}$`, "m");
    const match = value.match(markerPattern);
    if (!match || match.index === undefined) {
      return {
        hasMarker: false,
        before: value,
        after: "",
      };
    }

    const before = value.slice(0, match.index).trimEnd();
    const after = value.slice(match.index + match[0].length).trimStart();
    return { hasMarker: true, before, after };
  };
  const extractGroupedPreviews = (value: string) => {
    const { hasMarker, before, after } = splitPreviewSections(value);
    return {
      hasMarker,
      upper: extractPreviewLine(before),
      lower: hasMarker ? extractPreviewLine(after) : extractPreviewLine(before),
    };
  };
  const compareEntries = (left: entryType, right: entryType) => {
    const dateCompare = left.Date.localeCompare(right.Date);
    if (dateCompare !== 0) {
      return dateCompare;
    }

    const leftId = left.Id ?? 0;
    const rightId = right.Id ?? 0;
    return leftId - rightId;
  };
  const compareEntriesByRecentUpdate = (left: entryType, right: entryType) => {
    const leftUpdatedAt = left.UpdatedAt ?? left.CreatedAt ?? "";
    const rightUpdatedAt = right.UpdatedAt ?? right.CreatedAt ?? "";
    const updatedAtCompare = rightUpdatedAt.localeCompare(leftUpdatedAt);
    if (updatedAtCompare !== 0) {
      return updatedAtCompare;
    }

    return compareEntries(right, left);
  };
  const latestUpdatedEntry = (targetEntries: entryType[]) =>
    [...targetEntries].sort(compareEntriesByRecentUpdate)[0];
  const groupUpperPreview = (targetEntries: entryType[]) =>
    extractGroupedPreviews(latestUpdatedEntry(targetEntries)?.Value ?? "").upper;
  const displayEntries = $derived.by(() => {
    const source = resolvedSettings.useMonthFilter ? model.Lines : entries;
    const sorted = [...source].sort(compareEntries);
    return sortOrder === "asc" ? sorted : sorted.reverse();
  });
  const groupedEntriesByTag = $derived.by(() => {
    const groups = new Map<string, entryType[]>();

    for (const entry of displayEntries) {
      const tags = entry.Tags.length > 0 ? entry.Tags : [noTagLabel];
      for (const tag of tags) {
        if (!groups.has(tag)) {
          groups.set(tag, []);
        }
        groups.get(tag)?.push(entry);
      }
    }

    const grouped = [...groups.entries()]
      .map(([label, entries]) => ({ label, entries }) satisfies groupedEntriesByTagType)
      .sort((left, right) => left.label.localeCompare(right.label, "ja"));

    return sortOrder === "asc" ? grouped : grouped.reverse();
  });
  const groupedEntriesByOutline = $derived.by(() => {
    const groups = new Map<string, entryType[]>();

    for (const entry of displayEntries) {
      const outline = entry.Outline || noOutlineLabel;
      if (!groups.has(outline)) {
        groups.set(outline, []);
      }
      groups.get(outline)?.push(entry);
    }

    const grouped = [...groups.entries()].map(
      ([outline, entries]) => ({ outline, entries }) satisfies groupedEntriesByOutlineType,
    );

    grouped.sort((left, right) => left.outline.localeCompare(right.outline, "ja"));

    return sortOrder === "asc" ? grouped : grouped.reverse();
  });
  const hasEntries = $derived(displayEntries.length > 0);
  const calendarDate = $derived.by(() => {
    if (!resolvedSettings.useMonthFilter || !selectMonth) {
      return new Date();
    }

    const [year, month] = selectMonth.split("-").map(Number);
    if (!year || !month) {
      return new Date();
    }

    return new Date(year, month - 1, 1);
  });
  const calendarActiveDates = $derived.by(() =>
    resolvedSettings.useMonthFilter
      ? model.Lines.map((entry) => entry.Date)
      : [],
  );
  const calendarHolidayDates = $derived.by(() =>
    ($settingsStore?.Holidays ?? []).map((holiday) => holiday.Date),
  );
  const currentMonthIndex = $derived(model.WritedMonths.indexOf(selectMonth));
  const previousCalendarMonth = $derived.by(() => {
    if (currentMonthIndex < 0) {
      return undefined;
    }

    return model.WritedMonths[currentMonthIndex + 1];
  });
  const nextCalendarMonth = $derived.by(() => {
    if (currentMonthIndex <= 0) {
      return undefined;
    }

    return model.WritedMonths[currentMonthIndex - 1];
  });

  const showList = async () => {
    if (!categoryKey) {
      model = emptyListModel();
      entries = [];
      hasLoadedInitialList = true;
      return;
    }

    isLoading = true;
    try {
      if (resolvedSettings.useMonthFilter) {
        const month = selectMonth ?? "";
        const [entriesResponse, monthsResponse] = await Promise.all([
          fetch(`${apiPath(categoryKey)}?month=${month}`),
          fetch(`${apiPath(categoryKey)}?months=1`),
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
          CreatedAt: entry.CreatedAt,
          UpdatedAt: entry.UpdatedAt,
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

      const response = await fetch(apiPath(categoryKey));
      entries = (await response.json()) as entryType[];
    } finally {
      isLoading = false;
      hasLoadedInitialList = true;
    }
  };

  const listClick = (entry: entryType) => push({ path: detailPath(entry) });
  const calendarDateClick = (date: string) => {
    const targetEntry = model.Lines.find((entry) => entry.Date === date);
    if (targetEntry) {
      push({ path: detailPath(targetEntry) });
    }
  };
  const moveCalendarMonth = (month: string | undefined) => {
    if (month) {
      selectMonth = month;
    }
  };
  const toggleSortOrder = () =>
    (sortOrder = sortOrder === "asc" ? "desc" : "asc");
  const setViewMode = (nextMode: viewModeType) => (viewMode = nextMode);
  const outlineIconMap: Record<string, string> = {
    note: "note-sticky",
    tree: "folder-tree",
    book: "book",
    group: "layer-group",
    tag: "tags",
    calendar: "calendar-days",
    document: "file-lines",
    list: "list-ul",
    person: "user",
    user: "user",
    "address-card": "address-card",
    "note-sticky": "note-sticky",
    "folder-tree": "folder-tree",
    "layer-group": "layer-group",
    "calendar-days": "calendar-days",
    "file-lines": "file-lines",
  };
  const outlineIconName = $derived.by(() => {
    const rawIcon = (resolvedSettings.outlineIcon || "tree").trim().toLowerCase();
    return outlineIconMap[rawIcon] ?? outlineIconMap.note;
  });

  $effect(() => {
    if (typeof window === "undefined") {
      return;
    }

    const storedSortOrder = window.localStorage.getItem(sortOrderStorageKey);
    if (storedSortOrder === "asc" || storedSortOrder === "desc") {
      sortOrder = storedSortOrder;
    }

    const storedViewMode = window.localStorage.getItem(viewModeStorageKey);
    if (storedViewMode === "default" || storedViewMode === "tag" || storedViewMode === "outline") {
      viewMode = storedViewMode;
    }
  });

  $effect(() => {
    if (typeof window === "undefined") {
      return;
    }

    window.localStorage.setItem(sortOrderStorageKey, sortOrder);
    window.localStorage.setItem(viewModeStorageKey, viewMode);
  });

  $effect(() => {
    categoryKey;
    resolvedSettings.ready;
    resolvedSettings.showTags;
    selectMonth;

    if (!resolvedSettings.ready) {
      return;
    }

    if (!resolvedSettings.showTags && viewMode === "tag") {
      viewMode = "default";
    }

    showList();
  });
</script>

<div class="entry-list-page px-5 py-4">
  {#if hasLoadedInitialList && !isLoading}
    <div class="toolbar">
      <div class="toolbar-group">
        <div class="toolbar-controls">
          {#if resolvedSettings.useMonthFilter}
            <div class="select is-fullwidth-mobile toolbar-month-select">
              <select class="select" bind:value={selectMonth}>
                {#each model.WritedMonths as v}
                  <option value={v}>{v}</option>
                {/each}
              </select>
            </div>
          {/if}
          {#if hasEntries}
            <div class="toolbar-mode-row">
              <button
                class="button is-light sort-button"
                aria-label={sortOrder === "asc" ? "sort ascending" : "sort descending"}
                title={sortOrder === "asc" ? "昇順" : "降順"}
                onclick={toggleSortOrder}
              >
                <span class="sort-icon" aria-hidden="true">
                  <AppIcon name={sortOrder === "asc" ? "sort-asc" : "sort-desc"} />
                </span>
              </button>
              <div class="view-mode-picker" role="group" aria-label="表示モード">
                <span class="view-mode-label">表示モード</span>
                <div class="view-mode-options">
                  <button
                    class="view-mode-option"
                    class:is-active={viewMode === "default"}
                    type="button"
                    aria-pressed={viewMode === "default"}
                    onclick={() => setViewMode("default")}
                  >
                    <span class="icon"><AppIcon name="layer-group" /></span>
                    <span>通常</span>
                  </button>
                  {#if resolvedSettings.showTags}
                    <button
                      class="view-mode-option"
                      class:is-active={viewMode === "tag"}
                      type="button"
                      aria-pressed={viewMode === "tag"}
                      onclick={() => setViewMode("tag")}
                    >
                      <span class="icon"><AppIcon name="tags" /></span>
                      <span>{resolvedSettings.tagsLabel}</span>
                    </button>
                  {/if}
                <button
                  class="view-mode-option"
                  class:is-active={viewMode === "outline"}
                  type="button"
                  aria-pressed={viewMode === "outline"}
                  onclick={() => setViewMode("outline")}
                >
                  <span class="icon"><AppIcon name={outlineIconName} /></span>
                  <span>{resolvedSettings.outlineLabel}</span>
                </button>
                </div>
              </div>
            </div>
          {/if}
        </div>
      </div>
      <button
        class="button add-button add-button-primary is-hidden-mobile"
        aria-label={`add ${categoryKey}`}
        onclick={() => push({ path: addPath() })}
      >
        <span class="icon"><AppIcon name="plus" /></span>
        <span>新規</span>
      </button>
    </div>

    {#if resolvedSettings.useMonthFilter}
      <div class="calendar-panel is-hidden-mobile">
        <Calendar
          Value={calendarDate}
          ActiveDates={calendarActiveDates}
          HolidayDates={calendarHolidayDates}
          OnSelectDate={calendarDateClick}
        />
      </div>
    {/if}
  {/if}

  {#if !hasLoadedInitialList || isLoading}
    <div class="entry-empty has-text-grey">読み込み中...</div>
  {:else if !hasEntries}
    <div class="entry-empty">
      <p class="title is-6 mb-2">まだ記録がありません</p>
      <p class="has-text-grey mb-4">最初のエントリを作成すると、ここに一覧が表示されます。</p>
    </div>
  {:else}
    {#if viewMode === "default"}
      <div class="entry-list">
        {#each displayEntries as entry}
          <button class="entry-row" onclick={() => listClick(entry)}>
            <div class="entry-row-main">
              <div class="entry-row-head">
                <span class="entry-date">{entry.Date}</span>
                {#if entry.HasDetail}
                  <span class="icon has-text-grey-light">
                    <AppIcon name="note-sticky" />
                  </span>
                {/if}
              </div>
              <div class="entry-outline">{entry.Outline || noOutlineLabel}</div>
              {#if extractPreviewLine(entry.Value)}
                <div class="entry-preview">{extractPreviewLine(entry.Value)}</div>
              {/if}
            </div>
            {#if resolvedSettings.showTags && entry.Tags.length > 0}
              <div class="tags are-medium entry-tags entry-tags-badge">
                {#each entry.Tags as tag}
                  <span class="tag"><span class="icon"><AppIcon name="tags" /></span><span>{tag}</span></span>
                {/each}
              </div>
            {/if}
          </button>
        {/each}
      </div>
    {:else if viewMode === "tag"}
      <div class="entry-group-list">
        {#each groupedEntriesByTag as group}
          <section class="entry-group-card">
            <div class="entry-group-header">
              <div class="entry-group-heading">
                <div class="entry-group-title">
                  <span class="icon"><AppIcon name="tags" /></span>
                  <span>{group.label}</span>
                </div>
                {#if extractGroupedPreviews(latestUpdatedEntry(group.entries)?.Value ?? "").hasMarker &&
                  groupUpperPreview(group.entries)}
                  <div class="entry-group-preview">{groupUpperPreview(group.entries)}</div>
                {/if}
              </div>
              <div class="entry-group-meta">
                <span class="entry-group-count">{group.entries.length}件</span>
                <button
                  class="button is-small entry-group-add-button"
                  type="button"
                  aria-label={`${group.label} で新規作成`}
                  onclick={() =>
                    push({
                      path: addPathWithPreset({
                        tag: group.label === noTagLabel ? "" : group.label,
                        previousEntryId: latestUpdatedEntry(group.entries)?.Id,
                      }),
                    })}
                >
                  <span class="icon"><AppIcon name="plus" /></span>
                  <span>新規</span>
                </button>
              </div>
            </div>
            <div class="entry-list entry-list-nested">
              {#each group.entries as entry}
                <button class="entry-row" onclick={() => listClick(entry)}>
                  <div class="entry-row-main">
                    <div class="entry-row-head">
                      <span class="entry-date">{entry.Date}</span>
                      <span class="entry-date-outline">
                        {entry.Outline || noOutlineLabel}
                      </span>
                      {#if extractGroupedPreviews(entry.Value).lower}
                        <span class="entry-date-preview">
                          {extractGroupedPreviews(entry.Value).lower}
                        </span>
                      {/if}
                      {#if entry.HasDetail}
                        <span class="icon has-text-grey-light">
                          <AppIcon name="note-sticky" />
                        </span>
                      {/if}
                    </div>
                  </div>
                  {#if false}
                    <div class="tags are-medium entry-tags">
                      {#each entry.Tags as tag}
                        <span class="tag">{tag}</span>
                      {/each}
                    </div>
                  {/if}
                </button>
              {/each}
            </div>
          </section>
        {/each}
      </div>
    {:else}
      <div class="entry-group-list">
        {#each groupedEntriesByOutline as group}
          <section class="entry-group-card">
            <div class="entry-group-header">
              <div class="entry-group-heading">
                <div class="entry-group-title">
                  <span class="icon"><AppIcon name={outlineIconName} /></span>
                  <span>{group.outline}</span>
                </div>
                {#if extractGroupedPreviews(latestUpdatedEntry(group.entries)?.Value ?? "").hasMarker &&
                  groupUpperPreview(group.entries)}
                  <div class="entry-group-preview">{groupUpperPreview(group.entries)}</div>
                {/if}
              </div>
              <div class="entry-group-meta">
                <span class="entry-group-count">{group.entries.length}件</span>
                <button
                  class="button is-small entry-group-add-button"
                  type="button"
                  aria-label={`${group.outline} で新規作成`}
                  onclick={() =>
                    push({
                      path: addPathWithPreset({
                        outline: group.outline === noOutlineLabel ? "" : group.outline,
                        previousEntryId: latestUpdatedEntry(group.entries)?.Id,
                      }),
                    })}
                >
                  <span class="icon"><AppIcon name="plus" /></span>
                  <span>新規</span>
                </button>
              </div>
            </div>
            <div class="outline-entry-list">
              {#each group.entries as entry}
                <button class="outline-entry-row" onclick={() => listClick(entry)}>
                  <div class="entry-row-head">
                    <span class="entry-date">{entry.Date}</span>
                    {#if extractGroupedPreviews(entry.Value).lower}
                      <span class="entry-date-preview">
                        {extractGroupedPreviews(entry.Value).lower}
                      </span>
                    {/if}
                    {#if entry.HasDetail}
                      <span class="icon has-text-grey-light">
                        <AppIcon name="note-sticky" />
                      </span>
                    {/if}
                  </div>
                  {#if resolvedSettings.showTags && entry.Tags.length > 0}
                    <div class="tags are-medium entry-tags entry-tags-badge">
                      {#each entry.Tags as tag}
                        <span class="tag"><span class="icon"><AppIcon name="tags" /></span><span>{tag}</span></span>
                      {/each}
                    </div>
                  {/if}
                </button>
              {/each}
            </div>
          </section>
        {/each}
      </div>
    {/if}
  {/if}

  {#if hasLoadedInitialList && !isLoading}
    <div class="mobile-actions is-hidden-tablet">
      {#if resolvedSettings.useMonthFilter}
        <button
          class="button mobile-calendar-toggle-button"
          type="button"
          aria-expanded={isCalendarExpanded}
          aria-label={isCalendarExpanded ? "hide calendar" : "show calendar"}
          onclick={() => (isCalendarExpanded = !isCalendarExpanded)}
        >
          <span class="icon"><AppIcon name="calendar-days" /></span>
          <span class="mobile-calendar-toggle-label">カレンダー</span>
          <span class="icon mobile-calendar-toggle-chevron">
            <AppIcon name={isCalendarExpanded ? "chevron-up" : "chevron-down"} />
          </span>
        </button>
      {/if}

      <button
        class="button mobile-add-button add-button-primary"
        aria-label={`add ${categoryKey}`}
        onclick={() => push({ path: addPath() })}
      >
        <span class="icon"><AppIcon name="plus" /></span>
        <span>新規</span>
      </button>
    </div>

    {#if resolvedSettings.useMonthFilter && isCalendarExpanded}
      <div
        class="calendar-popup-backdrop is-hidden-tablet"
        role="button"
        tabindex="0"
        aria-label="close calendar popup"
        onclick={() => (isCalendarExpanded = false)}
        onkeydown={(event) => {
          if (event.key === "Enter" || event.key === " ") {
            isCalendarExpanded = false;
          }
        }}
      ></div>
      <div class="calendar-panel calendar-panel-mobile is-hidden-tablet">
        <div class="calendar-popup-card">
          <div class="calendar-popup-heading">
            <div class="calendar-popup-heading-side">
              {#if previousCalendarMonth}
                <button
                  class="button is-ghost calendar-month-nav"
                  type="button"
                  aria-label="previous month"
                  onclick={() => moveCalendarMonth(previousCalendarMonth)}
                >
                  <span class="icon"><AppIcon name="chevron-left" /></span>
                </button>
              {/if}
            </div>
            <div class="calendar-popup-heading-label">{format(calendarDate, "YYYY/MM", "ja")}</div>
            <div class="calendar-popup-heading-side">
              {#if nextCalendarMonth}
                <button
                  class="button is-ghost calendar-month-nav"
                  type="button"
                  aria-label="next month"
                  onclick={() => moveCalendarMonth(nextCalendarMonth)}
                >
                  <span class="icon"><AppIcon name="chevron-right" /></span>
                </button>
              {/if}
            </div>
          </div>
          <Calendar
            Value={calendarDate}
            ActiveDates={calendarActiveDates}
            HolidayDates={calendarHolidayDates}
            OnSelectDate={calendarDateClick}
          />
        </div>
      </div>
    {/if}
  {/if}
</div>

<style>
  .entry-list-page {
    max-width: 960px;
    margin: 0 auto;
  }

  .toolbar {
    display: flex;
    gap: 0.75rem;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    flex-wrap: wrap;
  }

  .toolbar-group {
    flex: 1 1 220px;
  }

  .toolbar-controls {
    display: flex;
    gap: 0.75rem;
    align-items: center;
    flex-wrap: wrap;
  }

  .toolbar-mode-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  .view-mode-picker {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    flex-wrap: wrap;
  }

  .view-mode-label {
    color: var(--bulma-text-weak);
    font-size: 0.86rem;
    font-weight: 600;
    white-space: nowrap;
  }

  .view-mode-options {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.22rem;
    border: 1px solid var(--bulma-border);
    border-radius: 999px;
    background: color-mix(in srgb, var(--bulma-scheme-main) 94%, black 6%);
  }

  .view-mode-option {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    min-width: 0;
    padding: 0.45rem 0.8rem;
    border: none;
    border-radius: 999px;
    background: transparent;
    color: var(--bulma-text-weak);
    font-size: 0.9rem;
    cursor: pointer;
  }

  .view-mode-option.is-active {
    background: color-mix(in srgb, #2d8f86 22%, var(--bulma-scheme-main));
    color: color-mix(in srgb, var(--bulma-text) 92%, white 8%);
    box-shadow: inset 0 0 0 1px color-mix(in srgb, #2d8f86 26%, transparent);
  }

  .add-button {
    min-width: 7rem;
  }

  .add-button-primary {
    background: color-mix(in srgb, #2d8f86 62%, var(--bulma-scheme-main));
    border: 1px solid color-mix(in srgb, #2d8f86 74%, black 26%);
    color: #edf8f6;
    box-shadow: 0 10px 24px rgba(10, 31, 29, 0.16);
  }

  .sort-button {
    min-width: 2.75rem;
    background: color-mix(in srgb, var(--bulma-border) 74%, var(--bulma-scheme-main));
    border: 1px solid color-mix(in srgb, var(--bulma-border) 86%, white 14%);
    color: color-mix(in srgb, var(--bulma-text) 88%, white 12%);
  }

  .sort-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 1.1rem;
    height: 1.1rem;
    line-height: 1;
  }

  .entry-list {
    display: flex;
    flex-direction: column;
    gap: 0.85rem;
  }

  .entry-list-nested {
    gap: 0.75rem;
  }

  .entry-group-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .entry-group-card {
    border: 1px solid var(--bulma-border);
    border-radius: 18px;
    background: color-mix(in srgb, var(--bulma-scheme-main) 95%, white 5%);
    padding: 0.9rem;
  }

  .entry-group-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.85rem;
  }

  .entry-group-title {
    display: inline-flex;
    align-items: center;
    gap: 0.45rem;
    font-size: 1rem;
    font-weight: 700;
    color: var(--bulma-text);
  }

  .entry-group-heading {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
    min-width: 0;
  }

  .entry-group-preview {
    color: var(--bulma-text-weak);
    font-size: 0.92rem;
    line-height: 1.45;
    min-width: 0;
  }

  .entry-group-count {
    color: var(--bulma-text-weak);
    font-size: 0.88rem;
    font-variant-numeric: tabular-nums;
  }

  .entry-group-meta {
    display: inline-flex;
    align-items: center;
    gap: 0.55rem;
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .entry-group-add-button {
    border-radius: 999px;
    border-color: color-mix(in srgb, #2d8f86 38%, var(--bulma-border));
    background: color-mix(in srgb, #2d8f86 16%, var(--bulma-scheme-main));
    color: color-mix(in srgb, var(--bulma-text) 92%, white 8%);
    font-weight: 600;
  }

  .outline-entry-list {
    display: flex;
    flex-direction: column;
  }

  .outline-entry-row {
    width: 100%;
    text-align: left;
    padding: 0.8rem 0;
    background: transparent;
    border: none;
    border-top: 1px solid color-mix(in srgb, var(--bulma-border) 78%, transparent);
  }

  .outline-entry-list .outline-entry-row:first-child {
    border-top: none;
    padding-top: 0;
  }

  .outline-entry-list .outline-entry-row:last-child {
    padding-bottom: 0;
  }

  .calendar-panel {
    margin-bottom: 1rem;
    overflow-x: auto;
  }

  .calendar-popup-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(8, 10, 16, 0.38);
    backdrop-filter: blur(6px);
    z-index: 18;
  }

  .calendar-panel-mobile {
    margin-bottom: 0;
    padding: 0 0.9rem;
    position: fixed;
    left: 0;
    right: 0;
    bottom: calc(4.75rem + env(safe-area-inset-bottom));
    z-index: 19;
    max-height: min(20rem, 55vh);
    overflow: auto;
    transform-origin: bottom center;
    animation: calendar-popup-slide-in 0.5s cubic-bezier(0.22, 1, 0.36, 1);
  }

  .calendar-popup-card {
    border: 1px solid color-mix(in srgb, var(--bulma-border) 82%, white 18%);
    border-radius: 18px;
    background: color-mix(in srgb, var(--bulma-scheme-main) 93%, black 7%);
    box-shadow:
      0 22px 48px rgba(0, 0, 0, 0.34),
      0 8px 18px rgba(0, 0, 0, 0.2);
    padding: 0.75rem;
  }

  .calendar-popup-heading {
    margin-bottom: 0.65rem;
    display: grid;
    grid-template-columns: 2.5rem 1fr 2.5rem;
    align-items: center;
    gap: 0.25rem;
    color: color-mix(in srgb, var(--bulma-text) 90%, white 10%);
  }

  .calendar-popup-heading-side {
    display: flex;
    justify-content: center;
  }

  .calendar-popup-heading-label {
    text-align: center;
    font-size: 0.95rem;
    font-weight: 700;
    letter-spacing: 0.04em;
  }

  .calendar-month-nav {
    width: 2.25rem;
    height: 2.25rem;
    border-radius: 999px;
    color: color-mix(in srgb, var(--bulma-text) 88%, white 12%);
    background: color-mix(in srgb, var(--bulma-scheme-main) 88%, black 12%);
  }

  .calendar-panel :global(table) {
    width: 100%;
    min-width: 20rem;
    background: var(--bulma-scheme-main);
    border-radius: 14px;
    overflow: hidden;
  }

  .mobile-actions {
    position: fixed;
    left: 0;
    right: 0;
    bottom: 0;
    padding: 0.85rem 0.9rem calc(0.85rem + env(safe-area-inset-bottom));
    z-index: 20;
    display: flex;
    gap: 0.65rem;
    align-items: center;
    background-color: var(--bulma-border);
  }

  .mobile-calendar-toggle-button {
    flex: 0 0 auto;
    min-width: 0;
    min-height: 3rem;
    padding-left: 0.9rem;
    padding-right: 0.9rem;
    gap: 0.25rem;
    border-radius: 999px;
    border-color: color-mix(in srgb, var(--bulma-border) 82%, white 18%);
    color: var(--bulma-text);
    background: color-mix(in srgb, var(--bulma-scheme-main) 92%, black 8%);
    box-shadow: 0 10px 24px rgba(10, 10, 10, 0.18);
  }

  .mobile-calendar-toggle-label {
    font-size: 0.92rem;
  }

  .mobile-calendar-toggle-chevron {
    margin-left: 0.1rem;
  }

  .entry-row {
    width: 100%;
    border: 1px solid var(--bulma-border);
    border-radius: 14px;
    background: var(--bulma-scheme-main);
    text-align: left;
    padding: 1rem 1.1rem;
    transition:
      border-color 0.15s ease,
      transform 0.15s ease,
      box-shadow 0.15s ease;
  }

  .entry-row:hover {
    border-color: var(--bulma-link);
    transform: translateY(-1px);
    box-shadow: 0 10px 24px rgba(10, 10, 10, 0.08);
  }

  .entry-row-main {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .entry-row-head {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--bulma-text-weak);
    font-size: 0.95rem;
    flex-wrap: wrap;
  }

  .entry-date {
    font-variant-numeric: tabular-nums;
  }

  .entry-date-outline {
    color: var(--bulma-text);
    font-weight: 600;
    line-height: 1.4;
  }

  .entry-date-preview {
    color: var(--bulma-text-weak);
    font-size: 0.9rem;
    line-height: 1.4;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    min-width: 0;
    flex: 1 1 10rem;
  }

  .entry-outline {
    font-size: 1.05rem;
    font-weight: 600;
    color: var(--bulma-text);
    line-height: 1.45;
  }

  .entry-preview {
    color: var(--bulma-text-weak);
    font-size: 0.92rem;
    line-height: 1.5;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .entry-tags {
    margin-top: 0.85rem;
    margin-bottom: 0;
  }

  .entry-tags-badge :global(.tag) {
    border-radius: 999px;
    background: transparent;
    color: var(--bulma-text-weak);
    font-size: 0.9rem;
    font-weight: 400;
    letter-spacing: 0.01em;
    padding-inline: 0;
    display: inline-flex;
    align-items: center;
    gap: 0.28rem;
  }

  .entry-tags-badge :global(.tag .icon) {
    margin-right: 0;
    color: var(--bulma-text-weak);
  }

  .entry-empty {
    padding: 3rem 1rem;
    text-align: center;
    border: 1px dashed var(--bulma-border);
    border-radius: 14px;
  }

  .mobile-add-button {
    flex: 1 1 auto;
    min-height: 3rem;
    box-shadow: 0 12px 28px rgba(10, 31, 29, 0.16);
  }

  @keyframes calendar-popup-slide-in {
    from {
      transform: translateY(48px);
    }

    to {
      transform: translateY(0);
    }
  }

  @media screen and (max-width: 768px) {
    .entry-list-page {
      padding-top: 0.35rem;
      padding-left: 0.75rem;
      padding-right: 0.75rem;
      padding-bottom: 5.5rem;
    }

    .toolbar {
      margin-bottom: 0.75rem;
      align-items: stretch;
    }

    .toolbar-group,
    .add-button {
      width: 100%;
    }

    .toolbar-controls {
      width: 100%;
    }

    .toolbar-month-select {
      width: 100%;
    }

    .toolbar-month-select :global(.select),
    .toolbar-month-select :global(select) {
      width: 100%;
    }

    .toolbar-mode-row {
      width: 100%;
      flex-wrap: nowrap;
      align-items: center;
    }

    .view-mode-picker {
      flex: 1 1 0;
      min-width: 0;
    }

    .view-mode-label {
      display: none;
    }

    .view-mode-options {
      width: 100%;
      justify-content: stretch;
    }

    .view-mode-option {
      flex: 1 1 0;
      justify-content: center;
      font-size: 0.74rem;
      padding: 0.38rem 0.3rem;
      gap: 0.16rem;
    }

    .view-mode-option :global(.icon) {
      display: none;
    }

    .entry-group-card {
      padding: 0.8rem;
    }

    .entry-group-header {
      align-items: flex-start;
      flex-direction: column;
      margin-bottom: 0.75rem;
    }

    .entry-group-meta {
      width: 100%;
      justify-content: space-between;
    }
  }
</style>
