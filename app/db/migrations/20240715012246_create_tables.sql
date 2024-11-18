-- +goose Up
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    `uid` VARCHAR(128) NOT NULL PRIMARY KEY COMMENT 'ユーザーID（cognito_sub）',
    `name` VARCHAR(255) NOT NULL COMMENT 'ユーザー名',
    `email` VARCHAR(255) NOT NULL COMMENT 'メールアドレス',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時'
) COMMENT 'ユーザー';

DROP TABLE IF EXISTS movies;
CREATE TABLE movies (
    `movie_id` INT NOT NULL PRIMARY KEY COMMENT '映画ID',
    `title` VARCHAR(255) NOT NULL COMMENT '映画タイトル',
    `overview` TEXT NOT NULL COMMENT '概要',
    `release_date` DATE NOT NULL COMMENT '公開日',
    `tmdb_image_url` VARCHAR(255) NOT NULL COMMENT 'TMDB画像URL',
    `s3_image_url` VARCHAR(255) NOT NULL COMMENT 'S3画像URL',
    `popularity` FLOAT NOT NULL COMMENT '人気度',
    `original_language` VARCHAR(10) NOT NULL COMMENT '言語',
    `vote_average` FLOAT NOT NULL COMMENT '平均評価',
    `vote_count` INT UNSIGNED NOT NULL COMMENT '評価数',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    KEY `title` (`title`),
    KEY `release_date` (`release_date`)
) COMMENT '映画';

DROP TABLE IF EXISTS genre_masters;
CREATE TABLE genre_masters (
    `genre_id` INT UNSIGNED NOT NULL PRIMARY KEY COMMENT 'ジャンルID',
    `name` VARCHAR(255) NOT NULL COMMENT 'ジャンル名',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時'
) COMMENT 'ジャンルマスタ';

DROP TABLE IF EXISTS movie_genres;
CREATE TABLE movie_genres (
    `movie_id` INT NOT NULL COMMENT '映画ID',
    `genre_id` INT UNSIGNED NOT NULL COMMENT 'ジャンルID',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    PRIMARY KEY (`movie_id`, `genre_id`),
    FOREIGN KEY (`movie_id`) REFERENCES movies(`movie_id`) ON DELETE CASCADE,
    FOREIGN KEY (`genre_id`) REFERENCES genre_masters(`genre_id`) ON DELETE CASCADE
) COMMENT '映画ジャンル';

DROP TABLE IF EXISTS user_reviews;
CREATE TABLE user_reviews (
    `uid` VARCHAR(128) NOT NULL COMMENT 'ユーザーID',
    `movie_id` INT NOT NULL COMMENT '映画ID',
    `content` TEXT NOT NULL COMMENT 'レビュー内容',
    `rating` FLOAT NOT NULL COMMENT '星評価',
    `review_date` DATETIME NOT NULL COMMENT 'レビュー日時',
    `is_publish` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '公開フラグ',
    `is_spoiler` BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'ネタバレフラグ',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    PRIMARY KEY (`uid`, `movie_id`),
    FOREIGN KEY (`uid`) REFERENCES users(`uid`) ON DELETE CASCADE,
    FOREIGN KEY (`movie_id`) REFERENCES movies(`movie_id`) ON DELETE CASCADE
) COMMENT 'ユーザーレビュー';

DROP TABLE IF EXISTS user_favorites;
CREATE TABLE user_favorites (
    `uid` VARCHAR(128) NOT NULL COMMENT 'ユーザーID',
    `movie_id` INT NOT NULL COMMENT '映画ID',
    `favorite_date` DATETIME NOT NULL COMMENT 'お気に入り日時',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    PRIMARY KEY (`uid`, `movie_id`),
    FOREIGN KEY (`uid`) REFERENCES users(`uid`) ON DELETE CASCADE,
    FOREIGN KEY (`movie_id`) REFERENCES movies(`movie_id`) ON DELETE CASCADE,
    KEY `favorite_date` (`favorite_date`) 
) COMMENT 'ユーザーお気に入り';

DROP TABLE IF EXISTS trending_movies;
CREATE TABLE trending_movies (
    `movie_id` INT NOT NULL COMMENT '映画ID',
    `record_date` DATE NOT NULL COMMENT '記録日',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    PRIMARY KEY (`movie_id`, `record_date`),
    FOREIGN KEY (`movie_id`) REFERENCES movies(`movie_id`) ON DELETE CASCADE
) COMMENT 'トレンド映画';

DROP TABLE IF EXISTS now_playing_movies;
CREATE TABLE now_playing_movies (
    `movie_id` INT NOT NULL COMMENT '映画ID',
    `record_date` DATE NOT NULL COMMENT '記録日',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    PRIMARY KEY (`movie_id`, `record_date`),
    FOREIGN KEY (`movie_id`) REFERENCES movies(`movie_id`) ON DELETE CASCADE
) COMMENT '上映中映画';

-- +goose Down
DROP TABLE IF EXISTS now_playing_movies;
DROP TABLE IF EXISTS trending_movies;
DROP TABLE IF EXISTS user_favorites;
DROP TABLE IF EXISTS user_reviews;
DROP TABLE IF EXISTS movie_genres;
DROP TABLE IF EXISTS genre_masters;
DROP TABLE IF EXISTS movies;
DROP TABLE IF EXISTS users;