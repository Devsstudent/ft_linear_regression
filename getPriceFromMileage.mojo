from python import Python

fn main():
	var mileageValue = 0
	try :
		let py_sys = Python.import_module('sys')
		var mileageValue = Python.evaluate("int(input('Enter a mielage: '))").to_float64()
		while mileageValue < 0:
			let str_err: Error = "Please enter a value greater than 0\n"
			mileageValue = Python.evaluate("int(input('Enter a mielage: '))").to_float64()
	except:
		return 
	try:
		var infoTheta = open("tetaInfo", "r")
		var letter: String = infoTheta.read(1)
		var theta0Str: String = ""
		var theta1Str: String = ""
		while (letter != ' ') :
			theta0Str = theta0Str.__add__(letter)
			letter = infoTheta.read(1)
		while (letter):
			theta1Str = theta1Str.__add__(letter)
			letter = infoTheta.read(1)
		print(theta0Str, theta1Str)
	except:
		return

	let theta0 = 0
	let theta1 = 0
	let estimatePrice = theta0 + (theta1 * mileageValue)
	print(estimatePrice)