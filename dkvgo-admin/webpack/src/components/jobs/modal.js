import React, { PropTypes } from 'react'
import { Form, Input, InputNumber, Radio, Modal } from 'antd'

const FormItem = Form.Item
const RadioGroup = Radio.Group

const formItemLayout = {
  labelCol: {
    span: 4
  },
  wrapperCol: {
    span: 18
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
    title: `${type === 'create' ? '新建任务' : '修改用户'}`,
    visible,
    onOk: handleOk,
    onCancel,
    width: 600,
    wrapClassName: 'vertical-center-modal'
  }

  return (
    <Modal {...modalOpts}>
      <Form horizontal>
        <FormItem label='作业名：' hasFeedback {...formItemLayout}>
          {getFieldDecorator('Name', {
            initialValue: item.Name,
            rules: [
              {
                required: true,
                message: '作业名不能为空'
              }
            ]
          })(<Input />)}
        </FormItem>
        <FormItem label='视频路径：' hasFeedback {...formItemLayout}>
          {getFieldDecorator('VideoDir', {
            initialValue: item.VideoDir,
            rules: [
              {
                required: true,
                message: '视频路径不能为空'
              }
            ]
          })(<Input autoComplete='off'/>)}
        </FormItem>
        <FormItem label='输出路径：' hasFeedback {...formItemLayout}>
          {getFieldDecorator('OutputDir', {
            initialValue: item.OutputDir,
            rules: [
              {
                required: true,
                message: '输出路径不能为空'
              }
            ]
          })(<Input autoComplete='off'/>)}
        </FormItem>
        <FormItem label='开始帧：' hasFeedback {...formItemLayout}>
          {getFieldDecorator('StartFrame', {
            initialValue: item.StartFrame,
            rules: [
              {
                required: true,
                message: '开始帧不能为空'
              }
            ]
          })(<InputNumber autoComplete='off'/>)}
        </FormItem>
        <FormItem label='结束帧：' hasFeedback {...formItemLayout}>
          {getFieldDecorator('OutputDir', {
            initialValue: item.EndFrame,
            rules: [
              {
                required: true,
                message: '结束帧不能为空'
              }
            ]
          })(<InputNumber autoComplete='off'/>)}
        </FormItem>
        <FormItem {...formItemLayout} label="相机类型">
          {getFieldDecorator('CameraType', { initialValue: item.CameraType || 'AURA'})(
            <RadioGroup>
              <Radio value="AURA">AURA</Radio>
              <Radio value="GOPRO">GOPRO</Radio>
            </RadioGroup>
          )}
        </FormItem>
        <FormItem {...formItemLayout} label="质量">
          {getFieldDecorator('Quality', { initialValue: item.Quality || '8k' })(
            <RadioGroup>
              <Radio value="8k">8K</Radio>
              <Radio value="6k">6K</Radio>
              <Radio value="4k">4K</Radio>
            </RadioGroup>
          )}
        </FormItem>
        <FormItem {...formItemLayout} label="生成顶">
          {getFieldDecorator('EnableTop', { initialValue: item.EnableTop || '1' })(
            <RadioGroup>
              <Radio value="1">是</Radio>
              <Radio value="0">否</Radio>
            </RadioGroup>
          )}
        </FormItem>
        <FormItem {...formItemLayout} label="生成底">
          {getFieldDecorator('EnableBottom', { initialValue: item.EnableBottom || '1' })(
            <RadioGroup>
              <Radio value="1">是</Radio>
              <Radio value="0">否</Radio>
            </RadioGroup>
          )}
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
