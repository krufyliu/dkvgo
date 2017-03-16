module.exports = {
  getJobStatus: function(status) {
    switch(status) {
    case 0:
        return { 'status': "default", 'text':"等待中" };
    case 1:
    case 2:
        return { 'status': "processing", 'text':"运行中" };
    case 3:
        return { 'status': "error", 'text':"终止中" };
    case 4:
        return { 'status': "error", 'text':"已终止" };
    case 5:
        return { 'status': "success", 'text':"完成" };
    case 6:
        return { 'status': "error", 'text':"失败" };
    default:
        return { 'status': "warning", 'text':"未知" };
    }
  },
  getProccessStatus: function(status) {
    switch(status) {
    case 0:
        return { 'status': "normal", 'text':"等待中" };
    case 1:
    case 2:
        return { 'status': "active", 'text':"运行中" };
    case 3:
        return { 'status': "normal", 'text':"终止中" };
    case 4:
        return { 'status': "normal", 'text':"已终止" };
    case 5:
        return { 'status': "success", 'text':"完成" };
    case 6:
        return { 'status': "exception", 'text':"失败" };
    default:
        return { 'status': "exception", 'text':"未知" };
    }
  },
  getActions: function(status) {
    var actions = {
        detail: { 'key': 'detail', 'text':'详情' },
        stop: { 'key': 'stop', 'text':'终止' },
        resume: { 'key': 'resume', 'text':'继续' },
    }
    var targetActions = [actions['detail']]
    switch(status) {
    case 0:
    case 1:
    case 2:
        targetActions.push(actions['stop'])
        return targetActions
    case 4:
    case 6:
        targetActions.push(actions['resume'])
        return targetActions
    default:
        return targetActions
    }

  }
}