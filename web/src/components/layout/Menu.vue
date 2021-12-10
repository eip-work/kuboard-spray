<i18n>
en:
  expand: Open
  collapse: Close
  clusters: Kubernetes Cluster
  settings: Settings
  resources: Resources Package
  kuboard: Kuboard Integration
zh:
  expand: 展开
  collapse: 收起
  clusters: 集群管理
  settings: 系统设置
  resources: 资源包管理
  kuboard: Kuboard 集成
</i18n>

<template>
  <div class="leftColumn">
    <div class="menuTitle noselect nowrap" :style="menuStyle" @click="$router.push('/')">
        Kuboard Spray
    </div>
    <el-menu :default-active="defaultActive" class="kuboard_namespaced_menu" unique-opened
      :collapse="isCollapse" background-color="white"
      ref="menu" router @open="handleOpen" @close="handleClose">
      <el-menu-item index="1" :route="`/clusters`" class="toplevel">
          <i class="el-icon-s-home"></i>
          <template #title>
            <span>{{$t('clusters')}}</span>
          </template>
      </el-menu-item>
      <el-submenu index="2">
        <template #title>
          <i class="el-icon-setting"></i>
          <span>{{$t('settings')}}</span>
        </template>
        <el-menu-item index="2-1" :route="`/settings/resources`">
          {{$t('resources')}}
        </el-menu-item>
        <el-menu-item index="2-2" :route="`/settings/kuboard`">
          {{$t('kuboard')}}
        </el-menu-item>
      </el-submenu>

    </el-menu>
    <div class="menu-toggler" @click="toggle">
      <i :class="isCollapse ? 'el-icon-d-arrow-right': 'el-icon-d-arrow-left'"/>
      {{isCollapse ? $t('expand') : $t('collapse')}}
    </div>
  </div>
</template>


<script>

export default {
  props: {
  },
  data () {
    let isCollapse = false
    if (localStorage.getItem('menu-collapse') === 'true') {
      isCollapse = true
    }
    return {
      defaultActive: '0',
      isCollapse: isCollapse,
    }
  },
  computed: {
    menuStyle () {
      if (this.isCollapse) {
        return 'width: 64px;'
      } else {
        return 'width: calc(6vw + 150px);'
      }
    }
  },
  watch: {
    '$route.path' () {
      this.$nextTick(() => {
        let active = this.refreshDefaultActive()
        this.defaultActive = active
      })
    },
  },
  mounted () {
    this.$nextTick(() => {
      let active = this.refreshDefaultActive()
      this.defaultActive = active
    })
  },
  methods: {
    refreshDefaultActive () {
      let active = '0'
      if (this.$refs.menu === undefined) {
        return active
      }

      let map = {
        'Home': '0',
        'Clusters': '1',
        'Cluster': '1',
        'Resources': '2-1',
        'Kuboard': '2-2'
      }
      if (map[this.$route.name]) {
        return map[this.$route.name]
      }

      return active
    },
    toggle() {
      this.isCollapse = !this.isCollapse
      localStorage.setItem('menu-collapse', this.isCollapse)
      this.$emit('collapse', this.isCollapse)
    },
    handleOpen(key, keyPath) {
      console.log(key, keyPath);
    },
    handleClose(key, keyPath) {
      console.log(key, keyPath);
    }
  }
}
</script>

<style lang="scss">
@import "./menu.scss";

.kuboard_global_menu {
  
  text-align: left;
  // width: 180px;
  background-color: rgb(244, 244, 245);
  padding: 10px 0 0 0;
  vertical-align: top;
  overflow: hidden;
  overflow-y: auto;
  border: none;
}
.kuboard_global_menu:not(.el-menu--collapse) {
  width: calc(8vw + 90px);
}
.kuboard_global_menu .el-menu-item-group__title {
  padding-left: 25px !important;
}
</style>
