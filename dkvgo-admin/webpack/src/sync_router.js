import React from 'react'
import {Router} from 'dva/router'
import App from './routes/app'
import Dashboard from './routes/dashboard'
import Users from './routes/users'
import Jobs from './routes/jobs'
import Err from './routes/error'

export default function ({history, app}) {
  app.model(require('./models/dashboard'))
  app.model(require('./models/users'))
  app.model(require('./models/jobs'))
  const routes = [
    {
      path: '/',
      component: App,
      indexRoute: {component: Dashboard},
      childRoutes: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: Dashboard
        }, {
          path: 'users',
          name: 'users',
          component: Users
        }, {
          path: 'jobs',
          name: 'jobs',
          component: Jobs
        }, {
          path: '*',
          name: 'error',
          component: Err
        }
      ]
    }
  ]

  return <Router history={history} routes={routes} />
}
