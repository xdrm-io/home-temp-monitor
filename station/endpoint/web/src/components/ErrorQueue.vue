<template>
	<div id='error-queue'>
		<transition-group>
		<div v-for='err of ordered' :key='err.uid' :class='err.error ? "error" : "info"'>
			{{ err.msg }}
		</div>
		</transition-group>
	</div>
</template>

<script lang="ts">
	import { Vue } from 'vue-class-component';
	import Queue, { Error, ErrorSubscriber } from '@/service/error';

	export default class ErrorQueue extends Vue implements ErrorSubscriber {
		protected queue_uid: string[] = [];
		protected queue_err: Error[]  = [];

		public mounted() {
			Queue.subscribe(this);
		}

		// implements the ErrorSubscriber
		public on(err: Error) {
			const MaxSize = 5;

			// already exists with this message, remove the other one
			const exists = this.queue_err.find( (v) => {
				return (v.msg == err.msg)
			})
			if( exists ){
				this.onHide(exists.uid)
			}

			this.queue_uid.push(err.uid);
			this.queue_err.push(err);

			setTimeout(() => this.onHide(err.uid), 3000);

			if( this.queue_uid.length > MaxSize ){
				this.queue_uid.shift();
				this.queue_err.shift();
			}
		}

		// returns the ordered list of visible errors
		public get ordered(): Error[] {
			const sorted = this.queue_err.sort( (a, b) => {
				return a.created.getTime() - b.created.getTime()
			})
			return sorted;
		}

		private onHide(uid: string) {
			const i = this.queue_uid.indexOf(uid);
			if( i < 0 ){
				return
			}
			this.queue_uid.splice(i, 1)
			this.queue_err.splice(i, 1)
		}
	}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
	#error-queue {
		display: flex;
		position: absolute;
			top: 0;
			left: 0;
			width: 100vw;
			height: 100vh;

		padding-bottom: 1em;

		flex-flow: column-reverse nowrap;
		justify-content: flex-start;
		align-items: flex-start;

		pointer-events: none;
		user-select: none;

		z-index: 1000;
	}

	.error, .info {
		display: block;
		position: relative;

		margin: .2em auto;
		padding: .4em 2em;

		font-size: 1.2em;
		font-weight: normal;
		color: #bb6565;
		background: #22262E;

		border: .1em solid lighten(#22262E, 5%);
		border-radius: .3rem / .3rem;
		box-shadow: 0 .5em 1em #111;

		pointer-events: fill;

		transform-origin: top center;
		transform: scaleY(1);

		opacity: 1;
		transition: opacity .2s ease-in-out,
					transform .2s ease-in-out;

		&:hover {
			opacity: .8;
		}
		&.v-enter-active{
			opacity: 0;
			transform: scaleY(0);
		}
		&.v-leave-active{
			opacity: 0;
			transform: scaleY(0);
		}

	}

	.info {
		color: #6581bb;
	}

</style>