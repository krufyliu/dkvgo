import React, { PropTypes } from 'react'
import { Form, Input, InputNumber, Radio, Modal,Row, Col, Button, Tooltip } from 'antd'

const FormItem = Form.Item
const RadioGroup = Radio.Group

const formItemLayout = {
  labelCol: {
    span: 5
  },
  wrapperCol: {
    span: 17
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
    validateFieldsAndScroll,
    getFieldsValue,
    getFieldValue,
    setFields
  }
}) => {
  function handleOk (e) {
    e.preventDefault()
    validateFieldsAndScroll((errors) => {
      if (errors) {
        return
      }
      var data = {
        ...getFieldsValue(),
        key: item.key
      }
      data.EnableBottom = data.EnableTop
      onOk(data)
    })
  }

  function checkStartEnd(rule, value, callback) {
    var startValue = getFieldValue('StartFrame')
    var endValue = getFieldValue('EndFrame')
    if (startValue && endValue) {
      if (parseInt(startValue) >= parseInt(endValue)) {
        callback('结束帧不能小于开始帧')
      }
    }
    callback()
  }

  function handleSwitch(e) {
    var videoPath = getFieldValue('VideoDir')
    if (typeof videoPath == 'undefined' || videoPath == "" ) {
      return 
    }
    const originVal = '',
          prePath = '\/dkvision',
          windowsPath = new RegExp(/\w:+?\\+?([^\:\?"\<\>\|\.\\\/]+?[\\\b]{0,1})+?/);
    if (windowsPath.test(videoPath)) {
      videoPath = videoPath.slice(videoPath.indexOf(':') + 1)
    } 
    videoPath = videoPath.split(/[\\/]/).length > 2 ? videoPath .split(/[\\/]/).join('\/') : videoPath
    videoPath = videoPath.indexOf(prePath) == 0 ? videoPath : prePath + videoPath
    setFields({
      VideoDir: { value: videoPath },
      OutputDir: { value:videoPath.replace('/data/', '/output/') }
    })
  }

  const modalOpts = {
    title: `${type === 'create' ? '新建任务' : '修改任务'}`,
    visible,
    onOk: handleOk,
    onCancel,
    width: 600,
    wrapClassName: 'vertical-center-modal'
  }

  return (
    <Modal {...modalOpts}>
      <Form layout="horizontal">
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
        <FormItem label='视频路径：' {...formItemLayout}>
          <Row>
            <Col span={22}>
            {getFieldDecorator('VideoDir', {
              initialValue: item.VideoDir,
              rules: [
                {
                  required: true,
                  message: '视频路径不能为空'
                }
              ]
            })(<Input size="large"/>)}
            </Col>
            <Col span={2}>
              <Tooltip title="Windows路径转换">
                <Button icon="retweet" size="large" onClick={handleSwitch}/>
              </Tooltip>
            </Col>
          </Row>
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
        <FormItem label='起始帧：' {...formItemLayout}>
          <Col span="6">
            <FormItem  {...formItemLayout}>
              {getFieldDecorator('StartFrame', {
                initialValue: item.StartFrame,
                rules: [
                  {
                    required: true,
                    message: '开始帧不能为空'
                  }, {
                    validator: checkStartEnd 
                  }
                ]
              })(<InputNumber min={0} autoComplete='off' />)}
            </FormItem>
          </Col>
          <Col span="1">
            <p className="ant-form-split">-</p>
          </Col>
          <Col span="6" push="1">
            <FormItem {...formItemLayout}>
              {getFieldDecorator('EndFrame', {
                initialValue: item.EndFrame,
                rules: [
                  {
                    required: true,
                    message: '结束帧不能为空'
                  }, {
                    validator: checkStartEnd
                  }
                ]
              })(<InputNumber min={0} autoComplete='off' />)}
            </FormItem>
          </Col>
        </FormItem>
        <FormItem {...formItemLayout} label="相机类型">
          {getFieldDecorator('CameraType', { initialValue: item.CameraType || 'AURA'})(
            <RadioGroup>
              <Radio value="AURA">AURA</Radio>
              <Radio value="GOPRO">GOPRO</Radio>
            </RadioGroup>
          )}
        </FormItem>
        <FormItem {...formItemLayout} label="分辨率">
          {getFieldDecorator('Quality', { initialValue: item.Quality || '8k' })(
            <RadioGroup>
              <Radio value="8k">8K</Radio>
              <Radio value="6k">6K</Radio>
              <Radio value="4k">4K</Radio>
            </RadioGroup>
          )}
        </FormItem>
        <FormItem {...formItemLayout} label='生成顶底'>
          {getFieldDecorator('EnableTop', { initialValue: item.EnableTop || '1' })(
            <RadioGroup>
              <Radio value="1">是</Radio>
              <Radio value="0">否</Radio>
            </RadioGroup>
          )}
        </FormItem>
        <FormItem {...formItemLayout} label='保留调试图片'>
          {getFieldDecorator('SaveDebugImg', { initialValue: item.SaveDebugImg || 'false' })(
            <RadioGroup>
              <Radio value="true">是</Radio>
              <Radio value="false">否</Radio>
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
