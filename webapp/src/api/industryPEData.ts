import request from './request'

export function getIndustryPEData() {
  return request({
    'url': '/industry_pe_data',
    'method': 'get'
  }
  )
}
