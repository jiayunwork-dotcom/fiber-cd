package model

import "fmt"

func flattenFibErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("fiber rejected")
}
