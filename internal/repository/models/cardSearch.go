package models

import "strconv"

type CardSearchModel interface {
	GetSubQuery() (string, []any)
}

type CardSearch struct {
	UID    string
	Name   string
	Limit  int
	Offset int
}

type cardSearchValue struct {
	item  string
	value string
}

// GetSubQuery возвращает подзапрос для поиска аккаунтов
// todo вынести общую логику с доугими моделями поиска куда-нибудь отдельно
func (cs *CardSearch) GetSubQuery() (string, []any) {
	str := ""
	values := make([]any, 0)
	items := make([]cardSearchValue, 0)
	if cs.UID != "" {
		items = append(items, cardSearchValue{item: "uid", value: cs.UID})
	}
	if cs.Name != "" {
		items = append(items, cardSearchValue{item: "name", value: cs.Name})
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

	if cs.Limit > 0 {
		str += " LIMIT $" + strconv.Itoa(len(values)+1)
		values = append(values, cs.Limit)
	}

	if cs.Offset > 0 {
		str += " OFFSET $" + strconv.Itoa(len(values)+1)
		values = append(values, cs.Offset)
	}

	return str, values
}
