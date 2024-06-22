package jsontime

import (
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

var (
	Default          = jsoniter.ConfigCompatibleWithStandardLibrary
	defaultExtension CustomTimeExtension
)

var (
	defaultFormat = time.RFC3339
	defaultLocale = time.Local
)

func init() {
	defaultExtension = CustomTimeExtension{
		Format:   time.RFC3339,
		Location: time.Local,
	}

	Default.RegisterExtension(&defaultExtension)
}

func SetDefaultTimeFormat(timeFormat string, timeLocation *time.Location) {
	defaultExtension.Format = timeFormat
	defaultExtension.Location = timeLocation
}

type CustomTimeExtension struct {
	jsoniter.DummyExtension
	Format   string
	Location *time.Location
}

func New(timeFormat string, timeLocation *time.Location) jsoniter.API {
	api := jsoniter.ConfigCompatibleWithStandardLibrary
	api.RegisterExtension(&CustomTimeExtension{
		Format:   timeFormat,
		Location: timeLocation,
	})
	return api
}

func (extension *CustomTimeExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		var typeErr error
		var isPtr bool
		typeName := binding.Field.Type().String()
		if typeName == "time.Time" {
			isPtr = false
		} else if typeName == "*time.Time" {
			isPtr = true
		} else {
			continue
		}

		binding.Encoder = &funcEncoder{fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			if typeErr != nil {
				stream.Error = typeErr
				return
			}

			var tp *time.Time
			if isPtr {
				tpp := (**time.Time)(ptr)
				tp = *(tpp)
			} else {
				tp = (*time.Time)(ptr)
			}

			if tp != nil {
				lt := tp.In(extension.Location)
				str := lt.Format(extension.Format)
				stream.WriteString(str)
			} else {
				stream.Write([]byte("null"))
			}
		}}
		binding.Decoder = &funcDecoder{fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if typeErr != nil {
				iter.Error = typeErr
				return
			}

			str := iter.ReadString()
			var t *time.Time
			if str != "" {
				var err error
				tmp, err := time.ParseInLocation(extension.Format, str, extension.Location)
				if err != nil {
					iter.Error = err
					return
				}
				t = &tmp
			} else {
				t = nil
			}

			if isPtr {
				tpp := (**time.Time)(ptr)
				*tpp = t
			} else {
				tp := (*time.Time)(ptr)
				if tp != nil && t != nil {
					*tp = *t
				}
			}
		}}
	}
}

type funcDecoder struct {
	fun jsoniter.DecoderFunc
}

func (decoder *funcDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	decoder.fun(ptr, iter)
}

type funcEncoder struct {
	fun         jsoniter.EncoderFunc
	isEmptyFunc func(ptr unsafe.Pointer) bool
}

func (encoder *funcEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	encoder.fun(ptr, stream)
}

func (encoder *funcEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	if encoder.isEmptyFunc == nil {
		return false
	}
	return encoder.isEmptyFunc(ptr)
}
