<template>
	<div id='header'>
		<span>Homenv</span>
	</div>
</template>


<script lang='ts'>
import { Vue } from 'vue-class-component';

export default class AppHeader extends Vue {
}
</script>

<style scoped lang='scss'>
#header {
	display: flex;
	position: relative;
		width: 100%;
		height: 4rem;

	padding-left: 1em;

	flex-flow: row nowrap;
	justify-content: flex-start;
	align-items: center;

	font-size: 1.5em;
	font-weight: 300;
	color: #fff;

	background: linear-gradient(to right bottom, #435770 0%, #16397d 100%);
}

</style>
