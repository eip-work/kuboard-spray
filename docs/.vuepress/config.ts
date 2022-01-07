import { defineUserConfig } from 'vuepress'
import type { DefaultThemeOptions } from 'vuepress'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

const { path } = require('@vuepress/utils')

export default defineUserConfig<DefaultThemeOptions>({
  clientAppEnhanceFiles: path.resolve(
    __dirname,
    './clientAppEnhance.js'
  ),

  bundler: '@vuepress/bundler-vite',
  // Vite 打包工具的配置项
  bundlerConfig: {
    plugins: [
      // ...
      AutoImport({
        resolvers: [ElementPlusResolver()],
      }),
      Components({
        resolvers: [ElementPlusResolver()],
      }),
    ],
  },

  // 站点配置
  lang: 'zh-CN',
  title: 'Kuboard Spray',
  description: '在图形界面种离线安装高可用 Kubernetes 集群',

  // 主题和它的配置
  theme: '@vuepress/theme-default',
  plugins: [
    [
      '@vuepress/plugin-google-analytics',
      {
        id: 'G-ECTJF3X9EJ',
      },
    ],
  ],
  themeConfig: {
    logo: 'https://kuboard.cn/favicon.png',
    repo: 'eip-work/kuboard-spray',
    docsDir: 'docs',
    navbar: [
      {
        text: '学习指南',
        link: '/guide/install-k8s.md',
      },
      {
        text: '技术支持',
        link: '/support/',
      },
    ],
    sidebar: {
      '/guide/': [
        {
          text: '集群安装',
          children: [
            '/guide/install-k8s.md',
          ]
        },
        {
          text: '集群维护',
          children: [
            '/guide/maintain/ha-mode.md',
            '/guide/maintain/add-replace-node.md',
            '/guide/maintain/upgrade.md'
          ]
        },
        {
          text: '资源包',
          children: [
            '/guide/resource-package/repo.md',
            '/guide/resource-package/make.md',
          ]
        },
        {
          text: '网络插件',
          children: [
            '/guide/network-plugin/calico.md',
            '/guide/network-plugin/flannel.md',
          ]
        }
      ],
      '/support/': [
        {
          text: '支持',
          link: '/support/',
        },
        {
          text: '变更记录',
          link: '/support/change-log/v1.md'
        }
      ]
    }
  },
})