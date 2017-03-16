import Ajax from 'robe-ajax'
import { getCookie } from './cookie'
import base64_decode from 'locutus/php/url/base64_decode'
import {message} from 'antd'

export default function request (url, options) {
  function getXsrfToken() {
    var xsrf, xsrflist
    xsrf = getCookie('_xsrf')
    xsrflist = xsrf.split("|")
    return xsrflist[0]
  }

  var method, headers, xsrf
  method = options.method || 'get'
  headers = options.headers || {} 
  xsrf = getXsrfToken()
  if (xsrf != "" && method.toUpperCase() !== 'GET') {
    headers['X-Xsrftoken'] = base64_decode(xsrf)
  }
  return Ajax.ajax({
    url: url,
    method: method,
    data: options.data || {},
    headers: headers,
    // processData: options.method === 'get',
    dataType: options.dataType || 'JSON'
  }).done((data) => {
    if (data && !data.success) {
      console.log(data.message)
      message.error(data.message, 5)
    }
    return data
  })
}
