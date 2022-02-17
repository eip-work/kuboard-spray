module.exports = {
  configureWebpack: (config, isServer) => {
  },

  markdown: {
    lineNumbers: true
  },

  // 站点配置
  lang: 'zh-CN',
  title: 'Kuboard Spray',
  description: '在图形界面种离线安装高可用 Kubernetes 集群',

  head: [
    ['meta', { name: 'theme-color', content: '#007af5' }],
  ],

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
    nav: [
      {
        text: '学习指南',
        link: '/guide/install-k8s.md',
      },
      {
        text: 'Kuboard',
        link: 'https://kuboard.cn',
      },
    ],
    sidebarDepth: 2,
    sidebar: {
      '/': [
        {
          title: 'TLDR',
          path: '/',
        },
        {
          title: '集群维护',
          children: [
            {
              title: '集群安装',
              path: '/guide/install-k8s.md',
            },
            {
              title: '安装选项',
              children: [
                '/guide/options/admission-plugins.md',
              ]
            },
            '/guide/maintain/ha-mode.md',
            '/guide/maintain/add-replace-node.md',
            '/guide/maintain/upgrade.md',
            '/guide/extra/speedup.md'
          ]
        },
        {
          title: '资源包',
          children: [
            '/support/',
            '/guide/resource-package/make.md',
          ]
        },
        {
          title: '网络插件',
          children: [
            '/guide/network-plugin/calico.md',
            '/guide/network-plugin/flannel.md',
          ],
        },
        {
          title: '可选组件',
          children: [
            '/guide/addons/install_addon.md',
            '/guide/addons/nodelocaldns.md',
            '/guide/addons/netchecker.md',
            '/guide/addons/metrics_server.md',
          ],
        },
        '/support/change-log/v1.md',
      ],
    }
  },
}