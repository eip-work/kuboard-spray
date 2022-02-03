<i18n>
en:
  title: Resource Packages
zh:
  title: 资源包列表
</i18n>

<template>
  <div>
    <Resources ref="resources"></Resources>
  </div>
</template>

<script>
import Resources from './Resources.vue'
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
      { label: this.$t('title') }
    ]
  },
  refresh () {
    this.$refs.resources.refresh()
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { Resources },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="css">

</style>
