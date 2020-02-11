// TODO: MAKE TEMPLATE
importScripts('react_rust_wasm.js');

let ready = false;

const loadBindgen = (cb) => {
	wasm_bindgen('react_rust_wasm_bg.wasm').then(() => {
		console.log('IN WORKER THREAD after WASM:');
		ready = true;
		cb();
	}).catch(console.error);
}

const sendReady = () => {
	postMessage({action: 'init', payload: {ready: ready}});
};

const doBigComputation = () => {
	wasm_bindgen.big_computation();
};

onmessage = (e) => {
	console.log("In WORKER THREAD -- HEARD:");
	console.log(e);

	if (!ready && e.data.action !== 'init') {
		// TODO: ERR OUT HERE.
	}

	switch(e.data.action) {
		case 'init':
			if (ready) {
				sendReady();
				break;
			}

			loadBindgen(sendReady);
			break;

		case 'do_wasm':
			if (!ready) {
				return;
			}

			doBigComputation();
			break;

		default:
			console.error(`Unknown action: ${e.data.action}`);
			break;
	}
}
