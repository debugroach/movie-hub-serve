DROP TABLE if EXISTS ratings;
DROP TABLE if EXISTS users;
DROP TABLE if EXISTS movies;
CREATE TABLE `users` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `username` VARCHAR(255) NOT NULL UNIQUE,
  `password` VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE `ratings` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `username` VARCHAR(255) NOT NULL,
  `movie_id` INT NOT NULL,
  `rating` INT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE `movies` (
  `id` INT PRIMARY KEY,
  `title` VARCHAR(255) NOT NULL,
  `poster_path` VARCHAR(255) NOT NULL,
  `backdrop_path` VARCHAR(255) NOT NULL,
  `vote_average` DOUBLE NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX `index_users_on_username` ON `users` (`username`);
CREATE INDEX `index_ratings_on_username` ON `ratings` (`username`);
CREATE INDEX `index_ratings_on_movie_id` ON `ratings` (`movie_id`);
ALTER TABLE `ratings`
ADD FOREIGN KEY (`username`) REFERENCES `users` (`username`);