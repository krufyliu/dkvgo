import Ajax from 'robe-ajax'
import { getCookie } from './cookie'
import base64_decode from 'locutus/php/url/base64_decode'

function getXsrfToken() {
  var xsrf, xsrflist
  xsrf = getCookie('_xsrf')
  xsrflist = xsrf.split("|")
  return xsrflist[0]
}

export default function request (url, options) {
  var method, data, xsrf
  method = options.method || 'get'
  data = options.data || {} 
  xsrf = getXsrfToken()
  if (xsrf != "") {
    data._xsrf = base64_decode(xsrf)
  }
  return Ajax.ajax({
    url: url,
    method: method,
    data: data,
    processData: options.method === 'get',
    dataType: options.dataType || 'JSON'
  }).done((data) => {
    console.log(data)
    return data
  })
}
