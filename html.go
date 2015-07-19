package webapp

import (
	"fmt"
	"time"
)

// Html is bounded in the ViewModel to allow helper functions to be called from the view.
type Html struct {
	ViewModel *ViewModel
	Context   *Context
}

// Human() display's a basic type to a more human readible format.
func (this *Html) Human(value interface{}) string {
	switch val := value.(type) {
	case time.Time:
		return val.Format("2006-01-02 15:04")
	default:
		return fmt.Sprintf("%v", value)
	}
	return "human"
}
