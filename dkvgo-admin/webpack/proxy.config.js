const mock = {}

// require('fs').readdirSync(require('path').join(`${__dirname}/mock`)).forEach((file) => {
//   Object.assign(mock, require('./mock/' + file))
// })
mock['GET /api/*'] = 'http://localhost:8080'
mock['POST /api/*'] = 'http://localhost:8080'
mock['DELETE /api/*'] = 'http://localhost:8080'
mock['PUT /api/*'] = 'http://localhost:8080'
mock['GET /ajax/user/:id'] = function(req, res) {
  res.json(req.params)
}

module.exports = mock
