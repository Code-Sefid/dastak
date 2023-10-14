package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/common"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/constants"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
	"github.com/soheilkhaledabdi/dastak/pkg/service_errors"

	"gorm.io/gorm"
)

type preload struct {
	string
}

type BaseService[T any, Tc any, Tu any, Tr any] struct {
	Database *gorm.DB
	Logger   logging.Logger
	Config   *config.Config
	Preloads []preload
}
type Updater interface {
    UpdateFields(updateData interface{}) error
}

func NewBaseService[T any, Tc any, Tu any, Tr any](cfg *config.Config) *BaseService[T, Tc, Tu, Tr] {
	return &BaseService[T, Tc, Tu, Tr]{
		Database: db.GetDb(),
		Logger:   logging.NewLogger(cfg),
	}
}

func (s *BaseService[T, Tc, Tu, Tr]) Create(ctx context.Context, req *Tc) (*Tr, error) {
	model, _ := common.TypeConverter[T](req)
	tx := s.Database.WithContext(ctx).Begin()
	err := tx.
		Create(model).
		Error
	if err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return nil, err
	}
	tx.Commit()
	bm, _ := common.TypeConverter[models.BaseModel](model)
	return s.GetByIdWithoutUserID(ctx, bm.ID)
}

func (s *BaseService[T, Tc, Tu, Tr]) CreateByUserId(ctx context.Context, req *Tc, userID int) (*Tr, error) {
	model, _ := common.TypeConverter[T](req)

	modelValue := reflect.ValueOf(model).Elem()
	userIDField := modelValue.FieldByName("UserID")
	if userIDField.IsValid() && userIDField.CanSet() {
		userIDField.SetInt(int64(userID))
	}

	tx := s.Database.WithContext(ctx).Begin()
	err := tx.Create(model).Error
	if err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return nil, err
	}
	tx.Commit()

	bm, _ := common.TypeConverter[models.BaseModel](model)
	return s.GetById(ctx, bm.ID,userID)
}


func (s *BaseService[T, Tc, Tu, Tr]) Update(ctx context.Context, id int , userID int , req *Tu) (*Tr, error) {

	updateMap, _ := common.TypeConverter[map[string]interface{}](req)
	snakeMap := map[string]interface{}{}
	for k, v := range *updateMap {
		snakeMap[common.ToSnakeCase(k)] = v
	}
	snakeMap["modified_by"] = &sql.NullInt64{Int64: int64(ctx.Value(constants.UserIdKey).(float64)), Valid: true}
	snakeMap["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}
	model := new(T)
	tx := s.Database.WithContext(ctx).Begin()
	if err := tx.Model(model).
		Where("id = ? and user_id = ? AND deleted_by is null", id,userID).
		Updates(snakeMap).
		Error; err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return nil, err
	}
	tx.Commit()
	return s.GetById(ctx, id,userID)

}


func (s *BaseService[T, Tc, Tu, Tr]) Delete(ctx context.Context, id int,userID int) error {
	tx := s.Database.WithContext(ctx).Begin()

	model := new(T)

	deleteMap := map[string]interface{}{
		"deleted_by": &sql.NullInt64{Int64: int64(ctx.Value(constants.UserIdKey).(float64)), Valid: true},
		"deleted_at": sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}

	if ctx.Value(constants.UserIdKey) == nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	if err := tx.
		Model(model).
		Where("id = ? AND user_id = ? AND deleted_by is null", id,userID).
		Updates(deleteMap).
		Error; err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Update, service_errors.RecordNotFound, nil)
		return err
	}
	tx.Commit()
	return nil
}

func (s *BaseService[T, Tc, Tu, Tr]) GetByUserId(ctx context.Context, userId int) (*Tr, error) {
	model := new(T)
	db := Preload(s.Database, s.Preloads)
	err := db.
		Where("user_id = ? and deleted_by is null", userId).
		First(model).
		Error
	if err != nil {
		return nil, err
	}
	return common.TypeConverter[Tr](model)
}

func (s *BaseService[T, Tc, Tu, Tr]) GetByIdWithoutUserID(ctx context.Context, id int) (*Tr, error) {
	model := new(T)
	db := Preload(s.Database, s.Preloads)
	err := db.
		Where("id = ? AND deleted_by is null", id).
		First(model).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return common.TypeConverter[Tr](model)
}

func (s *BaseService[T, Tc, Tu, Tr]) GetById(ctx context.Context, id int,userID int) (*Tr, error) {
	model := new(T)
	db := Preload(s.Database, s.Preloads)
	err := db.
		Where("id = ? AND user_id = ? AND deleted_by is null", id,userID).
		First(model).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return common.TypeConverter[Tr](model)
}

