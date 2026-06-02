// outbound (secondary adapter)

package ports

import (
	"context"

	"github.com/HtetAungKhant23/velora/internal/core/domain/user"
)

type UserRepository interface {
	ExistsByEmail(ctx context.Context, email user.Email) (bool, error)
	FindByEmail(ctx context.Context, email user.Email) (*user.User, error)
	FindByID(ctx context.Context, id user.UserID) (*user.User, error)
	Save(ctx context.Context, user *user.User) error
}
