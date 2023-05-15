package main

import ec "app/internal/exprcalc"

func main() {
	calculatorCircuitBoard := ec.NewCalculatorCircuitBoard()

	ec.PowerOn(calculatorCircuitBoard)
}
