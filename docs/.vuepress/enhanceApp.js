import Comp from './comp/index'
import VueClipboard from 'vue-clipboard2'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';

export default ({
  Vue, // VuePress 正在使用的 Vue 构造函数
  options, // 附加到根实例的一些选项
  router, // 当前应用的路由实例
  siteData, // 站点元数据
  isServer // 当前应用配置是处于 服务端渲染 或 客户端
}) => {
  Vue.use(ElementUI, { size: 'small', zIndex: 3000 })
  Vue.use(Comp)
  Vue.use(VueClipboard)
}