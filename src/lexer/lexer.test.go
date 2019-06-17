import "testing"

func TestNextToken(t *testing.T) {
	input := `five <- 5;ten <- 10

add <- function(x, y) x + y

result <- add(five, ten)`

	test := []struct {
		expect
	}
}
