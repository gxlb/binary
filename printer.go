package binary

import (
	"bytes"
	"reflect"
	"strconv"
)

// use "%##v" "%++v" to format a value with indented-multi-line style string
const (
	maxArrayElemInLine = 5
)

type ShowOption uint

const (
	OptionSingleLine ShowOption = 1 << iota
)

//ShowString shows more effective info of value than fmt(%#v)
func ShowString(data interface{}) string {
	var p Printer
	p.Init(0)
	return p.ShowString(data)
}

func ShowSingleLineString(data interface{}) string {
	var p Printer
	p.Init(OptionSingleLine)
	return p.ShowString(data)
}

type Printer struct {
	option ShowOption
	buff   bytes.Buffer
}

func (p *Printer) Init(option ShowOption) {
	p.option = option
	p.buff.Truncate(0)
}

func (p *Printer) ShowString(data interface{}) string {
	p.show(reflect.ValueOf(data), 0)
	return p.buff.String()
}

func (p *Printer) show(v reflect.Value, depth int) {
	t := v.Type()
	switch k := t.Kind(); k {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		p.buff.WriteString(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		p.buff.WriteString(strconv.FormatUint(v.Uint(), 10))
	case reflect.Bool:
		p.buff.WriteString(strconv.FormatBool(v.Bool()))
	case reflect.Float32:
		p.float32(v.Float())
	case reflect.Float64:
		p.float64(v.Float())
	case reflect.Complex64:
		c := v.Complex()
		p.buff.WriteString("complex(")
		p.float32(real(c))
		p.buff.WriteByte(',')
		p.float32(imag(c))
		p.buff.WriteByte(')')
	case reflect.Complex128:
		c := v.Complex()
		p.buff.WriteString("complex(")
		p.float64(real(c))
		p.buff.WriteByte(',')
		p.float64(imag(c))
		p.buff.WriteByte(')')
	case reflect.String:
		p.buff.WriteByte('`')
		p.buff.WriteString(v.String())
		p.buff.WriteByte('`')
	case reflect.Slice, reflect.Array:
		elemKind := t.Elem().Kind()
		p.buff.WriteString(t.String())
		p.buff.WriteByte('{')
		lines := 0
		for i, size := 0, v.Len(); i < size; i++ {
			if p.newLineInArray(elemKind, i, size) {
				lines++
				p.newLine(depth + 1)
			}
			p.show(v.Index(i), depth+1)
			p.buff.WriteByte(',')
			//p.buff.WriteByte(' ')
		}
		if lines > 0 {
			p.newLine(depth)
		}
		p.buff.WriteByte('}')

	case reflect.Map:
		p.buff.WriteString(t.String())
		p.buff.WriteByte('{')
		keys := v.MapKeys()
		for i, l := 0, len(keys); i < l; i++ {
			p.newLine(depth + 1)
			key := keys[i]
			p.show(key, depth+1)
			p.buff.WriteByte(':')
			p.buff.WriteByte(' ')
			p.show(v.MapIndex(key), depth+1)
			p.buff.WriteByte(',')
			//p.buff.WriteByte(' ')

		}
		p.newLine(depth)
		p.buff.WriteByte('}')

	case reflect.Struct:
		p.buff.WriteString(t.String())
		p.buff.WriteByte('{')
		for i, n := 0, v.NumField(); i < n; i++ {
			p.newLine(depth + 1)
			f := v.Field(i)
			ft := t.Field(i)
			p.buff.WriteString(ft.Name)
			p.buff.WriteByte(':')
			p.buff.WriteByte(' ')
			p.show(f, depth+1)
			p.buff.WriteByte(',')
			//p.buff.WriteByte(' ')
		}
		p.newLine(depth)
		p.buff.WriteByte('}')

	case reflect.Ptr:
		if v.IsNil() {
			p.buff.WriteString("nil")
		} else {
			pv := v.Elem()
			switch v.Kind() {
			case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
				p.buff.WriteByte('&')
				p.show(pv, depth+1)
			default:
			}
		}

	default:

	}
}

func (p *Printer) float32(f float64) {
	p.buff.WriteString(strconv.FormatFloat(f, 'f', 6, 32))
}
func (p *Printer) float64(f float64) {
	p.buff.WriteString(strconv.FormatFloat(f, 'f', 6, 64))
}

func (p *Printer) newLineInArray(elemKind reflect.Kind, index, size int) bool {
	switch elemKind {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct,
		reflect.String, reflect.Ptr, reflect.Complex64, reflect.Complex128:
		return true
	default:
		if index == 0 && size >= maxArrayElemInLine || index > 0 && index%maxArrayElemInLine == 0 {
			return true
		}
	}
	return false
}

func (p *Printer) newLine(depth int) {
	if p.option&OptionSingleLine != 0 { //single line
		return
	}
	p.buff.WriteByte('\n')
	for i := 0; i < depth; i++ {
		p.buff.WriteString("    ")
	}
}
