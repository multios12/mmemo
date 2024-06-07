import { writable } from 'svelte/store';
import { type settingType } from './models/memoModels.js';

export const settingsStore = writable(<settingType>{});