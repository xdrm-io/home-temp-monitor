<template>
	<div id='time-series'>
		<h1 ref='title'>
			Evolution from
			<input type='date' class='from' v-model='query.from' @change="onChange()"/>
			to
			<input type='date' class='to' v-model='query.to' @change="onChange()"/>
			by
			<select v-model='query.by' @change="onChange()">
				<option v-for='freq in Frequencies' :key='freq' :value='freq'>
					{{ freq }}
				</option>
			</select>
		</h1>
		<h2>5 min ago</h2>
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
	</div>
</template>


<script lang='ts'>
import { Vue } from 'vue-class-component';
import client from '../api/client'

enum Frequency {
	Second = 'second',
	Minute = 'minute',
	Hour   = 'hour',
	Day    = 'day',
}

interface Query {
	from?: Date;
	to?:   Date;
	by?:   Frequency;
	rooms: { [name:string]: boolean };
	ref?: string;
}

interface Series {
	name: string;
	data: {
		x: Date;
		y: number;
	}[];
}

export default class TimeSeries extends Vue {
	public readonly Frequencies: Frequency[] = [Frequency.Second, Frequency.Minute, Frequency.Hour, Frequency.Day];
	public rooms: string[] = [];

	public query: Query = {
		by: Frequency.Hour,
		to: new Date(),
		rooms: {},
	};
	public series: { [name:string]: Series } = {};

	public mounted() {
		client.getRoomNames()
		.then( (rooms) => { this.rooms = rooms; })
		.catch( console.error );
	}

	public onChange() {
		this.fetchCurrent()
	}

	private fetchCurrent() : string {
		const from  = this.query.from ? new Date(this.query.from) : undefined;
		const to    = this.query.to ? new Date(this.query.to) : undefined;
		const by    = this.query.by;
		const rooms = Object.keys(this.query.rooms).filter( (name) => this.query.rooms[name] === true ) ?? this.rooms;
		const ref   = this.query.ref || '';

		if( !by ){
			return "by is required";
		}

		const query = new URLSearchParams();
		if( from !== undefined ) {
			query.set('from', Math.round(from.getTime()/1000).toString());
		}
		if( to !== undefined ) {
			query.set('to', Math.round(to.getTime()/1000).toString());
		}
		if( rooms.length > 0 ) {
			for( const room of rooms ){
				query.append('rooms', room);
			}
		}
		if( ref ) {
			query.set('ref', ref);
		}
		query.set('by', by);
		console.debug(query.toString())
		return "";
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
