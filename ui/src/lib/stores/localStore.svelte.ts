import { browser } from '$app/environment';

export class LocalStore<T> {
	private v: T;
	private key: string;

	get value() {
		return this.v;
	}

	set value(val: T) {
		localStorage.setItem(this.key, JSON.stringify(val));
		this.v = val;
	}

	constructor(key: string, value: T) {
		this.key = key;
		this.v = $state<T>(value);

		if (browser) {
			const item = localStorage.getItem(key);
			if (item) this.v = JSON.parse(item);
		}
	}
}

export function localStore<T>(key: string, value: T) {
	return new LocalStore(key, value);
}

export const localSettings = {
	showOnlyAfterNow: localStore('showOnlyAfterNow', false),
};
