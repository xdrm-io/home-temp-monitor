<template>
	<div id='time-series'>
		<h1 ref='title'>
			Evolution from
			<input type='date' class='from' v-model='query.from' @change='onChange()'/>
			to
			<input type='date' class='to' v-model='query.to' @change="onChange()"/>
			by
			<select v-model='query.by' @change="onChange()">
				<option v-for='freq in Frequencies' :key='freq' :value='freq'>
					{{ freq }}
				</option>
			</select>
		</h1>
		<h2>{{ lastUpdateInMinutes(updatedAt) }} min ago</h2>
		<div class='filters'>
			<select v-model='query.ref' placeholder='reference' @change="onChange()">
				<option value=''>reference</option>
				<option v-for='name in rooms' :key='name' :value='name'>
					{{ name }}
				</option>
			</select>

			<div class='separator'></div>

			<div class='room' v-for='name in rooms' :key='name'>
				<input type='checkbox' :id='name' :value='name' v-model='query.rooms[name]' @change="onChange()"/>
				<label :for='name'>{{ name }}</label>
			</div>
		</div>

		<figure class="highcharts-figure">
			<div ref='temperatureChart' />
		</figure>
		<figure class="humidity-figure">
			<div ref='humidityChart' />
		</figure>
	</div>
</template>


<script lang='ts'>
import { Vue } from 'vue-class-component';
import client, {SeriesResponse} from '../api/client'
import Queue from '@/service/error';

import * as Highcharts from 'highcharts';
import highchartsMore from 'highcharts/highcharts-more';
highchartsMore(Highcharts);

enum Frequency {
	Second = 'second',
	Minute = 'minute',
	Hour   = 'hour',
	Day    = 'day',
}

interface Query {
	from?: string;
	to?:   string;
	by?:   Frequency;
	rooms: { [name:string]: boolean };
	ref?: string;
}

const LOAD_DELAY_MS = 2*1000;

export default class TimeSeries extends Vue {
	public readonly Frequencies: Frequency[] = [Frequency.Second, Frequency.Minute, Frequency.Hour, Frequency.Day];
	public rooms: string[] = [];

	public updatedAt?: Date;

	public query: Query = {
		by: Frequency.Hour,
		from: new Date(Date.now()-24*3600*1000).toISOString().substring(0,10),
		rooms: {},
	};
	public series: SeriesResponse = {};

	private timeout: number|undefined = undefined;

	private tChart?: Highcharts.Chart;
	private hChart?: Highcharts.Chart;


	public lastUpdateInMinutes(date: Date|undefined): string{
		if( date === undefined || !(date instanceof Date) ) {
			return '-';
		}
		return Math.round( (Date.now() - date.getTime())/1000 / 60 ).toString();
	}

	public mounted() {

		// init charts
		const baseConfig: Highcharts.Options = {
			legend: {
				layout: 'vertical',
				align: 'right',
				verticalAlign: 'middle'
			},
			xAxis: {
				type: 'datetime',
			},
			plotOptions: {
				line: {
					dataLabels: {
						style: { fontFamily: 'Outfit' },
						enabled: true
					},
					enableMouseTracking: true
				},
			},
		}
		this.tChart = Highcharts.chart(this.$refs.temperatureChart as HTMLElement, {
			...baseConfig,
			title: {
				style: { fontFamily: 'Outfit' },
				text: 'Temperature per room'
			},
			yAxis: {
				title: {
					style: { fontFamily: 'Outfit' },
					text: 'Temperature (°C)'
				},
			},
			tooltip: {
				style: { fontFamily: 'Outfit', fontWeight: '300' },
				shared: true,
				valueSuffix: '°C'
			},
		});
		this.hChart = Highcharts.chart(this.$refs.humidityChart as HTMLElement, {
			...baseConfig,
			title: {
				style: { fontFamily: 'Outfit' },
				text: 'Humidity per room'
			},
			yAxis: {
				title: {
					style: { fontFamily: 'Outfit' },
					text: 'Humidity (%)'
				},
			},
			tooltip: {
				style: { fontFamily: 'Outfit', fontWeight: '300' },
				shared: true,
				valueSuffix: '%'
			},
		});

		client.getRoomNames()
		.then( (rooms) => {
			this.rooms = rooms;
			this.query.rooms = {};
			for( const name of rooms ){
				this.query.rooms[name] = true;
			}

			// launch initial query
			this.fetchSeries().catch( err => Queue.raise(err) );
		})
		.catch( (err) => Queue.raise(err) );

	}

	public onChange() {
		this.delayFetch();
	}

