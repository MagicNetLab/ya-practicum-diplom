package client

import (
	"context"
	"errors"
	"testing"

	mockAccClient "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/mocks"
	accProto "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	mockAuthClient "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/mocks"
	authProto "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
	mockCardClient "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/mocks"
	cardProto "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
	mockFileClient "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/mocks"
	fileProto "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
	mockNoteClient "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/mocks"
	noteProto "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// FileReaderMock - это мок для интерфейса FileManager
type FileReaderMock struct {
	mock.Mock
}

// Read чтение содержание файла
func (f *FileReaderMock) Read(filename string) ([]byte, error) {
	ret := f.Called(filename)

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]byte, error)); ok {
		return rf(filename)
	}
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(filename)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TestAuth - тест функции Auth
func TestAuth(t *testing.T) {
	mockAuth := new(mockAuthClient.AuthClient)
	client := &AppClientImpl{authClient: mockAuth}

	t.Run("Успешная аутентификация", func(t *testing.T) {
		mockAuth.On("Auth", mock.Anything, &authProto.AuthRequest{Login: "testuser", Secret: "password"}).
			Return(&authProto.AuthResponse{Token: "valid_token"}, nil)

		token, err := client.Auth(context.Background(), "testuser", "password")
		assert.NoError(t, err)
		assert.Equal(t, "valid_token", token)
	})

	t.Run("Ошибка аутентификации: неверный логин или пароль", func(t *testing.T) {
		mockAuth.On("Auth", mock.Anything, &authProto.AuthRequest{Login: "wronguser", Secret: "wrongpassword"}).
			Return(nil, errors.New("invalid credentials"))

		token, err := client.Auth(context.Background(), "wronguser", "wrongpassword")
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

// TestRegister - тест функции Register
func TestRegister(t *testing.T) {
	mockAuth := new(mockAuthClient.AuthClient)
	client := &AppClientImpl{authClient: mockAuth}

	t.Run("Успешная регистрация", func(t *testing.T) {
		mockAuth.On("Register", mock.Anything, &authProto.RegRequest{Login: "newuser", Secret: "password"}).
			Return(&authProto.RegResponse{Token: "new_valid_token"}, nil)

		token, err := client.Register(context.Background(), "newuser", "password")
		assert.NoError(t, err)
		assert.Equal(t, "new_valid_token", token)
	})

	t.Run("Ошибка регистрации: логин уже существует", func(t *testing.T) {
		mockAuth.On("Register", mock.Anything, &authProto.RegRequest{Login: "existinguser", Secret: "password"}).
			Return(nil, errors.New("user already exists"))

		token, err := client.Register(context.Background(), "existinguser", "password")
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

// TestListAccounts - тест вывода списка аккаунтов
func TestListAccounts(t *testing.T) {
	t.Run("Успешное получение списка аккаунтов", func(t *testing.T) {
		mockAccounts := new(mockAccClient.AccountsClient)
		client := &AppClientImpl{accounts: mockAccounts}
		mockAccounts.On("Search", mock.Anything, mock.Anything).
			Return(&accProto.SearchAccountResponse{
				Acc: []*accProto.Account{
					{Id: "1", Login: "user1", Url: "url1", Description: "desc1"},
					{Id: "2", Login: "user2", Url: "url2", Description: "desc2"},
				},
			}, nil)

		accounts, err := client.ListAccounts(context.Background())
		assert.NoError(t, err)
		assert.Len(t, accounts, 2)
		assert.Equal(t, "1", accounts[0].ID)
		assert.Equal(t, "user1", accounts[0].Login)
	})

	t.Run("Ошибка получения списка аккаунтов", func(t *testing.T) {
		mockAccounts := new(mockAccClient.AccountsClient)
		client := &AppClientImpl{accounts: mockAccounts}
		mockAccounts.On("Search", mock.Anything, mock.Anything).
			Return(nil, errors.New("failed to fetch accounts"))

		accounts, err := client.ListAccounts(context.Background())
		assert.Error(t, err)
		assert.Nil(t, accounts)
	})
}

// TestAddAccount - тест добавления аккаунта
func TestAddAccount(t *testing.T) {
	t.Run("Успешное добавление аккаунта", func(t *testing.T) {
		mockAccounts := new(mockAccClient.AccountsClient)
		client := &AppClientImpl{accounts: mockAccounts}
		mockAccounts.On("Create", mock.Anything, mock.Anything).
			Return(&accProto.CreateAccountResponse{}, nil)

		err := client.AddAccount(context.Background(), "newuser", "password", "url", "meta")
		assert.NoError(t, err)
	})

	t.Run("Ошибка добавления аккаунта", func(t *testing.T) {
		mockAccounts := new(mockAccClient.AccountsClient)
		client := &AppClientImpl{accounts: mockAccounts}
		mockAccounts.On("Create", mock.Anything, mock.Anything).
			Return(nil, errors.New("failed to create account"))

		err := client.AddAccount(context.Background(), "newuser", "password", "url", "meta")
		assert.Error(t, err)
	})
}

// TestRemoveAccount - тест удаления аккаунта
func TestRemoveAccount(t *testing.T) {
	t.Run("Успешное удаление аккаунта", func(t *testing.T) {
		mockAccounts := new(mockAccClient.AccountsClient)
		client := &AppClientImpl{accounts: mockAccounts}
		mockAccounts.On("Remove", mock.Anything, mock.Anything).
			Return(&accProto.RemoveAccountResponse{}, nil)

		err := client.RemoveAccount(context.Background(), "1")
		assert.NoError(t, err)
	})

	t.Run("Ошибка удаления аккаунта", func(t *testing.T) {
		mockAccounts := new(mockAccClient.AccountsClient)
		client := &AppClientImpl{accounts: mockAccounts}
		mockAccounts.On("Remove", mock.Anything, mock.Anything).
			Return(nil, errors.New("account not found"))

		err := client.RemoveAccount(context.Background(), "1")
		assert.Error(t, err)
	})
}

// TestGetAccount - тест получения аккаунта
func TestGetAccount(t *testing.T) {
	t.Run("Успешное получение аккаунта", func(t *testing.T) {
		mockAccounts := new(mockAccClient.AccountsClient)
		client := &AppClientImpl{accounts: mockAccounts}
		mockAccounts.On("Get", mock.Anything, mock.Anything).
			Return(&accProto.GetAccountResponse{
				Account: &accProto.Account{Id: "1", Login: "user1", Url: "url1", Description: "desc1"},
			}, nil)

		account, err := client.GetAccount(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, "1", account.ID)
		assert.Equal(t, "user1", account.Login)
	})

	t.Run("Ошибка получения аккаунта", func(t *testing.T) {
		mockAccounts := new(mockAccClient.AccountsClient)
		client := &AppClientImpl{accounts: mockAccounts}
		mockAccounts.On("Get", mock.Anything, mock.Anything).
			Return(nil, errors.New("account not found"))

		account, err := client.GetAccount(context.Background(), "1")
		assert.Error(t, err)
		assert.Empty(t, account)
	})
}

// TestSearchAccount - тест поиска аккаунтов
func TestSearchAccount(t *testing.T) {
	t.Run("Успешный поиск аккаунтов", func(t *testing.T) {
		mockAccounts := new(mockAccClient.AccountsClient)
		client := &AppClientImpl{accounts: mockAccounts}
		mockAccounts.On("Search", mock.Anything, mock.Anything).
			Return(&accProto.SearchAccountResponse{
				Acc: []*accProto.Account{
					{Id: "1", Login: "user1", Url: "url1", Description: "desc1"},
					{Id: "2", Login: "user2", Url: "url2", Description: "desc2"},
				},
			}, nil)

		accounts, err := client.SearchAccount(context.Background(), "user")
		assert.NoError(t, err)
		assert.Len(t, accounts, 2)
		assert.Equal(t, "1", accounts[0].ID)
	})

	t.Run("Ошибка поиска аккаунтов", func(t *testing.T) {
		mockAccounts := new(mockAccClient.AccountsClient)
		client := &AppClientImpl{accounts: mockAccounts}
		mockAccounts.On("Search", mock.Anything, mock.Anything).
			Return(nil, errors.New("search failed"))

		accounts, err := client.SearchAccount(context.Background(), "user")
		assert.Error(t, err)
		assert.Nil(t, accounts)
	})
}

// TestAddCard - тест добавления карты
func TestAddCard(t *testing.T) {
	t.Run("Успешное добавление карты", func(t *testing.T) {
		mockCards := new(mockCardClient.CardClient)
		client := &AppClientImpl{cards: mockCards}
		mockCards.On("Create", mock.Anything, mock.Anything).
			Return(&cardProto.CreateCardResponse{}, nil)

		err := client.AddCard(context.Background(), CardData{
			Name:   "Test Card",
			Number: "1234567812345678",
			Month:  12,
			Year:   2025,
			CVC:    "123",
			PIN:    "1234",
			Meta:   "Test Meta",
		})
		assert.NoError(t, err)
	})

	t.Run("Ошибка добавления карты", func(t *testing.T) {
		mockCards := new(mockCardClient.CardClient)
		client := &AppClientImpl{cards: mockCards}
		mockCards.On("Create", mock.Anything, mock.Anything).
			Return(nil, errors.New("failed to create card"))

		err := client.AddCard(context.Background(), CardData{
			Name:   "Test Card",
			Number: "1234567812345678",
			Month:  12,
			Year:   2025,
			CVC:    "123",
			PIN:    "1234",
			Meta:   "Test Meta",
		})
		assert.Error(t, err)
	})
}

// TestRemoveCard - тест удаления карты
func TestRemoveCard(t *testing.T) {
	t.Run("Успешное удаление карты", func(t *testing.T) {
		mockCards := new(mockCardClient.CardClient)
		client := &AppClientImpl{cards: mockCards}
		mockCards.On("Delete", mock.Anything, mock.Anything).
			Return(&cardProto.DeleteCardResponse{}, nil)

		err := client.RemoveCard(context.Background(), "1")
		assert.NoError(t, err)
	})

	t.Run("Ошибка удаления карты", func(t *testing.T) {
		mockCards := new(mockCardClient.CardClient)
		client := &AppClientImpl{cards: mockCards}
		mockCards.On("Delete", mock.Anything, mock.Anything).
			Return(nil, errors.New("card not found"))

		err := client.RemoveCard(context.Background(), "1")
		assert.Error(t, err)
	})
}

// TestcardDetail - тест получения информации о карте
func TestCardDetail(t *testing.T) {
	t.Run("Успешное получение информации о карте", func(t *testing.T) {
		mockCards := new(mockCardClient.CardClient)
		client := &AppClientImpl{cards: mockCards}
		mockCards.On("Get", mock.Anything, mock.Anything).
			Return(&cardProto.GetCardResponse{
				Card: &cardProto.CardModel{ID: "1", Name: "Test Card", Number: "1234567812345678", Month: 12, Year: 2025, CVC: "123", PIN: "1234", Meta: "Test Meta"},
			}, nil)

		card, err := client.CardDetail(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, "1", card.ID)
		assert.Equal(t, "Test Card", card.Name)
	})

	t.Run("Ошибка получения информации о карте", func(t *testing.T) {
		mockCards := new(mockCardClient.CardClient)
		client := &AppClientImpl{cards: mockCards}
		mockCards.On("Get", mock.Anything, mock.Anything).
			Return(nil, errors.New("card not found"))

		card, err := client.CardDetail(context.Background(), "1")
		assert.Error(t, err)
		assert.Empty(t, card)
	})
}

// TestCardSearch - тест поиска карт
func TestCardSearch(t *testing.T) {
	t.Run("Успешный поиск карт", func(t *testing.T) {
		mockCards := new(mockCardClient.CardClient)
		client := &AppClientImpl{cards: mockCards}
		mockCards.On("Search", mock.Anything, mock.Anything).
			Return(&cardProto.SearchCardResponse{
				Cards: []*cardProto.ShortCardModel{
					{ID: "1", Name: "Test Card 1", Number: "1234567812345678"},
					{ID: "2", Name: "Test Card 2", Number: "8765432187654321"},
				},
			}, nil)

		cards, err := client.CardSearch(context.Background(), CardSearchData{Name: "Test Card"})
		assert.NoError(t, err)
		assert.Len(t, cards, 2)
		assert.Equal(t, "1", cards[0].ID)
	})

	t.Run("Ошибка поиска карт", func(t *testing.T) {
		mockCards := new(mockCardClient.CardClient)
		client := &AppClientImpl{cards: mockCards}
		mockCards.On("Search", mock.Anything, mock.Anything).
			Return(nil, errors.New("search failed"))

		cards, err := client.CardSearch(context.Background(), CardSearchData{Name: "Test Card"})
		assert.Error(t, err)
		assert.Nil(t, cards)
	})
}

// TestNotesList - тест получения списка заметок
func TestNotesList(t *testing.T) {
	t.Run("Успешное получение списка заметок", func(t *testing.T) {
		mockNotes := new(mockNoteClient.NoteClient)
		client := &AppClientImpl{notes: mockNotes}
		mockNotes.On("Search", mock.Anything, mock.Anything).
			Return(&noteProto.SearchNoteResponse{
				Notes: []*noteProto.NoteModel{
					{ID: "1", Title: "Note 1", Meta: "Meta 1"},
					{ID: "2", Title: "Note 2", Meta: "Meta 2"},
				},
			}, nil)

		notes, err := client.NotesList(context.Background())
		assert.NoError(t, err)
		assert.Len(t, notes, 2)
		assert.Equal(t, "1", notes[0].ID)
		assert.Equal(t, "Note 1", notes[0].Title)
	})

	t.Run("Ошибка получения списка заметок", func(t *testing.T) {
		mockNotes := new(mockNoteClient.NoteClient)
		client := &AppClientImpl{notes: mockNotes}
		mockNotes.On("Search", mock.Anything, mock.Anything).
			Return(nil, errors.New("failed to fetch notes"))

		notes, err := client.NotesList(context.Background())
		assert.Error(t, err)
		assert.Nil(t, notes)
	})
}

// TestNoteCreate - тест создания заметки
func TestNoteCreate(t *testing.T) {
	t.Run("Успешное создание заметки", func(t *testing.T) {
		mockNotes := new(mockNoteClient.NoteClient)
		client := &AppClientImpl{notes: mockNotes}
		mockNotes.On("Create", mock.Anything, mock.Anything).
			Return(&noteProto.CreateNoteResponse{}, nil)

		err := client.NoteCreate(context.Background(), NoteData{
			Title:   "New Note",
			Content: "This is a note.",
			Meta:    "Test Meta",
		})
		assert.NoError(t, err)
	})

	t.Run("Ошибка создания заметки", func(t *testing.T) {
		mockNotes := new(mockNoteClient.NoteClient)
		client := &AppClientImpl{notes: mockNotes}
		mockNotes.On("Create", mock.Anything, mock.Anything).
			Return(nil, errors.New("failed to create note"))

		err := client.NoteCreate(context.Background(), NoteData{
			Title:   "New Note",
			Content: "This is a note.",
			Meta:    "Test Meta",
		})
		assert.Error(t, err)
	})
}

// TestNoteDelete - тест удаления заметки
func TestNoteDelete(t *testing.T) {
	t.Run("Успешное удаление заметки", func(t *testing.T) {
		mockNotes := new(mockNoteClient.NoteClient)
		client := &AppClientImpl{notes: mockNotes}
		mockNotes.On("Remove", mock.Anything, mock.Anything).
			Return(&noteProto.RemoveNoteResponse{}, nil)

		err := client.NoteDelete(context.Background(), "1")
		assert.NoError(t, err)
	})

	t.Run("Ошибка удаления заметки", func(t *testing.T) {
		mockNotes := new(mockNoteClient.NoteClient)
		client := &AppClientImpl{notes: mockNotes}
		mockNotes.On("Remove", mock.Anything, mock.Anything).
			Return(nil, errors.New("note not found"))

		err := client.NoteDelete(context.Background(), "1")
		assert.Error(t, err)
	})
}

// TestNoteSearch - тест поиска заметок
func TestNoteSearch(t *testing.T) {
	t.Run("Успешный поиск заметок", func(t *testing.T) {
		mockNotes := new(mockNoteClient.NoteClient)
		client := &AppClientImpl{notes: mockNotes}
		mockNotes.On("Search", mock.Anything, mock.Anything).
			Return(&noteProto.SearchNoteResponse{
				Notes: []*noteProto.NoteModel{
					{ID: "1", Title: "Note 1", Meta: "Meta 1"},
					{ID: "2", Title: "Note 2", Meta: "Meta 2"},
				},
			}, nil)

		notes, err := client.NoteSearch(context.Background(), NoteSearchData{Search: "Note"})
		assert.NoError(t, err)
		assert.Len(t, notes, 2)
		assert.Equal(t, "1", notes[0].ID)
	})

	t.Run("Ошибка поиска заметок", func(t *testing.T) {
		mockNotes := new(mockNoteClient.NoteClient)
		client := &AppClientImpl{notes: mockNotes}
		mockNotes.On("Search", mock.Anything, mock.Anything).
			Return(nil, errors.New("search failed"))

		notes, err := client.NoteSearch(context.Background(), NoteSearchData{Search: "Note"})
		assert.Error(t, err)
		assert.Nil(t, notes)
	})
}

// TestNoteDetail - тест получения информации о заметке
func TestNoteDetail(t *testing.T) {
	t.Run("Успешное получение информации о заметке", func(t *testing.T) {
		mockNotes := new(mockNoteClient.NoteClient)
		client := &AppClientImpl{notes: mockNotes}
		mockNotes.On("Get", mock.Anything, mock.Anything).
			Return(&noteProto.GetNoteResponse{
				Note: &noteProto.NoteModel{ID: "1", Title: "Note 1", Content: "Content of Note 1", Meta: "Meta 1"},
			}, nil)

		note, err := client.NoteDetail(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, "1", note.ID)
		assert.Equal(t, "Note 1", note.Title)
	})

	t.Run("Ошибка получения информации о заметке", func(t *testing.T) {
		mockNotes := new(mockNoteClient.NoteClient)
		client := &AppClientImpl{notes: mockNotes}
		mockNotes.On("Get", mock.Anything, mock.Anything).
			Return(nil, errors.New("note not found"))

		note, err := client.NoteDetail(context.Background(), "1")
		assert.Error(t, err)
		assert.Empty(t, note)
	})
}

// TestFileList - тест получения списка файлов
func TestFileList(t *testing.T) {
	t.Run("Успешное получение списка файлов", func(t *testing.T) {
		mockFiles := new(mockFileClient.FilesClient)
		client := &AppClientImpl{files: mockFiles}
		mockFiles.On("List", mock.Anything, mock.Anything).
			Return(&fileProto.ListFilesResponse{
				Files: []*fileProto.FileModel{
					{Id: "1", Name: "File 1", Size: 1024, Meta: "Meta 1"},
					{Id: "2", Name: "File 2", Size: 2048, Meta: "Meta 2"},
				},
			}, nil)

		files, err := client.FileList(context.Background())
		assert.NoError(t, err)
		assert.Len(t, files, 2)
		assert.Equal(t, "1", files[0].ID)
	})

	t.Run("Ошибка получения списка файлов", func(t *testing.T) {
		mockFiles := new(mockFileClient.FilesClient)
		client := &AppClientImpl{files: mockFiles}
		mockFiles.On("List", mock.Anything, mock.Anything).
			Return(nil, errors.New("failed to fetch files"))

		files, err := client.FileList(context.Background())
		assert.Error(t, err)
		assert.Nil(t, files)
	})
}

// TestFileAdd - тест добавления файла
func TestFileAdd(t *testing.T) {
	t.Run("Успешное добавление файла", func(t *testing.T) {
		mockFiles := new(mockFileClient.FilesClient)
		mockFiles.On("Put", mock.Anything, mock.Anything).
			Return(&fileProto.PutFileResponse{}, nil)
		mockFileReader := new(FileReaderMock)
		mockFileReader.On("Read", mock.Anything, mock.Anything).Return([]byte("file content"), nil)
		client := &AppClientImpl{files: mockFiles, fileReader: mockFileReader}

		err := client.FileAdd(context.Background(), FileData{
			Name: "New File",
			Path: "path/to/file",
			Meta: "Test Meta",
		})
		assert.NoError(t, err)
	})

	t.Run("Ошибка чтения файла", func(t *testing.T) {
		mockFiles := new(mockFileClient.FilesClient)
		mockFileReader := new(FileReaderMock)
		mockFiles.On("Put", mock.Anything, mock.Anything).
			Return(&fileProto.PutFileResponse{}, nil)
		mockFileReader.On("Read", mock.Anything, mock.Anything).Return(nil, errors.New("failed to read file"))
		client := &AppClientImpl{files: mockFiles, fileReader: mockFileReader}

		err := client.FileAdd(context.Background(), FileData{
			Name: "New File",
			Path: "path/to/file",
			Meta: "Test Meta",
		})
		assert.Error(t, err)
	})

	t.Run("Ошибка добавления файла", func(t *testing.T) {
		mockFiles := new(mockFileClient.FilesClient)
		mockFiles.On("Put", mock.Anything, mock.Anything).
			Return(nil, errors.New("failed to add file"))
		mockFileReader := new(FileReaderMock)
		mockFileReader.On("Read", mock.Anything, mock.Anything).Return([]byte("file content"), nil)
		client := &AppClientImpl{files: mockFiles, fileReader: mockFileReader}

		err := client.FileAdd(context.Background(), FileData{
			Name: "New File",
			Path: "path/to/file",
			Meta: "Test Meta",
		})
		assert.Error(t, err)
	})
}

// TestFileRemove - тест удаления файла
func TestFileRemove(t *testing.T) {
	t.Run("Успешное удаление файла", func(t *testing.T) {
		mockFiles := new(mockFileClient.FilesClient)
		client := &AppClientImpl{files: mockFiles}
		mockFiles.On("Remove", mock.Anything, mock.Anything).
			Return(&fileProto.RemoveFileResponse{}, nil)

		err := client.FileRemove(context.Background(), "1")
		assert.NoError(t, err)
	})

	t.Run("Ошибка удаления файла", func(t *testing.T) {
		mockFiles := new(mockFileClient.FilesClient)
		client := &AppClientImpl{files: mockFiles}
		mockFiles.On("Remove", mock.Anything, mock.Anything).
			Return(nil, errors.New("file not found"))

		err := client.FileRemove(context.Background(), "1")
		assert.Error(t, err)
	})
}

// TestFileDownload - тест скачивания файла
func TestFileDownload(t *testing.T) {
	t.Run("Успешное скачивание файла", func(t *testing.T) {
		mockFiles := new(mockFileClient.FilesClient)
		client := &AppClientImpl{files: mockFiles}
		mockFiles.On("Download", mock.Anything, mock.Anything).
			Return(&fileProto.DownloadFileResponse{
				Name:    "File 1",
				Content: []byte("file content"),
				Size:    1024,
			}, nil)

		file, err := client.FileDownload(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, "File 1", file.Name)
	})

	t.Run("Ошибка скачивания файла", func(t *testing.T) {
		mockFiles := new(mockFileClient.FilesClient)
		client := &AppClientImpl{files: mockFiles}
		mockFiles.On("Download", mock.Anything, mock.Anything).
			Return(nil, errors.New("file not found"))

		file, err := client.FileDownload(context.Background(), "1")
		assert.Error(t, err)
		assert.Empty(t, file)
	})
}

// TestFileSearch - тест поиска файлов
func TestFileSearch(t *testing.T) {
	t.Run("Успешный поиск файлов", func(t *testing.T) {
		mockFiles := new(mockFileClient.FilesClient)
		client := &AppClientImpl{files: mockFiles}
		mockFiles.On("Search", mock.Anything, mock.Anything).
			Return(&fileProto.SearchFilesResponse{
				Files: []*fileProto.FileModel{
					{Id: "1", Name: "File 1", Size: 1024, Meta: "Meta 1"},
					{Id: "2", Name: "File 2", Size: 2048, Meta: "Meta 2"},
				},
			}, nil)

		files, err := client.FileSearch(context.Background(), "File")
		assert.NoError(t, err)
		assert.Len(t, files, 2)
		assert.Equal(t, "1", files[0].ID)
	})

	t.Run("Ошибка поиска файлов", func(t *testing.T) {
		mockFiles := new(mockFileClient.FilesClient)
		client := &AppClientImpl{files: mockFiles}
		mockFiles.On("Search", mock.Anything, mock.Anything).
			Return(nil, errors.New("search failed"))

		files, err := client.FileSearch(context.Background(), "File")
		assert.Error(t, err)
		assert.Nil(t, files)
	})
}
