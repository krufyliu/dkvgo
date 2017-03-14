import React, {PropTypes} from 'react'
import {connect} from 'dva'
import {Row, Col, Card} from 'antd'

const bodyStyle = {
  bodyStyle: {
    height: 432,
    background: '#fff'
  }
}

function Dashboard ({dashboard}) {
  return (
    <Row gutter={24}>
      <p>Welcome</p>
    </Row>
  )
}

Dashboard.propTypes = {
  dashboard: PropTypes.object
}

export default connect(({dashboard}) => ({dashboard}))(Dashboard)
