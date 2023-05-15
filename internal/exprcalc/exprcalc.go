package exprcalc

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	opBr byte = '('
	clBr byte = ')'
	dot  byte = '.'
	add  byte = '+'
	sub  byte = '-'
	mul  byte = '*'
	div  byte = '/'
	pow  byte = '^'
	mod  byte = '%'
)

type expressionCalculator interface {
	run()
}

func PowerOn(e expressionCalculator) {
	e.run()
}

// Engine.
type calculatorCircuitBoard struct {
	input  string
	result string

	errNo         int
	lowBound      int
	highBound     int
	openBrackets  int
	closeBrackets int

	constants map[string]float64
}

// Engine Constructor.
func NewCalculatorCircuitBoard() *calculatorCircuitBoard {
	constants := map[string]float64{}

	constants["e"] = math.E
	constants["pi"] = math.Pi
	constants["phi"] = math.Phi

	return &calculatorCircuitBoard{constants: constants}
}

func (c *calculatorCircuitBoard) run() {
	for {
		c.resetWorkValues()

		fmt.Print("> ")

		err := c.readInput()

		if err != nil {
			c.writeError()
			continue
		}

		c.filterInput()

		if c.isUserWantExit() {
			break
		}

		if c.isWrongCharsFound() ||
			c.isInequalOpenCloseBackets() {
			c.writeError()
			continue
		}

		c.changeConstantsNamesToValues()
		c.workWithExpr()
		c.writeResult()
	}

	fmt.Println("Exited.")
}

func (c *calculatorCircuitBoard) resetWorkValues() {
	c.errNo = 0
	c.lowBound = 0
	c.highBound = 0
	c.openBrackets = 0
	c.closeBrackets = 0
}

func (c *calculatorCircuitBoard) readInput() error {
	var err error
	reader := bufio.NewReader(os.Stdin)
	c.input, err = reader.ReadString('\n')

	if err != nil {
		c.errNo = 3
		return err
	}

	return nil
}

func (c *calculatorCircuitBoard) writeError() {
	var errorVerbose string

	switch c.errNo {
	case 1:
		errorVerbose = "entered wrong symbol"
	case 2:
		errorVerbose = "amount of the opening and closing brackets is not equal"
	case 3:
		errorVerbose = "input buffer error"
	case 4:
		errorVerbose = "wrong number notation"
	case 5:
		errorVerbose = "zero division is not allowed"
	default:
		errorVerbose = "no error"
	}

	fmt.Println("Error:", errorVerbose)
}

func (c *calculatorCircuitBoard) isUserWantExit() bool {
	if c.input == "q" || c.input == "quit" || c.input == "exit" {
		return true
	}

	return false
}

func (c *calculatorCircuitBoard) filterInput() {
	c.input = strings.Replace(c.input, " ", "", -1)
	c.input = strings.Replace(c.input, "\r", "", -1)
	c.input = strings.Replace(c.input, "\n", "", -1)
	c.input = strings.Replace(c.input, ",", ".", -1)
}

func (c *calculatorCircuitBoard) isWrongCharsFound() bool {
	var elem byte

	for i := 0; i < len(c.input); i++ {
		elem = c.input[i]

		if !((elem == dot) || (elem == opBr) || (elem == clBr) ||
			c.isAvailableDigit(elem) || c.isAvailableOp(elem)) {
			c.errNo = 1
			return true
		}
	}

	return false
}

func (c *calculatorCircuitBoard) isInequalOpenCloseBackets() bool {
	c.openBrackets = 0
	c.closeBrackets = 0

	for i := 0; i < len(c.input); i++ {
		elem := string(c.input[i])
		if elem == ")" {
			c.closeBrackets++
		} else if elem == "(" {
			c.openBrackets++
		}
	}

	if c.openBrackets != c.closeBrackets {
		c.errNo = 2
		return true
	}

	return false
}

func (c *calculatorCircuitBoard) workWithExpr() {
	for c.closeBrackets > 0 {
		c.performBracketOp()
		c.closeBrackets--
	}

	c.result = c.input
	c.calculateExpr()
	c.input = c.result
}

func (c *calculatorCircuitBoard) changeConstantsNamesToValues() {
	for name, val := range c.constants {
		var count int
		var nameLen, valLen int = len(name), len(c.input)

		for i := 0; i < valLen-(nameLen-1); i++ {
			var word string = c.input[i : i+nameLen]
			if word == name {
				count++
			}
		}

		for count > 0 {
			for i := 0; i < valLen-(nameLen-1); i++ {
				var word string = c.input[i : i+nameLen]
				if word == name {
					c.input = c.input[:i] + strconv.FormatFloat(val, 'f', -1, 64) +
						c.input[i+nameLen:]
					break
				}
			}
			count--
		}
	}
}

