package impl

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/doar"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	"github.com/training-of-new-employees/qon/internal/store"
	"github.com/training-of-new-employees/qon/internal/store/cache"
	mock_doar "github.com/training-of-new-employees/qon/mocks/pkg/doar"
	mock_jwttoken "github.com/training-of-new-employees/qon/mocks/pkg/jwttoken"
	mock_store "github.com/training-of-new-employees/qon/mocks/store"
	mock_cache "github.com/training-of-new-employees/qon/mocks/store/cache"
)

func Test_newUserService(t *testing.T) {
	type args struct {
		db         store.Storages
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		cache      cache.Cache
		jwtGen     jwttoken.JWTGenerator
		jwtVal     jwttoken.JWTValidator
		sender     doar.EmailSender
		host       string
	}
	tests := []struct {
		name string
		args args
		want *uService
	}{
		{
			"Stub creation",
			args{},
			&uService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUserService(tt.args.db, tt.args.secretKey, tt.args.aTokenTime, tt.args.rTokenTime, tt.args.cache, tt.args.jwtGen, tt.args.jwtVal, tt.args.sender, tt.args.host); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_WriteAdminToCache(t *testing.T) {
	type fields struct {
		userdb     *mock_store.MockRepositoryUser
		cache      *mock_cache.MockCache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     *mock_doar.MockEmailSender
	}
	type args struct {
		ctx context.Context
		val model.CreateAdmin
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *model.CreateAdmin
		wantErr bool
	}{
		{
			"User exist",
			func(f *fields) {
				u := &model.User{
					Email: "email@mail.com",
				}
				f.userdb.EXPECT().GetUserByEmail(nil, u.Email).Return(u, nil)
			},
			args{
				nil,
				model.CreateAdmin{
					Email: "email@mail.com",
				},
			},
			nil,
			true,
		},
		{
			"Can't set cache",
			func(f *fields) {
				email := "email@mail.com"
				f.userdb.EXPECT().GetUserByEmail(nil, email).Return(nil, errs.ErrUserNotFound)
				f.cache.EXPECT().Set(nil, gomock.Any(), gomock.Any()).Return(errs.ErrInternal)
			},
			args{
				nil,
				model.CreateAdmin{
					Email: "email@mail.com",
				},
			},
			nil,
			true,
		},
		{
			"Can't send code",
			func(f *fields) {
				email := "email@mail.com"
				f.userdb.EXPECT().GetUserByEmail(nil, email).Return(nil, errs.ErrUserNotFound)
				f.cache.EXPECT().Set(nil, gomock.Any(), gomock.Any()).Return(nil)
				f.sender.EXPECT().Mode().Return("api")
				f.sender.EXPECT().SendCode(email, gomock.Any()).Return(errs.ErrInternal)
			},
			args{
				nil,
				model.CreateAdmin{
					Email: "email@mail.com",
				},
			},
			&model.CreateAdmin{
				Email: "email@mail.com",
			},
			false,
		},
		{
			"Success write cache",
			func(f *fields) {
				email := "email@mail.com"
				f.userdb.EXPECT().GetUserByEmail(nil, email).Return(nil, errs.ErrUserNotFound)
				f.cache.EXPECT().Set(nil, gomock.Any(), gomock.Any()).Return(nil)
				f.sender.EXPECT().SendCode(email, gomock.Any()).Return(nil)
			},
			args{
				nil,
				model.CreateAdmin{
					Email: "email@mail.com",
				},
			},
			&model.CreateAdmin{
				Email: "email@mail.com",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{}
			f.userdb = mock_store.NewMockRepositoryUser(ctrl)
			f.cache = mock_cache.NewMockCache(ctrl)
			f.sender = mock_doar.NewMockEmailSender(ctrl)
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			storages := mockUserStorage(ctrl, f.userdb)
			u := &uService{
				db:         storages,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			got, err := u.WriteAdminToCache(tt.args.ctx, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.WriteAdminToCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Email != tt.want.Email {
				t.Errorf("uService.WriteAdminToCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_GetUserByID(t *testing.T) {
	type fields struct {
		userdb     *mock_store.MockRepositoryUser
		companydb  *mock_store.MockRepositoryCompany
		posdb      *mock_store.MockRepositoryPosition
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *model.UserInfo
		wantErr bool
	}{
		{
			"User Not Found",
			func(f *fields) {
				f.userdb.EXPECT().GetUserByID(nil, 1).Return(nil, errs.ErrUserNotFound)
			},
			args{nil, 1},
			nil,
			true,
		},
		{
			"Company Not Found",
			func(f *fields) {
				u := &model.User{
					Email:     "email@mail.com",
					CompanyID: 1,
				}
				f.userdb.EXPECT().GetUserByID(nil, 1).Return(u, nil)
				f.companydb.EXPECT().GetCompany(nil, 1).Return(nil, errs.ErrCompanyNotFound)
			},
			args{nil, 1},
			nil,
			true,
		},
		{
			"Position Not Found",
			func(f *fields) {
				u := &model.User{
					Email:      "email@mail.com",
					CompanyID:  1,
					PositionID: 1,
				}
				company := &model.Company{
					Name: "company",
				}
				f.userdb.EXPECT().GetUserByID(nil, 1).Return(u, nil)
				f.companydb.EXPECT().GetCompany(nil, 1).Return(company, nil)
				f.posdb.EXPECT().GetPositionInCompany(nil, 1, 1).Return(nil, errs.ErrPositionNotFound)
			},
			args{nil, 1},
			nil,
			true,
		},
		{
			"Get user success",
			func(f *fields) {
				u := &model.User{
					Email:      "email@mail.com",
					CompanyID:  1,
					PositionID: 1,
				}
				company := &model.Company{
					Name: "company",
				}
				pos := &model.Position{
					Name: "position",
				}
				f.userdb.EXPECT().GetUserByID(nil, 1).Return(u, nil)
				f.companydb.EXPECT().GetCompany(nil, 1).Return(company, nil)
				f.posdb.EXPECT().GetPositionInCompany(nil, 1, 1).Return(pos, nil)
			},
			args{nil, 1},
			&model.UserInfo{
				User: model.User{
					Email:      "email@mail.com",
					CompanyID:  1,
					PositionID: 1,
				},
				CompanyName:  "company",
				PositionName: "position",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{}
			f.userdb = mock_store.NewMockRepositoryUser(ctrl)
			f.companydb = mock_store.NewMockRepositoryCompany(ctrl)
			f.posdb = mock_store.NewMockRepositoryPosition(ctrl)
			if tt.prepare != nil {
				tt.prepare(f)
			}

			storages := mockUCPStorage(ctrl, f.userdb, f.companydb, f.posdb)

			u := &uService{
				db:         storages,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			got, err := u.GetUserByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_ArchiveUser(t *testing.T) {
	type fields struct {
		userdb     *mock_store.MockRepositoryUser
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx             context.Context
		id              int
		editorCompanyID int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *model.UserEdit
		wantErr bool
	}{
		{
			"User not found",
			func(f *fields) {
				f.userdb.EXPECT().GetUserByID(nil, 1).Return(nil, errs.ErrUserNotFound)
			},
			args{
				nil,
				1,
				2,
			},
			nil,
			true,
		},
		{
			"Editor company id is missing",
			func(f *fields) {
				u := &model.User{
					ID:        1,
					CompanyID: 1,
				}
				f.userdb.EXPECT().GetUserByID(nil, 1).Return(u, nil)
			},
			args{
				nil,
				1,
				2,
			},
			nil,
			true,
		},
		{
			"Can't edit user",
			func(f *fields) {
				u := &model.User{
					ID:        1,
					CompanyID: 1,
				}

				f.userdb.EXPECT().GetUserByID(nil, 1).Return(u, nil)
				f.userdb.EXPECT().EditUser(nil, gomock.Any()).Return(nil, errs.ErrInternal)
			},
			args{
				nil,
				1,
				1,
			},
			nil,
			true,
		},
		{
			"Can edit user",
			func(f *fields) {
				u := &model.User{
					ID:        1,
					CompanyID: 1,
				}
				edit := &model.UserEdit{
					ID: 1,
				}

				f.userdb.EXPECT().GetUserByID(nil, 1).Return(u, nil)
				f.userdb.EXPECT().EditUser(nil, gomock.Any()).Return(edit, nil)
			},
			args{
				nil,
				1,
				1,
			},
			&model.UserEdit{
				ID: 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{}
			f.userdb = mock_store.NewMockRepositoryUser(ctrl)
			if tt.prepare != nil {
				tt.prepare(f)
			}
			storages := mockUserStorage(ctrl, f.userdb)
			u := &uService{
				db:         storages,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}

			if err := u.ArchiveUser(tt.args.ctx, tt.args.id, tt.args.editorCompanyID); (err != nil) != tt.wantErr {
				t.Errorf("uService.ArchiveUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_uService_EditUser(t *testing.T) {
	type fields struct {
		userdb     *mock_store.MockRepositoryUser
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx             context.Context
		val             *model.UserEdit
		editorCompanyID int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *model.UserEdit
		wantErr bool
	}{
		{
			"User not found",
			func(f *fields) {
				f.userdb.EXPECT().GetUserByID(nil, 1).Return(nil, errs.ErrUserNotFound)
			},
			args{
				nil,
				&model.UserEdit{
					ID: 1,
				},
				2,
			},
			nil,
			true,
		},
		{
			"Editor company id is missing",
			func(f *fields) {
				u := &model.User{
					ID:        1,
					CompanyID: 1,
				}
				f.userdb.EXPECT().GetUserByID(nil, 1).Return(u, nil)
			},
			args{
				nil,
				&model.UserEdit{
					ID: 1,
				},
				2,
			},
			nil,
			true,
		},
		{
			"Can edit user",
			func(f *fields) {
				u := &model.User{
					ID:        1,
					CompanyID: 1,
				}
				edit := &model.UserEdit{
					ID: 1,
				}

				f.userdb.EXPECT().GetUserByID(nil, 1).Return(u, nil)
				f.userdb.EXPECT().EditUser(nil, gomock.Any()).Return(edit, nil)
			},
			args{
				nil,
				&model.UserEdit{
					ID: 1,
				},
				1,
			},
			&model.UserEdit{
				ID: 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{}
			f.userdb = mock_store.NewMockRepositoryUser(ctrl)
			if tt.prepare != nil {
				tt.prepare(f)
			}
			storages := mockUserStorage(ctrl, f.userdb)
			u := &uService{
				db:         storages,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			got, err := u.EditUser(tt.args.ctx, tt.args.val, tt.args.editorCompanyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.EditUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.EditUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_GetUsersByCompany(t *testing.T) {
	type fields struct {
		userdb     *mock_store.MockRepositoryUser
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx       context.Context
		companyID int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    []model.User
		wantErr bool
	}{
		{
			"Users not found",
			func(f *fields) {
				f.userdb.EXPECT().GetUsersByCompany(nil, 1).Return(nil, errs.ErrNotFound)
			},
			args{nil, 1},
			nil,
			true,
		},
		{
			"Users found",
			func(f *fields) {
				us := []model.User{
					{
						Email:     "u1@mail.com",
						CompanyID: 1,
					},
					{
						Email:     "u2@mail.com",
						CompanyID: 1,
					},
				}
				f.userdb.EXPECT().GetUsersByCompany(nil, 1).Return(us, nil)
			},
			args{nil, 1},
			[]model.User{
				{
					Email:     "u1@mail.com",
					CompanyID: 1,
				},
				{
					Email:     "u2@mail.com",
					CompanyID: 1,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{}
			f.userdb = mock_store.NewMockRepositoryUser(ctrl)
			if tt.prepare != nil {
				tt.prepare(f)
			}
			storages := mockUserStorage(ctrl, f.userdb)
			u := &uService{
				db:         storages,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			got, err := u.GetUsersByCompany(tt.args.ctx, tt.args.companyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.GetUsersByCompany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.GetUsersByCompany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_GenerateTokenPair(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      *mock_cache.MockCache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   *mock_jwttoken.MockJWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx       context.Context
		userId    int
		isAdmin   bool
		companyID int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *model.Tokens
		wantErr bool
	}{
		{
			"Error Generate access token",
			func(f *fields) {
				f.aTokenTime = 12 * time.Hour
				f.rTokenTime = 120 * time.Hour
				r := "refresh"

				f.tokenGen.EXPECT().GenerateToken(1, true, 1, "", f.rTokenTime).Return(r, nil)

				hasher := sha1.New()
				hasher.Write([]byte(r))
				hashedRefresh := hex.EncodeToString(hasher.Sum(nil))

				f.tokenGen.EXPECT().GenerateToken(1, true, 1, hashedRefresh, f.aTokenTime).Return("", errs.ErrInternal)
			},
			args{
				nil,
				1, true, 1},
			nil,
			true,
		},
		{
			"Error Generate refresh token",
			func(f *fields) {
				f.aTokenTime = 12 * time.Hour
				f.rTokenTime = 120 * time.Hour
				f.tokenGen.EXPECT().GenerateToken(1, true, 1, "", f.rTokenTime).Return("", errs.ErrInternal)
			},
			args{
				nil,
				1, true, 1},
			nil,
			true,
		},
		{
			"Success Generate token",
			func(f *fields) {
				f.aTokenTime = 12 * time.Hour
				f.rTokenTime = 120 * time.Hour
				a := "access"
				r := "refresh"
				f.tokenGen.EXPECT().GenerateToken(1, true, 1, "", f.rTokenTime).Return(r, nil)

				hasher := sha1.New()
				hasher.Write([]byte(r))
				hashedRefresh := hex.EncodeToString(hasher.Sum(nil))

				f.tokenGen.EXPECT().GenerateToken(1, true, 1, hashedRefresh, f.aTokenTime).Return(a, nil)
				f.cache.EXPECT().SetRefreshToken(gomock.Any(), hashedRefresh, r).Return(nil)
			},
			args{
				nil,
				1, true, 1},
			&model.Tokens{
				AccessToken:  "access",
				RefreshToken: "refresh",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := &fields{}
			f.tokenGen = mock_jwttoken.NewMockJWTGenerator(ctrl)
			f.cache = mock_cache.NewMockCache(ctrl)
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := &uService{
				db:         f.db,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			got, err := u.GenerateTokenPair(tt.args.ctx, tt.args.userId, tt.args.isAdmin, tt.args.companyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.GenerateTokenPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.GenerateTokenPair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_CreateUser(t *testing.T) {
	type fields struct {
		userdb     *mock_store.MockRepositoryUser
		cache      *mock_cache.MockCache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     *mock_doar.MockEmailSender
	}
	type args struct {
		ctx context.Context
		val model.UserCreate
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			"DB Create User Error",
			func(f *fields) {
				f.userdb.EXPECT().CreateUser(nil, gomock.Any()).Return(nil, errs.ErrEmailAlreadyExists)
			},
			args{
				nil,
				model.UserCreate{
					Email:    "user@mail.com",
					Password: "password",
				},
			},
			nil,
			true,
		},
		{
			"Sender Invite User Error",
			func(f *fields) {
				u := &model.User{
					ID:    1,
					Email: "user@mail.com",
				}
				f.userdb.EXPECT().CreateUser(nil, gomock.Any()).Return(u, nil)
				f.sender.EXPECT().InviteUser(u.Email, gomock.Any()).Return(errs.ErrInternal)
				f.cache.EXPECT().SetInviteCode(nil, gomock.Any(), gomock.Any()).Return(errs.ErrInternal)
				f.sender.EXPECT().Mode().Return("api")
			},
			args{
				nil,
				model.UserCreate{
					Email:    "user@mail.com",
					Password: "password",
				},
			},
			&model.User{
				ID:    1,
				Email: "user@mail.com",
			},
			false,
		},
		{
			"Sender Invite User success",
			func(f *fields) {
				u := &model.User{
					ID:    1,
					Email: "user@mail.com",
				}
				f.userdb.EXPECT().CreateUser(nil, gomock.Any()).Return(u, nil)
				f.sender.EXPECT().InviteUser(u.Email, gomock.Any()).Return(nil)
				f.cache.EXPECT().SetInviteCode(nil, gomock.Any(), gomock.Any()).Return(nil)
			},
			args{
				nil,
				model.UserCreate{
					Email:    "user@mail.com",
					Password: "password",
				},
			},
			&model.User{
				ID:    1,
				Email: "user@mail.com",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{}
			f.cache = mock_cache.NewMockCache(ctrl)
			f.userdb = mock_store.NewMockRepositoryUser(ctrl)
			f.sender = mock_doar.NewMockEmailSender(ctrl)
			if tt.prepare != nil {
				tt.prepare(f)
			}
			storages := mockUserStorage(ctrl, f.userdb)
			u := &uService{
				db:         storages,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			got, err := u.CreateUser(tt.args.ctx, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_GetAdminFromCache(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      *mock_cache.MockCache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx  context.Context
		code string
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *model.CreateAdmin
		wantErr bool
	}{
		{
			"Admin not found",
			func(f *fields) {
				f.cache.EXPECT().Get(nil, "register:admin:1234").Return(nil, errs.ErrNotFound)
			},
			args{
				nil,
				"1234",
			},
			nil,
			true,
		},
		{
			"Admin found",
			func(f *fields) {
				a := &model.CreateAdmin{
					Email:    "admin@mail.com",
					Password: "password",
					Company:  "admin",
				}
				f.cache.EXPECT().Get(nil, "register:admin:1234").Return(a, nil)
			},
			args{
				nil,
				"1234",
			},
			&model.CreateAdmin{
				Email:    "admin@mail.com",
				Password: "password",
				Company:  "admin",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{}
			f.cache = mock_cache.NewMockCache(ctrl)
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			u := &uService{
				db:         f.db,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			got, err := u.GetAdminFromCache(tt.args.ctx, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.GetAdminFromCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.GetAdminFromCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_DeleteAdminFromCache(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      *mock_cache.MockCache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		wantErr bool
	}{
		{
			"Delete not exist admin",
			func(f *fields) {
				f.cache.EXPECT().Delete(nil, "key").Return(errs.ErrNotFound)
			},
			args{
				nil,
				"key",
			},
			true,
		},
		{
			"Delete exist admin",
			func(f *fields) {
				f.cache.EXPECT().Delete(nil, "key").Return(nil)
			},
			args{
				nil,
				"key",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{}
			f.cache = mock_cache.NewMockCache(ctrl)
			if tt.prepare != nil {
				tt.prepare(f)
			}

			u := &uService{
				db:         f.db,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			if err := u.DeleteAdminFromCache(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("uService.DeleteAdminFromCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_uService_CreateAdmin(t *testing.T) {
	type fields struct {
		userdb     *mock_store.MockRepositoryUser
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx context.Context
		val model.CreateAdmin
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			"Error get user",
			func(f *fields) {
				f.userdb.EXPECT().GetUserByEmail(nil, gomock.Any()).Return(nil, errs.ErrInternal)
			},
			args{
				nil,
				model.CreateAdmin{
					Company:  "invalid",
					Email:    "invalid@mail.com",
					Password: "password",
				},
			},
			nil,
			true,
		},
		{
			"Email is exist",
			func(f *fields) {
				email := "exist@mail.com"
				u := &model.User{
					ID:        1,
					Email:     email,
					CompanyID: 2,
				}
				f.userdb.EXPECT().GetUserByEmail(nil, email).Return(u, nil)
			},
			args{
				nil,
				model.CreateAdmin{
					Company:  "exist",
					Email:    "exist@mail.com",
					Password: "password",
				},
			},
			nil,
			true,
		},
		{
			"Create admin fail",
			func(f *fields) {
				f.userdb.EXPECT().GetUserByEmail(nil, gomock.Any()).Return(nil, errs.ErrUserNotFound)
				f.userdb.EXPECT().CreateAdmin(nil, gomock.Any(), gomock.Any()).Return(nil, errs.ErrInternal)
			},
			args{
				nil,
				model.CreateAdmin{
					Company:  "notexist",
					Email:    "notexist@mail.com",
					Password: "password",
				},
			},
			nil,
			true,
		},
		{
			"Create admin success",
			func(f *fields) {
				admin := &model.User{
					ID:    2,
					Email: "admin@mail.com",
				}
				f.userdb.EXPECT().GetUserByEmail(nil, gomock.Any()).Return(nil, errs.ErrUserNotFound)
				f.userdb.EXPECT().CreateAdmin(nil, gomock.Any(), gomock.Any()).Return(admin, nil)
			},
			args{
				nil,
				model.CreateAdmin{
					Company:  "new company",
					Email:    "admin@mail.com",
					Password: "password",
				},
			},
			&model.User{
				ID:    2,
				Email: "admin@mail.com",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userdb := mock_store.NewMockRepositoryUser(ctrl)
			f := &fields{
				userdb: userdb,
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			storages := mockUserStorage(ctrl, f.userdb)

			u := &uService{
				db:         storages,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			got, err := u.CreateAdmin(tt.args.ctx, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.CreateAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.CreateAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_UpdatePasswordAndActivateUser(t *testing.T) {
	type fields struct {
		userdb     *mock_store.MockRepositoryUser
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		wantErr bool
	}{
		{
			"Email is not exist",
			func(f *fields) {
				f.userdb.EXPECT().GetUserByEmail(nil, "notexist@mail.com").Return(nil, errs.ErrUserNotFound)
			},
			args{
				nil,
				"notexist@mail.com",
				"",
			},
			true,
		},
		{
			"User Set Empty Password Error",
			func(f *fields) {
				u := &model.User{
					ID:    1,
					Email: "valid@mail.com",
				}
				f.userdb.EXPECT().GetUserByEmail(nil, u.Email).Return(u, nil)
				f.userdb.EXPECT().SetPasswordAndActivateUser(nil, u.ID, gomock.Any()).Return(errs.ErrInternal)
			},
			args{
				nil,
				"valid@mail.com",
				"",
			},
			true,
		},
		{
			"User Set Password",
			func(f *fields) {
				u := &model.User{
					ID:    1,
					Email: "valid@mail.com",
				}
				f.userdb.EXPECT().GetUserByEmail(nil, u.Email).Return(u, nil)
				f.userdb.EXPECT().SetPasswordAndActivateUser(nil, u.ID, gomock.Any()).Return(nil)
			},
			args{
				nil,
				"valid@mail.com",
				"password",
			},
			false,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		udb := mock_store.NewMockRepositoryUser(ctrl)
		f := fields{
			userdb: udb,
		}
		if tt.prepare != nil {
			tt.prepare(&f)
		}
		storages := mockUserStorage(ctrl, udb)
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         storages,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			if err := u.UpdatePasswordAndActivateUser(tt.args.ctx, tt.args.email, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("uService.UpdatePasswordAndActivateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_uService_ResetPassword(t *testing.T) {
	type fields struct {
		userdb     *mock_store.MockRepositoryUser
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     *mock_doar.MockEmailSender
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		wantErr bool
	}{
		{
			"Email is not exist",
			func(f *fields) {
				f.userdb.EXPECT().GetUserByEmail(nil, "notexist@mail.com").Return(nil, errs.ErrUserNotFound)
			},
			args{
				nil,
				"notexist@mail.com",
			},
			true,
		},
		{
			"User empty Email",
			func(f *fields) {
				f.userdb.EXPECT().GetUserByEmail(nil, "").Return(nil, errs.ErrUserNotFound)
			},
			args{
				nil,
				"",
			},
			true,
		},
		{
			"Update User Password error",
			func(f *fields) {
				u := &model.User{
					ID:    3,
					Email: "valid@mail.com",
				}
				f.userdb.EXPECT().GetUserByEmail(nil, "valid@mail.com").Return(u, nil)
				f.userdb.EXPECT().UpdateUserPassword(nil, 3, gomock.Any()).Return(errs.ErrInternal)
			},
			args{
				nil,
				"valid@mail.com",
			},
			true,
		},
		{
			"Send User Password Error",
			func(f *fields) {
				u := &model.User{
					ID:    3,
					Email: "valid@mail.com",
				}
				f.userdb.EXPECT().GetUserByEmail(nil, "valid@mail.com").Return(u, nil)
				f.userdb.EXPECT().UpdateUserPassword(nil, 3, gomock.Any()).Return(nil)
				f.sender.EXPECT().SendPassword("valid@mail.com", gomock.Any()).Return(errs.ErrInternal)
			},
			args{
				nil,
				"valid@mail.com",
			},
			true,
		},
		{
			"Send User Password",
			func(f *fields) {
				u := &model.User{
					ID:    3,
					Email: "valid@mail.com",
				}
				f.userdb.EXPECT().GetUserByEmail(nil, "valid@mail.com").Return(u, nil)
				f.userdb.EXPECT().UpdateUserPassword(nil, 3, gomock.Any()).Return(nil)
				f.sender.EXPECT().SendPassword("valid@mail.com", gomock.Any()).Return(nil)
			},
			args{
				nil,
				"valid@mail.com",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			udb := mock_store.NewMockRepositoryUser(ctrl)
			sender := mock_doar.NewMockEmailSender(ctrl)
			f := fields{userdb: udb, sender: sender}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			storages := mockUserStorage(ctrl, f.userdb)

			u := &uService{
				db:         storages,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			if err := u.ResetPassword(tt.args.ctx, tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("uService.ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_uService_EditAdmin(t *testing.T) {
	type fields struct {
		userdb     *mock_store.MockRepositoryUser
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx context.Context
		val model.AdminEdit
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		setargs func() args
		want    *model.AdminEdit
		wantErr bool
	}{
		{
			"Invalid Data",
			nil,
			func() args {
				email := "a"
				aedit := model.AdminEdit{
					ID:    1,
					Email: &email,
				}
				return args{
					nil,
					aedit,
				}
			},
			nil,
			true,
		},
		{
			"Valid data, db problem",
			func(f *fields) {
				email := "valid@mail.com"
				aedit := model.AdminEdit{
					ID:    1,
					Email: &email,
				}
				f.userdb.EXPECT().EditAdmin(gomock.Any(), aedit).Return(nil, errs.ErrNotFound)
			},
			func() args {
				email := "valid@mail.com"
				aedit := model.AdminEdit{
					ID:    1,
					Email: &email,
				}
				return args{
					nil,
					aedit,
				}
			},
			nil,
			true,
		},
		{
			"Valid Data",
			func(f *fields) {
				email := "valid@mail.com"
				aedit := model.AdminEdit{
					ID:    1,
					Email: &email,
				}
				f.userdb.EXPECT().EditAdmin(gomock.Any(), aedit).Return(&aedit, nil)

			},
			func() args {
				email := "valid@mail.com"
				aedit := model.AdminEdit{
					ID:    1,
					Email: &email,
				}
				return args{
					nil,
					aedit,
				}
			},
			&model.AdminEdit{
				ID: 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{}
			f.userdb = mock_store.NewMockRepositoryUser(ctrl)
			if tt.prepare != nil {
				tt.prepare(f)
			}
			storages := mockUserStorage(ctrl, f.userdb)
			u := &uService{
				db:         storages,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			args := tt.setargs()
			got, err := u.EditAdmin(args.ctx, args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.EditAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("uService.EditAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUService_ClearSession(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      *mock_cache.MockCache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   *mock_jwttoken.MockJWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx           context.Context
		hashedRefresh string
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		wantErr bool
	}{
		{
			"Error Clear Session",
			func(f *fields) {
				f.cache.EXPECT().DeleteRefreshToken(gomock.Any(), "hashed").Return(errs.ErrInternal)
			},
			args{nil, "hashed"},
			true,
		},
		{
			"Success Clear Session",
			func(f *fields) {
				f.cache.EXPECT().DeleteRefreshToken(gomock.Any(), "hashed").Return(nil)
			},
			args{nil, "hashed"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := &fields{}
			f.tokenGen = mock_jwttoken.NewMockJWTGenerator(ctrl)
			f.cache = mock_cache.NewMockCache(ctrl)
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := &uService{
				db:         f.db,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			err := u.ClearSession(tt.args.ctx, tt.args.hashedRefresh)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.ClearSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_uService_RegenerationInvitationLinkUser(t *testing.T) {
	type fields struct {
		userDB     *mock_store.MockRepositoryUser
		cache      *mock_cache.MockCache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     *mock_doar.MockEmailSender
		host       string
	}
	type args struct {
		ctx       context.Context
		email     string
		companyID int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    string
		wantErr bool
	}{
		{
			"Regeneration Invite Link User Error",
			func(f *fields) {
				f.userDB.EXPECT().GetUserByEmail(nil, "user@mail.com").Return(nil, errs.ErrUserNotFound)
			},
			args{
				ctx:       nil,
				email:     "user@mail.com",
				companyID: 1,
			},
			"",
			true,
		},
		{
			name: "Regeneration Invite Link User success",
			prepare: func(f *fields) {
				u := &model.User{
					ID:        1,
					Email:     "user@mail.com",
					CompanyID: 1,
					IsActive:  false,
				}

				f.cache.EXPECT().SetInviteCode(nil, gomock.Any(), gomock.Any()).Return(nil)
				f.userDB.EXPECT().GetUserByEmail(nil, "user@mail.com").Return(u, nil)
				f.sender.EXPECT().InviteUser(u.Email, gomock.Any()).Return(nil)

				f.host = "http://localhost"
			},
			args: args{
				ctx:       nil,
				email:     "user@mail.com",
				companyID: 1,
			},
			want:    "http://localhost/first-login\\?email=user@mail.com&invite=*",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{}
			f.cache = mock_cache.NewMockCache(ctrl)
			f.userDB = mock_store.NewMockRepositoryUser(ctrl)
			f.sender = mock_doar.NewMockEmailSender(ctrl)
			if tt.prepare != nil {
				tt.prepare(f)
			}
			storages := mockUserStorage(ctrl, f.userDB)
			u := &uService{
				db:         storages,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
				host:       f.host,
			}

			got, err := u.RegenerationInvitationLinkUser(tt.args.ctx, tt.args.email, tt.args.companyID)
			fmt.Println(err)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.RegenerationInvitationLinkUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(tt.want) > 0 {
				assert.Regexp(t, tt.want, got.Link)
			}
		})
	}
}

func mockUserStorage(ctrl *gomock.Controller, uStore *mock_store.MockRepositoryUser) *mock_store.MockStorages {
	posStore := mock_store.NewMockRepositoryPosition(ctrl)
	storages := mock_store.NewMockStorages(ctrl)
	storages.EXPECT().UserStorage().Return(uStore).AnyTimes()
	storages.EXPECT().PositionStorage().Return(posStore).AnyTimes()
	return storages
}

func mockUCPStorage(ctrl *gomock.Controller, uStore *mock_store.MockRepositoryUser, cStore *mock_store.MockRepositoryCompany, pStore *mock_store.MockRepositoryPosition) *mock_store.MockStorages {
	storages := mock_store.NewMockStorages(ctrl)
	storages.EXPECT().UserStorage().Return(uStore).AnyTimes()
	storages.EXPECT().CompanyStorage().Return(cStore).AnyTimes()
	storages.EXPECT().PositionStorage().Return(pStore).AnyTimes()
	return storages

}
