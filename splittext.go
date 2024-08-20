package gofpdf

import (
	"math"
	"strings"
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
	i := 0
	j := 0
	sep := -1
	l := 0
	enBool := true // 判断是否是连续英文，连续英文sep停留不动
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
			// l > wmax 字符串长度大于行的长度了，要换行
			// j是每一行start索引，sep是end索引，j<= lines < sep
			if sep == -1 {
				// sep = -1说明这一行写的全是英文
				sep = i
			}

			if j != sep {
				_addLines(&lines, string(s[j:sep]))
			} else {
				// j==sep时，说明是\n导致的，跳过该字符
				sep++
			}

			// 每一行写完要重置变量
			i, j = sep, sep
			sep = -1
			l = 0
			enBool = true // 判断是否是连续英文，连续英文sep停留不动
		} else {
			i++
		}
	}
	if i != j {
		_addLines(&lines, string(s[j:i]))
	}
	return lines
}

func _addLines(lines *[]string, line string) {
	// 如果整行等于一个空格
	if line != " " {
		if len(*lines) > 0 {
			// 首行缩进的空格保留，其它行的首尾空格去掉
			line = strings.Trim(line, " ")
		}
		//fmt.Println(line)
		*lines = append(*lines, line)
	}
}
