import { writable } from 'svelte/store';
import { type settingType } from './models/memoModels';

export const settingsStore = writable(<settingType>{});