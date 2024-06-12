export type settingType = {
    Diary: { Name: string }
    Categories: CategoryType[]
}
export type CategoryType = {
    Key: string,
    // 種類名
    Name: string,
    // 日付
    Date: Date,
    // 日付の使用・表示
    UseDate: boolean,
    // タグ表示 
    UseTag: boolean,
    // テンプレート
    Template: string,
    Fields: {}
}