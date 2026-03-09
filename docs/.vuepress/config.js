module.exports = {
  configureWebpack: (config, isServer) => {
  },

  markdown: {
    lineNumbers: true
  },

  // 站点配置
  lang: 'zh-CN',
  title: 'Pangee Cluster',
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
    // logo: '/favicon.png',
    repo: 'opencmit/pangee-cluster',
    docsDir: 'docs',
    nav: [
      {
        text: '学习指南',
        link: '/guide/install-k8s.md',
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
            '/guide/maintain/add-replace-node.md',
            '/guide/extra/speedup.md'
          ]
        },
        '/pangee-cluster/reset-password.md',
        '/setup-dev/dev.md',
        '/support/change-log/v2.md',
      ],
    }
  },
}