package textfmt

const (
	stateIndent = 0
	stateJump   = 1
	stateField  = 2
)

func Parse(fmt string) (fields []string, indents []string) {
	nf := 0
	ni := 0
	st := stateIndent

	states := [stateField + 1]func(c int32){}

	indents = append(indents, "")

	states[stateIndent] = func(c int32) {
		if rune(c) == '%' {
			st = stateJump
			return
		}
		indents[ni] += string([]rune{rune(c)})
	}

	states[stateJump] = func(c int32) {
		if rune(c) == '{' {
			st = stateField
			ni++
			fields = append(fields, "")
			return
		}
		indents[ni] += string([]rune{rune(c)})
	}

	states[stateField] = func(c int32) {
		if rune(c) == '}' {
			st = stateIndent
			nf++
			indents = append(indents, "")
			return
		}
		fields[nf] += string([]rune{rune(c)})
	}

	for _, c := range fmt {
		states[st](c)
	}

	return
}
