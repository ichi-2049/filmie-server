package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ichi-2049/filmie-server/db"
	domain "github.com/ichi-2049/filmie-server/internal/domain/models"
	"github.com/ichi-2049/filmie-server/internal/infrastructure/repositoryImpl"
	"github.com/ichi-2049/filmie-server/types"
	"github.com/ryanbradynd05/go-tmdb"
	"golang.org/x/time/rate"
)

const (
	maxPages   = 500  // TMDbの検索結果に対して取得できる最大ページ数
	startYear  = 1950 // 検索開始年度
	endYear    = 2025 // 検索終了年度
	maxWorkers = 10   // 同時に処理する年数
)

func main() {
	// 初期化処理
	ctx := context.Background()
	db := db.Init()
	repo := repositoryImpl.NewMovieRepositoryImpl(db)

	config := tmdb.Config{
		APIKey:   os.Getenv("TMDB_API_KEY"),
		Proxies:  nil,
		UseProxy: false,
	}
	tmdbClient := tmdb.Init(config)
	if tmdbClient == nil {
		panic("failed tmdb client")
	}

	// API全体で共有する単一のレートリミッター
	// TMDbのレート制限が秒間50リクエストのため、少し余裕を持って40に設定
	rateLimiter := rate.NewLimiter(rate.Every(time.Second/40), 40)

	var wg sync.WaitGroup
	yearCh := make(chan int)
	resultCh := make(chan string)

	// ワーカーの起動
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for year := range yearCh {
				if err := processYear(ctx, year, tmdbClient, repo, rateLimiter); err != nil {
					resultCh <- fmt.Sprintf("Error processing year %d: %v", year, err)
					continue
				}
				resultCh <- fmt.Sprintf("Successfully processed year %d", year)
			}
		}()
	}

	// 結果を受け取るゴルーチン
	go func() {
		for result := range resultCh {
			log.Println(result)
		}
	}()

	// 年の配布
	for year := startYear; year <= endYear; year++ {
		yearCh <- year
	}
	close(yearCh)

	// 全ワーカーの完了を待つ
	wg.Wait()
	close(resultCh)

	fmt.Println("All processing completed")
}

// 対象年度の作品情報を取得し、DBにUpsertする関数
func processYear(ctx context.Context, year int, tmdbClient *tmdb.TMDb, repo *repositoryImpl.MovieRepositoryImpl, rateLimiter *rate.Limiter) error {
	startDate := fmt.Sprintf("%d-01-01", year)
	endDate := fmt.Sprintf("%d-12-31", year)

	// TMDbの検索条件を設定
	options := map[string]string{
		"primary_release_date.gte": startDate,
		"primary_release_date.lte": endDate,
		"language":                 "ja-JP",           // 作品情報を日本語で取得
		"region":                   "JP",              // 日本で公開された作品を取得
		"sort_by":                  "popularity.desc", // 人気順でソートして取得
		"page":                     "1",
	}

	var allMovies []*domain.Movie
	currentPage := 1
	totalPages := maxPages

	for currentPage <= totalPages {
		// レート制限のチェックを行い、40リクエストを超える場合リクエスト可能になるまで待機
		if err := rateLimiter.Wait(ctx); err != nil {
			return fmt.Errorf("rate limit error: %v", err)
		}

		// ページを設定してTMDb APIを打鍵
		options["page"] = fmt.Sprintf("%d", currentPage)
		pageResult, err := tmdbClient.DiscoverMovie(options)
		if err != nil {
			return fmt.Errorf("failed to fetch page %d: %v", currentPage, err)
		}

		// 取得できる最大ページ数の制限があるため、総ページ数が500を超えていたらそれ以上取得しないようにする
		if currentPage == 1 {
			totalPages = pageResult.TotalPages
			if totalPages > maxPages {
				totalPages = maxPages
			}
		}

		// 取得した作品情報をモデルに格納してスライスに追加
		for _, m := range pageResult.Results {
			releaseDate, err := time.Parse("2006-01-02", m.ReleaseDate)
			if err != nil {
				continue
			}

			movie := &domain.Movie{
				MovieID:      m.ID,
				Title:        m.Title,
				Overview:     m.Overview,
				ReleaseDate:  *types.NewDate(releaseDate),
				TMDBImageURL: m.PosterPath,
				Popularity:   m.Popularity,
				VoteAverage:  m.VoteAverage,
				VoteCount:    m.VoteCount,
			}
			allMovies = append(allMovies, movie)
		}

		currentPage++
	}

	if len(allMovies) > 0 {
		// 1000件ごとにチャンクしてInsert
		const batchSize = 1000
		for i := 0; i < len(allMovies); i += batchSize {
			end := i + batchSize
			if end > len(allMovies) {
				end = len(allMovies)
			}
			if err := repo.BulkInsertMovies(allMovies[i:end]); err != nil {
				return fmt.Errorf("failed to insert movies batch: %v", err)
			}
		}
	}

	return nil
}
