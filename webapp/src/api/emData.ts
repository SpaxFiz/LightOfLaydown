import request from './request'

export function getEmAccountData(params: Object ) {
  return request({
    'url': '/em_account_data',
    'method': 'get'
  }
  )
}
