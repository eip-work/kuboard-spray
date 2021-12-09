import axios from 'axios'
import Cookies from 'js-cookie'

const kuboardSprayApi = axios.create({
  baseURL: `/kuboard-api/`,
  timeout: 30000,
  // httpAgent: new http.Agent({ keepAlive: true })
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
  // Cookies.remove('KuboardToken', { path: '/kuboard-api/' })
  // Cookies.remove('KuboardToken', { path: '/k8s-api/' })
  // Cookies.remove('KuboardToken', { path: '/k8s-ws/' })
  Cookies.remove('KuboardToken', { path: '/' })
  Cookies.remove('KuboardLogin', { path: '/' })
}
