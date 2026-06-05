import type { entryType } from "../models/entryModels.js";
import { apiPath } from "../basePath.js";

export const createEmptyEntry = (): entryType => ({
  Id: undefined,
  Outline: "",
  Date: new Date().toISOString().substring(0, 10),
  Value: "",
  HTML: "",
  Tags: [],
});

export const loadEntryApi = async (
  categoryKey: string,
  entryId: string,
): Promise<entryType> => {
  const response = await fetch(apiPath(`${categoryKey}/${entryId}`));
  if (!response.ok) {
    throw new Error("エントリを読み込めませんでした");
  }
  const nextEntry = (await response.json()) as entryType;
  nextEntry.Tags = nextEntry.Tags ?? [];
  nextEntry.HTML = nextEntry.HTML ?? "";
  nextEntry.Images = nextEntry.Images ?? [];
  return nextEntry;
};
