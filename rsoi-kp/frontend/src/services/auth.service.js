import axios from 'axios';

// const SESSION_API_URL = 'http://localhost:8081/session/api/v1/session/';
const SESSION_API_URL = '/session/api/v1/session/';

class AuthService {
  login(user) {
    return axios
      .post(SESSION_API_URL + 'signin', {
        login: user.login,
        password: user.password
      })
      .then(response => {
        if (response.data.access_token) {
          localStorage.setItem('user', JSON.stringify(response.data));
        }

        return response.data;
      });
  }

  logout() {
    localStorage.removeItem('user');
  }

  register(user) {
    return axios.post(SESSION_API_URL + 'signup', {
      login: user.login,
      password: user.password
    });
  }

  refreshToken() {
    console.log('refreshToken')
    let user = JSON.parse(localStorage.getItem('user'));
    return axios
      .post(SESSION_API_URL + 'token/refresh', {
        refresh_token: user.refresh_token,
      });
  }

}

export default new AuthService();
