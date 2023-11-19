function addSplitSymbol(str: string, symbol = '-'): string {
  var result = ''
  for (var i = 0; i < str.length; i++) {
    result += str.charAt(i);
    if ((i + 1) % 3 === 0 && i !== str.length - 1) {
      result += symbol
    }
  }
  return result
}

function delSplitSymbol(str: string, symbol = '-'): string {
  return str.replaceAll(symbol, '')
}

// Reference: https://developers.weixin.qq.com/community/develop/article/doc/00062eedbfcf50f19a6bc1abc56013
function isWechat() {
  return /MicroMessenger/i.test(window.navigator.userAgent)
}

export {
  addSplitSymbol,
  delSplitSymbol,
  isWechat,
}
