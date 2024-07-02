package fluentsql

import "fmt"

// Fetch clause
type Fetch struct {
	Fetch  int
	Offset int
}

func (f *Fetch) String() string {
	if f.Fetch > 0 || f.Offset > 0 {
		return fmt.Sprintf("OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", f.Offset, f.Fetch)
	}

	return ""
}
