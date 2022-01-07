import KuboardDemo from './KuboardDemo.vue'
import KuboardSprayResources from './kuboard-spray/Resources.vue'
import CopyToClipBoard from './CopyToClipBoard.vue'

let components = [
  { name: 'KuboardDemo', component: KuboardDemo },
  { name: 'KuboardSprayResources', component: KuboardSprayResources },
  { name: 'CopyToClipBoard', component: CopyToClipBoard },
]

export default function (app) {
  for (let component of components) {
    app.component(component.name, component.component)
    console.log('注册组件', component.name)
  }
}
