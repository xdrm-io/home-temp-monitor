const calendar = document.getElementById('calendar');
const picker = new DateRangePicker(calendar, {

});

const form = document.querySelector('form');
form.addEventListener('submit', (e) => {
    e.preventDefault();
    const from = picker.getDates()[0];
    const to   = picker.getDates()[1];
    if( !from || !to ){
        return;
    }


    // go the the data page with valid query parameters
    document.location.href = `/data?by=hour&from=${from.getTime()/1000}&to=${to.getTime()/1000}`;
}, false)