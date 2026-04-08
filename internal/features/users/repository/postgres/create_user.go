package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mirwinli/golang-todoapp/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.users (full_name,phone_number)
	VALUES ($1, $2)
	RETURNING id, version, full_name, phone_number;
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)

	var userModel UserModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	); err != nil {
		fmt.Printf("scan err: %v, unwrapped: %v\n", err, errors.Unwrap(err))
		return domain.User{}, fmt.Errorf("scan error %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
