package helpers

import "time"

// this type allows it to be parsed from yaml config file
// (it implements UnmarshalYAML)
// for example:
//     - your_yaml_field: 1d3h
type Duration time.Duration

func (d *Duration) UnmarshalYAML(unmarshal func(obj interface{}) error) error {
	var s string
	var err error
	if err = unmarshal(&s); err != nil {
		return err
	}
	var td time.Duration
	if td, err = time.ParseDuration(s); err != nil {
		return err
	}
	*d = Duration(td)
	return nil
}
