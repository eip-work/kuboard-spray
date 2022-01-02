import axios from 'axios'
import Cookies from 'js-cookie'

// let kuboardSprayId = 'default'
// let splitedPath = location.pathname.split('/')
// if (splitedPath[0] === 'kuboardspray' && splitedPath[1] !== undefined) {
//   kuboardSprayId = splitedPath[1]
// }

var kuboardSprayApi

const baseURL = `./api`

var vueapp

const comp = {
  install(app) {
    vueapp = app
    kuboardSprayApi = axios.create({
      baseURL: baseURL,
      timeout: 120000,
      headers: {
        Authorization: 'Bearer ' + Cookies.get('KuboardToken')
      }
    })
    kuboardSprayApi.interceptors.response.use(function (response) {
      return response;
    }, function (error) {
      console.log(error.response.status, error.response.data)
      if (error.response && error.response.status === 401) {
        window.VueAppComponent.$alert(window.VueAppComponent.$t('loginRequired'), window.VueAppComponent.$t('loginRequired'), {
          callback: () => {
            clearAllCookie()
            window.VueAppComponent.$router.push('/login')
          }
        })
      }
      return Promise.reject(error);
    });
    app.config.globalProperties.kuboardSprayApi = kuboardSprayApi
  }
}

export default comp

export {baseURL}

export function clearAllCookie() {
  Cookies.remove('KuboardToken', { path: location.pathname })
  Cookies.remove('KuboardLogin', { path: location.pathname })
  comp.install(vueapp)
}

export function setupCookie(token, expires) {
  console.log(token)
  Cookies.set('KuboardToken', token, { path: location.pathname, expires: expires } )
  comp.install(vueapp)
}