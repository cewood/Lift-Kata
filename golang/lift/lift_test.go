package lift

import (
	"reflect"
	"testing"
)

func TestNewRequest(t *testing.T) {
	var tests = []struct {
		name     string
		floor    int
		input    []int
		expected []int
	}{
		{
			"upwards mixed order",
			0,
			[]int{3, 1, 5, 2, 4},
			[]int{1, 2, 3, 4, 5},
		},
		{
			"upwards sorted order",
			0,
			[]int{4, 5, 6, 7, 8},
			[]int{4, 5, 6, 7, 8},
		},
		{
			"upwards reversed order",
			0,
			[]int{9, 8, 7, 6, 5},
			[]int{5, 6, 7, 8, 9},
		},
		{
			"upwards mixed order with reverse",
			3,
			[]int{4, 2, 1, 5, 0},
			[]int{4, 5, 2, 1, 0},
		},
		{
			"downwards mixed order",
			6,
			[]int{3, 1, 5, 2, 4},
			[]int{5, 4, 3, 2, 1},
		},
		{
			"downwards sorted order",
			7,
			[]int{2, 3, 4, 5, 6},
			[]int{6, 5, 4, 3, 2},
		},
		{
			"downwards reversed order",
			10,
			[]int{5, 6, 7, 8, 9},
			[]int{9, 8, 7, 6, 5},
		},
		{
			"downwards mixed order with reverse",
			6,
			[]int{4, 7, 5, 9, 8},
			[]int{5, 4, 7, 8, 9},
		},
		{
			"basement floors ordered",
			5,
			[]int{4, 3, 0, -1, -2},
			[]int{4, 3, 0, -1, -2},
		},
		{
			"basement floors unorder",
			5,
			[]int{4, -1, 3, -2, 0},
			[]int{4, 3, 0, -1, -2},
		},
	}

	for _, test := range tests {
		lift := Lift{test.name, test.floor, []int{}, true}
		for _, request := range test.input {
			lift.NewRequest(request)
		}

		if !reflect.DeepEqual(lift.Requests, test.expected) {
			t.Errorf("%s: wanted '%v' but got '%v'\n", test.name, test.expected, lift.Requests)
		}
	}
}

func TestFulfillRequest(t *testing.T) {
	var tests = []struct {
		name     string
		floor    int
		input    []int
		expected []int
	}{
		{
			"upwards ordered",
			1,
			[]int{1, 2, 3, 4, 5},
			[]int{2, 3, 4, 5},
		},
		{
			"upwards sorted change direction",
			4,
			[]int{4, 5, 6, 3, 2},
			[]int{5, 6, 3, 2},
		},
		{
			"downwards ordered",
			5,
			[]int{5, 4, 3, 2, 1},
			[]int{4, 3, 2, 1},
		},
		{
			"downwards sorted change direction",
			6,
			[]int{6, 5, 4, 7, 8},
			[]int{5, 4, 7, 8},
		},
	}

	for _, test := range tests {
		lift := Lift{test.name, test.floor, test.input, true}
		lift.FulfillRequest()

		if !reflect.DeepEqual(lift.Requests, test.expected) {
			t.Errorf("%s: wanted '%v' but got '%v'\n", test.name, test.expected, lift.Requests)
		}
	}
}

func TestGetDirection(t *testing.T) {
	var tests = []struct {
		name     string
		floor    int
		input    []int
		expected Direction
	}{
		{
			"upwards",
			0,
			[]int{1, 2},
			Up,
		},
		{
			"downwards",
			7,
			[]int{6, 3, 2},
			Down,
		},
		{
			"none",
			0,
			[]int{},
			None,
		},
	}

	for _, test := range tests {
		lift := Lift{test.name, test.floor, test.input, true}

		if !reflect.DeepEqual(lift.GetDirection(), test.expected) {
			t.Errorf("%s: wanted '%v' but got '%v'\n", test.name, test.expected, lift.GetDirection())
		}
	}
}
