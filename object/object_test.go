package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "Cristiano Ronaldo is the goat"}
	diff2 := &String{Value: "Cristiano Ronaldo is the goat"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("Strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("Strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("Strings with different content have same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	seven1 := &Integer{Value: 7}
	seven2 := &Integer{Value: 7}
	twentyTwo1 := &Integer{Value: 22}
	twentyTwo2 := &Integer{Value: 22}

	if seven1.HashKey() != seven2.HashKey() {
		t.Errorf("Integers with same content have different hash keys")
	}

	if twentyTwo1.HashKey() != twentyTwo2.HashKey() {
		t.Errorf("Integers with same content have different hash keys")
	}

	if seven1.HashKey() == twentyTwo1.HashKey() {
		t.Errorf("Integers with different content have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("Booleans with same content have different hash keys")
	}

	if false1.HashKey() != false2.HashKey() {
		t.Errorf("Booleans with same content have different hash keys")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Errorf("Booleans with different content have same hash keys")
	}
}
