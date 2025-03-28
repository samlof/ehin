const taxPercent = 25.5;
export function calculateTax(price: number) {
	return (price * taxPercent) / 100;
}

export function mwhToKwhPrice(price: number) {
	return (price / 10).toFixed(2);
}

export function formatPrice(price: number) {
	const tax = calculateTax(price);
	const totalPrice = price + tax;

	return mwhToKwhPrice(totalPrice);
}
