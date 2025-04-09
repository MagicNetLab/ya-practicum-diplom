package models

import "strconv"

// NoteSearchModel - интерфейс модели поиска заметок
type NoteSearchModel interface {
	GetSubQuery() (string, []any)
}

// NoteSearch - модель поиска заметок
type NoteSearch struct {
	UID    string
	Search string
	Limit  int32
	Offset int32
}

type searchNoteValue struct {
	item  string
	value string
}

// GetSubQuery - возвращает подзапрос для поиска заметок
func (n *NoteSearch) GetSubQuery() (string, []any) {
	str := ""
	values := make([]any, 0)

	items := make([]searchNoteValue, 0)
	if n.UID != "" {
		items = append(items, searchNoteValue{"uid", n.UID})
	}

	if n.Search != "" {
		items = append(items, searchNoteValue{"title", n.Search})
	}

	if len(items) > 0 {
		str += " WHERE "
		i := 1

		for _, v := range items {
			if i > 1 {
				str += " AND "
			}
			if v.item == "uid" {
				str += v.item + " = $" + strconv.Itoa(i)
				values = append(values, v.value)
			} else {
				str += v.item + " ilike $" + strconv.Itoa(i)
				values = append(values, "%"+v.value+"%")
			}
			i++
		}
	}

	if n.Limit > 0 {
		str += " LIMIT $" + strconv.Itoa(len(values)+1)
		values = append(values, n.Limit)
	}

	if n.Offset > 0 {
		str += " OFFSET $" + strconv.Itoa(len(values)+1)
		values = append(values, n.Offset)
	}

	return str, values
}
