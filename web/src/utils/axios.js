import axios from 'axios'
import Cookies from 'js-cookie'


var pangeeClusterApi

const baseURL = `./api`

var vueapp

let isLoginAlertShowing = false

const comp = {
  install(app) {
    vueapp = app
    pangeeClusterApi = axios.create({
      baseURL: baseURL,
      timeout: 120000,
      headers: {
        Authorization: 'Bearer ' + Cookies.get('PangeeClusterToken')
      }
    })
    pangeeClusterApi.interceptors.response.use(function (response) {
      return response;
    }, function (error) {
      if (error.response && error.response.status === 401) {
        if (error.response.request && error.response.request.responseURL && error.request.responseURL.indexOf('/api/login') >= 0) {
          return Promise.reject(error)
        }
        if (!isLoginAlertShowing) {
          isLoginAlertShowing = true
          window.VueAppComponent.$alert(window.VueAppComponent.$t('loginRequired'), window.VueAppComponent.$t('loginRequired'), {
            callback: () => {
              isLoginAlertShowing = false
              clearAllCookie()
              // window.VueAppComponent.$router.push('/login')
              console.log("loginrequired")
              location.href = "/#/login"
              location.reload()
            }
          })
        }
      } else {
        return Promise.reject(error);
      }
    });
    app.config.globalProperties.pangeeClusterApi = pangeeClusterApi
  }
}

export default comp

export { baseURL }

export function clearAllCookie() {
  Cookies.remove('PangeeClusterToken', { path: location.pathname })
  Cookies.remove('PangeeClusterLogin', { path: location.pathname })
  comp.install(vueapp)
}

export function setupCookie(token, expires) {
  console.log(token)
  Cookies.set('PangeeClusterToken', token, { path: location.pathname, expires: expires })
  comp.install(vueapp)
}