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
	UID    string
	Search string
	Limit  int
	Offset int
}

type searchValue struct {
	item  string
	value string
}

// GetSubQuery возвращает подзапрос и аргументы запроса
func (a AccountSearch) GetSubQuery() (string, []any) {
	str := ""
	values := make([]any, 0)
	items := make([]searchValue, 0)

	if a.Search != "" {
		items = append(items, searchValue{"url", a.Search})

	}

	if len(items) > 0 {
		str += " WHERE "
		i := 1

		for _, v := range items {
			if i > 1 {
				str += " AND "
			}

			str += v.item + " ilike $" + strconv.Itoa(i)
			values = append(values, "%"+v.value+"%")

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
