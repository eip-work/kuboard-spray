import axios from 'axios'
import Cookies from 'js-cookie'

let kuboardSprayId = 'default'
let splitedPath = location.pathname.split('/')
if (splitedPath[0] === 'kuboardspray' && splitedPath[1] !== undefined) {
  kuboardSprayId = splitedPath[1]
}

var kuboardSprayApi

const baseURL = `/kuboardspray/${kuboardSprayId}/api`

var vueapp

const comp = {
  install(app) {
    vueapp = app
    kuboardSprayApi = axios.create({
      baseURL: baseURL,
      timeout: 30000,
      headers: {
        Authorization: 'Bearer ' + Cookies.get('KuboardToken')
      }
    })
    kuboardSprayApi.interceptors.response.use(function (response) {
      return response;
    }, function (error) {
      if (error.response && error.response.status === 401) {
        if (error.response.data && error.response.data.kind === 'KuboardErrorResponse') {
          alert(window.VueAppComponent.$t('loginRequired'))
          clearAllCookie()
          location.href = '/login'
        }
      } else if (error.response && error.response.status === 424 && error.response.data.reason === 'AuditDbNotReady') {
        location.href = `/init-status-check/`
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