func (c *calculatorCircuitBoard) performBracketOp() {
	lenOfValue := len(c.input)

	for j := 0; j < lenOfValue; j++ {
		elem := string(c.input[j])
		if elem == ")" {
			c.closeBrackets = j

			for i := j - 1; i >= 0; i-- {
				elem := string(c.input[i])
				if elem == "(" {
					c.openBrackets = i
					break
				}
			}
			break
		}
	}

	c.result = c.input[c.openBrackets+1 : c.closeBrackets]
	c.calculateExpr()
	c.input = c.input[:c.openBrackets] + c.result + c.input[c.closeBrackets+1:]
}

func (c *calculatorCircuitBoard) calculateExpr() {
	var lowestOps, lowOps, highOps, highestOps int

	for i := 0; i < len(c.result); i++ {
		elem := c.result[i]
		if elem == pow {
			highestOps++
		} else if (elem == mul) || (elem == div) {
			highOps++
		} else if (elem == add) || (elem == sub) {
			lowOps++
		} else if elem == mod {
			lowestOps++
		}
	}

	for highestOps > 0 {
		c.performOp([]byte{pow})
		highestOps--
	}

	for highOps > 0 {
		c.performOp([]byte{mul, div})
		highOps--
	}

	for lowOps > 0 {
		c.performOp([]byte{add, sub})
		lowOps--
	}

	for lowestOps > 0 {
		c.performOp([]byte{mod})
		lowestOps--
	}
}

func (c *calculatorCircuitBoard) performOp(opsToCalc []byte) {
	var passToDoOp bool

	for j := 0; j < len(c.result); j++ {
		elem := c.result[j]

		for _, op := range opsToCalc {
			if elem == op {
				passToDoOp = true
			}
		}

		if passToDoOp {
			var i int

			for i = j - 1; i >= 0; i-- {
				elem := c.result[i]

				if c.isAvailableOp(elem) {
					break
				}
			}

			if i > 0 {
				c.lowBound = i + 1
			} else {
				c.lowBound = 0
			}

			for i = j + 1; i < len(c.result); i++ {
				elem := c.result[i]

				if c.isAvailableOp(elem) {
					break
				}
			}

			c.highBound = i - 1

			break
		}
	}

	binOp := c.doBinaryOp(c.result[c.lowBound : c.highBound+1])
	c.result = c.result[:c.lowBound] + binOp + c.result[c.highBound+1:]
}

func (c *calculatorCircuitBoard) doBinaryOp(binExpr string) string {
	var operatorChar byte
	var operatorPos int
	var resultf float64
	var err error

	for i := 0; i < len(binExpr); i++ {
		elem := binExpr[i]

		if c.isAvailableOp(elem) {
			operatorChar = elem
			operatorPos = i
			break
		}
	}

	operand1 := binExpr[:operatorPos]
	operand2 := binExpr[operatorPos+1:]

	operand1f, err := strconv.ParseFloat(operand1, 64)
	operand2f, err := strconv.ParseFloat(operand2, 64)

	if err != nil {
		c.errNo = 4
		return ""
	}

	switch operatorChar {
	case add:
		resultf = operand1f + operand2f
	case sub:
		resultf = operand1f - operand2f
	case mul:
		resultf = operand1f * operand2f
	case div:
		if operand2f != 0.0 {
			resultf = operand1f / operand2f
		} else {
			c.errNo = 5
			return ""
		}
	case pow:
		resultf = math.Pow(operand1f, operand2f)
	case mod:
		resultf = math.Mod(operand1f, operand2f)
	default:
		c.errNo = 1
		return ""
	}

	if math.Mod(resultf, 1.0) > 0.0 {
		return strconv.FormatFloat(resultf, 'f', 15, 64)
	}
	return fmt.Sprintf("%.0f", resultf)
}

func (c *calculatorCircuitBoard) isAvailableOp(symbol byte) bool {
	if (symbol == add) || (symbol == sub) || (symbol == mul) ||
		(symbol == div) || (symbol == pow) || (symbol == mod) {
		return true
	}

	return false
}

func (c *calculatorCircuitBoard) isAvailableDigit(symbol byte) bool {
	if (symbol >= '0') && (symbol <= '9') {
		return true
	}

	return false
}

func (c *calculatorCircuitBoard) writeResult() {
	if strings.Contains(c.input, ".") {
		for {
			if string(c.input[len(c.input)-1]) == "0" {
				c.input = c.input[:len(c.input)-1]
			} else {
				break
			}
		}
	}

	if c.input == "0." || c.input == "." {
		c.input = "0"
	}

	fmt.Println(c.input)
}
