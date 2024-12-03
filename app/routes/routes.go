package routes

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ichi-2049/filmie-server/db"
	graph "github.com/ichi-2049/filmie-server/graphql"
	"github.com/ichi-2049/filmie-server/graphql/resolver"
	"github.com/ichi-2049/filmie-server/graphql/resolver/container"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRoutes(e *echo.Echo) {
	// CORS設定の追加
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000", // Next.jsの開発サーバー
			// "https://hoge.com", // 本番環境のドメイン
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	// DB初期化
	db := db.Init()

	// DIコンテナの初期化
	container := container.NewContainer(db)

	// リゾルバの初期化
	resolver := resolver.NewResolver(container)

	graphqlHandler := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: resolver},
		),
	)
	playgroundHandler := playground.Handler("GraphQL", "/query")

	e.POST("/query", func(c echo.Context) error {
		graphqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET("/playground", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
}
