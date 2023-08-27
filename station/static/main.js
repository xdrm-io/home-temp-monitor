// fetch data from "GET http://192.168.1.10:8000/" in vanilla js
fetch('/api').then( (response) => response.json() ).then( (data) => {

    let tempSeries = []
    let humiditySeries = []
    for( const room in data ) {
        tempSeries.push({
            name: room,
            data: data[room].map( (d) => [d.ts*1000, d.temperature] )
        })
        humiditySeries.push({
            name: room,
            data: data[room].map( (d) => [d.ts*1000, d.humidity] )
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
        tooltip: {
            crosshairs: true
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
        series: humiditySeries,
    });


})

