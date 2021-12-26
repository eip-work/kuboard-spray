<i18n>
en:
  resourceList: Resources List
zh:
  resourceList: 资源包列表
</i18n>

<template>
  <div>
    <div class="app_block_title">
      {{ $t('obj.resource') }}
      
    </div>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'

export default {
  mixins: [mixin],
  percentage () {
    return 100
  },
  breadcrumb () {
    return [
      { label: this.$t('resourceList'), to: '/settings/resources' },
      { label: this.name },
    ]
  },
  refresh () {
    this.refresh()
  },
  props: {
    name: { type: String, required: true },
    mode: { type: String, required: false, default: 'view' },
  },
  data() {
    return {

    }
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
