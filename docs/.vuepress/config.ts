import { defineUserConfig } from 'vuepress'
import type { DefaultThemeOptions } from 'vuepress'

export default defineUserConfig<DefaultThemeOptions>({
  // 站点配置
  lang: 'zh-CN',
  title: 'Kuboard Spray',
  description: '在图形界面种离线安装高可用 Kubernetes 集群',

  // 主题和它的配置
  theme: '@vuepress/theme-default',
  themeConfig: {
    logo: 'https://kuboard.cn/favicon.png',
    sidebar: {
      '/guide/': [
        {
          text: '指南',
          children: [
            '/guide/install-k8s.md',
            '/guide/ha-mode.md',
          ]
        }
      ]
    }
  },
})