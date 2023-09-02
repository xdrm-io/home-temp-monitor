// proxy query parameters to the api
const params = new URLSearchParams(window.location.search)
let query = ""
for( const param of params ){
    const key   = param[0]
    const value = param[1]
    query += `${key}=${value}&`
}

// fetch data from "GET http://192.168.1.10:8000/" in vanilla js
fetch(`/api/?${query}`).then( (response) => response.json() ).then( (data) => {

    let startingDate = null

    let tempSeries = []
    let humiditySeries = []
    let colorIndex = 0;
    for( const room in data ) {
        colorIndex = (colorIndex+1) % Highcharts.getOptions().colors.length;

        if( data[room].length > 0 && data[room][0].t < startingDate ) {
            startingDate = data[room][0].t
        }
        tempSeries.push({
            name: `${room}`,
            data: data[room].map( (d) => [d.t*1000, d.tavg] ),
            color: Highcharts.getOptions().colors[colorIndex],
            zIndex: 1,
        })
        tempSeries.push({
            name: `${room} min/max`,
            data: data[room].map( (d) => [d.t*1000, d.tmin, d.tmax] ),
            type: 'arearange',
            lineWidth: 0,
            linkedTo: ':previous',
            color: Highcharts.getOptions().colors[colorIndex],
            fillOpacity: 0.2,
            zIndex: 0,
            marker: { enabled: false }
        })
        humiditySeries.push({
            name: `${room}`,
            data: data[room].map( (d) => [d.t*1000, d.havg] ),
            color: Highcharts.getOptions().colors[colorIndex],
            zIndex: 1,
        })
        humiditySeries.push({
            name: `${room} min/max`,
            data: data[room].map( (d) => [d.t*1000, d.hmin, d.hmax] ),
            type: 'arearange',
            lineWidth: 0,
            linkedTo: ':previous',
            color: Highcharts.getOptions().colors[colorIndex],
            fillOpacity: 0.2,
            zIndex: 0,
            marker: { enabled: false }
        })
    }


    const baseConfig = {
        legend: {
            layout: 'vertical',
            align: 'right',
            verticalAlign: 'middle'
        },
        chart: {
            type: 'line'
        },
        xAxis: {
            type: 'datetime',
        },
        plotOptions: {
            line: {
                dataLabels: {
                    enabled: true
                },
                enableMouseTracking: true
            },
        },
    }

    Highcharts.chart('temperatures', {
        ...baseConfig,
        title: {
            text: 'Temperature per room'
        },
        yAxis: {
            title: {
                text: 'Temperature (°C)'
            }
        },
        tooltip: {
            crosshairs: true,
            shared: true,
            valueSuffix: '°C'
        },
        series: tempSeries,
    });

    Highcharts.chart('humidity', {
        ...baseConfig,
        title: {
            text: 'Humidity per room'
        },
        yAxis: {
            title: {
                text: 'Humidity (%)'
            }
        },
        tooltip: {
            crosshairs: true,
            shared: true,
            valueSuffix: '%'
        },
        series: humiditySeries,
    });


}).catch( (error) => {
    document.write = `<h1>${error}</h1>`
})

