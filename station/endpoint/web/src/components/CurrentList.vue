<template>
	<div id='current-list'>
		<h1>Current</h1>
		<h2>{{ lastUpdateInMinutes(updatedAt) }} min ago</h2>
		<div class='cards'>
			<div class='card' v-for='(room, name) in rooms' :key='name'>
				<h3 :class="room.offline === true ? 'offline' : ''">{{ name }}</h3>
				<div class='grid'>
					<img src='@/assets/temperature.svg' alt='temperature' class='temperature-icon'/>
					<h4 class='temperature-label'>Temperature</h4>
					<div class='temperature-value'>
						{{ room.temperature }}°C
					</div>
					<div class='temperature-diff' :class="'temperature-diff ' + diffClass(room.diff.temperature)" >
						{{ Math.abs(room.diff.temperature) }}°C
					</div>

					<img src='@/assets/humidity.svg' alt='humidity' class='humidity-icon'/>
					<h4 class='humidity-label'>Humidity</h4>
					<div class='humidity-value'>
						{{ room.humidity }}%
					</div>
					<div class='humidity-diff' :class="'humidity-diff ' + diffClass(room.diff.humidity)" >
						{{ Math.abs(room.diff.humidity) }}%
					</div>
				</div>
			</div>
		</div>
	</div>
</template>


<script lang='ts'>
import { Vue } from 'vue-class-component';
import { CONFIG } from '../config'

interface Room {
	offline:     boolean|undefined;
	temperature: number;
	humidity:    number;
	diff: {
		temperature: number;
		humidity:    number;
	}
}

interface LocalStorageData {
	updatedAt: Date|undefined;
	rooms:     { [name:string]: Room };
}

const LSKEY = 'current.last';
const INITIAL_INTERVAL_MS = 1000;
const REFRESH_INTERVAL_MS = 60*1000;

export default class CurrentList extends Vue {
	public updatedAt: Date|undefined          = undefined;
	public rooms:     { [name:string]: Room } = {};

	private restore(): LocalStorageData|undefined {
		if( !window.localStorage ){
			return undefined;
		}
		const data = window.localStorage.getItem(LSKEY);
		if( data === null ) {
			return undefined;
		}
		try{
			return JSON.parse(data);
		} catch( error ) {
			return undefined;
		}
	}
	private save(data: LocalStorageData): void {
		if( !window.localStorage ){
			return;
		}
		window.localStorage.setItem(LSKEY, JSON.stringify(data));
	}

	public lastUpdateInMinutes(date: Date|undefined): string{
		if( date === undefined || !(date instanceof Date) ) {
			return '-';
		}
		return Math.round( (Date.now() - date.getTime())/1000 / 60 ).toString();
	}

	public diffClass(diff: number): string {
		if( diff > 0 ) { return 'up'; }
		if( diff < 0 ) { return 'down'; }
		return 'same';
	}

	public mounted(): void {
		// try to restore the last item from local storage
		const restored = this.restore();
		if( restored === undefined ){
			this.refresh();
			return;
		}

		this.updatedAt = restored.updatedAt ? new Date(restored.updatedAt) : undefined;
		this.rooms = restored.rooms;
		// clear diffs
		for( const room in this.rooms ){
			this.rooms[room].diff = { temperature: 0, humidity: 0 };
		}
		setTimeout(this.refresh, INITIAL_INTERVAL_MS);

	}

	private scheduleRefresh() {
		setTimeout(this.refresh, REFRESH_INTERVAL_MS);
	}

	private refresh() {
		this.fetchCurrent().
		then(this.scheduleRefresh)
		.catch( (error) => {
			if( error !== undefined ) {
				console.error(error);
			}
		});
	}

