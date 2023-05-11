package workers

import (
	"reflect"
	"testing"

	"github.com/sciensoft/fluenttests/fluent/contracts"
)

func TestCreateWorkerFactoryShouldProvideIWorkerInstance(t *testing.T) {
	// Arrange
	fluent := contracts.Fluent[IWorker](t)
	expectedType := reflect.TypeOf(&CsvProcessingWorker{})

	// Act
	worker := CreateWorkerFactory[*CsvProcessingWorker]()

	// Assert
	fluent.It(worker).
		Should().BeOfType(expectedType).
		And().HaveMethod("ExecuteWithContext")
}
