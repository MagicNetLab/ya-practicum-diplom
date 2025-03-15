package client

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	accpb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	authpb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
	cardpb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
)

type AccountData struct {
	ID       string
	Login    string
	Password string
	URL      string
	Meta     string
}

type CardData struct {
	ID     string
	Name   string
	Number string
	Mask   string
	Month  int32
	Year   int32
	CVC    string
	PIN    string
	Meta   string
}

type ShortCardData struct {
	ID     string
	Name   string
	Number string
	Meta   string
}

type CardSearchData struct {
	Name   string
	Number string
	Meta   string
	Year   string
}

type AppClient interface {
	Auth(ctx context.Context, username string, password string) (string, error)
	Register(ctx context.Context, username string, password string) (string, error)
	ListAccounts(ctx context.Context) ([]AccountData, error)
	AddAccount(ctx context.Context, username string, password string, url string, meta string) error
	RemoveAccount(ctx context.Context, username string) error
	GetAccount(ctx context.Context, id string) (AccountData, error)
	SearchAccount(ctx context.Context, login, url, meta string) ([]AccountData, error)
	ListCards(ctx context.Context) ([]ShortCardData, error)
	AddCard(ctx context.Context, data CardData) error
	RemoveCard(ctx context.Context, id string) error
	CardDetail(ctx context.Context, id string) (CardData, error)
	CardSearch(ctx context.Context, data CardSearchData) ([]ShortCardData, error)
}

func NewAppClient(cnf config.AppConfigurator) (AppClient, error) {
	connAddress := cnf.GetServerConf().GetHost() + ":" + cnf.GetServerConf().GetPort()
	connect, err := grpc.NewClient(connAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logger.Error("Failed to create grpc client", logger.StrArg("error", err.Error()), logger.StrArg("address", connAddress))
		return nil, err
	}

	return &AppClientImpl{
		cnf:        cnf,
		authClient: authpb.NewAuthClient(connect),
		accounts:   accpb.NewAccountsClient(connect),
		cards:      cardpb.NewCardClient(connect),
	}, nil
}

type AppClientImpl struct {
	cnf        config.AppConfigurator
	authClient authpb.AuthClient
	accounts   accpb.AccountsClient
	cards      cardpb.CardClient
}

// Auth - метод аутентификации пользователя
func (c *AppClientImpl) Auth(ctx context.Context, username string, password string) (string, error) {
	request := &authpb.AuthRequest{
		Login:  username,
		Secret: password,
	}

	resp, err := c.authClient.Auth(ctx, request)
	if err != nil {
		return "", err
	}

	return resp.GetToken(), nil
}

// Register - регистрация пользователя
func (c *AppClientImpl) Register(ctx context.Context, username string, password string) (string, error) {
	request := &authpb.RegRequest{
		Login:  username,
		Secret: password,
	}

	resp, err := c.authClient.Register(ctx, request)
	if err != nil {
		return "", err
	}

	return resp.GetToken(), nil
}

// ListAccounts - метод получения списка всех аккаунтов пользователя
func (c *AppClientImpl) ListAccounts(ctx context.Context) ([]AccountData, error) {
	req := &accpb.SearchAccountRequest{}

	resp, err := c.accounts.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	res := make([]AccountData, 0)

	for _, row := range resp.GetAcc() {
		r := AccountData{
			ID:    row.GetId(),
			Login: row.GetLogin(),
			URL:   row.GetUrl(),
			Meta:  row.GetDescription(),
		}
		res = append(res, r)
	}

	return res, nil
}

