import { createApp } from "vue";
import router from "./router/index.js";
import { initKuboardMfe } from "./micro-front-end.js";

import ElementPlus from "element-plus";
// import "element-plus/theme-chalk/src/index.scss";
import "./styles/element-variables.scss";
import "./styles/element-ui.css";
import icons from "./styles/el-icons.js";

import zhCn from "element-plus/es/locale/lang/zh-cn";

import store from "./store/index.js";
import i18n from "./i18n/index.js";
import initAxios from "./utils/axios.js";
import openUrlInBlank from "./utils/open-in-blank.js";
import validators from "./utils/validators.js";

import components from "./components/index.js";

import { VueClipboard } from "@soerenmartius/vue3-clipboard";

import App from "./App.vue";

import axios from "axios";

import i18nMixin from "./mixins/i18n-mixin.js";

// 自有组件库样式
// import "@paas-front/paas-css/css/index.css";
import "./styles/paas-index.css";

axios.get("./version.json?nocache=" + new Date().getTime()).then(resp => {
  window.PangeeCluster = { version: resp.data };
  window.PangeeCluster.version.trimed = window.PangeeCluster.version.version.slice(
    0,
    window.PangeeCluster.version.version.length - 6
  );
  window.PangeeCluster.version.arch = window.PangeeCluster.version.version.slice(window.PangeeCluster.version.version.length - 5);
  const app = createApp(App);
  app.use(ElementPlus, { size: "small", locale: zhCn, zIndex: 3000 });
  app.use(store);
  app.use(router);
  app.use(i18n);
  icons(app);
  app.use(components);
  app.use(initAxios);
  app.use(validators);
  app.use(openUrlInBlank);
  app.use(VueClipboard);
  app.mixin(i18nMixin);
  // app.config.unwrapInjectedRef = true;
  initKuboardMfe(app);
  app.mount("#app");

});
