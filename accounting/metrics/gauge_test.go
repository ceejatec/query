package metrics

import (
	"fmt"
	"testing"
)

func BenchmarkGuage(b *testing.B) {
	g := NewGauge()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Update(int64(i))
	}
}

func TestGauge(t *testing.T) {
	g := NewGauge()
	g.Update(int64(47))
	if v := g.Value(); 47 != v {
		t.Errorf("g.Value(): 47 != %v\n", v)
	}
}

func TestGetOrRegisterGauge(t *testing.T) {
	r := NewRegistry()
	NewRegisteredGauge("foo", r).Update(47)
	if g := GetOrRegisterGauge("foo", r); 47 != g.Value() {
		t.Fatal(g)
	}
}

func ExampleGetOrRegisterGauge() {
	m := "server.bytes_sent"
	g := GetOrRegisterGauge(m, nil)
	g.Update(47)
	fmt.Println(g.Value()) // Output: 47
}
