import React, {PropTypes} from 'react'
import {Table, Dropdown, Button, Menu, Icon, Modal, Badge, Progress} from 'antd'
import styles from './list.less'
import classnames from 'classnames'
import TableBodyWrapper from '../common/TableBodyWrapper'
import {getJobStatus, getProccessStatus} from '../../utils/jobStatus'

const confirm = Modal.confirm

function list ({ loading, dataSource, pagination, onPageChange, onDeleteItem, onEditItem, isMotion, location }) {
  const handleMenuClick = (record, e) => {
    if (e.key === '1') {
      onEditItem(record)
    } else if (e.key === '2') {
      confirm({
        title: '您确定要删除这条记录吗?',
        onOk () {
          onDeleteItem(record.id)
        }
      })
    }
  }

  const columns = [
    {
      title: '任务名',
      dataIndex: 'Name',
      key: 'Name'
    },
    {
      title: '源路径',
      dataIndex: 'VideoDir',
      key: 'VideoDir'
    },
    {
      title: '相机类型',
      dataIndex: 'CameraType',
      key: 'CameraType'
    },
    {
      title: '分辨率',
      dataIndex: 'Quality',
      key: 'Quality'
    },
    {
      title: '创建时间',
      dataIndex: 'CreateAt',
      key: 'CreateAt'
    },
    {
      title: '最近操作人',
      dataIndex: 'Operator.Username',
      key: 'Operator.Username'
    },
    {
      title: '进度',
      dataIndex: 'Progress',
      key: 'Progress',
      render: (text, record, index) => {
          const processStatus = getProccessStatus(record.Status)
          return (
              <Progress percent={Math.round(record.Progress)} strokeWidth={5} status={processStatus.status} />
          )
      }
    },
    {
      title: '状态',
      dataIndex: 'Status',
      key: 'Status',
      render: (text, record, index) => {
        const jobStatus = getJobStatus(record.Status);
        return (
            <Badge status={jobStatus.status} text={jobStatus.text} />
        )
      },
    },
    {
      title: '操作',
      key: 'operation',
      width: 100,
      render: (text, record) => {
        return (<Dropdown overlay={<Menu onClick={(e) => handleMenuClick(record, e)}>
          <Menu.Item key='1'>编辑</Menu.Item>
          <Menu.Item key='2'>删除</Menu.Item>
        </Menu>}>
          <Button style={{ border: 'none' }}>
            <Icon style={{ marginRight: 2 }} type='bars' />
            <Icon type='down' />
          </Button>
        </Dropdown>)
      }
    }
  ]

  const getBodyWrapperProps = {
    page: location.query.page,
    current: pagination.current
  }

  const getBodyWrapper = body => isMotion ? <TableBodyWrapper {...getBodyWrapperProps} body={body} /> : body

  return (
    <div>
      <Table
        className={classnames({[styles.table]: true, [styles.motion]: isMotion})}
        bordered
        scroll={{ x: 1200 }}
        columns={columns}
        dataSource={dataSource}
        loading={loading}
        onChange={onPageChange}
        pagination={pagination}
        simple
        rowKey={record => record.Id}
        getBodyWrapper={getBodyWrapper}
      />
    </div>
  )
}

list.propTypes = {
  loading: PropTypes.bool,
  dataSource: PropTypes.array,
  pagination: PropTypes.object,
  onPageChange: PropTypes.func,
  onDeleteItem: PropTypes.func,
  onEditItem: PropTypes.func,
  isMotion: PropTypes.bool,
  location: PropTypes.object
}

export default list
