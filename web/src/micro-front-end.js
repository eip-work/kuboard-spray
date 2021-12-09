function getParams() {
  var params = {};
  if (location.search.indexOf("?") == 0 && location.search.indexOf("=") > 1) {
      var paramArray = unescape(location.search).substring(1, location.search.length).split("&");
      if (paramArray.length > 0) {
          paramArray.forEach(function (currentValue) {
              params[currentValue.split("=")[0]] = currentValue.split("=")[1];
          });
      }
  }
  return params;
}

let isInFrame = self !== top

export function initKuboardMfe (app) {
  let KuboardMfe = {
    isInFrame: isInFrame,
    KUBOARD_FRAME_ID: getParams().KUBOARD_FRAME_ID,
    KUBOARD_FRAME_FIXED_HEIGHT: getParams().KUBOARD_FRAME_FIXED_HEIGHT,
    KUBOARD_ACCESS_ENDPOINT: getParams().KUBOARD_ACCESS_ENDPOINT,
    sendToParent (command, params) {
      if (window.parent) {
        window.parent.postMessage({
          kuboardCommandToParent: KuboardMfe.KUBOARD_FRAME_ID + ":" + command,
          params: params
        }, '*')
      }
    },
    messagesuccess (msg) {
      let message = {
        type: 'success',
        message: msg
      }
      if (isInFrame) {
        KuboardMfe.sendToParent('message', message)
      } else {
        window.VueAppComponent.$message(message)
      }
    },
    messageerror (msg) {
      let message = {
        type: 'error',
        message: msg
      }
      console.log('sendmessageerror', msg)
      if (isInFrame) {
        KuboardMfe.sendToParent('message', message)
      } else {
        window.VueAppComponent.$message(message)
      }
    },
    messageinfo (msg) {
      let message = {
        type: 'info',
        message: msg
      }
      console.log('sendmessageinfo', msg)
      if (isInFrame) {
        KuboardMfe.sendToParent('message', message)
      } else {
        window.VueAppComponent.$message(message)
      }
    }
  }
  
  window.KuboardMfe = KuboardMfe
  app.config.globalProperties.KuboardMfe = KuboardMfe
  
  console.log(KuboardMfe.KUBOARD_FRAME_ID, KuboardMfe.KUBOARD_FRAME_FIXED_HEIGHT)
  if (isInFrame && !KuboardMfe.KUBOARD_FRAME_FIXED_HEIGHT) {
    setInterval(() => {
      // console.log(window.document.body.clientHeight)
      KuboardMfe.sendToParent('iframeSize', {
        height: window.document.body.clientHeight,
      })
    }, 1000)
  }
}

