package dispersion

import "fmt"

func flattenZeroErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s", err.Error())
}
