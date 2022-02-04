import * as icons from '@element-plus/icons-vue'


function initIcons(app) {

  for (let key in icons) {
    app.component('ElIcon' + key, icons[key])
    console.log('ElIcon' + key)
  }
  // app.component('icon-box', Box)
  // app.component('icon-promotion', Promotion)
  // app.component('icon-check', Check)
  // app.component('icon-arrow-down-bold', ArrowDownBold)
  // app.component('icon-arrow-up-bold', ArrowUpBold)
}

export default initIcons