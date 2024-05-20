package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"strconv"
	"math"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"strings")

type Coord struct {
	x float64
	y float64
}

func costFunction(predictions []float64, points []Coord) float64 {
	var diff float64;
	for i := range predictions {
		d := predictions[i] - points[i].y;
		diff += math.Pow(d, 2.0);
	}
	loss := diff / float64(len(points));
	return loss;
}

func findMinMax(points []Coord) (minX, maxX, minY, maxY float64) {
	minX, maxX = points[0].x, points[0].x;
	minY, maxY = points[0].y, points[0].y;
	for _, point := range points {
		if point.x < minX {
			minX = point.x;
		}
		if point.x > maxX {
			maxX = point.x;
		}
		if point.y < minY {
			minY = point.y;
		}
		if point.y > maxY {
			maxY = point.y;
		}
	}
	return;
}

func minMaxNormalization(points []Coord) []Coord {
	minX, maxX, minY, maxY := findMinMax(points);

	normalizedPoints := make([]Coord, len(points))
	for i, point := range points {
		normalizedX := (point.x - minX) / (maxX - minX);
		normalizedY := (point.y - minY) / (maxY - minY);
		normalizedPoints[i] = Coord{x: normalizedX, y: normalizedY};
	}
	return normalizedPoints;
}

func denormalizeTheta(theta0, theta1 float64, originalPoints []Coord) (newTheta0, newTheta1 float64) {
	minX, maxX, minY, maxY := findMinMax(originalPoints);
	rangeX := maxX - minX;
	rangeY := maxY - minY;

	newTheta1 = theta1 * rangeY / rangeX;
	newTheta0 = theta0 * rangeY + minY - newTheta1 * minX;
	return
}

func calcGradientTeta1(points []Coord, teta0, teta1 float64) float64 {
	var diff float64 = 0.0;
	for i := 0; i < len(points); i++ {
		diff += ((teta0 + teta1 * points[i].y) - points[i].x) * points[i].y;
	}
	return (diff / float64(len(points)));
}

func calcGradientTeta0(points []Coord, teta0, teta1 float64) float64 {
	var diff float64 = 0.0;
	for i := 0; i < len(points); i++ {
		diff += (teta0 + teta1 * points[i].y) - points[i].x;
	}
	return (diff / float64(len(points)));
}

func f(points []Coord, teta0 float64, teta1 float64) []float64 {
	var result []float64
	for _, val := range points {
		y := teta1 * val.x + teta0;
		result = append(result, y);
	}
	return (result);
}

func main() {
	file, err := os.Open("./data.csv");
	if (err != nil) {
		log.Fatal(err);
	}
	defer file.Close()

	scanner := bufio.NewScanner(file);
	var points []Coord;
	var n = 0;
	var km, price float64;
	var pointValues []opts.ScatterData;
	for scanner.Scan() {
		if (scanner.Text() == "km,price") {
			continue ;
		}
		var line = strings.Split(scanner.Text(), ",");
		if (len(line) < 2) {
			log.Fatal("Wrong data format");
		}
		km, err = strconv.ParseFloat(line[0], 10);
		price, err = strconv.ParseFloat(line[1], 10);
		if (err != nil) {
			log.Fatal(err);
		}
		points = append(points, Coord{x: km, y: price});

		pointValues = append(pointValues, opts.ScatterData{Value: []interface{}{km, price} , SymbolSize: 10, XAxisIndex: int(km), YAxisIndex: int(price)});
		n++;
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err);
	}

	var iteration = 100000;
	var tetha0 float64= 0.0;
	var tetha1 float64= 0.0;
	var learningRate = 0.01;

	pointsNormed := minMaxNormalization(points);
	for i := 0; i < iteration; i++ {
		gradTeta0 := calcGradientTeta0(pointsNormed, tetha0, tetha1);
		gradTeta1 := calcGradientTeta1(pointsNormed, tetha0, tetha1);
		tetha0 -= gradTeta0 * learningRate;
		tetha1 -= gradTeta1 * learningRate;
	//	prediction := f(pointsNormed, tetha0, tetha1);
	//	cost := costFunction(prediction, pointsNormed);
	//	fmt.Println("COST :", cost);
	}

	a, b := denormalizeTheta(tetha0, tetha1, points);
	fmt.Println("TETA0", a, "TETA1", b);
	err = os.WriteFile("./tetaInfo", []byte(fmt.Sprintf("%f", a) + "," + fmt.Sprintf("%f", b)), 0666);
	if err != nil {
		log.Fatal(err)
	}
	
	line := charts.NewLine();
    line.SetGlobalOptions(
        charts.WithTitleOpts(opts.Title{Title: "Linear Regression"}),
        charts.WithXAxisOpts(opts.XAxis{Name: "Mileage"}),
        charts.WithYAxisOpts(opts.YAxis{Name: "Price"}),
    );
	scatter := 	charts.NewScatter();
	var lineInfo []opts.LineData;
	for _, point := range points {
		lineInfo = append(lineInfo, opts.LineData{Value: []interface{}{point.x, a + (b * point.x)}});
	}
     line.AddSeries("MyLine", lineInfo,
            charts.WithLineChartOpts(
                opts.LineChart{Smooth: true},
            ),
        )
	scatter.AddSeries("Points", pointValues);
    f, err := os.Create("line.html")
    if err != nil {
        panic(err)
    }
	line.Overlap(scatter);
    line.Render(f)
}
