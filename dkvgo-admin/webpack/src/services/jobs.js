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

export async function remove (params) {
  return request('/api/jobs', {
    method: 'delete',
    data: params
  })
}

export async function update (params) {
  return request('/api/jobs', {
    method: 'put',
    data: params
  })
}
