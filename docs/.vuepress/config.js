module.exports = {
  configureWebpack: (config, isServer) => {
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
        text: '技术支持',
        link: '/support/',
      },
    ],
    sidebar: {
      '/guide/': [
        {
          title: '集群安装',
          children: [
            '/guide/install-k8s.md',
          ]
        },
        {
          title: '集群维护',
          children: [
            '/guide/maintain/ha-mode.md',
            '/guide/maintain/add-replace-node.md',
            '/guide/maintain/upgrade.md'
          ]
        },
        {
          title: '资源包',
          children: [
            '/guide/resource-package/repo.md',
            '/guide/resource-package/make.md',
          ]
        },
        {
          title: '网络插件',
          children: [
            '/guide/network-plugin/calico.md',
            '/guide/network-plugin/flannel.md',
          ]
        }
      ],
      '/support/': [
        '',
        'change-log/v1.md',
      ]
    }
  },
}