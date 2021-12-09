<template>
  <div>
    <router-view v-if="inited" ref="view"></router-view>
    <IframeAdapter command="injectEnv" :handler="injectEnv"></IframeAdapter>
  </div>
</template>

<script>
import IframeAdapter from './IframeAdapter.vue'

export default {
  props: {
  },
  data() {
    return {
      inited: false
    }
  },
  computed: {
  },
  components: { IframeAdapter },
  mounted () {
    // this.injectEnv({ cluster: 'k8s-19' })
  },
  methods: {
    injectEnv (tokens) {
      console.log('tokens')
      window.KuboardEnv = tokens
      this.inited = true
    }
  }
}
</script>

<style scoped lang="scss">

</style>
