import axios from 'axios';
import AuthService from "./auth.service"
import { router  }from '../router'
import store from '../store'

function createAxiosResponseInterceptor() {
  const interceptor = axios.interceptors.response.use(
      response => response,
      error => {
          // Reject promise if usual error
          console.log("aaa")
          if (error.response.status !== 401) {
              return Promise.reject(error);
          }

          axios.interceptors.response.eject(interceptor);
          if (!error.config.__isRetryRequest) {
            return new Promise((resolve, reject) => {
              AuthService.refreshToken().then(response => {
                  // saveToken();
                  console.log("BBBB")
  
                  let user = JSON.parse(localStorage.getItem('user'));
                  user.access_token = response.data.access_token
                  user.refresh_token = response.data.refresh_token
                  localStorage.setItem('user', JSON.stringify(user));
                  error.config.__isRetryRequest = true
                  error.response.config.headers['Authorization'] = 'Bearer ' + response.data.access_token;
                  let req = axios(error.config);
                  console.log('req ', req)
                  resolve(req)
              }, error => {
                  console.log('CCCC')
                
                  store.dispatch('auth/logout');
                  router.push('/login');
                  reject(error);
              })
              createAxiosResponseInterceptor();
            })
          }
      }
  );
}

createAxiosResponseInterceptor()

// axios.interceptors.response.use(response => response, error => {
//     const status = error.response ? error.response.status : null
//     console.log('1 gh')
//     if (status === 401) {
//       console.log('2 gh')
//       // will loop if refreshToken returns 401
//       return AuthService.refreshToken().then(_ => {
//         // error.config.headers['Authorization'] = 'Bearer ' + store.state.auth.token;
//         // error.config.baseURL = undefined;
//         console.log("config ", error.config)
//         return axios.request(error.config);
//       })
//       .catch(err => err);
//     }
  
//     console.log('3 gh')
  
//     return Promise.reject(error);
//   });
  

class ErrorHandler {
    checkAuthError(error) {
        console.log('checkAuthError ' + error.toString())
        let res = ''
        if (error.response.status == 401) {
            AuthService.refreshToken().then(
                response => {
                    if (response.data.access_token) {
                        console.log('response.data.access_token')
                        let user = JSON.parse(localStorage.getItem('user'));
                        console.log(user)
                        user.access_token = response.data.access_token
                        user.refresh_token = response.data.refresh_token
                        console.log(user)
                        // localStorage.setItem('user', user);
                        localStorage.setItem('user', JSON.stringify(user));
                        console.log('Token was refreshed')
                    }
                  },
                  error => {
                    res =
                      (error.response && error.response.data && error.response.data.message) ||
                      error.message ||
                      error.toString()
                  }
                );
        } else {
            res = (error.response && error.response.data && error.response.data.message) || 
            error.message ||
            error.toString()
        }
        return res
    }
}


export default new ErrorHandler();

