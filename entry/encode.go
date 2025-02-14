package entry

import (
  "fmt"
)

func (entry Entry) Encode() string {
	return fmt.Sprintf(`%v,%s: %s`,
		entry.Time.Format(timeLayout),
		entry.Category,
		entry.Note,
	)
}