// AddAccount - метод добавления нового аккаунта
func (c *AppClientImpl) AddAccount(ctx context.Context, username string, password string, url string, meta string) error {
	req := &accpb.CreateAccountRequest{
		Login:       username,
		Password:    password,
		Url:         url,
		Description: meta,
	}

	_, err := c.accounts.Create(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

// RemoveAccount - метод удаления аккаунта
func (c *AppClientImpl) RemoveAccount(ctx context.Context, id string) error {
	req := &accpb.RemoveAccountRequest{Id: id}

	_, err := c.accounts.Remove(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// GetAccount - метод получения данных по конкретному аккаунту
func (c *AppClientImpl) GetAccount(ctx context.Context, id string) (AccountData, error) {
	req := &accpb.GetAccountRequest{Id: id}
	resp, err := c.accounts.Get(ctx, req)
	if err != nil {
		return AccountData{}, err
	}

	acc := resp.GetAccount()
	model := AccountData{
		ID:       acc.GetId(),
		Login:    acc.GetLogin(),
		Password: acc.GetPassword(),
		URL:      acc.GetUrl(),
		Meta:     acc.GetDescription(),
	}

	return model, nil
}

// SearchAccount - метод поиска аккаунтов по фильтрам
func (c *AppClientImpl) SearchAccount(ctx context.Context, login, url, meta string) ([]AccountData, error) {
	req := &accpb.SearchAccountRequest{
		Login:       login,
		Url:         url,
		Description: meta,
	}

	resp, err := c.accounts.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	res := make([]AccountData, 0)
	for _, row := range resp.GetAcc() {
		r := AccountData{
			ID:    row.GetId(),
			Login: row.GetLogin(),
			URL:   row.GetUrl(),
			Meta:  row.GetDescription(),
		}
		res = append(res, r)
	}
	return res, nil
}

// ListCards - метод получения списка всех карт пользователя
func (c *AppClientImpl) ListCards(ctx context.Context) ([]ShortCardData, error) {
	req := &cardpb.ListCardRequest{}

	resp, err := c.cards.List(ctx, req)
	if err != nil {
		return nil, err
	}

	res := make([]ShortCardData, 0)
	for _, row := range resp.GetCards() {
		r := ShortCardData{
			ID:     row.GetID(),
			Name:   row.GetName(),
			Number: row.GetNumber(),
			Meta:   row.GetMeta(),
		}

		res = append(res, r)
	}

	return res, nil
}

// AddCard - метод добавления новой карты
func (c *AppClientImpl) AddCard(ctx context.Context, data CardData) error {
	_, err := c.cards.Create(ctx, &cardpb.CreateCardRequest{
		Name:   data.Name,
		Number: data.Number,
		Meta:   data.Meta,
		Month:  data.Month,
		Year:   data.Year,
		CVC:    data.CVC,
		PIN:    data.PIN,
	})

	if err != nil {
		return err
	}

	return nil
}

// RemoveCard - метод удаления карты
func (c *AppClientImpl) RemoveCard(ctx context.Context, id string) error {
	_, err := c.cards.Delete(ctx, &cardpb.DeleteCardRequest{ID: id})
	if err != nil {
		return err
	}
	return nil
}

// CardDetail - метод получения информации по конкретной карте
func (c *AppClientImpl) CardDetail(ctx context.Context, id string) (CardData, error) {
	res, err := c.cards.Get(ctx, &cardpb.GetCardRequest{ID: id})
	if err != nil {
		return CardData{}, err
	}
	card := res.GetCard()
	return CardData{
		ID:     id,
		Name:   card.GetName(),
		Number: card.GetNumber(),
		Meta:   card.GetMeta(),
		Month:  card.GetMonth(),
		Year:   card.GetYear(),
		CVC:    card.GetCVC(),
		PIN:    card.GetPIN(),
	}, nil
}

// CardSearch - метод поиска карт по фильтрам
func (c *AppClientImpl) CardSearch(ctx context.Context, data CardSearchData) ([]ShortCardData, error) {
	req := &cardpb.SearchCardRequest{}
	req.Name = data.Name
	req.Number = data.Number
	req.Meta = data.Meta
	if data.Year != "" {
		y, err := strconv.Atoi(data.Year)
		if err == nil {
			req.Year = int32(y)
		}
	}
	res, err := c.cards.Search(ctx, req)

	if err != nil {
		return nil, err
	}

	resList := make([]ShortCardData, 0)
	for _, row := range res.GetCards() {
		r := ShortCardData{
			ID:     row.GetID(),
			Name:   row.GetName(),
			Number: row.GetNumber(),
			Meta:   row.GetMeta(),
		}
		resList = append(resList, r)
	}

	return resList, nil
}
