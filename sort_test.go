package hero

import (
	"reflect"
	"testing"
)

func TestAddVertices(t *testing.T) {
	sort := newSort()

	vertices := []string{"a", "b", "c", "d", "e"}
	for _, vertex := range vertices {
		sort.addVertices(vertex)
	}

	if len(sort.vertices) != len(vertices) {
		t.Fail()
	}

	for _, vertex := range vertices {
		if _, ok := sort.vertices[vertex]; !ok {
			t.Fail()
		}
	}
}

func TestAddEdge(t *testing.T) {
	sort := newSort()

	edges := map[string][]string{
		"a": {"1"},
		"b": {"2", "3"},
		"c": {"4", "5", "6"},
	}

	// init graph.
	for from, tos := range edges {
		for _, to := range tos {
			sort.addEdge(from, to)
		}
	}

	if len(sort.graph) != len(edges) {
		t.Fail()
	}

	for from, tos := range edges {
		if v, ok := sort.graph[from]; !ok || len(v) != len(tos) {
			t.Fail()
		}

		for _, to := range tos {
			if _, ok := sort.graph[from][to]; !ok {
				t.Fail()
			}
		}
	}
}

func TestCollect(t *testing.T) {
	sort := &sort{
		v: map[string]int{
			"a": 0,
			"b": 1,
			"c": 2,
		},
	}

	var queue []string
	sort.collect(&queue)

	if !reflect.DeepEqual(queue, []string{"a"}) ||
		!reflect.DeepEqual(sort.v, map[string]int{"b": 1, "c": 2}) {
		t.Fail()
	}
}

func TestSort(t *testing.T) {
	sort := newSort()

	vertices := []string{"a", "b", "c", "d", "e"}
	for _, vertex := range vertices {
		sort.addVertices(vertex)
	}

	edges := map[string][]string{
		"a": {"b"},
		"b": {"c"},
		"c": {"d"},
		"d": {"e"},
	}

	for from, tos := range edges {
		for _, to := range tos {
			sort.addEdge(from, to)
		}
	}

	if !reflect.DeepEqual(sort.sort(), vertices) {
		t.Fail()
	}
}
