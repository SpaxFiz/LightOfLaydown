import request from './request'

export function getFundInfo(code: number|string) {

  return request({
    'url': `/single_fund?code=${code}`,
    'method': 'get'
  }
  )
}
