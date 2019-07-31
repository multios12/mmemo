import axios, { AxiosPromise } from "axios";
const http = axios.create({ baseURL: `http://localhost:3000/api/` });
http.interceptors.request.use(
  config => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
      return config;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);
export default http;
