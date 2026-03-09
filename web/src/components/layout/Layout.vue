<i18n>
{
  "en": {
  },
  "cn": {
  }
}
</i18n>

<template>
  <div>
    <div class="main">
      <div class="main-header">
        <Header :pageLevel="pageLevel" :namespace="namespace" :percentage="percentage" :errorMsg="errorMsg" @collapse="isCollapse = $event"
          :breadcrumb="breadcrumb" :cluster="cluster" :switchToNamespaceHome="switchToNamespaceHome">
        </Header>
      </div>
      <div class="main-content">
        <Menu @collapse="isCollapse = $event"></Menu>
        <!-- <div class="rightColumn" > -->
        <div class="rightColumn" :style="contentStyle">
          <div class="content">
            <router-view :key="$route.path" v-slot="{ Component }">
              <transition name="el-fade-in" mode="out-in">
                <component :is="Component"></component>
              </transition>
            </router-view>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Menu from './Menu.vue'
import Header from './Header.vue'

export default {
  props: {
    cluster: { type: String, required: false, default: undefined },
    namespace: { type: String, required: false, default: undefined },
    breadcrumb: { type: Array, required: false, default: () => []},
    percentage: { type: Number, required: false, default: 100 },
    errorMsg: { type: String, required: false, default: undefined },
    pageLevel: { type: Boolean, required: false, default: false },
    switchToNamespaceHome: { type: Boolean, required: false, default: false },
  },
  data() {
    return {
      isCollapse: false,
    }
  },
  computed: {
    contentStyle () {
      let background = ''
      if (this.$route.meta.background) {
        background = 'background: ' + this.$route.meta.background + ';'
      }
      if (this.isCollapse) {
        // return 'width:  64px;' + background
        return 'width: calc(100vw - 64px);' + background
      } else {
        // return 'width: 220px;' + background
        return 'width: calc(100vw - 240px);' + background
      }
    }
  },
  components: { Menu, Header },
  mounted () {
    if (localStorage.getItem('menu-collapse') === 'true') {
      this.isCollapse = true
    } else {
      this.isCollapse = false
    }
  },
  methods: {

  }
}
</script>

<style scoped lang="css">
.main {
  padding: 0 0px;
  display: flex;
  background-color: #f1f4fa;
  flex-wrap: wrap; 
}
.content {
  display: inline-block;
  height: calc(100vh - 96px);
  width: 100%;
  overflow-x: auto;
  text-align: left;
  background-color: var(--el-color-white);
  /* padding: 0 20px; */
}
.footplaceholder {
  height: 30px;
}
</style>
