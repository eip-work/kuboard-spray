<i18n>
en:
  titleRepo: OS mirrors
  repoDescription1: You muse provide a mirror address (e.g. CentOS yum mirror, Ubuntu apt mirror) that all the nodes in your mirror could access.
  repoDescription2: Usually, a enterprise has its own os mirror, if not, Kuboard alos provide a wizard, so that you can create an os mirror quickly.
  addMirror: Add OS Mirror
  confirmToDelete: This is going to delete OS Mirror info, but OS Mirror server is not affected here.
  type: Type
  kind: Provision Method
  status: Status
  os: OS
  mirror_type_os: OS source
  mirror_type_docker: Docker source
zh:
  titleRepo: OS 软件源
  repoDescription1: 您必须定义您集群机器可以访问的操作系统软件地址（例如 CentOS 的 yum 源、Ubuntu 的 apt 源等）以供使用
  repoDescription2: 通常企业都有自己的常用操作系统的本地软件仓库，Kuboard 提供向导，帮助您快速设置一个操作系统软件源
  addMirror: 添加 OS 软件源
  confirmToDelete: 此操作将删除软件仓库信息，软件仓库本身并不受此影响
  type: 类 型
  kind: 提供方式
  status: 状 态
  os: 操作系统
  mirror_type_os: 操作系统软件源
  mirror_type_docker: docker 软件源
</i18n>

<template>
  <div>
    <div class="app_block_title">{{$t('titleRepo')}}</div>
    <el-alert :title="$t('titleRepo')" :closable="false" class="app_white_alert">
      <div class="description">
        <li>{{$t('repoDescription1')}}</li>
        <li>{{$t('repoDescription2')}}</li>
      </div>
    </el-alert>

    <div class="contentList">
      <!-- <el-card shadow="never" style="min-height: 234px;"> -->
        <el-skeleton v-if="loading" :rows="5" animated />
        <div v-else>
          <div style="text-align: right;" class="app_margin_bottom">
            <el-button type="primary" icon="el-icon-plus" @click="$refs.create.show()">{{$t('addMirror')}}</el-button>
          </div>
          <el-table v-if="mirrors" :data="mirrors" style="width: 100%;">
            <el-table-column prop="name" :label="$t('msg.name')">
              <template #default="scope">
                <router-link :to="`/settings/mirrors/${scope.row.name}`">
                  <el-icon :size="14" style="width: 14px; height: 14px; vertical-align: middle; margin-right: 5px;">
                    <el-icon-link></el-icon-link>
                  </el-icon>
                  {{ cleanName(scope.row.name) }}
                </router-link>
              </template>
            </el-table-column>
            <el-table-column prop="status.type" :label="$t('type')" width="150px">
              <template #default="scope">
                <el-tag v-if="scope.row.name.indexOf('docker_') === 0">{{$t('mirror_type_docker')}}</el-tag>
                <el-tag type="warning" v-else>{{$t('mirror_type_os')}}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status.type" :label="$t('os')" width="150px">
              <template #default="scope">
                <span v-if="scope.row.name.indexOf('docker_') === 0">{{ scope.row.name.split('_')[1].split('-')[0] }}</span>
                <span v-else>{{ scope.row.name.split('-')[0] }}</span>
              </template>
            </el-table-column>
            <!-- <el-table-column prop="status.kind" :label="$t('kind')" width="120px"></el-table-column> -->
            <el-table-column prop="status.url" label="url">
              <template #default="scope">
                <KuboardSprayLink v-if="scope.row.status" :href="scope.row.status.url" :size="12">{{scope.row.status.url}}</KuboardSprayLink>
              </template>
            </el-table-column>
            <el-table-column prop="status.status" :label="$t('status')" width="120px"></el-table-column>
            <el-table-column :label="$t('msg.operations')" width="200px">
              <template #default="scope">
                <el-button icon="el-icon-view" type="primary" @click="$router.push(`/settings/mirrors/${scope.row.name}`)">{{ $t('msg.view') }}</el-button>
                <el-popconfirm :confirm-button-text="$t('msg.ok')" :cancel-button-text="$t('msg.cancel')" icon="el-icon-warning" icon-color="red"
                  placement="bottom-end" :title="$t('confirmToDelete')" @confirm="deleteMirror(scope.row.name)">
                  <template #reference>
                    <el-button icon="el-icon-delete" type="danger">{{ $t('msg.delete') }}</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
          <!-- <div v-for="(item, index) in mirrors" :key="'mirror' + index" class="mirror">
            <div class="deleteButton">
            </div>
            <el-card shadow="never"
              @click="$router.push(`/settings/mirrors/${item}`)">
              <div class="noselect">
                {{item}}
              </div>
            </el-card>
          </div> -->
        </div>
      <!-- </el-card> -->
    </div>
    <MirrorCreate ref="create"></MirrorCreate>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import MirrorCreate from './MirrorCreate.vue'
import mirror_options from './mirror_options.js'

export default {
  mixins: [mixin],
  props: {
  },
  data () {
    return {
      loading: false,
      mirrors: undefined,
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
    cleanName(name) {
      for (let key in mirror_options) {
        if (name.indexOf(key) === 0) {
          return name.slice(key.length + 1)
        }
      }
      return name
    },
    async refresh () {
      this.loading = true
      this.mirrors = undefined
      await this.kuboardSprayApi.get('/mirrors').then(resp => {
        let temp = []
        for (let i in resp.data.data) {
          temp.push({ name: resp.data.data[i] })
        }
        this.mirrors = temp
        for (let item of this.mirrors) {
          this.loadMirror(item)
        }
      }).catch(e => {
        console.log(e)
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
    },
    loadMirror (mirror) {
      this.kuboardSprayApi.get('/mirrors/' + mirror.name).then(resp => {
        mirror.status = resp.data.data.status
      }).catch(e => {
        console.log(e)
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
  border-color: var(--el-color-primary);
}
.mirror:hover .deleteButton {
  overflow: visible;
}
</style>
