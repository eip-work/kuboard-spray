import axios from 'axios'
import Cookies from 'js-cookie'

let kuboardSprayId = 'default'
let splitedPath = location.pathname.split('/')
if (splitedPath[0] === 'kuboardspray' && splitedPath[1] !== undefined) {
  kuboardSprayId = splitedPath[1]
}

const kuboardSprayApi = axios.create({
  baseURL: `/kuboardspray/${kuboardSprayId}/api`,
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

export {kuboardSprayApi}

export function clearAllCookie() {
  Cookies.remove('KuboardToken', { path: location.pathname })
  Cookies.remove('KuboardLogin', { path: location.pathname })
}

export function setupCookie(token, expires) {
  Cookies.set('KuboardToken', token, { path: location.pathname, expires: expires } )
}