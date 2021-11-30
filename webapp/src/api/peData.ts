import request  from "./request"

export function getPEData(params: Object ) {
  return request({
    url: '/lg_pe_data',
    method: 'get',
    }
  )
}
