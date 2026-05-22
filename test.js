// const numbers = [1, 2, 3, 4];
// const cart = [
// 	{ name: "Áo", price: 200, quantity: 2 },
// 	{ name: "Quần", price: 300, quantity: 1 },
// ];
// const stringss = "abcxyz";
// let resultStr1 = "";
// for (let i = stringss.length - 1; i >= 0; i--) {
// 	resultStr = resultStr + stringss[i];
// }
// let resultStr2 = stringss.split("").reverse().join("");
// console.log(resultStr2);

// const doubled = numbers.map(x => x * 2);
// console.log(doubled);
// const evenNumber = numbers.filter(x => x % 2 == 0);
// console.log(evenNumber);
// const gtNumber = numbers.find(x => x > 2);
// console.log(gtNumber);
// //1. Tính tổng các số trong mảng (Kết quả là 1 con số)
// let total = numbers.reduce((sum, num) => num + num);
// console.log("total", total);
// //2. Tính tổng tiền giỏ hàng (Kết quả là số từ mảng Object)
// const totalPrice = cart.reduce(
// 	(total, item) => total + item.price * item.quantity,
// 	0,
// );
// console.log("totalPrice", totalPrice);
// Bài 3: Thêm chữ “Hello”
// const names = ["An", "Bình", "Cường"];
// const result = names.map(x => "Hello " + x);
// Bài 1: Lọc số chẵn
// Kết quả: [2, 4, 6]
// const arr = [5, 12, 8, 20, 3];
// const result = arr.filter(numbers => numbers % 2 == 0);
// console.log(result);
// Bài 2: Lọc số > 10
// const result = arr.filter(numbers => numbers > 10);
// console.log(result);
// Bài 3: Lọc tên dài hơn 3 ký tự
// const names = ["An", "Bình", "Cường", "Vy"];
// const result = names.filter(str => str.length > 3);
// console.log(result);
// Bài 1: Tìm số 10
// const arr = [1, 3, 5, 8, 10];
// const result = arr.find(x => x == 10);
// console.log(result);
// Bài 2: Tìm số chẵn đầu tiên
// const evenFirstNumber = arr.find(x => x % 2 === 0);
// console.log(evenFirstNumber);
// Bài 3: Tìm user có id = 3
// const users = [
// 	{ id: 1, name: "An" },
// 	{ id: 2, name: "Bình" },
// 	{ id: 3, name: "Cường" },
// ];

// const resultID = users.filter(id => id.id == 3);
// console.log(resultID);
// reduce()
// Bài 1: Tính tổng
// const arr = [1, 2, 3, 4];
// const sumArr = arr.reduce((sum, x) => sum + x, 0);
// console.log(sumArr);
// Bài 2: Tính tích
// const arr = [1, 2, 3, 4];
// const result = arr.reduce((h, x) => h * x, 1);
// console.log(result);
// Bài 3: Đếm tổng số ký tự
// const words = ["hi", "hello", "bye"];
// const result = words.reduce((count, str) => count + str.length, 0);
// console.log(result);
// Bài 4 (quan trọng): Tìm số lớn nhất
// const arr = [3, 9, 2, 15, 7];
// const result = arr.reduce((a, b) => Math.max(a, b));
// console.log(result);
// 1. Destructuring (phá cấu trúc)
// const person = { name: "Bình", city: "HCM" };
// let { name, city } = person;
// console.log(name);
// console.log(city);
// 👉 Lấy phần tử đầu và cuối
// const numbers = [5, 10, 15];
// const [first, ...rest] = numbers;
// const last = rest.pop();
// console.log("First", first);
// console.log("last", last);
// 2. Spread Operator (...)
// const a = [1, 2];
// const b = [3, 4];
// const mergeArr = [...a, ...b];
// console.log(mergeArr);
// const user = { name: "An", age: 20 };
// const newUser = { ...user, city: "HCM" };
// const coppyNewUser = { ...newUser, gender: "male" };
// console.log(coppyNewUser);
// "Xin chào, tôi là [name], tôi đến từ [city]"
// const name = "Quốc";
// const city = "Đồng Tháp";
// console.log(`Xin chào, tôi là ${name}, tôi đến từ ${city}`);

async function fetchAPI() {
	try {
		const res = await fetch("https://jsonplacesholder.typicode.com/posts");
		if (!res) {
			throw new Error("Lỗi");
		}
		const data = await res.json();
		data.forEach(post => {
			console.log(post.title);
		});
	} catch (error) {
		console.log("Đã bắt lỗi:", error.message);
	}
}
fetchAPI();
