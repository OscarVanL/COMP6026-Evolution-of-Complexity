package main

import (
	"fmt"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/optimisation"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	rastriginTest := [f.Rastrigin_n]uint16{}

	for i:=0; i<f.Rastrigin_n; i++ {
		rastriginTest[i] = uint16(rand.Int())
	}
	fmt.Println(rastriginTest)
	fmt.Println("Rastrigin: ", f.Rastrigin(rastriginTest))

	schwefelTest := [f.Schwefel_n]uint16{}
	for i:=0; i<f.Schwefel_n; i++ {
		schwefelTest[i] = uint16(rand.Int())
	}
	fmt.Println(schwefelTest)
	fmt.Println("Schwefel: ", f.Schwefel(schwefelTest))

	griewangkTest := [f.Griewangk_n]uint16{}
	for i:=0; i<f.Griewangk_n; i++ {
		griewangkTest[i] = uint16(rand.Int())
	}
	fmt.Println(griewangkTest)
	fmt.Println("Griewangk: ", f.Griewangk(griewangkTest))

	ackleyTest := [f.Ackley_n]uint16{}
	for i:=0; i<f.Ackley_n; i++ {
		ackleyTest[i] = uint16(rand.Int())
	}
	fmt.Println(ackleyTest)
	fmt.Println("Ackley: ", f.Ackley(ackleyTest))

	
}
