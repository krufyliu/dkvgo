const Watch = require('watchjs')
const config = require('./config')

module.exports = {}

if (typeof localStorage === 'undefined') {
    localStorage = {
      data: {},
      getItem: function(key) {
        return this.data[key]
      },
      setItem(key, value) {
        this.data[key] = value
      }
    }
}

module.exports.default = module.exports.mockStorge = function (name, defaultValue) {
  let key = config.prefix + name
  global[key] = localStorage.getItem(key)
    ? JSON.parse(localStorage.getItem(key))
    : defaultValue
  !localStorage.getItem(key) && localStorage.setItem(key, JSON.stringify(global[key]))
  Watch.watch(global[key], function () {
    localStorage.setItem(key, JSON.stringify(global[key]))
  })
  return key
}
