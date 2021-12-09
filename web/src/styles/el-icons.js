import {Box, Promotion, Check} from '@element-plus/icons'


function initIcons(app) {

  app.component('icon-box', Box)
  app.component('icon-promotion', Promotion)
  app.component('icon-check', Check)

}

export default initIcons