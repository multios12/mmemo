import axios, { AxiosPromise } from "axios";
const http = axios.create({ baseURL: `api/` });
export default http;