func (s *BaseService[T, Tc, Tu, Tr]) GetByFilter(ctx context.Context, req *dto.PaginationInputWithFilter, userId int) (*dto.PagedList[Tr], error) {
	res, err := Paginate[T, Tr](req, s.Preloads, s.Database, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewPagedList[T any](items *[]T, count int64, pageNumber int, pageSize int64) *dto.PagedList[T] {
	pl := &dto.PagedList[T]{
		PageNumber: pageNumber,
		TotalRows:  count,
		Items:      items,
	}
	pl.TotalPages = int(math.Ceil(float64(count) / float64(pageSize)))
	pl.HasNextPage = pl.PageNumber < pl.TotalPages
	pl.HasPreviousPage = pl.PageNumber > 1

	return pl
}

// Paginate
func Paginate[T any, Tr any](pagination *dto.PaginationInputWithFilter, preloads []preload, db *gorm.DB, userId int) (*dto.PagedList[Tr], error) {
	model := new(T)
	var items *[]T
	var rItems *[]Tr
	db = Preload(db, preloads)
	query := getQuery[T](&pagination.DynamicFilter)
	sort := getSort[T](&pagination.DynamicFilter)

	var totalRows int64 = 0

	db.
		Model(model).
		Where("user_id", userId).
		Where(query).
		Count(&totalRows)

	err := db.
		Where(query).
		Where("user_id", userId).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetPageSize()).
		Order(sort).
		Find(&items).
		Error

	if err != nil {
		return nil, err
	}
	rItems, err = common.TypeConverter[[]Tr](items)
	if err != nil {
		return nil, err
	}
	return NewPagedList(rItems, totalRows, pagination.PageNumber, int64(pagination.PageSize)), err

}

// getQuery
func getQuery[T any](filter *dto.DynamicFilter) string {
	t := new(T)
	typeT := reflect.TypeOf(*t)
	query := make([]string, 0)
	query = append(query, "deleted_by is null")
	if filter.Filter != nil {
		for name, filter := range filter.Filter {
			fld, ok := typeT.FieldByName(name)
			if ok {
				fld.Name = common.ToSnakeCase(fld.Name)
				switch filter.Type {
				case "contains":
					query = append(query, fmt.Sprintf("%s ILike '%%%s%%'", fld.Name, filter.From))
				case "notContains":
					query = append(query, fmt.Sprintf("%s not ILike '%%%s%%'", fld.Name, filter.From))
				case "startsWith":
					query = append(query, fmt.Sprintf("%s ILike '%s%%'", fld.Name, filter.From))
				case "endsWith":
					query = append(query, fmt.Sprintf("%s ILike '%%%s'", fld.Name, filter.From))
				case "equals":
					query = append(query, fmt.Sprintf("%s = '%s'", fld.Name, filter.From))
				case "notEqual":
					query = append(query, fmt.Sprintf("%s != '%s'", fld.Name, filter.From))
				case "lessThan":
					query = append(query, fmt.Sprintf("%s < %s", fld.Name, filter.From))
				case "lessThanOrEqual":
					query = append(query, fmt.Sprintf("%s <= %s", fld.Name, filter.From))
				case "greaterThan":
					query = append(query, fmt.Sprintf("%s > %s", fld.Name, filter.From))
				case "greaterThanOrEqual":
					query = append(query, fmt.Sprintf("%s >= %s", fld.Name, filter.From))
				case "inRange":
					if fld.Type.Kind() == reflect.String {
						query = append(query, fmt.Sprintf("%s >= '%s'", fld.Name, filter.From))
						query = append(query, fmt.Sprintf("%s <= '%s'", fld.Name, filter.To))
					} else {
						query = append(query, fmt.Sprintf("%s >= %s", fld.Name, filter.From))
						query = append(query, fmt.Sprintf("%s <= %s", fld.Name, filter.To))
					}
				}
			}
		}
	}
	return strings.Join(query, " AND ")
}

// getSort
func getSort[T any](filter *dto.DynamicFilter) string {
	t := new(T)
	typeT := reflect.TypeOf(*t)
	sort := make([]string, 0)
	if filter.Sort != nil {
		for _, tp := range *filter.Sort {
			fld, ok := typeT.FieldByName(tp.ColId)
			if ok && (tp.Sort == "asc" || tp.Sort == "desc") {
				fld.Name = common.ToSnakeCase(fld.Name)
				sort = append(sort, fmt.Sprintf("%s %s", fld.Name, tp.Sort))
			}
		}
	}
	return strings.Join(sort, ", ")
}

// Preload
func Preload(db *gorm.DB, preloads []preload) *gorm.DB {
	for _, item := range preloads {
		db = db.Preload(item.string)
	}
	return db
}
