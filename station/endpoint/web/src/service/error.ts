export interface Error {
	uid:     string;
	created: Date;

	msg:     string;
	error:   boolean;
}

export interface ErrorSubscriber {
	on(err: Error): void;
}

class ErrorCatcher {
	private sub?: ErrorSubscriber;

	public subscribe(sub?: ErrorSubscriber) {
		this.sub = sub;
	}

	public raise(err: string) {
		if( this.sub === undefined ){
			return;
		}
		this.sub.on({
			uid: `${new Date().toISOString()}-${err}`,
			created: new Date(),
			msg: err,
			error: true
		});
	}
	public info(err: string) {
		if( this.sub === undefined ){
			return;
		}
		this.sub.on({
			uid: `${new Date().toISOString()}-${err}`,
			created: new Date(),
			msg: err,
			error: false
		});
	}

}

export default new ErrorCatcher();