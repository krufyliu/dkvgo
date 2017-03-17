import React, { PropTypes } from 'react'
import { Form, Input, InputNumber, Radio, Modal } from 'antd'
const FormItem = Form.Item

const formItemLayout = {
  labelCol: {
    span: 6
  },
  wrapperCol: {
    span: 14
  }
}

const modal = ({
  visible,
  type,
  item = {},
  onOk,
  onCancel,
  form: {
    getFieldDecorator,
    validateFields,
    getFieldsValue
  }
}) => {
  function handleOk () {
    validateFields((errors) => {
      console.log("here")
      if (errors) {
        return
      }
      const data = {
        ...getFieldsValue(),
        key: item.key
      }
      onOk(data)
    })
  }

  const modalOpts = {
    title: `${type === 'create' ? '新建用户' : '修改用户'}`,
    visible,
    onOk: handleOk,
    onCancel,
    wrapClassName: 'vertical-center-modal'
  }

  return (
    <Modal {...modalOpts}>
      <Form layout="horizontal">
        <FormItem label='用户名：' hasFeedback {...formItemLayout}>
          {getFieldDecorator('Username', {
            initialValue: item.Username,
            rules: [
              {
                required: true,
                message: '用户名不能为空'
              }
            ]
          })(<Input />)}
        </FormItem>
        <FormItem label='邮箱：' hasFeedback {...formItemLayout}>
          {getFieldDecorator('Email', {
            initialValue: item.Email,
            rules: [
              {
                required: true,
                message: '邮箱不能为空'
              }
            ]
          })(<Input autoComplete='off'/>)}
        </FormItem>
        <FormItem label='密码：' hasFeedback {...formItemLayout}>
          {getFieldDecorator('Password', {
            initialValue: '',
            rules: [
              {
                required: true,
                min: 6,
                max: 20,
                message: '密码长度必须在6-20个字符'
              }
            ]
          })(<Input type="password"/>)}
        </FormItem>
      </Form>
    </Modal>
  )
}

modal.propTypes = {
  form: PropTypes.object.isRequired,
  visible: PropTypes.bool,
  type: PropTypes.string,
  item: PropTypes.object
}

export default Form.create()(modal)
