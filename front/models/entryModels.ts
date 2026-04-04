export type entryType = {
  Id: number | undefined;
  Outline: string;
  Date: string;
  Value: string;
  Tags: string[];
  HasDetail?: boolean;
};

export type listType = {
  WritedMonths: string[];
  Lines: entryType[];
};
