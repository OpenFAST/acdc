package main

import "testing"

func TestMinCostAssignment(t *testing.T) {

	C := [][]int{
		{20, 15, 18, 20, 25},
		{18, 20, 12, 14, 15},
		{21, 23, 25, 27, 25},
		{17, 18, 21, 23, 20},
		{18, 18, 16, 19, 20},
	}

	pathExp := [][2]int{{0, 1}, {1, 3}, {2, 0}, {3, 4}, {4, 2}}
	path, err := MinCostAssignment(C)
	if err != nil {
		t.Fatal(err)
	}
	if act, exp := len(path), len(pathExp); act != exp {
		t.Fatalf("len(path) = %v, expected %v", act, exp)
	}
	for i := range pathExp {
		if act, exp := path[i], pathExp[i]; act != exp {
			t.Fatalf("path[%d] = %v, expected %v", i, act, exp)
		}
	}

	C = [][]int{
		{9, 22, 58, 11, 19},
		{43, 78, 72, 50, 63},
		{41, 28, 91, 37, 45},
		{74, 42, 27, 49, 39},
		{36, 11, 57, 22, 25},
	}

	pathExp = [][2]int{{0, 3}, {1, 0}, {2, 1}, {3, 2}, {4, 4}}
	path, err = MinCostAssignment(C)
	if err != nil {
		t.Fatal(err)
	}
	if act, exp := len(path), len(pathExp); act != exp {
		t.Fatalf("len(path) = %v, expected %v", act, exp)
	}
	for i := range pathExp {
		if act, exp := path[i], pathExp[i]; act != exp {
			t.Fatalf("path[%d] = %v, expected %v", i, act, exp)
		}
	}
}
