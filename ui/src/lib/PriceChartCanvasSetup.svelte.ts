/* eslint-disable @typescript-eslint/no-unused-vars */
import { formatPrice } from '$lib/calcUtils';
import { formatDateDay, formatDateTime, isNow } from '$lib/dateUtils';
import type { PriceEntry } from '$lib/pricesApi';
import {
	LineController,
	LineElement,
	PointElement,
	CategoryScale,
	Chart,
	Colors,
	Legend,
	LinearScale,
	Tooltip,
	Filler,
	type ChartConfiguration,
	type ChartTypeRegistry,
	BarController,
	BarElement,
} from 'chart.js';
import ChartDataLabels from 'chartjs-plugin-datalabels';
import type { Context } from 'chartjs-plugin-datalabels';
import type {
	ScriptableLineSegmentContext,
	ScriptableContext,
	ChartDataset,
	LegendItem,
	TooltipItem,
} from 'chart.js';

const defaultColor = 'rgba(0, 200, 0, 0.5)';
const blueColor = 'rgba(0, 0, 200,0.5)';
const redColor = 'rgba(200, 0, 0,0.5)';
const greenColor = 'rgba(0, 200, 0,0.5)';
const transparentColor = 'rgba(54, 162, 235,0)';
const blackColor = 'rgba(100, 100, 100,0.5)';

function chooseColor(p: PriceEntry) {
	if (p.p > 100) {
		return redColor;
	} else if (p.p < 0) {
		return blueColor;
	} else if (p.p < 65) {
		return greenColor;
	}
	return defaultColor;
}

export type MyChartConfig = ChartConfiguration<'line' | 'bar', (number | string | null)[], string>;

export function chartConfig(prices: PriceEntry[], showOnlyAfterNow: boolean): MyChartConfig {
	const biggestTemp = formatPrice(Math.max(...prices.map((p) => p.p)));
	const biggest = +biggestTemp < 10 ? 10 : +biggestTemp;
	const pricesWithoutLast = prices.slice(0, prices.length - 2);
	const config = {
		type: 'line',
		plugins: [ChartDataLabels],
		options: {
			interaction: {
				mode: 'index',
				intersect: false,
			},
			animation: {
				duration: 1,
			},
			scales: {
				y: {
					beginAtZero: true,
				},
			},
			plugins: {
				datalabels: {
					labels: {
						title: null,
					},
				},
				legend: {
					labels: {
						filter(item: LegendItem) {
							return item.text === 'c/kWh';
						},
					},
				},
				tooltip: {
					filter(e: TooltipItem<'line' | 'bar'>) {
						return e.dataset.label === 'c/kWh';
					},
					displayColors: false,
					callbacks: {
						title(tooltipItems: TooltipItem<'line' | 'bar'>[]) {
							const p = prices[tooltipItems[0].dataIndex];
							return formatDateTime(p.s) + ' - ' + formatDateTime(p.e);
						},
						label(tooltipItem: TooltipItem<'line' | 'bar'>) {
							return formatPrice(prices[tooltipItem.dataIndex].p) + ' c/kWh';
						},
					},
				},
			},
		},
		data: {
			labels: prices.map((p) => formatDateTime(p.s)),
			datasets: [
				{
					label: 'c/kWh',
					type: 'line',
					data: prices.map((p) => +formatPrice(p.p)),
					borderColor: (ctx: ScriptableContext<'line' | 'bar'>) => {
						const segmentCtx = ctx as unknown as ScriptableLineSegmentContext;
						const p = prices[segmentCtx.p0DataIndex];
						if (p) return chooseColor(p);
						return 'rgba(0, 200, 0, 1)';
					},
					backgroundColor: 'rgba(0, 200, 0, 0.3)',
					borderWidth: 3,
					fill: 'origin',
					stepped: 'middle',
					pointRadius: 0,
					pointHoverRadius: 5,
					order: 7,
				} as ChartDataset<'line', number[]>,
				{
					label: 'Hover helper',
					type: 'bar',
					data: prices.map((p) => biggest),
					backgroundColor: transparentColor,
					order: 20,
					hoverBackgroundColor: 'rgba(0, 200, 0, 0.3)',
					grouped: false,
				} as ChartDataset<'bar', number[]>,
				{
					label: 'Day change',
					type: 'bar',
					backgroundColor: pricesWithoutLast.map((p) => {
						return formatDateTime(p.s) === '0' ? blackColor : transparentColor;
					}),
					data: pricesWithoutLast.map((p) => biggest),
					order: 1,
					grouped: false,
					categoryPercentage: 0.1,
					datalabels: {
						labels: {
							value: {
								color: 'black',
							},
						},
						formatter(value: unknown, context: Context) {
							return formatDateDay(prices[context.dataIndex].s);
						},
						display(context: Context) {
							const p = prices[context.dataIndex];
							return formatDateTime(p.s) === '0';
						},
					},
				} as ChartDataset<'bar', number[]>,
			],
		},
	};
	if (!showOnlyAfterNow) {
		config.data.datasets.push({
			label: 'Nyt',
			type: 'bar',
			data: prices.map((p) => (isNow(p) ? biggest : 0)),
			grouped: false,
			order: 3,
			categoryPercentage: 0.5,
			datalabels: {
				labels: {
					value: {
						color: 'black',
					},
				},
				formatter(value: unknown, context: Context) {
					const p = prices[context.dataIndex];
					return formatPrice(p.p) + ' c/kWh';
				},
				display(context: Context) {
					const p = prices[context.dataIndex];
					return isNow(p);
				},
				offset: -50,
				anchor: 'center',
				align: 'start',
			},
		} as ChartDataset<'bar', number[]>);
	}
	return config as unknown as MyChartConfig;
}

Chart.register(
	Colors,
	LineController,
	LineElement,
	PointElement,
	BarController,
	BarElement,
	CategoryScale,
	LinearScale,
	Filler,
	Legend,
	Tooltip,
	ChartDataLabels,
);

export function setupChart(elementId: string, config: () => MyChartConfig) {
	let chart: Chart<'line' | 'bar', (number | string | null)[], string>;
	$effect(() => {
		chart = new Chart(document.getElementById(elementId) as HTMLCanvasElement, config());
		return () => {
			if (chart) {
				chart.destroy();
			}
		};
	});
}
