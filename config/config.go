package config

import "time"

var DBSource = "root:secret@tcp(127.0.0.1:3306)/movie_hub?parseTime=true"
var TokenSymmetricKey = "12345678901234567890123456789012"
var Duration time.Duration = 15 * time.Minute
