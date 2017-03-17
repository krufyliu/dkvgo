import React, { PropTypes } from 'react'
import { routerRedux } from 'dva/router'
import { connect } from 'dva'
import JobList from '../components/jobs/list'
import JobModal from '../components/jobs/modal'
import JobSearch from '../components/jobs/search'

function Jobs ({ loading, location, dispatch, jobs }) {
  const { list, pagination, currentItem, modalVisible, modalType, isMotion } = jobs
  const { field, keyword } = location.query

  const jobModalProps = {
    item: modalType === 'create' ? {} : currentItem,
    type: modalType,
    visible: modalVisible,
    onOk (data) {
      dispatch({
        type: `jobs/${modalType}`,
        payload: data
      })
    },
    onCancel () {
      dispatch({
        type: 'jobs/hideModal'
      })
    }
  }

  const jobListProps = {
    dataSource: list,
    loading: loading.global,
    pagination: pagination,
    location,
    isMotion,
    onPageChange (page) {
      const { query, pathname } = location
      dispatch(routerRedux.push({
        pathname: pathname,
        query: {
          ...query,
          page: page.current,
          pageSize: page.pageSize
        }
      }))
    },
    onStop(id) {
      dispatch({
        type: 'jobs/stop',
        payload: id
      })
    },
    onResume(id) {
      dispatch({
        type: 'jobs/resume',
        payload: id
      })
    },
    onDeleteItem (id) {
      dispatch({
        type: 'jobs/delete',
        payload: id
      })
    },
    onEditItem (item) {
      dispatch({
        type: 'jobs/showModal',
        payload: {
          modalType: 'update',
          currentItem: item
        }
      })
    }
  }

  const jobSearchProps = {
    field,
    keyword,
    isMotion,
    onSearch (fieldsValue) {
      fieldsValue.keyword.length ? dispatch(routerRedux.push({
        pathname: '/jobs',
        query: {
          field: fieldsValue.field,
          keyword: fieldsValue.keyword
        }
      })) : dispatch(routerRedux.push({
        pathname: '/jobs'
      }))
    },
    onAdd () {
      dispatch({
        type: 'jobs/showModal',
        payload: {
          modalType: 'create'
        }
      })
    }
  }

  const JobModalGen = () =>
    <JobModal {...jobModalProps} />

  return (
    <div className='content-inner'>
      <JobSearch {...jobSearchProps} />
      <JobList {...jobListProps} />
      <JobModalGen />
    </div>
  )
}

Jobs.propTypes = {
  users: PropTypes.object,
  location: PropTypes.object,
  dispatch: PropTypes.func
}

export default connect(({jobs, loading}) => ({jobs, loading}))(Jobs)
