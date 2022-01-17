package tool

import (
	"fmt"
	"testing"
)

func Test_IteratorOut(t *testing.T) {
	arrStr := IteratorOut{IteratorSingleOut{Value: "a"}, IteratorSingleOut{Value: "b"}, IteratorSingleOut{Value: "c"}, IteratorSingleOut{Value: "d"}}
	arrInt := IteratorOut{IteratorSingleOut{Value: "1"}, IteratorSingleOut{Value: "2"}, IteratorSingleOut{Value: "3"}}
	arrChar := IteratorOut{IteratorSingleOut{Value: "一"}, IteratorSingleOut{Value: "二"}, IteratorSingleOut{Value: "三"}}
	iter := NewSliceIterator(arrStr).ToAssemblyIterator().BuildNewAssemblyIterator(NewSliceIterator(arrInt)).BuildNewAssemblyIterator(NewSliceIterator(arrChar))
	for {
		if value, ok := iter.Next(); ok {
			fmt.Println(value)
		} else {
			break
		}
	}
}
