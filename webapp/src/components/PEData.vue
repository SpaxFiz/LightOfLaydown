<script  lang="ts">
import { getPEData } from '@/api/peData'
import * as echarts from 'echarts'
import { defineComponent } from 'vue'

enum indexName {
  CSI1000 = '1000',
  CSI500 = '500',
  SHA50 = '50',
  CSI300 = '300',
}

enum seriesColor {
  CSI1000 = '#2db7f5',
  CSI500 = '#ff6600',
  SHA50 = '#808bc6',
  CSI300 = '#9ACC23',
}

enum mode {
  pe = 'pe',
  price = 'price',
  quantile = 'quantile',
}

enum title {
  pe = 'PE',
  quantile = 'PE分位数',
  price = '收盘价',
  weightedPE = '加权PE',
  weightedQuantile = '加权PE分位数',
}

const yAxis: echarts.YAXisComponentOption[] = [{
  id: 'PE',
  name: 'PE-TTM',
  min: 0,
  max: 150,
  position: 'right',
  axisLine: {
    show: 'auto',
    lineStyle: {
      color: '#188773'
    },
  }
}, {
  id: 'price',
  name: '收盘价',
  type: 'value',
  position: 'right',
  axisLine: {
    show: 'auto',
    lineStyle: {
      color: '#916400'
    }
  },
}, {
  id: 'quantile',
  name: 'PE分位',
  min: 0,
  max: 101,
  // offset: 80,
  position: 'right',
  axisLine: {
    show: 'auto',
    lineStyle: {
      color: '#a39552'
    }
  },
  axisLabel: {
    formatter: '{value}%'
  }
}]

declare type PEDataArray = Array<PEData>
declare type BackendData = Array<PEDataArray>
declare type PEArray = Array<(PEData | null)>
declare type SeriesData = Array<Array<number | null>>

type EChartsOption = echarts.EChartsOption

interface PEData {
  market_code: string,
  avg_pe_ttm: number,
  avg_pe_ttm_quantile: number,
  weighted_pe_ttm: number,
  weighted_pe_ttm_quantile: number,
  close_price: number,
  date: string,
  [param: string]: any,
}

interface GroupObject {
  pe300: PEArray,
  pe500: PEArray,
  pe50: PEArray,
  pe1000: PEArray,
  date: string[],
  series300?: SeriesData,
  series50?: SeriesData,
  series1000?: SeriesData,
  series500?: SeriesData,
}

interface iData {
  chartOptions: EChartsOption,
  mode: Object
  currentMode: mode,
  rawData: any,
  myChart?: echarts.ECharts,
  currentTitle: string
}

function padGroupObject(result: GroupObject) {
  result.pe500.push(null)
  result.pe1000.push(null)
  result.pe300.push(null)
  result.pe50.push(null)
}

function resolvePEData(data: PEArray): SeriesData {
  let result: SeriesData = [[], [], [], [], []]
  data.forEach(item => {
    if (item === null) {
      result[0].push(0)
      result[1].push(0)
      result[2].push(0)
      result[3].push(0)
      result[4].push(0)
    } else {
      result[0].push(item.avg_pe_ttm)
      result[1].push(parseFloat((item.avg_pe_ttm_quantile * 100).toFixed(2)))
      result[2].push(item.close_price)
      result[3].push(item.weighted_pe_ttm)
      result[4].push(parseFloat((item.weighted_pettm_quantile * 100).toFixed(2)))
    }
  })
  return result
}

function resolveGroupSeries(group: GroupObject) {
  group.series1000 = resolvePEData(group.pe1000)
  group.series300 = resolvePEData(group.pe300)
  group.series50 = resolvePEData(group.pe50)
  group.series500 = resolvePEData(group.pe500)
}

function genSeries(data: SeriesData, prefix: indexName, m: mode): any {
  let color: string = ''
  switch (prefix) {
    case indexName.CSI1000:
      color = seriesColor.CSI1000
      break;
    case indexName.CSI500:
      color = seriesColor.CSI500
      break;
    case indexName.CSI300:
      color = seriesColor.CSI300
      break;
    case indexName.SHA50:
      color = seriesColor.SHA50
      break;
    default:
      break;
  }
  switch (m) {
    case mode.pe:
      return [{
        name: prefix,
        data: data[0].map((item) => {
          return {
            value: item,
            tooltip: {
              formatter: prefix + '等权 '
            }
          }
        }),
        type: 'line',
        yAxisId: 'PE',
        smooth: true,
        lineStyle: {
          color: color,
          // type: 'dashed'
        },
        symbol: 'none'
      }, {
        name: prefix,
        data: data[3].map((item) => {
          return {
            value: item,
            tooltip: {
              formatter: prefix + '加权 '
            }
          }
        }),
        type: 'line',
        yAxisId: 'PE',
        smooth: true,
        lineStyle: {
          color: color,
          type: 'dashed'
        },
        symbol: 'none'
      }]
    case mode.quantile:
      return [{
        name: prefix,
        data: data[1].map((item) => {
          return {
            value: item,
            tooltip: {
              formatter: prefix + '等权(﹦) '
            }
          }
        }),
        type: 'line',
        yAxisId: 'quantile',
        lineStyle: {
          color: color,
          // type: 'dotted'
        },
        symbol: 'none'
      }, {
        name: prefix,
        data: data[4].map((item) => {
          return {
            value: item,
            tooltip: {
              formatter: prefix + '加权(➕) '
            }
          }
        }),
        type: 'line',
        yAxisId: 'quantile',
        lineStyle: {
          color: color,
          type: 'dashed'
        },
        symbol: 'none'
      }]
    case mode.price:
      return {
        name: prefix,
        data: data[2].map((item) => {
          return {
            value: item,
            tooltip: {
              formatter: prefix + '收盘价 '
            }
          }
        }),
        type: 'line',
        yAxisId: 'price',
        lineStyle: {
          color: color
        },
        symbol: 'none'
      }
  }
}

