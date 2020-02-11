package tmpstr

const CargoToml = `[package]
name = "{{.ProjectName}}"
version = "0.1.0"
edition= "2018"

[lib]
crate-type = ["cdylib"]

[dependencies]
wasm-bindgen = "0.2"`

const WebWorker = `importScripts('{{.ProjectName}}.js');

let ready = false;

const loadBindgen = (cb) => {
	wasm_bindgen('{{.ProjectName}}_bg.wasm').then(() => {
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
			console.error('Unknown action: ' + e.data.action);
			break;
	}
}`

const PackageJson = `{
  "name": "{{.ProjectName}}",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1",
    "build-wasm": "cargo build --target wasm32-unknown-unknown",
    "build-bindgen": "wasm-bindgen target/wasm32-unknown-unknown/debug/{{.ProjectName}}.wasm --out-dir dist --no-modules --no-typescript",
    "build": "npm run build-wasm && npm run build-bindgen && npx webpack",
    "start": "webpack-dev-server"
  },
  "keywords": [],
  "author": "",
  "license": "ISC",
  "dependencies": {
    "react": "^16.12.0",
    "react-dom": "^16.12.0"
  },
  "devDependencies": {
    "@babel/core": "^7.8.4",
    "@babel/preset-env": "^7.8.4",
    "@babel/preset-react": "^7.8.3",
    "babel-loader": "^8.0.6",
    "webpack": "^4.41.5",
    "webpack-cli": "^3.3.10",
    "webpack-dev-server": "^3.10.3"
  }
}`
