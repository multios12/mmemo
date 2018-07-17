import axios from 'axios';

/** 認証オブジェクト */
var auth = {
    /** 認証が完了し、ログイン済みであればtrueを返す */
    loggedIn: false,
    /** 認証とログイン処理を実行する */
    login: async function (username: string, password: string) {
        var url = `api/login`;
        var body = { username: username, password: password }
        console.log("auth");

        var response = await axios.post(url, body);

        console.log("auth=" + response.status);

        this.loggedIn = response.status == 200
        return this.loggedIn;
    },
    /** ログアウト処理を実行する */
    logout: function () {
        this.loggedIn = false
    }
};

export default auth;