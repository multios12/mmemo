export type settingType = {
    Diary: { Name: string }
    Holidays: HolidayType[]
    Categories: CategoryType[]
}

export type HolidayType = {
    Date: string,
    Name: string,
}

export type TemplateType = {
    Name: string,
    Value: string,
    Tags?: string[],
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
    // 同一日付の複数登録可否
    AllowMultipleEntriesPerDate: boolean,
    // テンプレート
    Templates?: TemplateType[],
    Fields: {
        Name?: string,
        Date?: string,
        Tags?: string,
        Value?: string,
        Outline?: string,
        OutlineIcon?: string,
    }
}
