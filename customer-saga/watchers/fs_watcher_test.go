package watchers

import (
	"regexp"
	"testing"

	"github.com/sciensoft/fluenttests/fluent/integers"
	"github.com/sciensoft/fluenttests/fluent/strings"
)

func TestCsvDateExpressionShouldMatch(t *testing.T) {
	// Arrange
	fluent := strings.Fluent(t)
	path := "/mnt/c/Projects/.samples/GData-Pipeline/customer-saga/.files/customer_20230509.csv"
	exp := `/customer_(?P<date>\d{6,8}).csv`

	// Act
	// Assert
	fluent.It(path).
		Should().NotBeEmpty().
		And().Match(exp)
}

func TestCsvDateRegexShouldExtractNamedGroup(t *testing.T) {
	// Arrange
	fluent := integers.Fluent[int](t)
	path := "/mnt/c/Projects/.samples/GData-Pipeline/customer-saga/.files/customers_20230507.csv"
	exp := `^(?P<path>[\w\W]*)/customers_(?P<date>\d{6,8}).csv$`
	regex := regexp.MustCompile(exp)

	// Act
	names := regex.SubexpNames()
	matches := regex.FindStringSubmatch(path)
	groups := make(map[string]string)

	for i, n := range names {
		groups[n] = matches[i]
	}

	// Assert
	fluent.It(len(groups)).
		Should().BePositive()
}
