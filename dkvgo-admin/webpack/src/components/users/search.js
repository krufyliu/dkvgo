import React, { PropTypes } from 'react'
import { Form, Button, Row, Col, Switch } from 'antd'
import SearchGroup from '../ui/search'

const search = ({
  field,
  keyword,
  onSearch,
  onAdd,
  isMotion,
  switchIsMotion,
  form: {
    getFieldDecorator,
    validateFields,
    getFieldsValue
  }
}) => {
  const searchGroupProps = {
    field,
    keyword,
    size: 'large',
    select: true,
    selectOptions: [{ value: 'username', name: '用户名' }, { value: 'email', name: '邮箱地址' }],
    selectProps: {
      defaultValue: field || '用户名'
    },
    onSearch: (value) => {
      onSearch(value)
    }
  }

  return (
    <Row gutter={24}>
      <Col lg={8} md={12} sm={16} xs={24} style={{marginBottom: 16}}>
        <SearchGroup {...searchGroupProps} />
      </Col>
      <Col lg={{offset: 8, span: 8}} md={12} sm={8} xs={24} style={{marginBottom: 16, textAlign: 'right'}}>
        <Button size='large' type='ghost' onClick={onAdd}>添加</Button>
      </Col>
    </Row>
  )
}

search.propTypes = {
  form: PropTypes.object.isRequired,
  onSearch: PropTypes.func,
  onAdd: PropTypes.func,
  field: PropTypes.string,
  keyword: PropTypes.string,
  isMotion: PropTypes.bool,
  switchIsMotion: PropTypes.func
}

export default Form.create()(search)
