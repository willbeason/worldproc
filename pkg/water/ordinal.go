package water

type Ordinal struct {
	Index int
	Height float64
}

type OrdinalList []Ordinal

func (l *OrdinalList) insert(ordinal Ordinal, idx int) {
	al := *l
	al = append(al, Ordinal{})
	copy(al[idx+1:], al[idx:])
	al[idx] = ordinal
	*l = al
}

func (l *OrdinalList) Insert(ordinal Ordinal) {
	al := *l
	if len(al) == 0 {
		*l = []Ordinal{ordinal}
		return
	}

	left := 0
	right := len(al)

	if ordinal.Height < al[0].Height {
		l.insert(ordinal, 0)
		return
	} else if ordinal.Height > al[right-1].Height {
		*l = append(al, ordinal)
		return
	}

	for left + 1 < right {
		mid := (left + right) / 2
		if al[mid].Height < ordinal.Height {
			left = mid
		} else {
			right = mid
		}
	}
	l.insert(ordinal, right)
}

func (l *OrdinalList) Pop() *Ordinal {
	if len(*l) == 0 {
		return nil
	}
	result := (*l)[0]
	*l = (*l)[1:]
	return &result
}
