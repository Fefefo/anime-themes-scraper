# Golang scraper for the [/r/AnimeThemes](https://www.reddit.com/r/AnimeThemes/)

A scraping tool for Gophers which scrapes [/r/AnimeThemes](https://www.reddit.com/r/AnimeThemes/), a simple and consistent repository of anime opening and ending themes.

It uses [Colly](https://github.com/gocolly/colly).


## Installation

`go get github.com/Fefefo/anime-themes-scraper`


## Project using this repo

-  [fefegobot](https://github.com/Fefefo/goBot) Telegram bot where you can inline search any opening by english or japanese title


## Structs

```go
// Anime is the struct which cointains the information about an anime
type Anime struct {
	IDAnime int
	Year    int
	NameJap string
	NameEng string
	MalLink string
	//ImageURL string
	Songs []animeSong
}

// animeSong is the stroct which contains the information about a theme
type animeSong struct {
	Title    string
	Link     string
	Version  string
	Episodes string
	Notes    string
}

// List is the struct wihich contains all the animes
type List []Anime
```


## Example

```go
func main() {
  //To put into list all the animes and their songs
  list := scraper.GetAnimeList()
  
  //To select only the animes which contains "Assassination Classroom"
  list = list.SelectByBothNames("Assassination Classroom")

  //To output the link of the first theme of the first anime selected
  fmt.Println(list[0].Songs[0].Link)
}
  /*
    OUTPUT

    https://animethemes.moe/video/AnsatsuKyoushitsu-OP1.webm
    
  */
```
