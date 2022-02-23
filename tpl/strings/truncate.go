// Copyright 2016 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strings

import (
	"errors"
	"github.com/spf13/cast"
	"html"
	"html/template"
	"regexp"
	"unicode"
	"unicode/utf8"
)

var (
	tagRE        = regexp.MustCompile(`^<(/)?([^ ]+?)(?:(\s*/)| .*?)?>`)
	htmlSinglets = map[string]bool{
		"br": true, "col": true, "link": true,
		"base": true, "img": true, "param": true,
		"area": true, "hr": true, "input": true,
	}
)

type htmlTag struct {
	name    string
	pos     int
	openTag bool
}

// Truncate truncates a given string to the specified length.
func (ns *Namespace) Truncate(a interface{}, coverageMap map[int][]bool, options ...interface{}) (template.HTML, error) {
	length, err := cast.ToIntE(a)
	coverageMap[0][0] = true
	if err != nil {
		coverageMap[0][1] = true
		return "", err
	}
	coverageMap[0][2] = true
	var textParam interface{}
	var ellipsis string
	coverageMap[1][0] = true
	coverageMap[2][0] = true
	coverageMap[3][0] = true
	switch len(options) {
	case 0:
		coverageMap[1][1] = true
		coverageMap[2][2] = true
		coverageMap[3][2] = true
		return "", errors.New("truncate requires a length and a string")
	case 1:
		coverageMap[2][1] = true
		coverageMap[1][2] = true
		coverageMap[3][2] = true
		textParam = options[0]
		ellipsis = " …"
	case 2:
		coverageMap[3][1] = true
		coverageMap[1][2] = true
		coverageMap[2][2] = true
		textParam = options[1]
		ellipsis, err = cast.ToStringE(options[0])
		coverageMap[4][0] = true
		if err != nil {
			coverageMap[4][1] = true
			return "", errors.New("ellipsis must be a string")
		}
		coverageMap[4][2] = true
		coverageMap[5][0] = true
		if _, ok := options[0].(template.HTML); !ok {
			coverageMap[5][1] = true
			ellipsis = html.EscapeString(ellipsis)
		}
		coverageMap[5][2] = true
	default:
		coverageMap[1][2] = true
		coverageMap[2][2] = true
		coverageMap[3][2] = true
		return "", errors.New("too many arguments passed to truncate")
	}
	coverageMap[6][0] = true
	if err != nil {
		coverageMap[6][1] = true
		return "", errors.New("text to truncate must be a string")
	}
	coverageMap[6][2] = true
	text, err := cast.ToStringE(textParam)
	coverageMap[7][0] = true
	if err != nil {
		coverageMap[7][1] = true
		return "", errors.New("text must be a string")
	}
	coverageMap[7][2] = true

	_, isHTML := textParam.(template.HTML)
	coverageMap[8][0] = true
	if utf8.RuneCountInString(text) <= length {
		coverageMap[8][1] = true
		coverageMap[9][0] = true
		if isHTML {
			coverageMap[9][1] = true
			return template.HTML(text), nil
		}
		coverageMap[9][2] = true
		return template.HTML(html.EscapeString(text)), nil
	}
	coverageMap[8][2] = true

	tags := []htmlTag{}
	var lastWordIndex, lastNonSpace, currentLen, endTextPos, nextTag int
	coverageMap[10][0] = true
	coverageMap[10][2] = true
	for i, r := range text {
		coverageMap[10][1] = true
		coverageMap[11][0] = true
		if i < nextTag {
			coverageMap[11][1] = true
			continue
		}
		coverageMap[11][2] = true
		coverageMap[12][0] = true
		if isHTML {
			coverageMap[12][1] = true
			// Make sure we keep tag of HTML tags
			slice := text[i:]
			m := tagRE.FindStringSubmatchIndex(slice)
			coverageMap[13][0] = true
			coverageMap[14][0] = true
			if len(m) > 0 && m[0] == 0 {
				coverageMap[13][1] = true
				coverageMap[14][1] = true
				nextTag = i + m[1]
				tagname := slice[m[4]:m[5]]
				lastWordIndex = lastNonSpace
				_, singlet := htmlSinglets[tagname]
				coverageMap[15][0] = true
				coverageMap[16][0] = true
				if !singlet && m[6] == -1 {
					coverageMap[15][1] = true
					coverageMap[16][1] = true
					tags = append(tags, htmlTag{name: tagname, pos: i, openTag: m[2] == -1})
				}
				if singlet && m[6] == -1 {
					coverageMap[15][2] = true
				}
				if !singlet && m[6] != -1 {
					coverageMap[16][2] = true
				}

				continue
			}
			coverageMap[13][2] = true
			coverageMap[14][2] = true
		}
		coverageMap[12][2] = true

		currentLen++
		coverageMap[17][0] = true
		coverageMap[18][0] = true
		if unicode.IsSpace(r) {
			coverageMap[17][1] = true
			coverageMap[18][2] = true
			lastWordIndex = lastNonSpace
		} else if unicode.In(r, unicode.Han, unicode.Hangul, unicode.Hiragana, unicode.Katakana) {
			coverageMap[18][1] = true
			coverageMap[17][2] = true
			lastWordIndex = i
		} else {
			coverageMap[17][2] = true
			coverageMap[18][2] = true
			lastNonSpace = i + utf8.RuneLen(r)
		}
		coverageMap[19][0] = true
		if currentLen > length {
			coverageMap[19][1] = true
			coverageMap[20][0] = true
			if lastWordIndex == 0 {
				coverageMap[20][1] = true
				endTextPos = i
			} else {
				coverageMap[20][2] = true
				endTextPos = lastWordIndex
			}
			out := text[0:endTextPos]
			coverageMap[21][0] = true
			if isHTML {
				coverageMap[21][1] = true
				out += ellipsis
				// Close out any open HTML tags
				var currentTag *htmlTag
				coverageMap[22][0] = true
				for i := len(tags) - 1; i >= 0; i-- {
					tag := tags[i]
					coverageMap[22][1] = true
					coverageMap[23][0] = true
					coverageMap[24][0] = true
					if tag.pos >= endTextPos || currentTag != nil {
						if tag.pos >= endTextPos && currentTag == nil {
							coverageMap[23][1] = true
						} else {
							coverageMap[24][1] = true
						}
						coverageMap[25][0] = true
						coverageMap[26][0] = true
						if currentTag != nil && currentTag.name == tag.name {
							coverageMap[25][1] = true
							coverageMap[26][1] = true
							currentTag = nil
						} else if currentTag != nil && currentTag.name != tag.name {
							coverageMap[25][1] = true
							coverageMap[26][2] = true
						} else if currentTag != nil {
							coverageMap[25][2] = true
						} else {
							coverageMap[25][2] = true
							coverageMap[26][2] = true
						}
						continue
					}
					coverageMap[23][2] = true
					coverageMap[24][2] = true

					coverageMap[27][0] = true
					if tag.openTag {
						coverageMap[27][1] = true
						out += ("</" + tag.name + ">")
					} else {
						coverageMap[27][2] = true
						currentTag = &tag
					}
				}
				coverageMap[22][2] = true

				return template.HTML(out), nil
			}
			return template.HTML(html.EscapeString(out) + ellipsis), nil
		}
		coverageMap[19][2] = true
	}
	coverageMap[28][0] = true
	if isHTML {
		coverageMap[28][1] = true
		return template.HTML(text), nil
	}
	coverageMap[28][2] = true
	return template.HTML(html.EscapeString(text)), nil
}
