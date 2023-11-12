<template>
	<error-queue/>
	<app-header/>
	<current-list/>
	<time-series/>
</template>


<script lang='ts'>
import { Options, Vue } from 'vue-class-component';
import AppHeader   from './components/AppHeader.vue';
import CurrentList from './components/CurrentList.vue';
import TimeSeries  from './components/TimeSeries.vue';
import ErrorQueue  from './components/ErrorQueue.vue';

@Options({
	components: {
		AppHeader,
		CurrentList,
		TimeSeries,
		ErrorQueue,
	}
})

export default class App extends Vue {}
</script>

<style lang='scss'>
*, *:before, *:after{
		box-sizing: border-box;
}

html, body, div, span, object, iframe, figure, h1, h2, h3, h4, h5, h6, p, blockquote, pre, a, code, em, img, small, strike, strong, sub, sup, tt, b, u, i, ol, ul, li, fieldset, form, label, table, caption, tbody, tfoot, thead, tr, th, td, main, canvas, embed, footer, header, nav, section, video{
		margin: 0;
		padding: 0;
		border: 0;
		font-size: 100%;
		font: inherit;
		vertical-align: baseline;
		text-rendering: optimizeLegibility;
		-webkit-font-smoothing: antialiased;
		text-size-adjust: none;
}

footer, header, nav, section, main{
		display: block;
}

body{
		line-height: 1;
}

ol, ul{
		list-style: none;
}

blockquote, q{
		quotes: none;
}

blockquote:before, blockquote:after, q:before, q:after{
		content: '';
		content: none;
}

table{
		border-collapse: collapse;
		border-spacing: 0;
}

input{
		-webkit-appearance: none;
		border-radius: 0;
}
</style>

<style lang='scss'>
body {
	display: block;
	position: absolute;
		width: 100%;
		min-height: 100%;

	background: #fff;

	font-family: 'Outfit';
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
	font-size: 16px;
	color: #000;

	overflow: hidden;
	overflow-y: auto;
}

</style>
