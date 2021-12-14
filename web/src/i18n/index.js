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
      addon: 'Addon: {name}',
    },
    node: {
      kube_control_plane: 'control plane',
      kube_node: 'worker node',
      etcd: 'etcd node',
      k8s_cluster: 'K8s Cluster',
    },
    msg: {
      ok: 'OK',
      save: 'Save',
      close: 'Close',
      cancel: 'Cancel',
      upload: 'Upload',
      edit: 'Edit',
      save_succeeded: 'Succeeded in saving.',
      save_failed: 'Failed in saving. {msg}',
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
      addon: '可选组件：{name}',
    },
    node: {
      kube_control_plane: '控制节点',
      kube_node: '工作节点',
      etcd: 'ETCD节点',
      k8s_cluster: 'K8s 集群',
    },
    msg: {
      ok: '确 定',
      save: '保 存',
      close: '关 闭',
      cancel: '取 消',
      upload: '上 传',
      edit: '编 辑',
      save_succeeded: '保存成功',
      save_failed: '保存失败： {msg}',
    }
  },
}

const i18n = createI18n({
  fallbackLocale: 'zh',
  globalInjection:true,
  legacy: true,
  silentFallbackWarn: true,
  silentTranslationWarn: true,
  locale: language.split("-")[0] || "zh",
  messages,
});

export default i18n