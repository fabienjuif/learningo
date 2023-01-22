module github.com/fabienjuif/learningo/client

go 1.19

require (
	github.com/fabienjuif/learningo/libs/env v0.0.0
	github.com/fabienjuif/learningo/libs/com v0.0.0
	github.com/fabienjuif/learningo/libs/utils v0.0.0
	github.com/joho/godotenv v1.4.0
)

replace github.com/fabienjuif/learningo/libs/env v0.0.0 => ../libs/env
replace github.com/fabienjuif/learningo/libs/com v0.0.0 => ../libs/com
replace github.com/fabienjuif/learningo/libs/utils v0.0.0 => ../libs/utils