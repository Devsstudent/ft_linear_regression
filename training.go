package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"strconv"
	"strings")

//resume of the linear regression :
//It is the form Y = a + bX
//To get a : (somme de x) / n - b * (somme de y) / n
//To get b : (some des combine xy (en moyenne)) - somme des combine / n
//Diviser par la somme des x^2 - somme des x^2 / n

//Il nous faut donc la somme des x^2

//La somme de x, et de y

//La somme de x par y

func main() {
	file, err := os.Open("./data.csv");
	if (err != nil) {
		log.Fatal(err);
	}
	defer file.Close()
	scanner := bufio.NewScanner(file);
	var Sxy, Sx2, Sx, Sy float64;
	var n = 0;
	for scanner.Scan() {
		if (scanner.Text() == "km,price") {
			continue ;
		}
		var line = strings.Split(scanner.Text(), ",");
		if (len(line) < 2) {
			log.Fatal("Wrong data format");
		}
		var km, price float64;
		km, err = strconv.ParseFloat(line[0], 10);
		price, err = strconv.ParseFloat(line[1], 10);
		if (err != nil) {
			log.Fatal(err);
		}
		Sxy += km * price;
		Sx2 += km * km;
		Sx += km;
		Sy += price;
		n++;
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err);
	}
	var b = (Sxy - (Sx * Sy) / float64(n)) / (Sx2 - (Sx * Sx) / float64(n));
	var a = Sy / float64(n) - b * Sx / float64(n);

	err = os.WriteFile("./tetaInfo", []byte(fmt.Sprintf("%f", a) + "," + fmt.Sprintf("%f", b)), 0666);
	if err != nil {
		log.Fatal(err)
	}
}
