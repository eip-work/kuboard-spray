import Comp from './comp/index'
import { defineClientAppEnhance } from '@vuepress/client'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import { VueClipboard } from '@soerenmartius/vue3-clipboard'

export default defineClientAppEnhance(({ app, router, siteData }) => {
  app.use(Comp)
  app.use(ElementPlus, { size: 'small', zIndex: 3000 })
  app.use(VueClipboard)
})