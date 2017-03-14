const qs = require('qs')
const Mock = require('mockjs')
const ms = require('../src/utils/mockStorge')

let dataKey = ms.mockStorge('JobsList', Mock.mock({
  'data|100': [
    {
        'Id|+1': 1,
        Name:              "@title",
        Priority:          100,
        Progress:          "@float(0, 100)",
        Status:            "@integer(0, 7)",
        StartFrame:        1000,
        EndFrame:          2500,
        CameraType:        "AURA",
        Algorithm:         "3D_AURA",
        VideoDir:          "/data/videos/record0001",
        OutputDir:         "/data/output/record0001",
        EnableBottom:      "1",
        EnableTop:         "1",
        Quality:           "8k",
        EanbleColorAdjust: "1",
        SaveDebugImg:      "@boolean",
        CreateAt: '@now',
        UpdateAt: '@now',
        Creator: {
            Username: '@cname'
        },
        Operator: {
            Username: '@cname'
        }
    }
  ],
  page: {
    total: 100,
    current: 1
  }
}))

let jobsListData = global[dataKey]

module.exports = {
  'GET /api/jobs' (req, res) {
    const page = qs.parse(req.query)
    const pageSize = page.pageSize || 10
    const currentPage = page.page || 1

    let data
    let newPage

    let newData = jobsListData.data.concat()

    if (page.field) {
      const d = newData.filter(function (item) {
        return item[page.field].indexOf(decodeURI(page.keyword)) > -1
      })

      data = d.slice((currentPage - 1) * pageSize, currentPage * pageSize)

      newPage = {
        current: currentPage * 1,
        total: d.length
      }
    } else {
      data = jobsListData.data.slice((currentPage - 1) * pageSize, currentPage * pageSize)
      jobsListData.page.current = currentPage * 1
      newPage = jobsListData.page
    }
    newPage.pageSize = Number(pageSize)
    res.json({success: true, data, page: newPage})
  },

  'POST /api/jobs' (req, res) {
    const newData = qs.parse(req.body)
    newData.CreateAt = Mock.mock('@now')
    newData.UpdateAt = Mock.mock('@now')

    newData.Id = jobsListData.data.length + 1
    jobsListData.data.unshift(newData)

    jobsListData.page.total = jobsListData.data.length
    jobsListData.page.current = 1

    global[dataKey] = jobsListData

    res.json({success: true, data: jobsListData.data, page: jobsListData.page})
  },

  'DELETE /api/jobs' (req, res) {
    const deleteItem = req.body

    jobsListData.data = jobsListData.data.filter(function (item) {
      if (item.id === deleteItem.id) {
        return false
      }
      return true
    })

    jobsListData.page.total = jobsListData.data.length

    global[dataKey] = jobsListData

    res.json({success: true, data: jobsListData.data, page: jobsListData.page})
  },

  'PUT /api/jobs' (req, res) {
    const editItem = req.body

    editItem.createTime = Mock.mock('@now')
    editItem.avatar = Mock.Random.image('100x100', Mock.Random.color(), '#757575', 'png', editItem.nickName.substr(0, 1))

    jobsListData.data = jobsListData.data.map(function (item) {
      if (item.id === editItem.id) {
        return editItem
      }
      return item
    })

    global[dataKey] = jobsListData
    res.json({success: true, data: jobsListData.data, page: jobsListData.page})
  }

}
