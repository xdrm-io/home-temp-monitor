import { CONFIG } from "@/config";

type RoomsResponse = string[];
type SeriesResponse = {
	[room_name: string]: {
		t: number;
		tavg: number;
		tmin: number;
		tmax: number;
		havg: number;
		hmin: number;
		hmax: number;
	}[]
}

interface Current {
	lastUpdate?: Date;
	rooms: {
		[room_name: string]: {
			temperature: number;
			humidity: number;
		}
	}
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
		// fetch last 1 hour by second
		const params = new URLSearchParams()
		params.append('from', (Math.round(Date.now() / 1000) - 3600).toString())
		params.append('by', 'second')
		const url = `${CONFIG.api_url}/series?${params.toString()}`;

		return new Promise( (resolve, reject) => {
			fetch(url).then( (response) => response.json() ).then( (data: SeriesResponse) => {
				const current: Current = {
					lastUpdate: undefined,
					rooms: {}
				};
				for( const room in data ) {
					const series = data[room].sort( (a, b) => a.t - b.t );
					if( series.length === 0 ) {
						continue
					}

					const last = series[series.length-1];
					if( current.lastUpdate === undefined ){
						current.lastUpdate = new Date(last.t*1000);
					} else {
						current.lastUpdate = new Date(Math.max(current.lastUpdate.getTime(), last.t*1000));
					}

					current.rooms[room] = {
						temperature: Math.floor(100*last.tavg)/100,
						humidity:    Math.floor(100*last.havg)/100,
					}
				}
				resolve(current);
			}).catch( (error) => {
				reject(error);
			})
		});
	}
}

export default new ClientClass();