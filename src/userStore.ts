import axios from "./axiosForApi";

/** 認証オブジェクト */
const auth = {
  token: () => localStorage.getItem("token"),
  /** 認証とログイン処理を実行する */
  // tslint:disable-next-line:object-literal-sort-keys
  login: async (username: string, password: string) => {
    const url = `login`;
    const body = { username, password };
    const response = await axios.post(url, body);
    localStorage.removeItem("token");

    const loggedIn = response.data.token !== undefined;

    if (loggedIn) {
      localStorage.setItem("token", response.data.token);
    }

    return loggedIn;
  },
  /** ログアウト処理を実行する */
  logout: () => {
    localStorage.removeItem("token");
  }
};

export default auth;
