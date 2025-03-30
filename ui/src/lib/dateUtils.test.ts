import { formatSeconds } from '$lib/dateUtils';
import { describe, expect, test } from 'vitest';

describe('lib/dateUtils.ts', () => {
	test('formatMillis works', () => {
		expect(formatSeconds(1000 * 60 * 5)).toBe('00:05:00');

		expect(formatSeconds(1000 * 60 * 5 + 1000 * 15)).toBe('00:05:15');

		expect(formatSeconds(1000 * 60 * 60 * 1 + 1000 * 60 * 5 + 1000 * 15)).toBe('01:05:15');
		expect(formatSeconds(1000 * 60 * 60 * 1 + 1000 * 60 * 5 + 1000 * 15 + 422)).toBe('01:05:15');
	});
});
