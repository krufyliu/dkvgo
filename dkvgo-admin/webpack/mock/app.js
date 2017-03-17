const qs = require('qs')
const Cookie = require('js-cookie')
const ms = require('../src/utils/mockStorge')

let dataKey = ms.mockStorge('AdminUsers', [
  {
    Id: 1,
    Email: 'admin@visiondk.com',
    Username: 'admin',
    Password: 'visiondk'
  }
])

module.exports = {
  'POST /api/auth' (req, res) {
    const userItem = qs.parse(req.body)
    const response = {
      success: false,
      message: ''
    }
    const d = global[dataKey].filter(function (item) {
      return item.Email === userItem.Email
    })
    if (d.length) {
      if (d[0].Password === userItem.Password) {
        const now = new Date()
        now.setDate(now.getDate() + 1)
        Cookie.set('user_session', now.getTime(), { path: '/' })
        Cookie.set('user_id', d[0].Id, { path: '/' })
        response.message = '登录成功'
        response.success = true
        response.user = d[0]
      } else {
        response.message = '密码不正确'
      }
    } else {
      response.message = '用户不存在'
    }
    res.json(response)
  },

  'GET /api/auth' (req, res) {
    const response = {
      success: Boolean(Cookie.get('user_id') != 'undefined' && Cookie.get('user_session') > new Date().getTime()),
      message: 'unauthorized'
    }
    var users = global[dataKey].filter(function (item) {
      return item.Id == parseInt(Cookie.get('user_id'))
    })
    if (users.length) {
      var user = users[0]
      delete(user['Password'])
      response['user'] = user
      response['message'] = "success"
    }
    res.json(response)
  },

  'DELETE /api/auth' (req, res) {
    Cookie.remove('user_session', { path: '/' })
    Cookie.remove('user_id', { path: '/' })
    res.json({
      success: true,
      message: '退出成功'
    })
  }
}
