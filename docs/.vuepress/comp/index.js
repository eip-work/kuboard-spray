let components = [
  { name: 'PangeeClusterResources', component: () => import('./pangee-cluster/Resources.vue') },
  { name: 'CopyToClipBoard', component: () => import('./CopyToClipBoard.vue') },
  { name: 'InstallAddon', component: () => import('./InstallAddon.vue') },
]

export default function (app) {
  for (let component of components) {
    app.component(component.name, component.component)
    console.log('注册组件', component.name)
  }
}
