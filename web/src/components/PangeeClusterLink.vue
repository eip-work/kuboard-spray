<template>
  <el-link target="blank" :href="href">
    <el-icon :style="`vertical-align: middle; width: ${size}px; height: ${size}px; margin-right: 5px;`">
      <component :is="icon" :style="`width: ${size}px; height: ${size}px;`"></component>
    </el-icon>
    <slot>
      <span :style="`font-size: ${size}px;`">{{ $t('msg.help') }}</span>
    </slot>
  </el-link>
</template>

<script>
export default {
  props: {
    href: { type: String, required: true },
    size: { Number, required: false, default: 14 },
    icon: { String, required: false, default: 'el-icon-link' }
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
