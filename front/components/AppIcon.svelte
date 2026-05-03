<script lang="ts">
  import ArrowDownWideNarrow from "lucide-svelte/icons/arrow-down-wide-narrow";
  import ArrowLeft from "lucide-svelte/icons/arrow-left";
  import ArrowUpNarrowWide from "lucide-svelte/icons/arrow-up-narrow-wide";
  import Bold from "lucide-svelte/icons/bold";
  import Book from "lucide-svelte/icons/book";
  import CalendarDays from "lucide-svelte/icons/calendar-days";
  import ChevronDown from "lucide-svelte/icons/chevron-down";
  import ChevronLeft from "lucide-svelte/icons/chevron-left";
  import ChevronRight from "lucide-svelte/icons/chevron-right";
  import ChevronUp from "lucide-svelte/icons/chevron-up";
  import CloudUpload from "lucide-svelte/icons/cloud-upload";
  import Code from "lucide-svelte/icons/code";
  import FileText from "lucide-svelte/icons/file-text";
  import FolderTree from "lucide-svelte/icons/folder-tree";
  import Grip from "lucide-svelte/icons/grip";
  import IdCard from "lucide-svelte/icons/id-card";
  import Image from "lucide-svelte/icons/image";
  import Italic from "lucide-svelte/icons/italic";
  import Layers from "lucide-svelte/icons/layers";
  import Link from "lucide-svelte/icons/link";
  import List from "lucide-svelte/icons/list";
  import ListOrdered from "lucide-svelte/icons/list-ordered";
  import NotebookPen from "lucide-svelte/icons/notebook-pen";
  import Pen from "lucide-svelte/icons/pen";
  import Plus from "lucide-svelte/icons/plus";
  import Quote from "lucide-svelte/icons/quote";
  import RotateCw from "lucide-svelte/icons/rotate-cw";
  import Save from "lucide-svelte/icons/save";
  import SquarePen from "lucide-svelte/icons/square-pen";
  import StickyNote from "lucide-svelte/icons/sticky-note";
  import Strikethrough from "lucide-svelte/icons/strikethrough";
  import Tags from "lucide-svelte/icons/tags";
  import Trash from "lucide-svelte/icons/trash";
  import User from "lucide-svelte/icons/user";
  import X from "lucide-svelte/icons/x";

  interface Props {
    name: string;
    size?: number;
    strokeWidth?: number;
    class?: string;
    ariaHidden?: boolean;
  }

  let {
    name,
    size = 18,
    strokeWidth = 2,
    class: className = "",
    ariaHidden = true,
  }: Props = $props();

  type IconComponent = typeof NotebookPen;

  const iconMap: Record<string, IconComponent> = {
    "address-card": IdCard,
    "arrow-left": ArrowLeft,
    "sort-asc": ArrowUpNarrowWide,
    "sort-desc": ArrowDownWideNarrow,
    bold: Bold,
    book: Book,
    calendar: CalendarDays,
    "calendar-days": CalendarDays,
    "chevron-down": ChevronDown,
    "chevron-left": ChevronLeft,
    "chevron-right": ChevronRight,
    "chevron-up": ChevronUp,
    "cloud-arrow-up": CloudUpload,
    code: Code,
    document: FileText,
    "file-lines": FileText,
    "floppy-disk": Save,
    "folder-tree": FolderTree,
    grip: Grip,
    "grip-lines": Grip,
    group: Layers,
    image: Image,
    italic: Italic,
    "layer-group": Layers,
    link: Link,
    list: List,
    "list-ol": ListOrdered,
    "list-ul": List,
    note: StickyNote,
    "note-sticky": StickyNote,
    pen: Pen,
    "pen-to-square": SquarePen,
    person: User,
    plus: Plus,
    "quote-left": Quote,
    "arrows-rotate": RotateCw,
    "strikethrough": Strikethrough,
    tag: Tags,
    tags: Tags,
    template: Book,
    tree: FolderTree,
    trash: Trash,
    user: User,
    xmark: X,
  };

  const normalizeIconName = (value: string) =>
    value.trim().toLowerCase().replace(/^fa-/, "");
  const ResolvedIcon = $derived(iconMap[normalizeIconName(name)] ?? NotebookPen);
</script>

<ResolvedIcon
  {size}
  {strokeWidth}
  class={className}
  aria-hidden={ariaHidden}
/>
