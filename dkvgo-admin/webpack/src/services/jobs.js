import { request } from '../utils'

export async function query (params) {
  return request('/api/jobs', {
    method: 'get',
    data: params
  })
}

export async function create (params) {
  return request('/api/jobs', {
    method: 'post',
    data: params
  })
}

export async function stop(params) {
  return request('/api/jobs/' + params.id + '/action/stop', {
    method: 'post'
  })  
}

export async function resume(params) {
  return request('/api/jobs/' + params.id + '/action/resume', {
    method: 'post'
  })  
}

export async function remove (params) {
  return request('/api/jobs/' + params.id, {
    method: 'delete',
  })
}

export async function update (params) {
  return request('/api/jobs/' + params.id, {
    method: 'put',
  })
}
