import { writable } from 'svelte/store';
import { type settingType } from './models/settingType.js';
import { getMemosSettingApi } from './models/apiUrl.js';

export const initSetting = async (routes: any) => {
  const r = await getMemosSettingApi(routes);
  let settings = <settingType>await r.json();
  settingsStore.update((s) => settings);
};

export const settingsStore = writable(<settingType>{});