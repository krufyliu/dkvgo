import { create, stop, resume, remove, update, query } from '../services/jobs'
import { parse } from 'qs'

export default {
  namespace: 'jobs',

  state: {
    list: [],
    currentItem: {},
    modalVisible: false,
    modalType: 'create',
    isMotion: false,
    pagination: {
      showSizeChanger: true,
      showQuickJumper: true,
      showTotal: total => `共 ${total} 条`,
      current: 1,
      pageSize: 10,
      total: null
    }
  },

  subscriptions: {
    setup ({ dispatch, history}) {
      var timer = null
      function refresh(location) {
        timer = setInterval(() => {
          console.log('refresh')
          // if (jobs != null) {
          //   const filtered = jobs.list.filter((item) => item.Status == 1 || item.Status == 2)
          //   if (filtered.length === 0) {
          //     clearInterval(timer)
          //     timer = null
          //   }
          // }
          dispatch({
            type: 'query',
            payload: location.query
          })
        }, 30000)
      }
      history.listen(location => {
        if (location.pathname === '/jobs') {
          dispatch({
            type: 'query',
            payload: location.query
          })
          // refresh(location)
        } else {
          if (timer != null) {
            clearInterval(timer)
          }
        }
      })
    }
  },

  effects: {
    *query ({ payload }, { call, put}) {
      const data = yield call(query, parse(payload))
      if (data) {
        yield put({
          type: 'querySuccess',
          payload: {
            list: data.data,
            pagination: data.page
          }
        })
      }
    },
    *'delete' ({ payload }, { call, put }) {
      const data = yield call(remove, { id: payload })
      if (data && data.success) {
        yield put({
          type: 'querySuccess',
          payload: {
            list: data.data,
            pagination: {
              total: data.page.total,
              current: data.page.current
            }
          }
        })
      }
    },
    *create ({ payload }, { call, put, select }) {
      yield put({ type: 'hideModal' })
      const data = yield call(create, payload)
      if (data && data.success) {
        window.location.reload()
      }
    },
    *update ({ payload }, { select, call}) {
      yield put({ type: 'hideModal' })
      const id = yield select(({ jobs }) => jobs.currentItem.Id)
      const data = yield call(update, newUser)
      if (data && data.success) {
        yield put({
          type: 'updateSuccess',
          payload: {
            job: data.data,
          }
        })
      }
    },
    *stop({payload}, {call, put}) {
      const data = yield call(stop, {id: payload})      
      if (data && data.success) {
        yield put({
          type: 'updateSuccess',
          payload: {
            job: data.data
          }
        })
      }
    },
    *resume({payload}, {call, put}) {
      const data = yield call(resume, {id: payload})      
      if (data && data.success) {
        yield put({
          type: 'updateSuccess',
          payload: {
            job: data.data
          }
        })
      }
    },
    *switchIsMotion ({
      payload
    }, {put}) {
      yield put({
        type: 'handleSwitchIsMotion'
      })
    }
  },

  reducers: {
    querySuccess (state, action) {
      const {list, pagination} = action.payload
      return { ...state,
        list,
        pagination: {
          ...state.pagination,
          ...pagination
        }}
    },
    createSuccess (state, action) {
      const job = action.payload.job
      list = [...state.list]
      list.unshift(job)
      return { ...state,
        list,
        pagination: {
          ...state.pagination,
          total: state.pagination.total + 1
        }
      }
    },
    updateSuccess(state, action) {
      console.log(action.payload)
      const job = action.payload.job
      const list = state.list.map((item) => {
        if (item.Id == job.Id) {
          return job
        }
        return item
      }) 
      return { ...state,
        list,
      }
    },
    showModal (state, action) {
      return { ...state, ...action.payload, modalVisible: true }
    },
    hideModal (state) {
      return { ...state, modalVisible: false }
    },
    handleSwitchIsMotion (state) {
      localStorage.setItem('antdAdminUserIsMotion', !state.isMotion)
      return { ...state, isMotion: !state.isMotion }
    }
  }

}