function prepareLegend(): Array<any> {
  let result: echarts.LegendComponentOption[] = []
  let cat = ['收盘价', 'PE分位数', 'PE']
  let o = Object.assign(indexName)
  let color = Object.assign(seriesColor)
  cat.forEach((c: string) => {
    Object.keys(o).forEach((key: string) => {
      result.push({
        data: [{
          name: o[key] + c,
          textStyle: {
            color: color[key]
          },
          icon: 'circle',
        }]
      })
    })
  })
  return result
}

function prepareOptions(data: BackendData): GroupObject {
  let result: GroupObject = {
    pe300: [],
    pe500: [],
    pe50: [],
    pe1000: [],
    date: []
  }
  for (let i = 0; i < data.length; i++) {
    let group = data[i]
    result.date.push(group[0].date)
    padGroupObject(result)
    group.forEach(item => {
      switch (item.market_code) {
        case 'PE50':
          result.pe50[i] = item
          break
        case 'PE500':
          result.pe500[i] = item
          break
        case 'PE1000':
          result.pe1000[i] = item
          break
        case 'PE300':
          result.pe300[i] = item
          break
        default:
          break
      }
    })
  }
  return result
}

export default defineComponent({
  methods: {
    setup() {
      getPEData({})
        .then(res => {
          this.rawData = JSON.parse(res.data)
          this.chartOptions = this.constructData(this.rawData)
          this.render()
        })
        .catch(err => {
          console.log(err)
        })
    },
    constructData(data: BackendData): EChartsOption {
      let result = prepareOptions(data)

      resolveGroupSeries(result)
      let series: echarts.SeriesOption[] = []
      series = series.concat(genSeries(result.series1000 || [], indexName.CSI1000, this.currentMode))
      series = series.concat(genSeries(result.series50 || [], indexName.SHA50, this.currentMode))
      series = series.concat(genSeries(result.series300 || [], indexName.CSI300, this.currentMode))
      series = series.concat(genSeries(result.series500 || [], indexName.CSI500, this.currentMode))

      let yAxisOpt: echarts.YAXisComponentOption = {}
      switch (this.currentMode) {
        case mode.pe:
          yAxisOpt = yAxis[0]
          break;
        case mode.quantile:
          yAxisOpt = yAxis[2]
          break;
        case mode.price:
          yAxisOpt = yAxis[1]
          break;
      }

      let options: EChartsOption = {
        title: {
          text: this.currentTitle,
          left: '9%',
        },
        toolbox: {
          feature: {
            dataZoom: {
              yAxisIndex: 'none'
            },
            restore: {}
          }
        },
        dataZoom: [
          {
            type: 'inside',
            start: 0,
            end: 100
          },
          {
            start: 0,
            end: 100
          }
        ],
        tooltip: {
          trigger: 'axis',
          formatter: function (params: any): string {
            return [params[0]['axisValue']].concat(params.map((item: any) => {
              let color = ''
              switch (item['seriesName']) {
                case '1000':
                  color = seriesColor.CSI1000
                  break
                case '500':
                  color = seriesColor.CSI500
                  break
                case '50':
                  color = seriesColor.SHA50
                  break
                case '300':
                  color = seriesColor.CSI300
                  break
              }
              return '<span style="color:' + color + '">' + item['data']['tooltip']['formatter'] + item['data']['value'] + '</span>'
            })).join('<br>')

          }
        },
        // color: [seriesColor.CSI1000, seriesColor.SHA50, seriesColor.CSI300, seriesColor.CSI500],
        legend: {
          align: 'left',
          right: '10%',
          data: [
            { name: indexName.CSI1000, itemStyle: { color: seriesColor.CSI1000 } },
            { name: indexName.CSI300, itemStyle: { color: seriesColor.CSI300 } },
            { name: indexName.CSI500, itemStyle: { color: seriesColor.CSI500 } },
            { name: indexName.SHA50, itemStyle: { color: seriesColor.SHA50 } },
          ]
        },
        xAxis: {
          splitLine: { show: false },
          type: 'category',
          data: result.date,
          axisTick: {
            // alignWithLabel: true
          },
        },
        yAxis: yAxisOpt,
        series: series
      }
      return options
    },
    render() {
      var chartDom = document.getElementById('chartPE')!
      this.myChart = echarts.init(chartDom)
      this.myChart.setOption(this.chartOptions)
      var c = this.myChart
      window.onresize = function () {
        c.resize()
      }
    },
    switchMode(m: string) {
      switch (m) {
        case 'pe':
          this.currentMode = mode.pe
          this.currentTitle = title.pe
          break;
        case 'price':
          this.currentMode = mode.price
          this.currentTitle = title.price
          break;
        case 'quantile':
          this.currentTitle = title.quantile
          this.currentMode = mode.quantile
          break;
        default:
          break;
      }
      var newOption = this.constructData(this.rawData)
      this.myChart?.setOption(newOption, true, false)
    },
  },
  mounted() {
    this.setup()
  },
  data(): iData {
    return {
      chartOptions: {},
      currentMode: mode.pe,
      mode: mode,
      rawData: {},
      currentTitle: title.pe,
    }
  }
})

</script>

<template>
  <el-button-group>
    <el-button type="primary" size="small" @click="switchMode('pe')">PE</el-button>
    <el-button type="primary" size="small" @click="switchMode('price')">收盘价</el-button>
    <el-button type="primary" size="small" @click="switchMode('quantile')">分位数</el-button>
  </el-button-group>
  <div id="chartPE"></div>
</template>

<style scoped>
#chartPE {
  height: 400px;
}
</style>
