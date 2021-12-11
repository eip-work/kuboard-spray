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
      resources: 'Resource Package'
    },
  },
  zh: { 
    kuboardspray: 'Kuboard Spray',
    loginRequired: '请先登录。',
    isRequiredField: '为必填字段',
    field: field.zh,
    pkg: pkg.zh,
    obj: {
      resources: '资源包'
    },
  },
}

const i18n = createI18n({
  fallbackLocale: 'zh',
  // globalInjection:true,
  // legacy: false,
  silentFallbackWarn: true,
  locale: language.split("-")[0] || "zh",
  messages,
});

export default i18n