
import request from './request'

export function checkCipher(params: Record<string, string> ) {
  return request({
    'url': '/cipher',
    'method': 'post',
    'data': params
  }
  )
}