	private fetchCurrent() : Promise<void> {
		// proxy query parameters to the api
		const params = new URLSearchParams()
		// fetch last 1 hour
		params.append('from', (Math.round(Date.now() / 1000) - 3600).toString())
		// by second
		params.append('by', 'second')
		return new Promise( (resolve, reject) => {

			interface Response {
				[key: string]: {
					t: number;
					tavg: number;
					havg: number;
				}[]
			}

			fetch(`${CONFIG.api_url}/?${params.toString()}`).then( (response) => response.json() ).then( (data: Response) => {
				// get last data for each room
				let lastUpdate: Date|undefined;

				const current: { [name:string]: Room } = {};
				for( const room in data ) {
					const series = data[room].sort( (a, b) => a.t - b.t );
					if( series.length === 0 ) {
						continue
					}

					const last = series[series.length-1];
					if( lastUpdate === undefined ){
						lastUpdate = new Date(last.t*1000);
					} else {
						lastUpdate = new Date(Math.max(lastUpdate.getTime(), last.t*1000));
					}

					current[room] = {
						offline: true,
						temperature: Math.floor(100*last.tavg)/100,
						humidity:    Math.floor(100*last.havg)/100,
						diff: { temperature: 0, humidity: 0 }
					}
				}

				// update and calculate diff
				for( const room in current ){
					if( this.rooms[room] === undefined ) {
						this.rooms[room] = current[room];
						continue;
					}
					this.rooms[room].offline = false;
					this.rooms[room].temperature = current[room].temperature;
					this.rooms[room].humidity    = current[room].humidity;
					this.rooms[room].diff = {
						temperature: Math.round(100* (current[room].temperature - this.rooms[room].temperature)) / 100,
						humidity:    Math.round(100* (current[room].humidity    - this.rooms[room].humidity)) / 100
					}
				}

				// mark missing rooms as offline
				for( const room in this.rooms ){
					if( current[room] === undefined ) {
						this.rooms[room].offline = true;
					}
				}

				this.updatedAt = lastUpdate;

				// store in local storage if we refresh the page
				this.save({
					updatedAt: lastUpdate,
					rooms: this.rooms
				});

				resolve();
			}).catch( (error) => {
				reject(error);
			})
		});
	}

}
</script>

<style scoped lang='scss'>
#current-list {
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
}


.cards {
	display: flex;
	position: relative;
		width: 100%;
		height: auto;

	flex-flow: row wrap;
	justify-content: flex-start;
	align-items: flex-start;


	.card {
		display: block;
		position: relative;

		margin-right: 1.5em;
		margin-bottom: 1.5em;
		&:last-child {
			margin-right: 0;
		}

		border-radius: 1rem / 1rem;
		background: linear-gradient(to right bottom, #fcfcfc, #eeeeee);

		cursor: pointer;
		transition: .2s ease-in-out box-shadow,
					.2s ease-in-out transform;
		&:hover {
			box-shadow: 0 0 2em rgba(0,0,0,.2);
			transform: scale(1.05);
		}

		h3 {
			position: absolute;
			display: block;
				top: 0;
				left: 0;

			padding: .5em .8em;

			border-radius: 1rem 0 1rem 0;
			background: #b1c7ff;

			color: #000;
			font-size: 1em;
			font-weight: 400;

			z-index: 101;

			&.offline { background: #eee; }
		}

		.grid {
			display: grid;
			position: relative;

			padding: 1em;
			z-index: 100;

			grid-template-columns: 1fr auto 1fr;
			grid-template-rows: 2fr 1fr 1fr 2fr 1fr 1fr 2fr;

			.temperature-icon, .humidity-icon {
				grid-column: 1;
				width: 2em;
				margin: 0 2em;
			}
			.temperature-icon { grid-row: 2 / 4; }
			.humidity-icon { grid-row: 5 / 7; }

			h4 {
				grid-column: 2;

				color: #969696;
				font-size: 1.3em;
				font-weight: 300;
			}
			h4.temperature-label { grid-row: 2; }
			h4.humidity-label { grid-row: 5; }

			.temperature-value, .humidity-value {
				grid-column: 2;

				color: #000;
				font-size: 1.5em;
				font-weight: 400;
			}
			.temperature-value { grid-row: 3; }
			.humidity-value { grid-row: 6; }

			.temperature-diff, .humidity-diff {
				grid-column: 3;

				display: block;
				position: relative;

				color: #969696;
				font-size: 1em;
				font-weight: 300;

				&::before {
					display: inline-block;
					position: relative;
					content: '';

					width: 0;
					height: 0;

					border-radius: .1rem / .1rem;

					margin-right: .5em;
				}
				&.up::before {
					border-left: .5em solid transparent;
					border-right: .5em solid transparent;
					border-bottom: .5em solid #3069fe;
				}
				&.down::before {
					border-left: .5em solid transparent;
					border-right: .5em solid transparent;
					border-top: .5em solid #e73130;
				}
				&.same::before {
					width: .7em;
					margin-bottom: .2em;
					height: .2em;
					background: #969696;
				}
			}
			.temperature-diff { grid-row: 3; }
			.humidity-diff { grid-row: 6; }
		}
	}
}

</style>
