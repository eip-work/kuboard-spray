let components = [
  { name: 'KuboardDemo', component: () => import('./KuboardDemo.vue') },
  { name: 'KuboardSprayResources', component: () => import('./kuboard-spray/Resources.vue') },
  { name: 'CopyToClipBoard', component: () => import('./CopyToClipBoard.vue') },
  { name: 'InstallAddon', component: () => import('./InstallAddon.vue') },
]

export default function (app) {
  for (let component of components) {
    app.component(component.name, component.component)
    console.log('注册组件', component.name)
  }
}
