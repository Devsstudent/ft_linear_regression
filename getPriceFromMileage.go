package main

import (
	"fmt"
	"os"
	"strconv"
	"log"
	"bufio"
	"strings"
)

func main() {
	file, err := os.Open("./tetaInfo");
	if (err != nil) {
		log.Fatal(err);
	}
	defer file.Close();
	var theta0, theta1 float64;

	scanner := bufio.NewScanner(file);
	var i = 0;
	for scanner.Scan() {
		if (i == 1) {
			log.Fatal("too much line in tetaInfo");
		}
		var line = strings.Split(scanner.Text(), ",");
		if (len(line) < 2) {
			log.Fatal("Wrong tetaInfo");
		}
		theta0, err = strconv.ParseFloat(line[0], 10);
		theta1, err = strconv.ParseFloat(line[1], 10);
		if (err != nil) {
			log.Fatal("Error theta Format");
		}
		i++;
	}
	fmt.Printf("Enter a mileage :");
	var input string
	_, err = fmt.Scanf("%s\n", &input)
	if (err != nil) {
		log.Fatal("error reading input");
	}
	fmt.Println(theta0, theta1, input);
	var mileage float64;
	mileage, err = strconv.ParseFloat(input, 10);
	fmt.Println(theta0 + (theta1 * mileage));
	//var a = theta0 + theta1;
	return ;
}
