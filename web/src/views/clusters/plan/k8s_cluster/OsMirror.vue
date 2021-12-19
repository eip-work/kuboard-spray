<i18n>
en:
  label: OS Mirror
zh:
  label: OS 镜像
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled>
    ansible_python_interpreter    /usr/bin/python3  for ubuntu
    <FieldString :holder="cluster.inventory.all.vars" fieldName="yum_repo"></FieldString>
    <FieldString :holder="cluster.inventory.all.vars" fieldName="ubuntu_repo"></FieldString>
  </ConfigSection>
</template>

<script>
import ConfigSection from '../ConfigSection.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    enabled: {
      get () {return true},
      set () {},
    },
  },
  components: { ConfigSection },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">

</style>
