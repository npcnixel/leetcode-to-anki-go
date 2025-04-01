package tests

import (
	"regexp"
	"strings"
	"testing"

	leetcode "github.com/npcnixel/leetcode-to-anki-go/pkg/leetcode_to_anki"
)

func TestContainerWithMostWater(t *testing.T) {
	// The exact HTML from your example
	html := `<div class="view-lines monaco-mouse-cursor-text" role=presentation aria-hidden=true data-mprt=7 style='position:absolute;font-family:Menlo,Monaco,"Courier New",monospace;font-weight:normal;font-size:13px;font-feature-settings:"liga"0,"calt"0;line-height:20px;letter-spacing:0px;width:732px;height:626px'><div style=top:8px;height:20px class=view-line><span><span class=mtk4>class</span><span class=mtk1>&nbsp;</span><span class=mtk10>Solution</span><span class=mtk1>:</span></span></div><div style=top:28px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;</span><span class=mtk4>def</span><span class=mtk1>&nbsp;</span><span class=mtk11>maxArea</span><span class=mtk1>(</span><span class=mtk14>self</span><span class=mtk1>,&nbsp;</span><span class=mtk14>height</span><span class=mtk1>:&nbsp;List[</span><span class=mtk10>int</span><span class=mtk1>])&nbsp;-&gt;&nbsp;</span><span class=mtk10>int</span><span class=mtk1>:</span></span></div><div style=top:48px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;l&nbsp;=&nbsp;</span><span class=mtk7>0</span></span></div><div style=top:68px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;r&nbsp;=&nbsp;</span><span class=mtk11>len</span><span class=mtk1>(height)-</span><span class=mtk7>1</span></span></div><div style=top:88px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;maxArea&nbsp;=&nbsp;</span><span class=mtk7>0</span></span></div><div style=top:108px;height:20px class=view-line><span><span></span></span></div><div style=top:128px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span><span class=mtk13>while</span><span class=mtk1>&nbsp;l&lt;r:&nbsp;</span></span></div><div style=top:148px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;minH&nbsp;=&nbsp;</span><span class=mtk11>min</span><span class=mtk1>(height[l],&nbsp;height[r])</span></span></div><div style=top:168px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;currentArea&nbsp;=&nbsp;minH&nbsp;*&nbsp;(r-l)</span></span></div><div style=top:188px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;maxArea&nbsp;=&nbsp;</span><span class=mtk11>max</span><span class=mtk1>(maxArea,&nbsp;currentArea)</span></span></div><div style=top:208px;height:20px class=view-line><span><span></span></span></div><div style=top:228px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span><span class=mtk13>if</span><span class=mtk1>&nbsp;height[l]&nbsp;&lt;&nbsp;height[r]:</span></span></div><div style=top:248px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;l&nbsp;+=</span><span class=mtk7>1</span></span></div><div style=top:268px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span><span class=mtk13>else</span><span class=mtk1>:</span></span></div><div style=top:288px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;r&nbsp;-=</span><span class=mtk7>1</span></span></div><div style=top:308px;height:20px class=view-line><span><span class=mtk1>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span><span class=mtk13>return</span><span class=mtk1>&nbsp;maxArea</span></span></div></div>`
	expected := `class Solution:
    def maxArea(self, height: List[int]) -> int:
        l = 0
        r = len(height)-1
        maxArea = 0

        while l<r:
            minH = min(height[l], height[r])
            currentArea = minH * (r-l)
            maxArea = max(maxArea, currentArea)

            if height[l] < height[r]:
                l +=1
            else:
                r -=1
        return maxArea
`

	code, err := leetcode.ExtractCodeFromHTML(html)
	if err != nil {
		t.Fatalf("Failed to extract code: %v", err)
	}

	// Normalize whitespace in both strings
	normalizeWhitespace := func(s string) string {
		// Replace all whitespace sequences with a single space
		s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
		// Normalize indentation
		lines := strings.Split(s, "\n")
		for i, line := range lines {
			lines[i] = strings.TrimSpace(line)
		}
		return strings.Join(lines, "\n")
	}

	normalizedExpected := normalizeWhitespace(expected)
	normalizedActual := normalizeWhitespace(code)

	// Compare line by line for better error reporting
	expectedLines := strings.Split(strings.TrimSpace(normalizedExpected), "\n")
	actualLines := strings.Split(strings.TrimSpace(normalizedActual), "\n")

	if len(expectedLines) != len(actualLines) {
		t.Errorf("Line count mismatch.\nExpected %d lines:\n%s\n\nGot %d lines:\n%s",
			len(expectedLines), expected, len(actualLines), code)
		return
	}

	for i, expectedLine := range expectedLines {
		if i >= len(actualLines) {
			t.Errorf("Missing line %d: %q", i+1, expectedLine)
			continue
		}
		if expectedLine != actualLines[i] {
			t.Errorf("Line %d mismatch:\nExpected: %q\nGot: %q", i+1, expectedLine, actualLines[i])
		}
	}
}
