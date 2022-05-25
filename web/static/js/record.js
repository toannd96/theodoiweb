window.recorder = {
	events: [],
	rrweb: undefined,
	runner: undefined,
	session: {
		genId(length) {
			const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
			let result = "";
			const charactersLength = characters.length;
			for (let i = 0; i < length; i++) {
				result += characters.charAt(Math.floor(Math.random() * charactersLength));
			}
			return result;
		},
		get() {
			let session = window.sessionStorage.getItem('rrweb');
			if (session) return JSON.parse(session);
			session = {
				session_id: window.recorder.session.genId(64),
			};
			window.sessionStorage.setItem('rrweb', JSON.stringify(session));
			return session;
		},
		receive(data) {
			const session = window.recorder.session.get();
			window.sessionStorage.setItem('rrweb', JSON.stringify(Object.assign({}, session, data)));
		},
		clear() {
			window.sessionStorage.removeItem('rrweb')
		}
	},
	setSession: function () {
		const session = window.recorder.session.get();
		session.session_id = window.recorder.session.genId(64);
		window.recorder.session.receive(session)
		return window.recorder;
	},
	stop() {
		clearInterval(window.recorder.runner);
	},
	start() {
		window.recorder.runner = setInterval(function receive() {
			const session = window.recorder.session.get();
			fetch('https://theodoi-web.herokuapp.com/session/receive', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(Object.assign({}, { events: window.recorder.events }, session)),
			});
			window.recorder.events = []; // cleans-up events for next cycle
		}, 5 * 1000);
	},
	close() {
		clearInterval();
		window.recorder.session.clear();
	}
};
new Promise((resolve, reject) => {
	const script = document.createElement('script');
	script.src = 'https://cdn.jsdelivr.net/npm/rrweb@latest/dist/rrweb.min.js';
	script.addEventListener('load', resolve);
	script.addEventListener('error', e => reject(e.error));
	document.head.appendChild(script);
}).then(() => {
	window.recorder.rrweb = rrweb;
	rrweb.record({
		emit(event) {
			window.recorder.events.push(event);
		}
	});
	window.recorder.start();
}).catch(console.err);
