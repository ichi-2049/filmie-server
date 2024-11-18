package container

import (
	"sync"

	"github.com/ichi-2049/filmie-server/internal/application/services"
	"github.com/ichi-2049/filmie-server/internal/infrastructure/repositoryImpl"
	"gorm.io/gorm"
)

type Container struct {
	db           *gorm.DB
	movieService *services.MovieService
	once         sync.Once
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
	})
	return instance
}

func (c *Container) GetMovieService() *services.MovieService {
	c.once.Do(func() {
		movieRepo := repositoryImpl.NewMovieRepositoryImpl(c.db)
		c.movieService = services.NewMovieService(movieRepo)
	})
	return c.movieService
}
