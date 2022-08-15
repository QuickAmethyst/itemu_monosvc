package config

import (
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
)

var decoderConfig = func(m *mapstructure.DecoderConfig) {
	m.Metadata = nil
	m.WeaklyTypedInput = true
	m.DecodeHook = mapstructure.ComposeDecodeHookFunc(
		func(layout string) mapstructure.DecodeHookFunc {
			return func(
				f reflect.Type,
				t reflect.Type,
				data interface{},
			) (interface{}, error) {
				if f.Kind() != reflect.String || t != reflect.TypeOf(time.Time{}) {
					return data, nil
				}

				// Convert it by parsing
				strTime := data.(string)
				if strTime != "" {
					return time.Parse(layout, data.(string))
				}

				return time.Time{}, nil
			}
		}(time.RFC3339),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)
}
