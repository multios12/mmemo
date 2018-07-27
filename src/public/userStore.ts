import axios from 'axios';

/** 認証オブジェクト */
var auth = {
    token: function() {
        return localStorage.getItem('token');
    },
    /** 認証とログイン処理を実行する */
    login: async function (username: string, password: string) {
        var url = `api/login`;
        var body = { username: username, password: password }
        var response = await axios.post(url, body);
        localStorage.removeItem('token');

        var loggedIn = response.data.token != undefined;

        if (loggedIn) {
            localStorage.setItem('token', response.data.token);
            console.log(`token:${response.data.token}`)
        }

        return loggedIn;
    },
    /** ログアウト処理を実行する */
    logout: function () {
        localStorage.removeItem('token');
    },
};

export default auth;