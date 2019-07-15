package scanner

import (
	"fmt"
	"testing"
)

func TestScanner(t *testing.T) {
	source := `var x = 3; while(true) continue;`
	scanner := NewScanner(source)
	scanner.scan()
	fmt.Println(source)
	fmt.Println(scanner.Tokens)
}
