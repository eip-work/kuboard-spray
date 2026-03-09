<template>
  <div v-if="cluster && cluster.resourcePackage && cluster.inventory">
    <RemoteConfigComponent :resource-package="cluster.resourcePackage" v-model="inventory">
    </RemoteConfigComponent>
  </div>
</template>

<script lang="ts" setup>
import { loadModule } from 'vue3-sfc-loader';
import * as Vue from 'vue';
import { defineAsyncComponent, computed, provide } from "vue";
import axios from "axios";
import yaml from "js-yaml";
import { useI18n } from 'vue-i18n';

import compareVersions from 'compare-versions'

const props = defineProps<{
  cluster: any;
}>();

const inventory = computed({
  get() {
    return props.cluster.inventory;
  },
  set(val) {
    props.cluster.inventory = val;
  }
})

const i18n = useI18n({
  inheritLocale: true,
  useScope: 'local'
})

provide("t", i18n.t);
provide("locale", i18n.locale);
provide("compareVersions", compareVersions);

let options = {
  moduleCache: { vue: Vue },
  async getFile(url: string) {
    const res = await axios.get(url);
    return res.data;
  },
  addStyle(textContent: any) {
    const style = Object.assign(document.createElement('style'), { textContent });
    const ref = document.head.getElementsByTagName('style')[0] || null;
    document.head.insertBefore(style, ref);
  },
  customBlockHandler(block: any, filename: string, options: any) {

    if (block.type !== 'i18n')
      return

    const messages = yaml.load(block.content);
    for (let locale in messages)
      i18n.mergeLocaleMessage(locale, messages[locale]);
  }
}
let url = `/resource-package/${props.cluster.resourcePackage.metadata.version}/content/vue/index.vue`;
const RemoteConfigComponent = defineAsyncComponent(() => {
  return loadModule(url, options as any) as any
}
);

</script>