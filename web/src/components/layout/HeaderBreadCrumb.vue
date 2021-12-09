<template>
  <div class="bread-crumb">
    <!-- <div class="kind">{{kind}}</div> -->
    <div class="wrapper noselect" :style="hasNext ? '' : 'cursor: auto'">
      <router-link v-if="to" class="label" :to="to" style="cursor: pointer;">{{label}}</router-link>
      <span v-else :class="hasNext ? 'label' : 'label isLeaf'" style="cursor: pointer;" @click="$emit('click')">{{label}}</span>
      <slot></slot>
      <i v-if="hasNext" class="el-icon-arrow-right icon"></i>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    kind: { type: String, required: false },
    label: { type: String, required: false },
    to: { type: String, required: false, default: undefined },
    hasNext: { type: Boolean, required: false, default: true },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.bread-crumb {
  height: 70px;
  line-height: 70px;
  display: inline-block;
}
.wrapper {
  height: 70px;
  line-height: 70px;
}
.kind {
  opacity: 0.6;
  font-size: 13px;
  line-height: 22px;
  color: $--color-text-placeholder;
  transition: 0.5s;
}
.bread-crumb:hover .kind {
  opacity: 1;
}
.label {
  color: $--color-primary;
  font-weight: bold;
  font-size: 14px;
  border-bottom: dashed 1px $--color-primary;
}
.isLeaf {
  color: $--color-primary;
  font-weight: bold;
  font-size: 14px;
  border-bottom: none;
}
.icon {
  margin: 0 5px;
  color: $--color-primary-light-5;
  font-weight: 800;
}
</style>
