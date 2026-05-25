import type { SvelteComponentTyped } from "svelte";

export type ToolbarIcon = typeof SvelteComponentTyped<any, any, any>;

export interface ToolbarButtonContext {
  value: string;
  textarea: HTMLTextAreaElement | null;
  selectionStart: number;
  selectionEnd: number;
  updateValue: (value: string) => void;
  focus: () => void;
  openImageUploadModal: () => void;
  syncToolbarState: () => void;
}

export interface ToolbarButton {
  key: string;
  ariaLabel: string;
  icon?: ToolbarIcon;
  label?: string;
  title?: string;
  active?: boolean;
  disabled?: boolean;
  onClick?: (context: ToolbarButtonContext) => void;
}
