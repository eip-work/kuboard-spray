import axios from 'axios'

async function init () {
    axios.get(`https://addons.kuboard.cn/kuboard-spray-latest-version.json?r=${Math.random()}`, {
      timeout: 2000,
      withCredentials: false,
    }).then(resp => {
      let installed = window.KuboardSpray.version.version
      let arch = installed.slice(installed.lastIndexOf('-') + 1)
      let channel = 'latest'
      if (installed.indexOf('beta') > 0) channel = 'beta'
      if (installed.indexOf('alpha') > 0) channel = 'alpha'
      let currentVersion = resp.data[arch].channel[channel]
      
      let newVersion = currentVersion.version
      console.log('installed version: ' + installed)
      console.log('      new version: ' + newVersion)
      if (installed === 'v1.0.0-dev-amd64' || installed === 'v1.0.0-dev-arm64') {
        console.log('ignore new version in dev mode.')
        return
      }
      if (installed !== newVersion) {
        if (localStorage.getItem('IgnoreThisVersion') === newVersion) {
          console.log('  Ignored version: ' + newVersion)
          return
        }
        setTimeout(() => {
            window.$kuboardVersionCheckNotification = window.VueAppComponent.$notify({
              title: '检测到新版本',
              dangerouslyUseHTMLString: true,
              message: `<a href="https://kuboard-spray.cn/support/change-log/v1.html" target=_blank class="app_notify" style="width: 100%;" onclick="window.$kuboardVersionCheckNotification.close()">
                        <div class="app_text_mono">已经安装的版本是 <span style="color: var(--el-color-success); font-weight: bold;">${installed} </span></div>
                        <div class="app_text_mono">当前最新的版本是 <span style="color: var(--el-color-danger); font-weight: bold;">${newVersion} </span></div>
                        <div style="text-align: right; margin-top: 10px; color: var(--el-color-primary); font-weight: bold;">
                          <i class="el-icon" style="margin-right: 5px;"><svg class="icon" width="200" height="200" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" data-v-042ca774=""><path fill="currentColor" d="M511.552 128c-35.584 0-64.384 28.8-64.384 64.448v516.48L274.048 570.88a94.272 94.272 0 00-112.896-3.456 44.416 44.416 0 00-8.96 62.208L332.8 870.4A64 64 0 00384 896h512V575.232a64 64 0 00-45.632-61.312l-205.952-61.76A96 96 0 01576 360.192V192.448C576 156.8 547.2 128 511.552 128zM359.04 556.8l24.128 19.2V192.448a128.448 128.448 0 11256.832 0v167.744a32 32 0 0022.784 30.656l206.016 61.76A128 128 0 01960 575.232V896a64 64 0 01-64 64H384a128 128 0 01-102.4-51.2L101.056 668.032A108.416 108.416 0 01128 512.512a158.272 158.272 0 01185.984 8.32L359.04 556.8z"></path></svg></i>
                          点击此处前往升级
                        </div>
                        </a>
                        <div onclick="ignoreThisVersion('${newVersion}')"
                          style="margin-left: -30px; padding: 10px 20px 0px 0;" class="app_text_description app_notify">
                          ${window.VueAppComponent.$t('upgrade.ignoreThisVersion')}
                        </div>`,
              position: 'bottom-left',
              type: 'info',
              duration: 0,
              'custom-class': 'kuboard-spray-upgrade-notification',
            });
        }, 5000)
      }
    }).catch(e => {
      console.log('查询新版本失败: ', e)
    })

}

export default {init}

window.ignoreThisVersion = function(version) {
  localStorage.setItem('IgnoreThisVersion', version)
  window.VueAppComponent.$message.success(window.VueAppComponent.$t('upgrade.ignored', {version}))
  window.$kuboardVersionCheckNotification.close()
}

