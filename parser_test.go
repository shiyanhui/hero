package hero

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
)

const rootDir = "/tmp/gohero"

const indexHTML = `
<!DOCTYPE html>
<html>
  <head>
  </head>
  <body>
    <%@ body { %>
    <% } %>
  </body>
</html>
`

const itemHTML = `
<div>
    <a href="/user/<%= user %>">
        <%== user %>
    </a>
</div>
`

const listHTML = `
<%: func UserList(userList []string) *bytes.Buffer %>

<%!
	func Add(a, b int) int {
		return a + b
	}
%>

<%~ "index.html" %>

<%@ body { %>
    <%# this is note %>
    <% for _, user := range userList { %>
        <%+ "item.html" %>
    <% } %>
<% } %>
`

func init() {
	_, err := os.Stat(rootDir)
	if os.IsExist(err) {
		os.RemoveAll(rootDir)
	}
	os.Mkdir(rootDir, os.ModePerm)

	items := []struct {
		name    string
		content string
	}{
		{"index.html", indexHTML},
		{"item.html", itemHTML},
		{"list.html", listHTML},
	}

	for _, item := range items {
		ioutil.WriteFile(
			path.Join(rootDir, item.name),
			[]byte(item.content),
			os.ModePerm,
		)
	}
}

func TestNewNode(t *testing.T) {
	var n *node

	cases := []struct {
		t     uint8
		chunk []byte
	}{
		{TypeCode, []byte{1, 2}},
		{TypeExtend, []byte{3, 4}},
		{TypeEscapedValue, []byte{5, 6}},
	}

	for _, c := range cases {
		n = newNode(c.t, c.chunk)
		if n.t != c.t || !reflect.DeepEqual(n.chunk.Bytes(), c.chunk) ||
			n.children == nil || len(n.children) != 0 {
			t.Fail()
		}
	}
}

func TestSplitByEndBlock(t *testing.T) {
	cases := []struct {
		in  []byte
		bs1 []byte
		bs2 []byte
	}{
		{in: []byte("hello<% } %>world"), bs1: []byte("hello"), bs2: []byte("world")},
		{in: []byte(" a <% } %> b "), bs1: []byte(" a "), bs2: []byte(" b ")},
	}

	for _, c := range cases {
		if bs1, bs2 := splitByEndBlock(c.in); !reflect.DeepEqual(c.bs1, bs1) ||
			!reflect.DeepEqual(c.bs2, bs2) {
			t.Fail()
		}
	}
}

func buildTree() *node {
	root := newNode(TypeRoot, nil)
	root.insert("/tmp/gohero", "list.html", []byte(listHTML))
	return root
}

func testList(root *node, t *testing.T) {
	var (
		child   *node
		content string
	)

	child = root.children[0]
	content = strings.TrimSpace(child.chunk.String())
	if child.t != TypeDefinition ||
		content != "func UserList(userList []string) *bytes.Buffer" {
		t.Fail()
	}

	child = root.children[1]
	content = strings.TrimSpace(child.chunk.String())
	if child.t != TypeImport || content != `func Add(a, b int) int {
		return a + b
	}` {
		t.Fail()
	}

	child = root.children[2]
	content = strings.TrimSpace(child.chunk.String())
	if child.t != TypeExtend || content != "/tmp/gohero/index.html" {
		t.Fail()
	}

	child = root.children[3]
	content = strings.TrimSpace(child.chunk.String())
	if child.t != TypeBlock || content != "body" || len(child.children) != 4 {
		t.Fail()
	}

	include := child.children[2]
	if include.t != TypeInclude ||
		include.chunk.String() != path.Join(rootDir, "item.html") ||
		len(include.children) != 0 {
		t.Fail()
	}
}

func TestInsert(t *testing.T) {
	root := buildTree()
	testList(root, t)
}

func TestChildrenByType(t *testing.T) {
	root := buildTree()

	var children []*node

	children = root.childrenByType(TypeHTML)
	if len(children) > 0 {
		t.Fail()
	}

	children = root.childrenByType(TypeDefinition)
	content := strings.TrimSpace(children[0].chunk.String())
	if len(children) != 1 || children[0].t != TypeDefinition ||
		content != "func UserList(userList []string) *bytes.Buffer" {
		t.Fail()
	}
}

func TestFindBlockByName(t *testing.T) {
	root := buildTree()

	cases := []struct {
		name    string
		existed bool
	}{
		{"head", false},
		{"body", true},
	}

	for _, c := range cases {
		if block := root.findBlockByName(c.name); !(block != nil == c.existed) {
			t.Fail()
		}
	}
}

func TestParseFile(t *testing.T) {
	root := parseFile(rootDir, "list.html")
	testList(root, t)

	if len(dependencies.vertices) != 3 {
		t.Fail()
	}

	pathIndex := path.Join(rootDir, "index.html")
	pathItem := path.Join(rootDir, "item.html")
	pathList := path.Join(rootDir, "list.html")

	vertices := map[string]struct{}{
		pathIndex: struct{}{},
		pathItem:  struct{}{},
		pathList:  struct{}{},
	}

	if !reflect.DeepEqual(vertices, dependencies.vertices) {
		t.Fail()
	}

	graph := map[string]map[string]struct{}{
		pathIndex: map[string]struct{}{
			pathList: struct{}{},
		},
		pathItem: map[string]struct{}{
			pathList: struct{}{},
		},
	}

	if !reflect.DeepEqual(graph, dependencies.graph) {
		t.Fail()
	}
}

func testRebuild(root *node, t *testing.T) {
	if len(root.children) != 5 {
		t.Fail()
	}

	var (
		child   *node
		content string
	)

	child = root.children[0]
	content = strings.TrimSpace(child.chunk.String())
	if child.t != TypeImport || content != `func Add(a, b int) int {
		return a + b
	}` {
		t.Fail()
	}

	child = root.children[1]
	content = strings.TrimSpace(child.chunk.String())
	if child.t != TypeDefinition ||
		content != "func UserList(userList []string) *bytes.Buffer" {
		t.Fail()
	}

	child = root.children[2]
	content = strings.TrimSpace(child.chunk.String())
	if child.t != TypeHTML ||
		content != `<!DOCTYPE html>
<html>
  <head>
  </head>
  <body>` {
		t.Fail()
	}

	child = root.children[3]
	content = strings.TrimSpace(child.chunk.String())
	if child.t != TypeBlock || content != "body" || len(child.children) != 4 {
		t.Fail()
	}

	include := child.children[2]
	if include.t != TypeInclude ||
		include.chunk.String() != path.Join(rootDir, "item.html") ||
		len(include.children) != 5 {
		t.Fail()
	}

	child = root.children[4]
	content = strings.TrimSpace(child.chunk.String())
	if child.t != TypeHTML || content != `</body>
</html>` {
		t.Fail()
	}
}

func TestRebuild(t *testing.T) {
	paths := []string{
		"index.html",
		"item.html",
		"list.html",
	}

	for _, p := range paths {
		parsedNodes[path.Join(rootDir, p)] = parseFile(rootDir, p)
	}

	root := parsedNodes[path.Join(rootDir, "list.html")]
	root.rebuild()

	testRebuild(root, t)
}

func TestParseDir(t *testing.T) {
	dependencies.graph = make(map[string]map[string]struct{})
	dependencies.vertices = make(map[string]struct{})

	parseDir(rootDir)

	root := parsedNodes[path.Join(rootDir, "list.html")]
	testRebuild(root, t)
}
