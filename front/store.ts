import { writable } from "svelte/store";
import type { settingType } from "./models/settingType.js";

export const settingsStore = writable({} as settingType);
