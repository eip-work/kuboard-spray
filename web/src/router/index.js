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
    path: '/tail/:ownerType/:ownerName/history/:pid/:file',
    component: () => import('../views/logs/TailFile.vue'),
    props: route => ({ ownerType: route.params.ownerType, ownerName: route.params.ownerName, pid: route.params.pid, file: route.params.file })
  },
  {
    path: '/ssh/:ownerType/:ownerName/:nodeName',
    component: () => import('../views/logs/Ssh.vue'),
    props: route => ({ ownerType: route.params.ownerType, ownerName: route.params.ownerName, nodeName: route.params.nodeName })
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
        path: 'resources/:name',
        name: 'Resource',
        component: () => import('../views/resources/Resource.vue'),
        props: route => ({ name: route.params.name, mode: route.query.mode })
      },
      {
        path: 'resources/:name/on_air',
        name: 'ResourceOnAir',
        component: () => import('../views/resources/ResourceOnAir.vue'),
        props: route => ({ name: route.params.name, mode: route.query.mode })
      },
      {
        path: 'mirrors',
        name: 'Mirrors',
        component: () => import('../views/os_mirror/Mirrors.vue')
      },
      {
        path: 'mirrors/:name',
        name: 'Mirror',
        component: () => import('../views/os_mirror/Mirror.vue'),
        props: route => ({ name: route.params.name, mode: route.query.mode })
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