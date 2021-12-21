<i18n>
en:

zh:


</i18n>

<template>
  <div>
    <FieldString :holder="inventory.all.children.target.vars" prop="inventory.all.children.target.vars" fieldName="apt_mirror_apache_root_dir" required></FieldString>
  </div>
</template>

<script>
export default {
  props: {
    os_mirror: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    inventory: {
      get () { return this.os_mirror.inventory },
      set () {},
    }
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
