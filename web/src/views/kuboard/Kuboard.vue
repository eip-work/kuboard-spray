<i18n>
en:
  hi: kuboard
zh:
  hi: kuboard
</i18n>

<template>
  <div>
    <div class="app_block_title">正在建设 ...</div>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'

export default {
  mixins: [mixin],
  props: {
  },
  percentage () {
    return 100
  },
  breadcrumb () {
    return [
      { label: 'Kuboard 集成' }
    ]
  },
  computed: {
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">

</style>