	private delayFetch() {
		if( this.timeout !== undefined ) {
			clearTimeout(this.timeout);
		}
		this.timeout = setTimeout( () => {
			this.timeout = undefined;
			this.fetchSeries().catch( err => Queue.raise(err) );
		}, LOAD_DELAY_MS);
	}

	private fetchSeries() : Promise<void> {
		const from  = this.query.from ? new Date(this.query.from) : undefined;
		const to    = this.query.to ? new Date(this.query.to) : undefined;
		const by    = this.query.by;
		const rooms = Object.keys(this.query.rooms).filter( (name) => this.query.rooms[name] === true ) ?? this.rooms;
		const ref   = this.query.ref || '';

		return new Promise<void>( (resolve, reject) => {
			if( from === undefined ){
				reject(new Error("from is required"));
				return;
			}
			if( by === undefined ){
				reject(new Error("by is required"));
				return;
			}
			client.getSeries({ from, to, by, rooms, ref })
				.then( (series) => {
					this.series = series;
					this.configureChart(series);
					resolve();
				})
				.catch(reject);
		});
	}

	private configureChart(data: SeriesResponse) {
		let startTimestamp: number|undefined;

		const tSeries:     Highcharts.SeriesOptionsType[] = []
		const hSeries: Highcharts.SeriesOptionsType[] = []
		let   colorIndex                                     = 0;
		const palette                                        = Highcharts.getOptions().colors!;
		for( const room in data ) {
			colorIndex = (colorIndex+1) % palette.length;

			if( data[room].length > 0 && (startTimestamp === undefined || data[room][0].t < startTimestamp) ){
				startTimestamp = data[room][0].t
			}
			tSeries.push({
				type: 'line',
				name: `${room}`,
				data: data[room].map( (d) => [d.t*1000, Math.round(100*d.tavg)/100] ),
				color: palette[colorIndex],
				zIndex: 1,
			})
			tSeries.push({
				name: `${room} range`,
				data: data[room].map( (d) => [d.t*1000, Math.round(100*d.tmin)/100, Math.round(100*d.tmax)/100] ),
				type: 'arearange',
				lineWidth: 0,
				linkedTo: ':previous',
				color: palette[colorIndex],
				fillOpacity: 0.2,
				zIndex: 0,
				marker: { enabled: false }
			})
			hSeries.push({
				type: 'line',
				name: `${room}`,
				data: data[room].map( (d) => [d.t*1000, Math.round(100*d.havg)/100] ),
				color: palette[colorIndex],
				zIndex: 1,
			})
			hSeries.push({
				name: `${room} range`,
				data: data[room].map( (d) => [d.t*1000, Math.round(100*d.hmin)/100, Math.round(100*d.hmax)/100] ),
				type: 'arearange',
				lineWidth: 0,
				linkedTo: ':previous',
				color: palette[colorIndex],
				fillOpacity: 0.2,
				zIndex: 0,
				marker: { enabled: false }
			})
		}

		this.tChart?.update({ series: tSeries, }, true, true)
		this.hChart?.update({ series: hSeries, }, true, true)

		this.updatedAt = new Date();
		Queue.info('Time series updated')
	}

}
</script>

<style scoped lang='scss'>
#time-series {
	display: block;
	position: relative;
		width: 100%;
		height: auto;

	padding: 1.5em;

	color: #000;
	background: none;

	h1 {
		color: #000;
		font-size: 2em;
		font-weight: 400;

		margin-bottom: .1em;
	}
	h2 {
		color: #707070;
		font-size: 1em;
		font-weight: 300;

		margin-bottom: .5em;
	}


	input[type=date],
	select {
		padding: .4em;

		font-family: 'Outfit';
		font-weight: 300;
		font-size: 1.2rem;

		border: none;
		border-radius: .3rem / .3rem;
		background: #f6f6f6;
	}

	.filters {
		display: flex;
		position: relative;
		padding: 1em 1.3em;

		flex-flow: row wrap;
		justify-content: flex-start;
		align-items: center;

		background: #f9f9f9;

		.room {
			input[type=checkbox] {
				display: none;
			}

			label {
				display: inline-block;
				position: relative;
				padding: .5em .8em;
				margin: .2em .5em;

				font-size: 1.2rem;
				color: #666;

				background: #fff;
				border: none;
				border-radius: .3rem / .3rem;

				user-select: none;
				cursor: pointer;
			}

			input[type=checkbox]:checked + label {
				color: #fff;
				background: #3069fe;
			}
		}

		select {
			background: #fff;
		}

		.separator {
			display: block;
			position: relative;
			width: 1px;
			height: 1.5em;

			margin: 0 1em;

			background: #ccc;
		}
	}

}

</style>
