package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/constant"
	"github.com/nozzlium/eniqilo_store/internal/model"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(
	db *pgx.Conn,
) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repository *UserRepository) Save(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	query := `
    insert into users
    (
      id, 
      phone_number,
      password,
      name
    ) values 
    (
      $1,
      $2,
      $3,
      $4 
    );
  `

	_, err := repository.db.Exec(
		ctx,
		query,
		user.ID,
		user.PhoneNumber,
		user.Password,
		user.Name,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (repository *UserRepository) FindByPhoneNumber(
	ctx context.Context,
	phone_number string,
) (model.User, error) {
	query := `
    select
      id,
      phone_number,
      name,
      password
    from users 
    where 
      phone_number = $1;
  `
	user := model.User{}
	err := repository.db.QueryRow(ctx, query, phone_number).
		Scan(&user.ID, &user.PhoneNumber, &user.Name, &user.Password)
	if err != nil {
		log.Println(err)
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return user, constant.ErrNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (repository *UserRepository) FindByPhoneNumberAndID(
	ctx context.Context,
	id string,
	phone_number string,
) (model.User, error) {
	query := `
    select
      id,
      phone_number,
      name,
      password
    from users 
    where 
      id = $1 and phone_number = $2;
  `
	user := model.User{}
	err := repository.db.QueryRow(ctx, query, id, phone_number).
		Scan(&user.ID, &user.PhoneNumber, &user.Name, &user.Password)
	if err != nil {
		log.Println(err)
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return user, constant.ErrNotFound
		}
		return model.User{}, err
	}

	return user, nil
}
