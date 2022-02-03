import {Box, Promotion, Check, ArrowDownBold, ArrowUpBold } from '@element-plus/icons-vue'


function initIcons(app) {

  app.component('icon-box', Box)
  app.component('icon-promotion', Promotion)
  app.component('icon-check', Check)
  app.component('icon-arrow-down-bold', ArrowDownBold)
  app.component('icon-arrow-up-bold', ArrowUpBold)
}

export default initIcons