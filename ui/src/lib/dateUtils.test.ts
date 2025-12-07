import { formatSeconds, isNextDayVisible } from '$lib/dateUtils';
import { describe, expect, test } from 'vitest';

describe('lib/dateUtils.ts', () => {
	test('formatMillis works', () => {
		expect(formatSeconds(60 * 5)).toBe('00:05:00');

		expect(formatSeconds(60 * 5 + 15)).toBe('00:05:15');

		expect(formatSeconds(60 * 60 * 1 + 60 * 5 + 15)).toBe('01:05:15');
	});

	describe('isNextDayVisible', () => {
		test('returns false when no prices exist for next day', () => {
			const now = new Date('2025-12-07T10:00:00');
			const prices = [
				{ s: new Date('2025-12-07T08:00:00'), e: new Date('2025-12-07T09:00:00'), p: 10 },
				{ s: new Date('2025-12-07T09:00:00'), e: new Date('2025-12-07T10:00:00'), p: 11 },
			];
			expect(isNextDayVisible(prices, now)).toBe(false);
		});

		test('returns false when next day prices are before 03:00', () => {
			const now = new Date('2025-12-07T10:00:00');
			const prices = [
				{ s: new Date('2025-12-07T08:00:00'), e: new Date('2025-12-07T09:00:00'), p: 10 },
				{ s: new Date('2025-12-08T02:00:00'), e: new Date('2025-12-08T03:00:00'), p: 11 },
			];
			expect(isNextDayVisible(prices, now)).toBe(false);
		});

		test('returns true when next day prices are at or after 03:00', () => {
			const now = new Date('2025-12-07T10:00:00');
			const prices = [
				{ s: new Date('2025-12-07T08:00:00'), e: new Date('2025-12-07T09:00:00'), p: 10 },
				{ s: new Date('2025-12-08T03:00:00'), e: new Date('2025-12-08T04:00:00'), p: 11 },
			];
			expect(isNextDayVisible(prices, now)).toBe(true);
		});

		test('returns true when prices are much later in future', () => {
			const now = new Date('2025-12-07T10:00:00');
			const prices = [
				{ s: new Date('2025-12-07T08:00:00'), e: new Date('2025-12-07T09:00:00'), p: 10 },
				{ s: new Date('2025-12-09T10:00:00'), e: new Date('2025-12-09T11:00:00'), p: 11 },
			];
			expect(isNextDayVisible(prices, now)).toBe(true);
		});

		test('returns false when year changes but month and day are earlier', () => {
			const now = new Date('2025-12-31T10:00:00');
			const prices = [
				{ s: new Date('2025-12-31T08:00:00'), e: new Date('2025-12-31T09:00:00'), p: 10 },
				{ s: new Date('2026-01-01T10:00:00'), e: new Date('2026-01-01T11:00:00'), p: 11 },
			];
			expect(isNextDayVisible(prices, now)).toBe(true);
		});
	});
});
