import { describe, test, expect } from 'vitest';
import { calculateTax } from './calcUtils';

describe('lib/calcUtils.ts', () => {
	test('should render h1', () => {
		const price = 6.73;
		expect(price + calculateTax(price)).toBeCloseTo(8.45);
	});
});
