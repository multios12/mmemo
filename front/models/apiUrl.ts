export const getMonthApi = (month: string, route: any) => {
    let location = window.location.href.replace(route.result.path.original, "");
    let url = `${location}/api/diary/${month}`;
    return fetch(url)
}
export const getDetailApi = (year: string, month: string, day: string, route: any) => {
    let location = window.location.href.replace(route.result.path.original, "");
    let url = `${location}/api/diary/${year}/${month}/${day}`;
    return fetch(url)
}
export const postDetailApi = (year: string, month: string, day: string, body: string, route: any) => {
    let location = window.location.href.replace(route.result.path.original, "");
    let url = `${location}/api/diary/${year}/${month}/${day}`;
    let options = { method: "post", body };
    return fetch(url, options)
}
export const deleteDetailApi = (year: string, month: string, day: string, route: any) => {
    let location = window.location.href.replace(route.result.path.original, "");
    let url = `${location}/api/diary/${year}/${month}/${day}`;
    let options = { method: "delete" };
    return fetch(url, options)
}

export const getMemosSettingApi = (route: any) => {
    let location = window.location.href.replace(route.result.path.original, "");
    let url = `${location}/api/memos`;
    return fetch(url);
}
export const getMemosApi = (category: string, route: any) => {
    let location = window.location.href.replace(route.result.path.original, "");
    category = category.replace("/", "");
    let url = `${location}/api/memos/${category}`;
    return fetch(url);
}
export const getMemosIdApi = (category: string, id: string, route: any) => {
    let location = window.location.href.replace(route.result.path.original, "");
    category = category.replace("/", "");
    let url = `${location}/api/memos/${category}/${id}`;
    return fetch(url);
}
export const putMemosApi = (category: string, memo: string, route: any) => {
    let location = window.location.href.replace(route.result.path.original, "");
    let url = `${location}/api/memos/${category}`;
    let options = { method: "put", body: memo };
    return fetch(url, options);
}
export const postMemosApi = (category: string, id: string, memo: string, route: any) => {
    let location = window.location.href.replace(route.result.path.original, "");
    let url = `${location}/api/memos/${category}/${id}`;
    let options = { method: "post", body: memo };
    return fetch(url, options);
}
export const deleteMemosIdApi = (category: string, id: string, route: any) => {
    let location = window.location.href.replace(route.result.path.original, "");
    let url = `${location}/api/memos/${category}/${id}`;
    let options = { method: "delete" };
    return fetch(url, options);
}
