// http.ts
import axios, { AxiosRequestConfig, AxiosResponse } from 'axios'

const showStatus = (status: number) => {
    let message = ''

    switch (status) {
    case 400:
      message = '请求错误(400)'
      break
    case 401:
      message = '未授权，请重新登录(401)'
      break
    case 403:
      message = '拒绝访问(403)'
      break
    case 404:
      message = '请求出错(404)'
      break
    case 408:
      message = '请求超时(408)'
      break
    case 500:
      message = '服务器错误(500)'
      break
    case 501:
      message = '服务未实现(501)'
      break
    case 502:
      message = '网络错误(502)'
      break
    case 503:
      message = '服务不可用(503)'
      break
    case 504:
      message = '网络超时(504)'
      break
    case 505:
      message = 'HTTP版本不受支持(505)'
      break
    default:
      message = `连接出错(${status})!`
    }
    return `${message}`
  },

  service = axios.create({
    'baseURL': '/api',
    'headers': {
      'Content-Type': 'application/json;charset=utf-8'
    },
    'withCredentials': false,
    'timeout': 30000,
    'transformRequest': [(data) => {
      data = JSON.stringify(data)
      return data
    }],
    validateStatus() {
      return true
    },
    'transformResponse': [(data) => {
      if (typeof data === 'string' && data.startsWith('{')) {
        data = JSON.parse(data)
      }
      return data
    }]
  })

// service.interceptors.request.use((config: AxiosRequestConfig) => {
//   let token = localStorage.getItem('token')
//   if(token){
//     config.headers.Authorization = `${token}`;
//   }
//   return config
// }, (error) => {
//   error.data = {}
//   error.data.msg = 'internal error'
//   return Promise.resolve(error)
// })

service.interceptors.response.use((response: AxiosResponse) => {
  const status = response.status
  let msg = ''

  if (status < 200 || status >= 300) {

    msg = showStatus(status)
    if (typeof response.data === 'string') {
      response.data = { msg }
    } else {
      response.data.msg = msg
    }
  }
  return response
}, (error) => {
  if (axios.isCancel(error)) {
    console.log('repeated request: ' + error.message)
  } else {
    console.log(error)
  }
  return Promise.reject(error)
})

export default service
