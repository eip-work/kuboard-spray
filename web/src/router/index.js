import { createRouter, createWebHashHistory } from 'vue-router';
import Cookies from 'js-cookie'

const constantRouterMap = [
  {
    path: '/',
    component: () => import('../components/layout/Layout.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('../views/home/Home.vue')
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/Login.vue')
  },
  {
    path: '/clusters',
    component: () => import('../components/layout/Layout.vue'),
    children: [
      {
        path: '',
        name: 'Clusters',
        component: () => import('../views/clusters/Clusters.vue')
      },
      {
        path: ':name',
        name: 'Cluster',
        component: () => import('../views/clusters/Cluster.vue'),
        props: route => ({ name: route.params.name, mode: route.query.mode })
      }
    ],
  },
  {
    path: '/clusters/:name/history/:pid/tail/:file',
    component: () => import('../views/logs/TailFile.vue'),
    props: route => ({ cluster: route.params.name, pid: route.params.pid, file: route.params.file })
  },
  {
    path: '/settings',
    component: () => import('../components/layout/Layout.vue'),
    children: [
      {
        path: 'resources',
        name: 'Resources',
        component: () => import('../views/resources/Resources.vue')
      },
      {
        path: 'kuboard',
        name: 'Kuboard',
        component: () => import('../views/kuboard/Kuboard.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes: constantRouterMap
})

router.beforeEach((to, from, next) => {
  let token = Cookies.get('KuboardToken')
  if (to.name == 'Login') {
    next()
  } else if (token == null || token == undefined || token == '') {
    next({ name: 'Login' })
  } else {
    next()
  }
})

export default router