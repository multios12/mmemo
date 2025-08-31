const apiUrl = import.meta.env.VITE_API_URL;

export const getMonthApi = (month: string) => {
  let url = `${apiUrl}diary/${month}`;
  return fetch(url)
}
export const getDetailApi = (year: string, month: string, day: string) => {
  let url = `${apiUrl}diary/${year}/${month}/${day}`;
  return fetch(url)
}
export const postDetailApi = (year: string, month: string, day: string, body: string) => {
  let url = `${apiUrl}diary/${year}/${month}/${day}`;
  let options = { method: "post", body };
  return fetch(url, options)
}
export const deleteDetailApi = (year: string, month: string, day: string) => {
  let url = `${apiUrl}diary/${year}/${month}/${day}`;
  let options = { method: "delete" };
  return fetch(url, options)
}

export const getMemosSettingApi = () => {
  let url = `${apiUrl}memos`;
  return fetch(url);
}
export const getMemosApi = (category: string) => {
  category = category.replace("/", "");
  let url = `${apiUrl}memos/${category}`;
  return fetch(url);
}
export const getMemosIdApi = (category: string, id: string) => {
  category = category.replace("/", "");
  let url = `${apiUrl}memos/${category}/${id}`;
  return fetch(url);
}
export const putMemosApi = (category: string, memo: string) => {
  let url = `${apiUrl}memos/${category}`;
  let options = { method: "put", body: memo };
  return fetch(url, options);
}
export const postMemosApi = (category: string, id: string, memo: string) => {
  let url = `${apiUrl}memos/${category}/${id}`;
  let options = { method: "post", body: memo };
  return fetch(url, options);
}
export const deleteMemosIdApi = (category: string, id: string) => {
  let url = `${apiUrl}memos/${category}/${id}`;
  let options = { method: "delete" };
  return fetch(url, options);
}
