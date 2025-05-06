import { onMount, onCleanup, createEffect } from 'solid-js';
import * as echarts from 'echarts';

export type Point = {
  timestamp: string;
  value: number;
};

export type ChartData = Record<string, Point[]>;


// function filterLastMinutes(points: Point[], n: number): Point[] {
//   const now = Date.now();
//   const fiveMinutesAgo = now - n * 60 * 1000;
//   return points.filter(p => new Date(p.timestamp).getTime() >= fiveMinutesAgo);
// }

// function buildSeries(data: ChartData): Series[] {
//   return Object.entries(data).map(([name, points]) => ({
//     name,
//     type: 'line',
//     data: filterLastMinutes(points, 1)
//       .map((p) => [echarts.time.parse(p.timestamp), p.value] as [Date, number])
//       .sort((a, b) => new Date(a[0]!).getTime() - new Date(b[0]!).getTime())
//   }));
// }


interface ChartProps {
  name: string;
  data: Record<string, Point[]>;
}

export function Chart(props: ChartProps) {
  let chartRef: HTMLDivElement;
  let chartInstance: echarts.ECharts | undefined;

  onMount(() => {
    chartInstance = echarts.init(chartRef);

    window.addEventListener('resize', () => chartInstance?.resize());

    onCleanup(() => {
      chartInstance?.dispose();
    });
  });

  const updateChart = () => {
    if (!chartInstance) return;

    const mappedData: Record<number, number> = {};

    Object.values(props.data)
      .forEach((ps) =>
        ps.forEach((p) => mappedData[p.value] = (mappedData[p.value] || 0) + 1)
      );

    const uniqueVoteNumbers = Object.keys(mappedData).map((n) => parseInt(n));
    const seriesData = uniqueVoteNumbers.map((n) => mappedData[n]);

    chartInstance.setOption({
      title: {
        text: `Estimate ${props.name}!`,
        textStyle: {
          color: "#fafafa"
        },
      },
      xAxis: {
        max: 'dataMax',
        axisLabel: {
          show: false,
        },
        splitLine: {
          lineStyle: {
            width: 1,
            color: '#161619',
          },
        },
      },

      yAxis: {
        axisLabel: {
          color: '#818188',
        },
        axisTick: {
          show: false,
        },
        axisLine: {
          show: false,
        },
        type: 'category',
        inverse: true,
        data: uniqueVoteNumbers,
      },

      color: ['#2662d9', '#e23670', '#e88c30', '#af57db', '#2eb88a'],

      series: [{
        roundCap: true,
        colorBy: 'data',
        itemStyle: {
          borderCap: 'round',
        },
        label: {
          show: true,
          position: 'right',
          valueAnimation: true
        },
        realtimeSort: true,
        data: seriesData,
        type: 'bar',
      }],

      backgroundColor: '#09090b',
    });
  };

  createEffect(() => {
    updateChart();
  });

  return <div ref={el => (chartRef = el)} class="w-full h-96" />;
}
