import axios, { AxiosPromise } from 'axios';

// const http = axios.create({ baseURL: `${process.env.API_URL}/api/v1/` });
const http = axios.create();
http.interceptors.request.use((config) => {
    var token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
        return config
    }
    return config
}, function (error) {
    return Promise.reject(error)
})
export default http;
