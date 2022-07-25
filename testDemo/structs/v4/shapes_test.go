package main

import "testing"

//表格驱动测试在我们要创建一系列相同测试方式的测试用例时很有用
// func TestArea(t *testing.T) {
// 	areaTests := []struct {
// 		shape Shape
// 		want  float64
// 	}{
// 		// {Rectangle{12, 6}, 72.0},
// 		// {Circle{10}, 314.1592653589793},
// 		// {Triangle{12, 6}, 36.0},
// 		{shape: Rectangle{Width: 12, Height: 6}, want: 72.0},
// 		{shape: Circle{Radius: 10}, want: 314.1592653589793},
// 		{shape: Triangle{Base: 12, Height: 6}, want: 36.0},
// 	}

// 	for _, tt := range areaTests {
// 		got := tt.shape.Area()
// 		if got != tt.want {
// 			t.Errorf("%#v got %g want %g", tt.shape, got, tt.want)
// 		}
// 	}

// }

func TestArea(t *testing.T) {

	areaTests := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		{name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0},
		{name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
		{name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
	}

	for _, tt := range areaTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.hasArea {
				t.Errorf("%#v got %.2f want %.2f", tt.shape, got, tt.hasArea)
			}
		})

	}

}
