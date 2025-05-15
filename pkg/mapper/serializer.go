package mapper

import (
	"github.com/viant/xunsafe"
	"reflect"
	"sync"
	"unsafe"
)

// caches para evitar recalcular estruturas
var (
	structCache = make(map[reflect.Type]*xunsafe.Struct)
	structMux   sync.RWMutex

	fieldCache = make(map[string]*xunsafe.Field)
	fieldMux   sync.RWMutex
)

func getXStruct(t reflect.Type) *xunsafe.Struct {
	structMux.RLock()
	xs, ok := structCache[t]
	structMux.RUnlock()
	if ok {
		return xs
	}

	structMux.Lock()
	defer structMux.Unlock()
	// Double-check locking
	if xs, ok = structCache[t]; ok {
		return xs
	}
	xs = xunsafe.NewStruct(t)
	structCache[t] = xs
	return xs
}

func getXField(structType reflect.Type, fieldName string) *xunsafe.Field {
	key := structType.PkgPath() + "." + structType.Name() + "." + fieldName

	fieldMux.RLock()
	xf, ok := fieldCache[key]
	fieldMux.RUnlock()
	if ok {
		return xf
	}

	field, found := structType.FieldByName(fieldName)
	if !found {
		return nil
	}
	xf = xunsafe.NewField(field)

	fieldMux.Lock()
	fieldCache[key] = xf
	fieldMux.Unlock()

	return xf
}

// Serialize converte uma interface em map[string]interface{} ou []map[string]interface{}
func Serialize(data interface{}) interface{} {
	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		length := val.Len()
		serializedList := make([]map[string]interface{}, length)
		for i := 0; i < length; i++ {
			serializedList[i] = serializeStruct(val.Index(i))
		}
		return serializedList
	case reflect.Struct, reflect.Ptr:
		return serializeStruct(val)
	default:
		return nil
	}
}

func serializeStruct(val reflect.Value) map[string]interface{} {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	structType := val.Type()
	xStruct := getXStruct(structType)

	var structPtr unsafe.Pointer

	// Verifica se val é endereçável para evitar cópias desnecessárias
	if val.CanAddr() {
		structPtr = unsafe.Pointer(val.UnsafeAddr())
	} else {
		// Caso contrário, cria uma cópia endereçável
		copyVal := reflect.New(structType).Elem()
		copyVal.Set(val)
		structPtr = unsafe.Pointer(copyVal.UnsafeAddr())
	}

	result := make(map[string]interface{}, len(xStruct.Fields))
	for _, field := range xStruct.Fields {
		fieldValue := field.ValuePointer(structPtr)
		result[field.Name] = reflect.NewAt(field.Type, fieldValue).Elem().Interface()
	}

	return result
}

// Deserialize preenche uma struct a partir de um map[string]interface{} usando xunsafe
func Deserialize(input map[string]interface{}, target interface{}) {
	val := reflect.ValueOf(target).Elem()
	typeOfVal := val.Type()
	structPtr := unsafe.Pointer(val.UnsafeAddr())

	for key, v := range input {
		if field, ok := typeOfVal.FieldByName(key); ok {
			xField := getXField(typeOfVal, key)

			if field.Type.Kind() == reflect.Struct {
				if subMap, ok := v.(map[string]interface{}); ok {
					subStruct := reflect.New(field.Type).Elem()
					Deserialize(subMap, subStruct.Addr().Interface())
					xField.SetValue(structPtr, subStruct.Interface())
				}
			} else {
				xField.SetValue(structPtr, v)
			}
		}
	}
}
