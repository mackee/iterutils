package iterutils_test

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/mackee/iterutils"
)

func ExampleFromNexter() {
	s := "foo\nbar\nbaz\n"
	scanner := bufio.NewScanner(strings.NewReader(s))
	nexter := iterutils.NewNexterWithT(scanner, func(scanner *bufio.Scanner) bool {
		return scanner.Scan()
	})
	it := iterutils.FromNexter(nexter, func(nexter iterutils.NexterWithT[*bufio.Scanner]) string {
		return nexter.T().Text()
	})
	for s := range it {
		fmt.Println(s)
	}
	// Output:
	// foo
	// bar
	// baz
}
