package gofpdf

import (
	"math"
)

// SplitText splits UTF-8 encoded text into several lines using the current
// font. Each line has its length limited to a maximum width given by w. This
// function can be used to determine the total height of wrapped text for
// vertical placement purposes.
func (f *Fpdf) SplitText(txt string, w float64) (lines []string) {
	cw := f.currentFont.Cw
	wmax := int(math.Ceil((w - 2*f.cMargin) * 1000 / f.fontSize))
	s := []rune(txt) // Return slice of UTF-8 runes
	nb := len(s)
	for nb > 0 && s[nb-1] == '\n' {
		nb--
	}
	s = s[0:nb]
	sep := -1
	i := 0
	j := 0
	l := 0
	enBool := false // 判断是否是连续英文，连续英文sep停留不动
	for i < nb {
		c := s[i]
		l += cw[c]

		// 非英文，只有全英文的时候sep才等于-1
		if !isEnglish(c) {
			sep = i
			enBool = false
		} else {
			if !enBool {
				sep = i
			}
			enBool = true
		}

		if c == '\n' || l > wmax {
			// 字符串长度大于行的长度了，要换行
			if sep == -1 {
				// 全是英文正常换行
				lines = append(lines, string(s[j:i]))
				j = i
			} else {
				// 包含中文，为了把英文单词放在一起
				lines = append(lines, string(s[j:sep]))
				// 换完行之后，指针移动
				j = sep
				i = sep
			}
			sep = -1
			l = 0
		} else {
			i++
		}
	}
	if i != j {
		lines = append(lines, string(s[j:i]))
	}
	return lines
}
