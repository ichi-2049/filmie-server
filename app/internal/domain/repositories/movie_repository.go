package repositories

import domain "github.com/ichi-2049/filmie-server/internal/domain/models"

type MovieRepository interface {
	/*
	映画作品情報をカーソルページネーションで取得する関数
	タイトル（optional）で検索をかけ、人気順（降順）と映画ID（昇順）でソートする
	*/
	GetMovieConnection(first int, after *string, title *string) (*domain.MovieConnection, error)
}
