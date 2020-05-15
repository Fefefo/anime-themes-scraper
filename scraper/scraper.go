package scraper

import (
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

// Anime is the struct which cointains the information about an anime
type Anime struct {
	IDAnime int
	NameJap string
	NameEng string
	MalLink string
	//ImageURL string
	Songs []animeSong
}

type animeSong struct {
	Title    string
	Link     string
	Version  string
	Episodes string
	Notes    string
}

func GetAnimeList() []Anime {
	c := colly.NewCollector()
	d := colly.NewCollector()

	count := 0

	var wg sync.WaitGroup
	var animeList []Anime

	c.OnHTML("div.md.wiki", func(e *colly.HTMLElement) {

		for i := 1; i < e.DOM.Find("p").Length(); i++ {
			link, _ := e.DOM.Find("p").Eq(i).Children().Eq(0).Attr("href")
			linkAbs := "https://reddit.com" + strings.TrimSuffix(strings.Split(link, "#")[0], "/")

			if vis, _ := d.HasVisited(linkAbs); !vis {
				wg.Add(1)
				go d.Visit(linkAbs)
			}
		}
	})

	d.OnHTML("div.md.wiki", func(e *colly.HTMLElement) {
		element := e.DOM.Find("h3")
		for i := 0; i < element.Length(); i++ {
			el := element.Eq(i)
			var temp Anime
			count++

			temp.IDAnime = count
			temp.NameJap = el.Text()
			temp.MalLink, _ = el.Attr("href")

			el = el.Next()
			if el.Is("p") {
				temp.NameEng = el.Text()
				el = el.Next()
			}
			el = el.Find("tbody")
			tempNameSong := ""
			for j := 0; j < el.Find("tr").Length(); j++ {
				var songs animeSong

				row := el.Find("tr").Eq(j).Children()
				title := row.Eq(0)
				//fmt.Println(title.Text(), "is a song of", temp.NameEng)
				if title.Text() != "" {
					songs.Title = substr(title.Text(), strings.Index(title.Text(), `"`)+1) //len(title.Text())-1)
					tempNameSong = songs.Title
				} else {
					songs.Title = tempNameSong
				}
				songs.Version = strings.Replace(title.Text(), `"`+songs.Title+`"`, "", -1)
				link := title.Next()
				songs.Link, _ = link.Children().Attr("href")
				eps := link.Next()
				songs.Episodes = eps.Text()
				notes := eps.Next()
				songs.Notes = notes.Text()

				temp.Songs = append(temp.Songs, songs)
			}
			animeList = append(animeList, temp)
		}

		wg.Done()
		//fmt.Println(count)
	})

	c.Visit("https://reddit.com/r/AnimeThemes/wiki/anime_index")

	wg.Wait()
	return animeList
}

func substr(a string, b int) string {
	d := []rune(a)
	// fmt.Printf("%d  %d \n", len(d), len(a))
	// fmt.Println(string(d), string(a))
	return string(d[b : len(d)-1])
}
