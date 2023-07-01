package consts

type FeedType string

const (
	Recommend FeedType = "homefeed_recommend"
	Fasion    FeedType = "homefeed.fashion_v3"
	Food      FeedType = "homefeed.food_v3"
	Cosmetics FeedType = "homefeed.cosmetics_v3"
	Movie     FeedType = "homefeed.movie_and_tv_v3"
	Career    FeedType = "homefeed.career_v3"
	Emotion   FeedType = "homefeed.love_v3"
	Household FeedType = "homefeed.household_product_v3"
	Game      FeedType = "homefeed.gaming_v3"
	Travel    FeedType = "homefeed.travel_v3"
	Fitness   FeedType = "homefeed.fitness_v3"
)

type NoteType string

const (
	Normal NoteType = "normal"
	Video  NoteType = "video"
)

type SearchSortType string

const (
	General     SearchSortType = "general"
	MostPopular SearchSortType = "popularity_descending"
	Latest      SearchSortType = "time_descending"
)

type SearchNoteType string

const (
	All    SearchNoteType = "0"
	Videos SearchNoteType = "1"
	Image  SearchNoteType = "2"
)
