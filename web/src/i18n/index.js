import { createI18n } from 'vue-i18n' 
import field from './field.js'
import pkg from './pkg.js';

const language = (
  (navigator.language ? navigator.language : navigator.userLanguage) || "zh"
).toLowerCase();

const messages = {
  en: {
    kuboardspray: 'Kuboard Spray',
    loginRequired: 'Login Required.',
    isRequiredField: 'is required.',
    field: field.en,
    pkg: pkg.en,
    obj: {
      resources: 'Resource Package',
      bastion: 'Bastion',
      localhost: 'KuboardSpray',
    },
    node: {
      kube_control_plane: 'control plane',
      kube_node: 'worker node',
      etcd: 'etcd node',
      k8s_cluster: 'K8s Cluster',
    }
  },
  zh: { 
    kuboardspray: 'Kuboard Spray',
    loginRequired: '请先登录。',
    isRequiredField: '为必填字段',
    field: field.zh,
    pkg: pkg.zh,
    obj: {
      resources: '资源包',
      bastion: '堡垒机',
      localhost: 'KuboardSpray',
    },
    node: {
      kube_control_plane: '控制节点',
      kube_node: '工作节点',
      etcd: 'ETCD节点',
      k8s_cluster: 'K8s 集群',
    }
  },
}

const i18n = createI18n({
  fallbackLocale: 'zh',
  globalInjection:true,
  legacy: true,
  silentFallbackWarn: true,
  locale: language.split("-")[0] || "zh",
  messages,
});

export default i18n