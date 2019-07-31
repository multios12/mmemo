import fs from "fs";
import path from "path";
class Settings {
  /** データ保存パス */
  private datasPath: string;

  /** データ保存パスを返す */
  public DatasPath = () => {
    if (!this.datasPath) {
      this.datasPath = process.env.DATAS_PATH
        ? process.env.DATAS_PATH
        : path.join(__dirname, "data");
      if (!fs.existsSync(this.datasPath)) {
        fs.mkdirSync(this.datasPath);
      }
    }
    return this.datasPath;
  };

  /** 指定されたIDのデータ保存パスを返す */
  public CreateIdPath = (id: string) => path.join(this.DatasPath(), id);

  /** 指定されたID・ファイル名のデータ保存パスを返す */
  public CreateImagePath = (id: string, fileName: string) =>
    path.join(this.DatasPath(), id, fileName);

  /** memos.dbファイルのパスを返す */
  public MemosDBPath = () => path.join(this.DatasPath(), "memos.db");
}

export = new Settings();
