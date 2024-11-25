// Package kongpleter generates a yaml description of the command line as described by Kong.
// This yaml can then be used by github.com/miekg/gompletely to generate specific completions.
package kongpleter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
)

func walk(node *kong.Node, visit func(*kong.Node)) {
	if node == nil {
		return
	}

	queue := []*kong.Node{node}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		visit(current)

		save := current.Name
		for _, a := range current.Aliases {
			current.Name = a
			visit(current)
			current.Name = save
		}

		for _, child := range current.Children {
			queue = append(queue, child)
		}
	}
	return
}

// Walk walks kong and prints out the completion yaml according to the kong Parser.
// Hidden flags are skipped, and the extra tag "completion" is used to generate completions.
//
// completion may hold:
//   - <file> (or any other bash compgen *action*)
//   - a string not starting with <, which is interpreted as a shell command.
//
// Basic usage:
//
//	func (c Completion) BeforeReset(ctx *kong.Context, p *kong.Kong) error {
//		out, _ := kongpleter.Walk(p)
//		println(string(out))
//	  ...
//
// Where BeforeReset is used to have a flag do something special, in this case output the
// completion yaml.
func Walk(kong *kong.Kong) []byte {
	c := &comp{b: &bytes.Buffer{}}
	c.completeFunc(kong.Model.Node)

	for _, n := range kong.Model.Children {
		walk(n, c.completeFunc)
	}
	return c.b.Bytes()
}

type comp struct {
	b *bytes.Buffer
}

func (c *comp) completeFunc(n *kong.Node) {
	name := Path(n)
	newline := "\n"
	if name == n.Name {
		newline = ""
	}
	fmt.Fprintf(c.b, "%s%s:\n", newline, name)
	c.Subcommands(n)
	c.FlagSimple(n)
	c.Positional(n)

	c.FlagComplexDetail(n)
}

func (c *comp) Subcommands(n *kong.Node) {
	for _, n1 := range n.Children {
		if n1.Type == kong.CommandNode {
			fmt.Fprintf(c.b, "- S,%s,\n", n1.Name)
			for _, a := range n1.Aliases {
				fmt.Fprintf(c.b, "- S,%s,\n", a)
			}
		}
	}
}

func (c *comp) Positional(n *kong.Node) {
	for i, p := range n.Positional {
		fmt.Fprintf(c.b, "- %d,%s\n", i+1, p.Name)
	}
}

func (c *comp) FlagSimple(n *kong.Node) {
	for _, f := range n.Flags {
		if f.Hidden {
			continue
		}
		fmt.Fprintf(c.b, "- --%s[%s]\n", f.Name, f.Help)
		if f.Tag.Negatable != "" {
			fmt.Fprintf(c.b, "- --no-%s[%s]\n", f.Name, f.Help)
			continue
		}
		if f.Short > 0 {
			fmt.Fprintf(c.b, "- -%s[%s]\n", string(f.Short), f.Help)
			continue
		}
	}
}

func (c *comp) FlagComplexDetail(n *kong.Node) {
	name := Path(n)
	for _, f := range n.Flags {
		if f.Hidden {
			continue
		}
		comp := f.Tag.Get("completion")
		if comp != "" {
			fmt.Fprintf(c.b, "\n%s*--%s:\n", name, f.Name)
			if strings.HasPrefix(comp, "<") {
				fmt.Fprintf(c.b, "- %s[%s]\n", comp, f.Help)
			} else {
				fmt.Fprintf(c.b, "- $(%s)[%s]\n", comp, f.Help)
			}
		}
		if f.Enum != "" {
			fmt.Fprintf(c.b, "\n%s*--%s:\n", name, f.Name)
			for _, e := range f.EnumSlice() {
				fmt.Fprintf(c.b, "- %s[%s]\n", e, f.Help)
			}
		}
		if len(f.Envs) > 0 {
			fmt.Fprintf(c.b, "\n%s*--%s:\n", name, f.Name)
			envs := make([]string, len(f.Envs))
			for i := range envs {
				envs[i] = "'$" + envs[i] + "'"
			}
			format := strings.Repeat("%s\n", len(envs))
			list := strings.Join(envs, " ")
			fmt.Fprintf(c.b, "- $(printf %q %s)[%s]\n", format, list, f.Help)
		}
	}
}

func Path(n *kong.Node) string {
	name := n.Name
	for n.Parent != nil {
		n = n.Parent
		name = n.Name + " " + name
	}
	return name
}
