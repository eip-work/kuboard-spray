function openUrlInBlank (url) {
  let a = document.createElement('a')
  document.body.appendChild(a)
  a.setAttribute('href', url)
  a.setAttribute('target', '_blank')
  a.click()
  document.body.removeChild(a)
}


export default {
  install (app) {
    app.config.globalProperties.openUrlInBlank = openUrlInBlank
  }
}