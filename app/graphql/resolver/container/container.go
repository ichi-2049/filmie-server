package container

import (
	"sync"

	"github.com/ichi-2049/filmie-server/internal/application/services"
	"github.com/ichi-2049/filmie-server/internal/domain/repositories"
	"github.com/ichi-2049/filmie-server/internal/infrastructure/repositoryImpl"
	"gorm.io/gorm"
)

type Container struct {
	once         sync.Once
	db           *gorm.DB
	movieRepo    repositories.MovieRepository
	userRepo     repositories.UserRepository
	movieService *services.MovieService
	userService  *services.UserService
}

var (
	instance *Container
	once     sync.Once
)

func NewContainer(db *gorm.DB) *Container {
	once.Do(func() {
		instance = &Container{
			db: db,
		}
		instance.initialize()
	})
	return instance
}

func (c *Container) initialize() {
	c.once.Do(func() {
		c.movieRepo = repositoryImpl.NewMovieRepositoryImpl(c.db)
		c.userRepo = repositoryImpl.NewUserRepositoryImpl(c.db)
		c.movieService = services.NewMovieService(c.movieRepo)
		c.userService = services.NewUserService(c.userRepo)
	})
}

func (c *Container) GetMovieService() *services.MovieService {
	return c.movieService
}

func (c *Container) GetUserService() *services.UserService {
	return c.userService
}
