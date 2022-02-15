<i18n>
en:
  execute_scan: Re-Execute CIS Scan
  expected_result: Expected Result
  remediation: Remediation
  manual: manual
  actual_value: Actual value
  audit: Audit command
  AuditEnv: Audit Env
  AuditConfig: AuditConfig
  expected_and_remediation: Expected and remediation
  last_run_time: Scaned at
zh:
  execute_scan: 重新执行 CIS 扫描
  expected_result: 预期结果
  remediation: 问题修复
  manual: 人工
  actual_value: 实测结果
  audit: 审计命令
  AuditEnv: 审计命令环境变量
  AuditConfig: 审计命令配置文件
  expected_and_remediation: 预期及修复
  last_run_time: 检查时间
</i18n>

<template>
  <div>
    <el-scrollbar height="calc(100vh - 220px)">
      <div class="app_block_title">
        Kube-bench
        <el-button type="primary" icon="el-icon-promotion" style="margin-left: 20px;"
          @click="execute_scan('ignore_cache')" :loading="loading">
          {{ $t('execute_scan') }}
        </el-button>
        <div v-if="loadingPercentage < 100" style="display: inline-block; width: 200px; margin-left: 20px; vertical-align: middle;">
          <el-progress :percentage="loadingPercentage" :stroke-width="20" text-inside></el-progress>
        </div>
        <span v-else style="font-size: 13px; font-weight: normal; margin-left: 20px;">
          {{ $t('last_run_time') }}
          <KuboardSprayTime :time="last_run_time"></KuboardSprayTime>
        </span>
        <KuboardSprayLink style="float: right;" href="https://kuboard-spray.cn/guide/cis-scan/kube-bench.html" target="blank">Kubebench CIS Scan</KuboardSprayLink>
      </div>
      <div>
      </div>
      <el-table ref="table" :data="tableRows" style="width: 100%; margin-bottom: 20px" row-key="id"
        :expand-row-keys="expanedRowKeys"
        :default-sort="{ prop: 'id', order: 'ascending' }">
        <el-table-column prop="id" label="id" sortable width="420" fixed>
          <template #default="scope">
            <span class="row_title" @click="$refs.table.toggleRowExpansion(scope.row)">
              [{{scope.row.id.trim()}}]
              {{scope.row.text}}
              {{scope.row.desc}}
              {{scope.row.test_desc}}
            </span>
          </template>
        </el-table-column>
        <el-table-column width="120" fixed>
          <template #header>
            {{$t('expected_and_remediation')}}
          </template>
          <template #default="scope">
            <el-tooltip v-if="scope.row.expected_result || scope.row.remediation" class="box-item" effect="dark"
              placement="top-end" trigger="click">
              <el-button>{{$t('expected_and_remediation')}}</el-button>
              <template #content>
                <div style="width: 650px;">
                  <template v-if="scope.row.test_info && false">
                    <div class="test_title" >{{ $t('test_info') }}</div>
                    <pre v-for="(info, index) in scope.row.test_info" :key="scope.row.id + 'info' + index">{{ info }}</pre>
                  </template>
                  <template v-if="scope.row.AuditConfig">
                    <div class="test_title">{{$t('AuditConfig')}}</div>
                    <pre>{{scope.row.AuditConfig}}</pre>
                  </template>
                  <template v-if="scope.row.AuditEnv">
                    <div class="test_title">{{$t('AuditEnv')}}</div>
                    <pre>{{scope.row.AuditEnv}}</pre>
                  </template>
                  <template v-if="scope.row.audit">
                    <div class="test_title">{{$t('audit')}}</div>
                    <pre>{{scope.row.audit}}</pre>
                  </template>
                  <template v-if="scope.row.expected_result">
                    <div class="test_title" >{{ $t('expected_result') }}</div>
                    <pre>{{ scope.row.expected_result }}</pre>
                  </template>
                  <div class="test_title" >{{ $t('remediation') }}</div>
                  <pre>{{ scope.row.remediation }}</pre>
                </div>
              </template>
            </el-tooltip>
          </template>
        </el-table-column>
        <template v-for="(node, name) in mergedResult" :key="'re' + name">
          <el-table-column :label="name" width="180">
            <template #header>
              {{name}}
            </template>
            <template #default="scope">
              <span v-if="loadingPercentage < 100 && scope.row.id === '1'" style="display: none;">{{loadingPercentage}}</span>
              <div v-if="node[scope.row.id] && node[scope.row.id].total_pass !== undefined" 
                class="row_click" @click="$refs.table.toggleRowExpansion(scope.row)">
                <el-tag v-if="node[scope.row.id].total_pass" type="success" effect="dark">{{node[scope.row.id].total_pass}}</el-tag>
                <el-tag v-if="node[scope.row.id].total_fail" type="danger" effect="dark">{{node[scope.row.id].total_fail}}</el-tag>
                <el-tag v-if="node[scope.row.id].total_warn" type="warning" effect="dark">{{node[scope.row.id].total_warn}}</el-tag>
                <el-tag v-if="node[scope.row.id].total_info" type="info" effect="light">{{node[scope.row.id].total_info}}</el-tag>
              </div>
              <div v-if="node[scope.row.id] && node[scope.row.id].pass !== undefined"
                class="row_click" @click="$refs.table.toggleRowExpansion(scope.row)">
                <el-tag v-if="node[scope.row.id].pass" type="success" effect="dark">{{node[scope.row.id].pass}}</el-tag>
                <el-tag v-if="node[scope.row.id].fail" type="danger" effect="dark">{{node[scope.row.id].fail}}</el-tag>
                <el-tag v-if="node[scope.row.id].warn" type="warning" effect="dark">{{node[scope.row.id].warn}}</el-tag>
                <el-tag v-if="node[scope.row.id].info" type="info" effect="light">{{node[scope.row.id].info}}</el-tag>
              </div>
              <template v-if="node[scope.row.id] && node[scope.row.id].test_number !== undefined">
                <template v-if="node[scope.row.id].type === 'manual'">
                  <el-tag type="info">{{$t('manual')}}</el-tag>
                </template>
                <template v-else>
                  <el-tag v-if="node[scope.row.id].status === 'PASS'" type="success">{{node[scope.row.id].status}}</el-tag>
                  <el-tag v-if="node[scope.row.id].status === 'WARN'" type="warning">{{node[scope.row.id].status}}</el-tag>
                  <el-tag v-if="node[scope.row.id].status === 'FAIL'" type="danger">{{node[scope.row.id].status}}</el-tag>
                  <el-tag v-if="node[scope.row.id].status === 'INFO'" type="info">{{node[scope.row.id].status}}</el-tag>
                </template>
                <el-tooltip v-if="node[scope.row.id].actual_value && node[scope.row.id].status !== 'PASS'" class="box-item" effect="dark"
                  placement="top-end" trigger="click">
                  <template #content>
                    <pre style="width: 650px;">{{node[scope.row.id].actual_value}}</pre>
                  </template>
                  <el-tag class="reason" type="info">{{$t('actual_value')}}</el-tag>
                </el-tooltip>
                <el-tooltip v-if="node[scope.row.id].reason" class="box-item" effect="dark"
                  placement="top-end" trigger="click">
                  <template #content>
                    <pre style="width: 650px;">{{node[scope.row.id].reason.replaceAll('\\n', '\n')}}</pre>
                  </template>
                  <el-tag class="reason" type="info">{{$t('reason')}}</el-tag>
                </el-tooltip>
              </template>
            </template>
          </el-table-column>
        </template>
      </el-table>
    </el-scrollbar>
  </div>
