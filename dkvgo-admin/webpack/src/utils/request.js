import Ajax from 'robe-ajax'

export default function request (url, options) {
    return Ajax.ajax({
      url: url,
      method: options.method || 'get',
      data: options.data || {},
      processData: options.method === 'get',
      dataType: options.dataType || 'JSON'
    }).done((data) => {
      console.log(data)
      return data
    })
}
