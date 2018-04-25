import fs from 'fs';
import path from 'path';
class settings {
    /** データ保存パス */
    private _datasPath: string;

    /** データ保存パスを返す */
    public DatasPath = () => {
        if (!this._datasPath) {
            this._datasPath = process.env.DATAS_PATH ? process.env.DATAS_PATH : path.join(__dirname, 'data');
            if (!fs.existsSync(this._datasPath)) {
                fs.mkdirSync(this._datasPath);
            }
        }
        return this._datasPath;
    };

    /** 指定されたIDのデータ保存パスを返す */
    public CreateIdPath = (id: string) => path.join(this.DatasPath(), id);

    /** 指定されたID・ファイル名のデータ保存パスを返す */
    public CreateImagePath = (id: string, fileName: string) => path.join(this.DatasPath(), id, fileName);

    /** memos.dbファイルのパスを返す */
    public MemosDBPath = () => path.join(this.DatasPath(), 'memos.db');
}

export = new settings;
