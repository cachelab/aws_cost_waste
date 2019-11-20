//+build !test

package main

import (
	"log"
	"os"
)

const name = "aws_cost_waste"
const version = "1.0.0"

func main() {
	svc := NewService()

	err := svc.Init()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}

	err = svc.RunEbsReport()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}

	err = svc.RunElbReport()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}
}
