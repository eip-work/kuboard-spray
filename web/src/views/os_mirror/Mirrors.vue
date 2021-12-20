<i18n>
en:
  titleRepo: OS mirrors
  repoDescription1: You muse provide a mirror address (e.g. CentOS yum mirror, Ubuntu apt mirror) that all the nodes in your mirror could access.
  repoDescription2: Usually, a enterprise has its own os mirror, if not, Kuboard alos provide a wizard, so that you can create an os mirror quickly.
  addMirror: Add OS Mirror
  confirmToDelete: This is going to delete OS Mirror info, but OS Mirror server is not affected here.
zh:
  titleRepo: OS 镜像源
  repoDescription1: 您必须定义您集群机器可以访问的操作系统镜像地址（例如 CentOS 的 yum 源、Ubuntu 的 apt 源等）以供使用
  repoDescription2: 通常企业都有自己的常用操作系统的本地镜像仓库，Kuboard 提供向导，帮助您快速设置一个操作系统镜像源
  addMirror: 添加 OS 镜像源
  confirmToDelete: 此操作将删除镜像仓库信息，镜像仓库本身并不受此影响
</i18n>

<template>
  <div>
    <div class="app_block_title">{{$t('titleRepo')}}</div>
    <el-alert :title="$t('titleRepo')" type="default" :closable="false">
      <div class="description">
        <li>{{$t('repoDescription1')}}</li>
        <li>{{$t('repoDescription2')}}</li>
      </div>
    </el-alert>

    <div class="contentList">
      <el-card shadow="none" style="min-height: 234px;">
        <el-skeleton v-if="loading" :rows="5" animated />
        <div v-else style="display: flex; flex-wrap: wrap;">
          <div v-for="(item, index) in mirrors" :key="'mirror' + index" class="mirror">
            <div class="deleteButton">
              <el-popconfirm :confirm-button-text="$t('msg.ok')" :cancel-button-text="$t('msg.cancel')" icon="el-icon-warning" icon-color="red"
                placement="right" :title="$t('confirmToDelete')" @confirm="deleteMirror(item)">
                <template #reference>
                  <el-button icon="el-icon-delete" circle type="danger"></el-button>
                </template>
              </el-popconfirm>
            </div>
            <el-card shadow="none"
              @click="$router.push(`/settings/mirrors/${item}`)">
              <div class="noselect">
                {{item}}
              </div>
            </el-card>
          </div>
          <el-button type="primary" size="large" icon="el-icon-plus" @click="$refs.create.show()">{{$t('addMirror')}}</el-button>
        </div>
      </el-card>
    </div>
    <MirrorCreate ref="create"></MirrorCreate>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import MirrorCreate from './MirrorCreate.vue'

export default {
  mixins: [mixin],
  props: {
  },
  data () {
    return {
      loading: false,
      mirrors: [],
    }
  },
  refresh () {
    this.refresh()
  },
  percentage () {
    return 100
  },
  breadcrumb () {
    return [
      { label: this.$t('titleRepo') },
    ]
  },
  computed: {
  },
  components: { MirrorCreate },
  mounted () {
    this.refresh()
  },
  methods: {
    async refresh () {
      this.loading = true
      await this.kuboardSprayApi.get('/mirrors').then(resp => {
        this.mirrors = resp.data.data
      }).catch(e => {
        this.$message.error('Error: ' + e)
      })
      this.loading = false
    },
    deleteMirror(mirror) {
      this.kuboardSprayApi.delete('/mirrors/' + mirror).then(() => {
        this.refresh()
        this.$message.success(this.$t('msg.delete_succeeded'))
      }).catch(e => {
        this.$message.error(this.$t('msg.delete_failed', {msg: e.response.data.message }))
      })
    }
  }
}
</script>

<style scoped lang="scss">
.description {
  line-height: 28px;
}
.contentList {
  margin: 10px 0;
}
.mirror {
  margin-right: 10px;
  width: 200px;
  border-radius: 6px;
  cursor: pointer;
  .deleteButton {
    height: 0px;
    overflow: hidden;
    position: relative;
    top: -5px;
    left: 5px;
    text-align: right;
  }
}
.mirror:hover {
  border-color: $--color-primary;
}
.mirror:hover .deleteButton {
  overflow: visible;
}
</style>
