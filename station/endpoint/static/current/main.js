// proxy query parameters to the api
const params = new URLSearchParams()
// 1 hour ago
params.append('from', parseInt(new Date().getTime() / 1000) - 3600)
// by minute
params.append('by', 'second')

// fetch data from "GET http://192.168.1.10:8000/" in vanilla js
fetch(`/api/?${params.toString()}`).then( (response) => response.json() ).then( (data) => {
    // get last data for each room

    // type Room = {
    //      t:           number,
    //      temperature: number,
    //      humidity:    number,
    // };
    const rooms = {};
    for( const room in data ) {
        const series = data[room].sort( (a, b) => a.t - b.t );
        if( series.length === 0 ) {
            continue
        }

        const last = series[series.length-1];

        rooms[room] = {
            t: last.t,
            temperature: parseInt(100*last.tavg)/100,
            humidity:    parseInt(100*last.havg)/100,
        }
    }

    // build DOM
    const container = document.getElementById('current');
    for( const room in rooms ){

        const name = document.createElement('div');
        name.classList.add('name');
        name.innerText = room;
        container.appendChild(name);

        const temp = document.createElement('div');
        temp.classList.add('temp');
        temp.innerText = `${rooms[room].temperature}`;
        container.appendChild(temp);

        const humidity = document.createElement('div');
        humidity.classList.add('humidity');
        humidity.innerText = `${rooms[room].humidity}`;
        container.appendChild(humidity);
    }


}).catch( (error) => {
    document.write = `<h1>${error}</h1>`
})

