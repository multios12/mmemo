export type entryType = {
  Id: number | undefined;
  Outline: string;
  Date: string;
  Value: string;
  HTML?: string;
  Images?: {
    Id: string;
    Src: string;
    Alt: string;
    Markdown?: string;
  }[];
  Tags: string[];
  HasDetail?: boolean;
  CreatedAt?: string;
  UpdatedAt?: string;
};

export type listType = {
  WritedMonths: string[];
  Lines: entryType[];
};
