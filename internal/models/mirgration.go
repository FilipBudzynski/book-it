package models

var MigrateModels = []any{
	&User{},
	&UserBook{},
	&Book{},
	&ReadingProgress{},
	&DailyProgressLog{},
	&ExchangeRequest{},
	&OfferedBook{},
	&ExchangeMatch{},
	&Genre{},
    &Location{},
}
