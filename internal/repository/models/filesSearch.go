package models

import "strconv"

type FilesSearchModel interface {
	GetUID() string
	GetSubQuery() (string, []any)
}

type FilesSearch struct {
	UID    string
	Name   string
	Meta   string
	Limit  int32
	Offset int32
}

func (s *FilesSearch) GetUID() string {
	return s.UID
}

func (s *FilesSearch) GetSubQuery() (string, []any) {
	str := ""
	values := make([]any, 0)
	items := make([]cardSearchValue, 0)
	if s.UID != "" {
		items = append(items, cardSearchValue{item: "uid", value: s.UID})
	}
	if s.Name != "" {
		items = append(items, cardSearchValue{item: "name", value: s.Name})
	}
	if s.Meta != "" {
		items = append(items, cardSearchValue{item: "meta", value: s.Meta})
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

	if s.Limit > 0 {
		str += " LIMIT $" + strconv.Itoa(len(values)+1)
		values = append(values, s.Limit)
	}

	if s.Offset > 0 {
		str += " OFFSET $" + strconv.Itoa(len(values)+1)
		values = append(values, s.Offset)
	}

	return str, values
}
