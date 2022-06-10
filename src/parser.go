package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func ParseView(result_html io.Reader, id string, l string) (res Result, err error) {
	var r Result
	r.Id, _ = strconv.Atoi(id)
	doc, err := html.Parse(result_html)
	dfsv(doc, &r, l)
	return r, err
}

func dfsv(n *html.Node, in *Result, l string) *html.Node {
	// Get the Hangul
	if CheckClass(n, "word_head") {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c != nil {
				if c.Type == html.TextNode {
					in.Hangul = c.Data
				} else if c.Type == html.ElementNode {
					in.HomonymNumber, _ = strconv.Atoi(GetTextAll(c))
					break
				}
			}
		}
		// Get the Hanja
	} else if CheckClass(n, "chi_info ml5") {
		hanja := MatchBetween(GetTextAll(n), "(", " )")
		in.Hanja = hanja

	} else if CheckClass(n, "search_sub") {
		// Get the pronounciation and audio file
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode && strings.TrimSpace(c.Data) != "" && c.Data[0:1] == "[" {
				in.Pronounciation = strings.TrimSpace(c.Data) + "]"
			} else if c.Data == "a" && c.Attr[1].Val == "sound" {
				for _, a := range c.Attr {
					if a.Key == "href" {
						in.Audio = MatchBetween(a.Val, "'", "');")
						break
					}
				}
				break
			}
		}

	} else if CheckClass(n, "word_att_type1") {
		// Get the Korean word type
		match := MatchBetween(GetTextAll(n.FirstChild), "「", "」")
		in.TypeKr = strings.ToValidUTF8(match, "")

	} else if CheckClass(n, fmt.Sprintf("manyLang%s ml5", l)) {
		// Get the English word type
		in.TypeEng = strings.TrimSpace(GetTextSingle(n.FirstChild))

	} else if CheckClass(n, "ri-star-s-fill") {
		// Get the level of the word
		in.Level++

	} else if CheckClass(n, "printArea") {
		// Get the Inflections
		in.Inflections = (GetTextSingle(n))
		inf := InitInflectionLinks()
		ind := 0
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "a" {
				in.InflectionLinks = append(in.InflectionLinks, inf)
				re := regexp.MustCompile("[0-9]+")
				id := c.Attr[1].Val
				id = re.FindAllString(id, -1)[0]
				in.InflectionLinks[ind].Id, _ = strconv.Atoi(id)
				in.InflectionLinks[ind].Hangul = GetTextSingle(c) + "<sup>" + GetTextSingle(c.NextSibling) + "</sup>"
				ind++
			}
		}

	} else if CheckClass(n, fmt.Sprintf("multiTrans manyLang%s mb10 printArea", l)) {
		// Get the english translation
		s := InitSense()
		in.Senses = append(in.Senses, s)
		l := len(in.Senses)
		in.Senses[l-1].Translation = cleanStringSpecial([]byte(GetTextAll(n)))

	} else if CheckClass(n, "senseDef ml20 printArea") {
		// Get the korean definition
		l := len(in.Senses)
		in.Senses[l-1].KrDefinition = GetTextSingle(n.LastChild)

	} else if CheckClass(n, fmt.Sprintf("multiSenseDef manyLang%s ml20 printArea", l)) {
		// Get the english definition
		in.Senses[len(in.Senses)-1].Definition = GetTextAll(n)

	} else if CheckClass(n, "dot") || CheckClass(n, "dot printArea") {
		var ex example
		// Get the examples
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "li" {
				in.Senses[len(in.Senses)-1].Examples = append(in.Senses[len(in.Senses)-1].Examples, ex)
				v := GetContent(c, "b")
				t := GetTextAll(c)
				in.Senses[len(in.Senses)-1].Examples[len(in.Senses[len(in.Senses)-1].Examples)-1].Value = v
				in.Senses[len(in.Senses)-1].Examples[len(in.Senses[len(in.Senses)-1].Examples)-1].Type = t
			} else if c.Parent.Data != "ul" && c.Parent.Data != "li" {
				break
			}
		}

	} else if CheckClass(n, "heading_wrap printArea") {
		getRef(n, in, l)
	}

	// Traverse the tree of nodes vi depth-first search
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// Skip all commment nodes or nodes whose type is "script"
		if c.Type == html.CommentNode || c.Data == "script" {
			continue
		} else {
			dfsv(c, in, l)
		}
	}
	return n
}
