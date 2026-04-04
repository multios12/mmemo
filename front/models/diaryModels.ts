export type listType = { WritedMonths: string[]; Lines: lineType[] };
export type lineType = {
  Id: number | undefined;
  Day: string;
  Outline: string;
  Tags: string[];
  IsDetail: boolean;
  HCount: number;
};
