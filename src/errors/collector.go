package errors

import "strings"

func NewCollector() *Collector {
	return &Collector{}
}

// Collector helps to collect errors without interruption and check them afterwords.
// A common use-case for it is a loop where you don't want to interrupt the process but have to return
// errors if any occurred in the loop.
type Collector struct {
	errors []error
}

// Add collects an error
func (c *Collector) Add(err error) {
	c.errors = append(c.errors, err)
}

// HasErrors returns true if at least one error was collected
func (c *Collector) HasErrors() bool {
	for _, err := range c.errors {
		if err != nil {
			return true
		}
	}

	return false
}

// Error joins all collected errors together and returns as one string
func (c *Collector) Error() string {
	var result []string
	for _, err := range c.errors {
		if err != nil {
			result = append(result, err.Error())
		}
	}

	return strings.Join(result, " | ")
}

// Errors returns a slice of all collected errors
func (c *Collector) Errors() []error {
	return c.errors
}
