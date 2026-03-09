<template>
  <div :style="`display: inline-block; disabled: ${loading}`" class="noselect">
    <div v-if="to" class="kb-button router" :style="`width: ${size}; height: ${size}; ${loading ? 'cursor: not-allowed;' : ''}`">
      <router-link :to="to"  class="router">
          <div class="kb-content" :style="`width: ${contentSize}; height: ${contentSize}`">
            <slot v-if="!loading">
              <el-icon :size="contentSize" class="kb-icon" :style="`width: ${contentSize}; height: ${contentSize}; font-size: ${contentSize}`">
                <component :is="icon"></component>
              </el-icon>
            </slot>
            <el-icon v-else :size="contentSize" class="kb-icon" :style="`width: ${contentSize}; height: ${contentSize}; font-size: ${contentSize}`">
              <component :is="icon"></component>
            </el-icon>
          </div>
      </router-link>
    </div>
    <div v-else class="kb-button" @click.prevent.stop="_ => {if (!loading) {$emit('click')}}" :style="`width: ${size}; height: ${size}; ${loading ? 'cursor: not-allowed;' : ''}`">
      <div class="kb-content" :style="`width: ${contentSize}; height: ${contentSize}`">
        <slot v-if="!loading">
          <el-icon :size="contentSize" class="kb-icon" :style="`width: ${contentSize}; height: ${contentSize}; font-size: ${contentSize}`">
            <component :is="icon"></component>
          </el-icon>
        </slot>
        <el-icon v-else :size="contentSize" class="kb-icon" :style="`width: ${contentSize}; height: ${contentSize}; font-size: ${contentSize}`">
          <component :is="icon"></component>
        </el-icon>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    icon: { type: String, required: true },
    to: { type: String, required: false, default: undefined },
    size: { type: String, required: false, default: '30px' },
    contentSize: { type: Number, required: false, default: 15 },
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
  color: var(--el-color-white);
  transition: 0.2s;
}
.router:hover {
  color: var(--el-color-primary);
}
.kb-button {
  display: inline-block;
  width: 40px;
  height: 40px;
  color: var(--text-3);
  /* border-radius: 20px; */
  /* background: var(--el-color-primary); */
  vertical-align: middle;
  cursor: pointer;
  /* transition: 0.2s; */

  border-radius: 4px;
  background: #FFFFFF;
  border: 1px solid #EAEEEF;

}
.kb-button:hover {
  /* background-color: var(--el-color-primary-light-3); */
  color: var(--el-color-primary);
  /* transform:rotate(45deg);
  -ms-transform:rotate(45deg); 	
  -moz-transform:rotate(45deg); 	
  -webkit-transform:rotate(45deg); 
  -o-transform:rotate(45deg);  */

  border-radius: 4px;
  background: #FFFFFF;
  border: 1px solid var(--el-color-primary);

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
