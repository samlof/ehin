/* eslint-disable @typescript-eslint/no-unused-vars */
import { formatPrice } from '$lib/calcUtils';
import { formatDateDay, formatDateTime, isNow } from '$lib/dateUtils';
import type { PriceEntry } from '$lib/pricesApi';
import {
	BarController,
	BarElement,
	CategoryScale,
	Chart,
	Colors,
	Legend,
	LinearScale,
	Tooltip,
	type ChartConfiguration,
	type ChartTypeRegistry,
} from 'chart.js';
import ChartDataLabels from 'chartjs-plugin-datalabels';

const defaultColor = 'rgba(54, 162, 235,0.5)';
const redColor = 'rgba(200, 0, 0,0.5)';
const greenColor = 'rgba(0, 200, 0,0.5)';
const transparentColor = 'rgba(54, 162, 235,0)';
const blackColor = 'rgba(100, 100, 100,0.5)';

function chooseColor(p: PriceEntry) {
	if (p.p > 100) {
		return redColor;
	} else if (p.p < 65) {
		return greenColor;
	}
	return defaultColor;
}

export type MyChartConfig = ChartConfiguration<keyof ChartTypeRegistry, string[], string>;

export function chartConfig(prices: PriceEntry[]): MyChartConfig {
	const biggestTemp = formatPrice(Math.max(...prices.map((p) => p.p)));
	const biggest = +biggestTemp < 10 ? '10' : biggestTemp;
	return {
		type: 'bar',
		plugins: [ChartDataLabels],
		options: {
			animation: {
				duration: 1,
			},
			scales: {},
			plugins: {
				datalabels: {
					labels: {
						title: null,
					},
				},
				legend: {
					labels: {
						filter(item, data) {
							return item.text === 'c/kWh';
						},
					},
				},
				tooltip: {
					filter(e, index, array, data) {
						return e.dataset.label === 'Hover helper' || e.dataset.label === 'c/kWh';
					},
					displayColors: false,
					callbacks: {
						title(tooltipItems) {
							const p = prices[tooltipItems[0].dataIndex];
							return formatDateTime(p.s) + ' - ' + formatDateTime(p.e);
						},
						label(tooltipItem) {
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
					data: prices.map((p) => formatPrice(p.p)),
					backgroundColor: prices.map(chooseColor),
					order: 7,
					grouped: false,
				},
				{
					label: 'Nyt',
					data: prices.map((p) => (isNow(p) ? biggest : '0')),
					grouped: false,
					order: 3,
					categoryPercentage: 0.5,
					datalabels: {
						labels: {
							value: {
								color: 'black',
							},
						},
						formatter(value, context) {
							const p = prices[context.dataIndex];
							return formatPrice(p.p) + ' c/kWh';
						},
						display(context) {
							const p = prices[context.dataIndex];
							return isNow(p);
						},
						offset: -50,
						anchor: 'center',
						align: 'start',
					},
				},
				{
					label: 'Hover helper',
					data: prices.map((p) => biggest),
					backgroundColor: transparentColor,
					order: 20,
					hoverBackgroundColor: defaultColor,
					grouped: false,
				},
				{
					label: 'Day change',
					backgroundColor: prices.map((p) =>
						formatDateTime(p.s) === '0' ? blackColor : transparentColor,
					),
					data: prices.map((p) => (formatDateTime(p.s) === '0' ? biggest : '0')),
					order: 1,
					grouped: false,
					categoryPercentage: 0.1,
					datalabels: {
						labels: {
							value: {
								color: 'black',
							},
						},
						formatter(value, context) {
							return formatDateDay(prices[context.dataIndex].s);
						},
						display(context) {
							const p = prices[context.dataIndex];
							return formatDateTime(p.s) === '0';
						},
					},
				},
			],
		},
	};
}

Chart.register(
	Colors,
	BarController,
	BarElement,
	CategoryScale,
	LinearScale,
	Legend,
	Tooltip,
	ChartDataLabels,
);

export function setupChart(elementId: string, config: () => MyChartConfig) {
	let chart: Chart<keyof ChartTypeRegistry, string[], string>;
	$effect(() => {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		chart = new Chart(document.getElementById(elementId) as any, config());
		return () => {
			if (chart) {
				chart.destroy();
			}
		};
	});
}
