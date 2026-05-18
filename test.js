const numbers = [1, 2, 3, 4];
const cart = [
	{ name: "Áo", price: 200, quantity: 2 },
	{ name: "Quần", price: 300, quantity: 1 },
];
const doubled = numbers.map(x => x * 2);
console.log(doubled);
const evenNumber = numbers.filter(x => x % 2 == 0);
console.log(evenNumber);
const gtNumber = numbers.find(x => x > 2);
console.log(gtNumber);
//1. Tính tổng các số trong mảng (Kết quả là 1 con số)
let total = numbers.reduce((sum, num) => num + num);
console.log("total", total);
//2. Tính tổng tiền giỏ hàng (Kết quả là số từ mảng Object)
const totalPrice = cart.reduce(
	(total, item) => total + item.price * item.quantity,
	0,
);
console.log("totalPrice", totalPrice);
//3. Đếm số lần xuất hiện của các phần tử (Kết quả là 1 Object)