</template>

<script>
import clone from 'clone'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      first_check_started: false,
      kube_control_plane: undefined,
      kube_node: undefined,
      etcd: undefined,
      metadata: undefined,
      expanedRowKeys: ['1', '2', '3', '4', '5'],
      last_run_time: undefined,
    }
  },
  computed: {
    loadingPercentage () {
      let count = 10
      if (this.kube_control_plane) {
        count += 30
      }
      if (this.kube_node) {
        count += 30
      }
      if (this.etcd) {
        count += 30
      }
      return count
    },
    loading () {
      return this.first_check_started && (this.kube_control_plane === undefined || this.kube_node === undefined || this.etcd === undefined)
    },
    tableRows () {
      let result = []
      if (this.metadata) {
        for (let item of this.metadata) {
          let temp = clone(item)
          temp.children = temp.tests
          for (let i in temp.children) {
            let test = temp.children[i]
            test.children = test.results
            test.id = test.section + ' '
            for (let i in test.children) {
              test.children[i].id = test.children[i].test_number
            }
          }
          result.push(temp)
        }
      }
      return result
    },
    mergedResult () {
      let result = {}
      if (this.kube_control_plane) {
        for (let nodename in this.kube_control_plane) {
          result[nodename] = result[nodename] || {}
          for (let key in this.kube_control_plane[nodename]) {
            result[nodename][key] =  this.kube_control_plane[nodename][key]
          }
        }
      }
      if (this.kube_node) {
        for (let nodename in this.kube_node) {
          result[nodename] = result[nodename] || {}
          for (let key in this.kube_node[nodename]) {
            result[nodename][key] =  this.kube_node[nodename][key]
          }
        }
      }
      if (this.etcd) {
        for (let nodename in this.etcd) {
          result[nodename] = result[nodename] || {}
          for (let key in this.etcd[nodename]) {
            result[nodename][key] =  this.etcd[nodename][key]
          }
        }
      }
      return result
    }
  },
  components: { },
  mounted () {
    this.retrieve_metadata()
    this.execute_scan('force_cache')
  },
  methods: {
    execute_scan(cache_mode) {
      this.do_scan('kube_control_plane', cache_mode)
      this.do_scan('kube_node', cache_mode)
      this.do_scan('etcd', cache_mode)
      this.first_check_started = true
    },
    do_scan(target, cache_mode) {
      this[target] = undefined
      this.kuboardSprayApi.post(`/clusters/${this.cluster.name}/cis_scan`, { target: target, cache_mode }).then(resp => {
        this[target] = {}
        this.last_run_time = resp.data.last_run_time
        for (let index in resp.data.data.plays[0].tasks[0].hosts) {
          let host = resp.data.data.plays[0].tasks[0].hosts[index]
          let hostresult = JSON.parse(host.stdout)
          let temp = {}
          for (let item of hostresult.Controls) {
            temp[item.id] = item
            for (let test of item.tests) {
              temp[test.section + ' '] = test
              for (let k in test.results) {
                let r = test.results[k]
                temp[r.test_number] = r
              }
              delete test.results
            }
            delete item.tests
          }
          this[target][index] = this[target][index] || temp
        }
      }).catch(e => {
        console.log(e)
        if (e.response && e.response.data.message) {
          if (e.response.data.message === 'no cache found') {
            this[target] = {}
          } else {
            this.$message.error('failed to execute cis_scan: ', e.response.data.message)
          }
        }
      })
    },
    retrieve_metadata() {
      this.kuboardSprayApi.post(`/clusters/${this.cluster.name}/cis_scan`, { target: 'kube_control_plane[0]', cache_mode: 'prefer_cache' }).then(resp => {
        this.metadata = {}
        for (let index in resp.data.data.plays[0].tasks[0].hosts) {
          let host = resp.data.data.plays[0].tasks[0].hosts[index]
          let hostresult = JSON.parse(host.stdout)
          // console.log(hostresult)
          this.metadata = hostresult.Controls
        }
      }).catch(e => {
        if (e.response && e.response.data.message) {
          this.$message.error('failed to execute cis_scan: ', e.response.data.message)
        }
      })
    }
  }
}
</script>

<style scoped lang="scss">
.test_title {
  font-weight: bolder;
  font-size: 14px;
  color: white;
  line-height: 24px;
}
.reason {
  margin-left: 10px;
  cursor: pointer;
}
.row_title {
  cursor: pointer;
}
.row_click {
  cursor: pointer;
}
</style>
