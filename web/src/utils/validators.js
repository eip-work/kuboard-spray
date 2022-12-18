const comp = {
  install(app) {
    app.config.globalProperties.$validators = {}
    app.config.globalProperties.$validators.dnsName = function (rule, value, callback) {
      if (!value) {
        return callback()
      }
      if (value.indexOf('_') >= 0) {
        callback(new Error('不能包含下划线'))
      }
      if (parseInt(value[0]) >= 0 && parseInt(value[0]) <= 9) {
        callback(new Error('不能以数字开头，必须以字母开头'))
      }
      var reg = /^[0-9a-zA-Z\-.]+$/
      if (reg.test(value)) {
        return callback()
      } else {
        return callback(window.VueAppComponent.$t('validator.dnsName'))
      }
    }
    app.config.globalProperties.$validators.ipv4 = function (rule, value, callback) {
      if (!value) {
        return callback()
      }
      let reg = /^(?:(?:1[0-9][0-9]\.)|(?:2[0-4][0-9]\.)|(?:25[0-5]\.)|(?:[1-9][0-9]\.)|(?:[0-9]\.)){3}(?:(?:1[0-9][0-9])|(?:2[0-4][0-9])|(?:25[0-5])|(?:[1-9][0-9])|(?:[0-9]))$/
      if (reg.test(value)) {
        return callback()
      } else {
        return callback('非法的 ipv4 地址')
      }
    }
  }
}

export default comp