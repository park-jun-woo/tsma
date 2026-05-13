//ff:func feature=cli type=helper control=sequence
//ff:what Prints the aggregate progress summary with percentage breakdown
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

func showOverallStatus(sess *model.Session) error {
	s := sess.Summary
	total := float64(s.Total)
	if total == 0 {
		fmt.Println("No endpoints found.")
		return nil
	}

	fmt.Printf("%d endpoints\n", s.Total)
	fmt.Printf("DONE:    %3d (%.1f%%)\n", s.Done, float64(s.Done)/total*100)
	fmt.Printf("PARTIAL: %3d (%.1f%%)\n", s.Partial, float64(s.Partial)/total*100)
	fmt.Printf("TODO:    %3d (%.1f%%)\n", s.Todo, float64(s.Todo)/total*100)

	return nil
}
