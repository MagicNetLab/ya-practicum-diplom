package models

import (
	"strconv"
)

// AccountSearchModel интерфейс для модели данных поиска аккаунтов
type AccountSearchModel interface {
	GetSubQuery() (string, []any)
}

// AccountSearch модель данных поиска аккаунтов
type AccountSearch struct {
	UID         string
	Login       string
	URL         string
	Description string
	Limit       int
	Offset      int
}

type searchValue struct {
	item  string
	value string
}

func (a AccountSearch) GetSubQuery() (string, []any) {
	str := ""
	values := make([]any, 0)
	items := make([]searchValue, 0)
	if a.UID != "" {
		items = append(items, searchValue{"uid", a.UID})
	}

	if a.Login != "" {
		items = append(items, searchValue{"login", a.Login})
	}

	if a.URL != "" {
		items = append(items, searchValue{"url", a.URL})
	}

	if a.Description != "" {
		items = append(items, searchValue{"description", a.Description})
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

	if a.Limit > 0 {
		str += " LIMIT $" + strconv.Itoa(len(values)+1)
		values = append(values, a.Limit)
	}

	if a.Offset > 0 {
		str += " OFFSET $" + strconv.Itoa(len(values)+1)
		values = append(values, a.Offset)
	}

	return str, values
}
