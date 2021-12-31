<template>
  <div :style="`display: inline-block; disabled: ${loading}`" class="noselect">
    <div v-if="to" class="kb-button router" :style="`width: ${size}; height: ${size}; ${loading ? 'cursor: not-allowed;' : ''}`">
      <router-link :to="to"  class="router">
          <div class="kb-content" :style="`width: ${contentSize}; height: ${contentSize}`">
            <slot v-if="!loading">
              <i :class="`${icon} kb-icon`" :style="`width: ${contentSize}; height: ${contentSize}; font-size: ${contentSize}`"></i>
            </slot>
            <i v-else class="el-icon-loading kb-icon" :style="`width: ${contentSize}; height: ${contentSize}; font-size: ${contentSize}`"></i>
          </div>
      </router-link>
    </div>
    <div v-else class="kb-button" @click.prevent.stop="_ => {if (!loading) {$emit('click')}}" :style="`width: ${size}; height: ${size}; ${loading ? 'cursor: not-allowed;' : ''}`">
      <div class="kb-content" :style="`width: ${contentSize}; height: ${contentSize}`">
        <slot v-if="!loading">
          <i :class="`${icon} kb-icon`" :style="`width: ${contentSize}; height: ${contentSize}; font-size: ${contentSize}`"></i>
        </slot>
        <i v-else class="el-icon-loading kb-icon" :style="`width: ${contentSize}; height: ${contentSize}; font-size: ${contentSize}`"/>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    icon: { type: String, required: true },
    to: { type: String, required: false, default: undefined },
    size: { type: String, required: false, default: '30px'},
    contentSize: { type: String, required: false, default: '15px'},
    loading: { type: Boolean, required: false, default: false },
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
.router a {
  color: $--color-white;
  transition: 0.2s;
}
.router:hover {
  color: $--color-primary;
}
.kb-button {
  display: inline-block;
  width: 40px;
  height: 40px;
  border-radius: 20px;
  // background: rgba(#FFF, 0.2);
  background: $--color-primary;
  vertical-align: middle;
  cursor: pointer;
  transition: 0.2s;
}
.kb-button:hover {
  // background-color: $--color-white;
  // color: $--color-primary;
  // background-color: rgba(#FFF, 0.2);
  background-color: $--color-primary-light-2;
  color: $--color-white;
  text-align: center;
}
.kb-content {
  width: 24px;
  height: 24px;
  padding: 8px;
  font-size: 24px;
  
}
.kb-icon {
  vertical-align: top;
}
</style>
