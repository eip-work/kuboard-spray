import * as icons from '@element-plus/icons-vue'


function initIcons(app) {

  for (let key in icons) {
    app.component('ElIcon' + key, icons[key])
    // console.log('ElIcon' + key)
  }
}

export default initIcons