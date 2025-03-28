import { browser } from '$app/environment';

export type breakpointVals = 'xs' | 's' | 'm' | 'l' | 'xl' | 'xxl';

let breakpointState: breakpointVals = $state('xxl');

if (browser) {
	const breakpoints = [
		{ value: 'xs', mediaquery: window.matchMedia('(max-width:  479px)') },
		{ value: 's', mediaquery: window.matchMedia('(min-width:  480px) and (max-width:  719px)') },
		{ value: 'm', mediaquery: window.matchMedia('(min-width:  720px) and (max-width:  959px)') },
		{ value: 'l', mediaquery: window.matchMedia('(min-width:  960px) and (max-width: 1439px)') },
		{ value: 'xl', mediaquery: window.matchMedia('(min-width: 1440px) and (max-width: 1919px)') },
		{ value: 'xxl', mediaquery: window.matchMedia('(min-width: 1920px)') },
	] as const;

	for (const b of breakpoints) {
		//set the current breakpoint
		if (b.mediaquery.matches === true) {
			// EventBus.$emit("breakpoint", breakpoint.value);
			breakpointState = b.value;
		}
		b.mediaquery.addEventListener('change', (event) => {
			if (event.matches === true) {
				breakpointState = b.value;
			}
		});
	}
}

export const breakpoint = () => breakpointState;
