import { CONFIG } from "@/config";

type RoomsResponse = string[];

interface Current {
	lastUpdate?: Date;
	rooms: {
		[room_name: string]: {
			temperature: number;
			humidity: number;
		}
	}
}

export enum Frequency {
	Second = 'second',
	Minute = 'minute',
	Hour   = 'hour',
	Day    = 'day',
}
export const Frequencies: Frequency[] = [Frequency.Second, Frequency.Minute, Frequency.Hour, Frequency.Day];

export interface SeriesRequest {
	from: Date;
	to?: Date;
	by: Frequency;
	rooms?: string[];
	ref?: string;
}

export type LastSeriesResponse = {
	[room_name: string]: {
		t:           number;
		temperature: number;
		humidity:    number;
	}
}
export type SeriesResponse = {
	[room_name: string]: {
		t:    number;
		tmin: number;
		tavg: number;
		tmax: number;
		hmin: number;
		havg: number;
		hmax: number;
	}[]
}


class ClientClass {

	public getRoomNames(): Promise<string[]> {
		const url = `${CONFIG.api_url}/rooms`;
		return new Promise<string[]>( (resolve, reject) => {
			fetch(url).then( (response) => response.json() ).then( (res: RoomsResponse) => {
					resolve(res);
				})
				.catch( reject )
		});
	}

	public getCurrent() : Promise<Current> {
		const url = `${CONFIG.api_url}/last`;
		return new Promise( (resolve, reject) => {
			fetch(url).then( (response) => response.json() ).then( (data: LastSeriesResponse) => {
				const current: Current = {
					lastUpdate: undefined,
					rooms: {}
				};
				for( const room in data ) {
					const entry = data[room];
					if( current.lastUpdate === undefined ){
						current.lastUpdate = new Date(entry.t*1000);
					} else {
						current.lastUpdate = new Date(Math.max(current.lastUpdate.getTime(), entry.t*1000));
					}

					current.rooms[room] = {
						temperature: Math.floor(100*entry.temperature)/100,
						humidity:    Math.floor(100*entry.humidity)/100,
					}
				}
				resolve(current);
			}).catch( (error) => {
				reject(error);
			})
		});
	}
	public getSeries(req: SeriesRequest) : Promise<SeriesResponse> {
		const params = new URLSearchParams()
		params.append('from', Math.floor(req.from.getTime()/1000).toString())
		if( req.to !== undefined ) {
			params.append('to', Math.floor(req.to.getTime()/1000).toString())
		}
		if( req.rooms !== undefined && req.rooms.length > 0 ){
			for( const room of req.rooms ){
				params.append('rooms', room)
			}
		}
		params.append('by', req.by)
		if( req.ref !== undefined ) {
			params.append('ref', req.ref)
		}


		const url = `${CONFIG.api_url}/series?${params.toString()}`;
		return new Promise( (resolve, reject) => {
			fetch(url).then( (response) => response.json() ).then( (res: SeriesResponse) => resolve(res) )
			.catch( (error) => {
				reject(error);
			})
		});
	}
}

export default new ClientClass();