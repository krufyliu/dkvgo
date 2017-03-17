import { request } from '../utils'

export async function login (params) {
  return request('/api/auth', {
    method: 'post',
    data: params
  })
}

export async function logout (params) {
  return request('/api/auth', {
    method: 'delete',
    data: params
  })
}

export async function userInfo (params) {
  return request('/api/auth', {
    method: 'get',
    data: params
  })
}
