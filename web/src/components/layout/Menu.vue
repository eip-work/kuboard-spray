<i18n locale="en" lang="yaml">
  expand: Open
  collapse: Close
  clusters: Kubernetes Cluster
  settings: Settings
  resources: Resources Package
  mirrors: OS Mirrors
  profile: Change Password
</i18n>
<i18n locale="zh" lang="yaml">
  expand: 展开
  collapse: 收起
  clusters: 集群管理
  settings: 系统设置
  resources: 资源包管理
  mirrors: OS 软件源
  profile: 修改密码
</i18n>

<template>
  <div class="leftColumn" :style="menuStyle" >
    <!-- <div class="menuTitle noselect nowrap" :style="menuStyle" @click="$router.push('/home')">
      Pangee Cluster
    </div> -->
    <el-menu :default-active="defaultActive" class="pangeecluster_namespaced_menu" unique-opened :collapse="isCollapse"
      background-color="white" ref="menu" router>
      <el-menu-item index="1" :route="`/clusters`" class="toplevel">
        <el-icon>
          <el-icon-home-filled></el-icon-home-filled>
        </el-icon>
        <template #title>
          <span>{{ t('clusters') }}</span>
        </template>
      </el-menu-item>
      <el-sub-menu index="2">
        <template #title>
          <el-icon>
            <el-icon-setting></el-icon-setting>
          </el-icon>
          <span>{{ t('settings') }}</span>
        </template>
        <el-menu-item index="2-1" :route="`/settings/resources`">
          {{ t('resources') }}
        </el-menu-item>
        <el-menu-item index="2-4" :route="`/settings/profile`">
          {{ t('profile') }}
        </el-menu-item>
      </el-sub-menu>

    </el-menu>
    <div class="menu-toggler" @click="toggle">
      <el-icon style="vertical-align: middle;" :size="14">
        <el-icon-d-arrow-right v-if="isCollapse"></el-icon-d-arrow-right>
        <el-icon-d-arrow-left v-else></el-icon-d-arrow-left>
      </el-icon>
      {{ isCollapse ? t('expand') : t('collapse') }}
    </div>
  </div>
</template>


<script>

export default {
  props: {
  },
  data() {
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
    menuStyle() {
      if (this.isCollapse) {
        return 'width: 64px;'
      } else {
        return 'width: 240px;'
        // return 'width: calc(6vw + 220px);'
      }
    }
  },
  watch: {
    '$route.path'() {
      this.$nextTick(() => {
        let active = this.refreshDefaultActive()
        this.defaultActive = active
      })
    },
  },
  mounted() {
    this.$nextTick(() => {
      let active = this.refreshDefaultActive()
      this.defaultActive = active
    })
  },
  methods: {
    refreshDefaultActive() {
      let active = '0'
      if (this.$refs.menu === undefined) {
        return active
      }

      let map = {
        'Home': '0',
        'Clusters': '1',
        'Cluster': '1',
        'Resources': '2-1',
        'Resource': '2-1',
        'ResourceOnAir': '2-1',
        'Mirrors': '2-2',
        'Mirror': '2-2',
        'Profile': '2-4',
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
  }
}
</script>

<style lang="css">
@import "./menu.scss";

.pangeecluster_global_menu {
  text-align: left;
  background-color: rgb(244, 244, 245);
  padding: 10px 0 0 0;
  vertical-align: top;
  overflow: hidden;
  overflow-y: auto;
  border: none;
}

.pangeecluster_global_menu:not(.el-menu--collapse) {
  width: calc(8vw + 90px);
}

.pangeecluster_global_menu .el-menu-item-group__title {
  padding-left: 25px !important;
}
</style>
