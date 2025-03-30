package client

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	accpb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	authpb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
	cardpb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
	filepb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
	notepb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	Name string
}

type NoteData struct {
	ID      string
	Title   string
	Content string
	Meta    string
}

type FileData struct {
	ID      string
	Name    string
	Path    string
	Content []byte
	Size    int64
	Meta    string
}

type NoteSearchData struct {
	Search string
}

type AppClient interface {
	Auth(ctx context.Context, username string, password string) (string, error)
	Register(ctx context.Context, username string, password string) (string, error)
	ListAccounts(ctx context.Context) ([]AccountData, error)
	AddAccount(ctx context.Context, username string, password string, url string, meta string) error
	RemoveAccount(ctx context.Context, username string) error
	GetAccount(ctx context.Context, id string) (AccountData, error)
	SearchAccount(ctx context.Context, search string) ([]AccountData, error)
	ListCards(ctx context.Context) ([]ShortCardData, error)
	AddCard(ctx context.Context, data CardData) error
	RemoveCard(ctx context.Context, id string) error
	CardDetail(ctx context.Context, id string) (CardData, error)
	CardSearch(ctx context.Context, data CardSearchData) ([]ShortCardData, error)
	NotesList(ctx context.Context) ([]NoteData, error)
	NoteCreate(ctx context.Context, data NoteData) error
	NoteDelete(ctx context.Context, id string) error
	NoteSearch(ctx context.Context, data NoteSearchData) ([]NoteData, error)
	NoteDetail(ctx context.Context, id string) (NoteData, error)
	FileList(ctx context.Context) ([]FileData, error)
	FileAdd(ctx context.Context, data FileData) error
	FileRemove(ctx context.Context, id string) error
	FileDownload(ctx context.Context, id string) (FileData, error)
	FileSearch(ctx context.Context, title string) ([]FileData, error)
}

// NewAppClient - конструктор объекта реализующего интерфейс AppClient
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
		notes:      notepb.NewNoteClient(connect),
		files:      filepb.NewFilesClient(connect),
		fileReader: &FileReader{},
	}, nil
}

// AppClientImpl - реализация интерфейса AppClient
type AppClientImpl struct {
	cnf        config.AppConfigurator
	authClient authpb.AuthClient
	accounts   accpb.AccountsClient
	cards      cardpb.CardClient
	notes      notepb.NoteClient
	files      filepb.FilesClient
	fileReader FileManager
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
func (c *AppClientImpl) SearchAccount(ctx context.Context, search string) ([]AccountData, error) {
	req := &accpb.SearchAccountRequest{Search: search}

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

// NotesList - метод получения списка всех заметок пользователя
func (c *AppClientImpl) NotesList(ctx context.Context) ([]NoteData, error) {
	resp, err := c.notes.Search(ctx, &notepb.SearchNoteRequest{})
	if err != nil {
		return nil, err
	}
	res := make([]NoteData, 0)
	for _, row := range resp.GetNotes() {
		r := NoteData{
			ID:    row.GetID(),
			Title: row.GetTitle(),
			Meta:  row.GetMeta(),
		}
		res = append(res, r)
	}

	return res, nil
}

// NoteCreate - метод создания новой заметки
func (c *AppClientImpl) NoteCreate(ctx context.Context, data NoteData) error {
	_, err := c.notes.Create(ctx, &notepb.CreateNoteRequest{
		Title:   data.Title,
		Meta:    data.Meta,
		Content: data.Content,
	})
	if err != nil {
		return err
	}
	return nil
}

// NoteDelete - метод удаления заметки
func (c *AppClientImpl) NoteDelete(ctx context.Context, id string) error {
	_, err := c.notes.Remove(ctx, &notepb.RemoveNoteRequest{ID: id})
	if err != nil {
		return err
	}
	return nil
}

// NoteSearch - метод поиска заметок по фильтрам
func (c *AppClientImpl) NoteSearch(ctx context.Context, data NoteSearchData) ([]NoteData, error) {
	resp, err := c.notes.Search(ctx, &notepb.SearchNoteRequest{Search: data.Search})

	if err != nil {
		return nil, err
	}

	resList := make([]NoteData, 0)
	for _, row := range resp.GetNotes() {
		r := NoteData{
			ID:    row.GetID(),
			Title: row.GetTitle(),
			Meta:  row.GetMeta(),
		}
		resList = append(resList, r)
	}
	return resList, nil
}

// NoteDetail - метод получения детальной информации по заметки
func (c *AppClientImpl) NoteDetail(ctx context.Context, id string) (NoteData, error) {
	res, err := c.notes.Get(ctx, &notepb.GetNoteRequest{ID: id})
	if err != nil {
		return NoteData{}, err
	}

	note := res.GetNote()
	return NoteData{
		ID:      note.GetID(),
		Title:   note.GetTitle(),
		Meta:    note.GetMeta(),
		Content: note.GetContent(),
	}, nil
}

// FileList - метод получения списка всех файлов пользователя
func (c *AppClientImpl) FileList(ctx context.Context) ([]FileData, error) {
	resp, err := c.files.List(ctx, &filepb.ListFilesRequest{})
	if err != nil {
		return nil, err
	}

	resList := make([]FileData, 0)
	for _, row := range resp.GetFiles() {
		f := FileData{
			ID:   row.GetId(),
			Name: row.GetName(),
			Size: row.GetSize(),
			Meta: row.GetMeta(),
		}
		resList = append(resList, f)
	}
	return resList, nil
}

// FileAdd - метод добавления нового файла
func (c *AppClientImpl) FileAdd(ctx context.Context, data FileData) error {
	fileContent, err := c.fileReader.Read(data.Path)
	if err != nil {
		return err
	}

	req := &filepb.PutFileRequest{
		Name:    data.Name,
		Meta:    data.Meta,
		Content: fileContent,
	}

	_, err = c.files.Put(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// FileRemove - метод удаления файла
func (c *AppClientImpl) FileRemove(ctx context.Context, id string) error {
	_, err := c.files.Remove(ctx, &filepb.RemoveFileRequest{Id: id})
	if err != nil {
		return err
	}
	return nil
}

// FileDownload - метод скачивания файла
func (c *AppClientImpl) FileDownload(ctx context.Context, id string) (FileData, error) {
	resp, err := c.files.Download(ctx, &filepb.DownloadFileRequest{Id: id})
	if err != nil {
		return FileData{}, err
	}

	return FileData{Name: resp.GetName(), Content: resp.GetContent(), Size: resp.GetSize()}, nil
}

// FileSearch - метод поиска файлов по фильтрам
func (c *AppClientImpl) FileSearch(ctx context.Context, title string) ([]FileData, error) {
	resp, err := c.files.Search(ctx, &filepb.SearchFilesRequest{Name: title})
	if err != nil {
		return nil, err
	}

	resList := make([]FileData, 0)
	for _, row := range resp.GetFiles() {
		r := FileData{
			ID:   row.GetId(),
			Name: row.GetName(),
			Size: row.GetSize(),
			Meta: row.GetMeta(),
		}
		resList = append(resList, r)
	}

	return resList, nil
}
