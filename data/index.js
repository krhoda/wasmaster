import React, {useEffect, useReducer} from "react";
import ReactDOM from "react-dom";

const wasmWorker = window.wasmWorker;

const doSomeWasm = () => {
	wasmWorker.postMessage({'action': 'do_wasm'});
}

const App = () => {
	const wasmReducer = (state, action) => {
		console.log("???")
		console.log(state, action);
		if (action.type === 'init') {
			console.log('state is true');
			return true;
		}
	};

	let [wasmLoaded, wasmDispatch] = useReducer(wasmReducer, false);
	console.log('WASM LOADED  IS');

	console.log(wasmLoaded)

	useEffect(() => {
		console.log("In USE EFFECT");
		wasmWorker.onmessage = (e) => {
			console.log("In UI THREAD -- HEARD:");
			console.log(e);

			let {data} = e;
			if (!data) {
				console.error('Could not find data in onmessage event!');
				console.error(e);
				return
			}

			let {action} = data;

			console.log("In UI THREAD -- ABOUT TO DISPATCH:");
			console.log(action);

			wasmDispatch({type: action})

			/* wasmWorker.postMessage({'action': 'do_wasm'}); */
		}

		wasmWorker.postMessage(
			{
				'action': 'init',
				'payload': {}
			}
		);
	}, [])

	let loaded = <p>Loading...</p>;
	if (wasmLoaded) {
		loaded = <button onClick={doSomeWasm}>WASM NOW!</button>;
	}

	return (
		<div>
			<h1>Hi there</h1>
			{loaded}
		</div>
	);
};

ReactDOM.render(<App />, document.getElementById("root"));
