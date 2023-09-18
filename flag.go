package hcpvaultengine

import "strconv"

type boolValue struct {
	target *bool
}

func (v *boolValue) String() string {
	return strconv.FormatBool(*v.target)
}

func (v *boolValue) Set(s string) error {
	value, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*v.target = value
	return nil
}
