const dayMs = 24 * 60 * 60 * 1000;

const toUtcDay = (value: Date) =>
  Date.UTC(value.getFullYear(), value.getMonth(), value.getDate());

export const addDay = (value: Date) =>
  new Date(value.getFullYear(), value.getMonth(), value.getDate() + 1);

export const diffDays = (left: Date, right: Date) =>
  Math.round((toUtcDay(left) - toUtcDay(right)) / dayMs);

export const monthStart = (value: Date) =>
  new Date(value.getFullYear(), value.getMonth(), 1);

export const monthEnd = (value: Date) =>
  new Date(value.getFullYear(), value.getMonth() + 1, 0);

export const weekStart = (value: Date) =>
  new Date(value.getFullYear(), value.getMonth(), value.getDate() - value.getDay());

const weekdayFormatter = new Intl.DateTimeFormat("ja-JP", { weekday: "short" });
const dayFormatter = new Intl.DateTimeFormat("ja-JP", { day: "numeric" });
const monthLabelFormatter = new Intl.DateTimeFormat("ja-JP", { month: "numeric" });

export const format = (value: Date, pattern: string) => {
  if (pattern === "YYYY-MM-DD") {
    const year = String(value.getFullYear());
    const month = String(value.getMonth() + 1).padStart(2, "0");
    const day = String(value.getDate()).padStart(2, "0");
    return `${year}-${month}-${day}`;
  }

  if (pattern === "D") {
    return dayFormatter.format(value);
  }

  if (pattern === "ddd") {
    return weekdayFormatter.format(value);
  }

  if (pattern === "M月") {
    return `${monthLabelFormatter.format(value)}月`;
  }

  return value.toISOString();
};
