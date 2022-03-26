import { ElNotification } from 'element-plus'

function entry () {
  setTimeout(() => {
    showTip()
  }, 5000)
}

window.kuboardspray_github_star_action = function () {
  window.kuboardspray_github_star_ref.close()
  localStorage.setItem('kuboardspray_github_star_visited', 'true')
}

function showTip() {
  if (localStorage.getItem('kuboardspray_github_star_visited') === 'true') {
    return
  }
  window.kuboardspray_github_star_ref = ElNotification({
    title: '加关注，不迷路',
    dangerouslyUseHTMLString: true,
    // message: '<a onClick="kuboardspray_open_github_star()">This is <i>HTML</i> string</a>',
    message: `<div class="app_text_mono">加 Github star / watch，可以随时获取本项目的源代码、文档</div>
              <a href="https://github.com/eip-work/kuboard-spray" target=_blank class="app_notify" style="width: 100%;" onclick="window.kuboardspray_github_star_action()">
              <div style="text-align: right; margin-top: 10px; color: var(--el-color-primary); font-weight: bold;">
                <i class="el-icon" style="margin-right: 5px;"><svg class="icon" width="200" height="200" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" data-v-042ca774=""><path fill="currentColor" d="M511.552 128c-35.584 0-64.384 28.8-64.384 64.448v516.48L274.048 570.88a94.272 94.272 0 00-112.896-3.456 44.416 44.416 0 00-8.96 62.208L332.8 870.4A64 64 0 00384 896h512V575.232a64 64 0 00-45.632-61.312l-205.952-61.76A96 96 0 01576 360.192V192.448C576 156.8 547.2 128 511.552 128zM359.04 556.8l24.128 19.2V192.448a128.448 128.448 0 11256.832 0v167.744a32 32 0 0022.784 30.656l206.016 61.76A128 128 0 01960 575.232V896a64 64 0 01-64 64H384a128 128 0 01-102.4-51.2L101.056 668.032A108.416 108.416 0 01128 512.512a158.272 158.272 0 01185.984 8.32L359.04 556.8z"></path></svg></i>
                点击前往 Github
              </div>
              </a>`,
    showClose: false,
    position: 'bottom-left',
    type: 'warning',
    duration: 0,
    'custom-class': 'kuboard-spray-upgrade-notification',
  })
}

export default entry
