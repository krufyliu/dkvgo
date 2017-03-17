const color = require('../src/utils/theme')
const Mock = require('mockjs')
const ms = require('../src/utils/mockStorge')

let dataKey = ms.mockStorge('Dashboard', Mock.mock({
  'cpu': {
    'usage|50-600': 1,
    space: 825,
    'cpu|40-90': 1,
    'data|20': [
      {
        'cpu|20-80': 1
      }
    ]
  },
  'browser': [
    {
      name: 'Google Chrome',
      percent: 43.3,
      status: 1
    },
    {
      name: 'Mozilla Firefox',
      percent: 33.4,
      status: 2
    },
    {
      name: 'Apple Safari',
      percent: 34.6,
      status: 3
    },
    {
      name: 'Internet Explorer',
      percent: 12.3,
      status: 4
    },
    {
      name: 'Opera Mini',
      percent: 3.3,
      status: 1
    },
    {
      name: 'Chromium',
      percent: 2.53,
      status: 1
    }
  ],
  user: {
    name: 'zuiidea',
    email: 'zuiiidea@gmail.com',
    sales: 3241,
    sold: 3556,
    avatar: 'http://tva4.sinaimg.cn/crop.0.0.996.996.180/6ee6a3a3jw8f0ks5pk7btj20ro0rodi0.jpg'
  },
  'completed|12': [
    {
      'name|+1': 2008,
      'Task complete|200-1000': 1,
      'Cards Complete|200-1000': 1
    }
  ],
  numbers: [
    {
      icon: 'pay-circle-o',
      color: color.green,
      title: 'Online Review',
      number: 2781
    }, {
      icon: 'team',
      color: color.blue,
      title: 'New Customers',
      number: 3241
    }, {
      icon: 'message',
      color: color.purple,
      title: 'Active Projects',
      number: 253
    }, {
      icon: 'shopping-cart',
      color: color.red,
      title: 'Referrals',
      number: 4324
    }
  ]
}))

module.exports = {
  'GET /api/dashboard' (req, res) {
    res.json(global[dataKey])
  }
}
