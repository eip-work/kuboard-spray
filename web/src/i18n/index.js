import { createI18n } from 'vue-i18n' 

const language = (
  (navigator.language ? navigator.language : navigator.userLanguage) || "zh"
).toLowerCase();

const i18n = createI18n({
  fallbackLocale: 'zh',
  // globalInjection:true,
  // legacy: false,
  locale: language.split("-")[0] || "zh",
  messages,
});

const messages = {
  en: {
    kuboardspray: 'Kuboard Spray',
    loginRequired: 'Login Required.'
  },
  zh: { 
    kuboardspray: 'Kuboard Spray',
    loginRequired: '请先登录。'
  },
}

export default i18n