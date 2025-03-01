package models

import "time"

// AccountSearchModel интерфейс для модели данных поиска аккаунтов
type AccountSearchModel interface {
	GetUID() string
	GetLogin() string
	GetURL() string
	GetDescription() string
	GetCreatedFrom() time.Time
	GetCreatedTo() time.Time
	GetUpdatedFrom() time.Time
	GetUpdatedTo() time.Time
	GetLimit() int
	GetOffset() int
}

// AccountSearch модель данных поиска аккаунтов
type AccountSearch struct {
	uid         string
	login       string
	url         string
	description string
	createdFrom time.Time
	createdTo   time.Time
	updatedFrom time.Time
	updatedTo   time.Time
	limit       int
	offset      int
}

// GetUID возвращает идентификатор владельца аккаунта
func (a AccountSearch) GetUID() string {
	return a.uid
}

// GetLogin возвращает логин аккаунта
func (a AccountSearch) GetLogin() string {
	return a.login
}

// GetURL возвращает URL аккаунта
func (a AccountSearch) GetURL() string {
	return a.url
}

// GetDescription возвращает описание для аккаунта
func (a AccountSearch) GetDescription() string {
	return a.description
}

// GetCreatedFrom возвращает минимальную дату создания аккаунта
func (a AccountSearch) GetCreatedFrom() time.Time {
	return a.createdFrom
}

// GetCreatedTo возвращает максимальную дату создания аккаунта
func (a AccountSearch) GetCreatedTo() time.Time {
	return a.createdTo
}

// GetUpdatedFrom возвращает минимальную дату обновления аккаунта
func (a AccountSearch) GetUpdatedFrom() time.Time {
	return a.updatedFrom
}

// GetUpdatedTo возвращает максимальную дату обновления аккаунта
func (a AccountSearch) GetUpdatedTo() time.Time {
	return a.updatedTo
}

// GetLimit возвращает лимит поиска
func (a AccountSearch) GetLimit() int {
	return a.limit
}

// GetOffset возвращает смещение поиска
func (a AccountSearch) GetOffset() int {
	return a.offset
}
