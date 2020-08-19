package scraper

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

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

type animeSong struct {
	Title    string
	Link     string
	Version  string
	Episodes string
	Notes    string
}

// List is the struct wihich contains all the animes
type List []Anime

//GetAnimeList returns an array of the Anime struct with all anime loaded
func GetAnimeList() List {
	c := colly.NewCollector()
	d := colly.NewCollector()

	count := 0

	var wg sync.WaitGroup
	var animeList List

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
		year := []rune(e.DOM.Find("h2").Eq(1).Text())
		element := e.DOM.Find("h3")
		for i := 0; i < element.Length(); i++ {
			el := element.Eq(i)
			var temp Anime
			count++

			temp.IDAnime = count
			temp.NameJap = el.Text()
			temp.MalLink, _ = el.Attr("href")
			temp.Year, _ = strconv.Atoi(string(year[0:4]))
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
					songs.Title = getTitle(title.Text()) //len(title.Text())-1)
					tempNameSong = songs.Title
				} else {
					songs.Title = tempNameSong
				}
				if row.Find("a").Length() != 0 {
					songs.Version = strings.Replace(title.Text(), `"`+songs.Title+`"`, "", -1)
					link := title.Next()
					songs.Link, _ = link.Children().Attr("href")
					eps := link.Next()
					songs.Episodes = eps.Text()
					notes := eps.Next()
					songs.Notes = notes.Text()

					temp.Songs = append(temp.Songs, songs)
				}
			}
			animeList = append(animeList, temp)
		}

		wg.Done()
		fmt.Println(count)
	})

	c.Visit("https://reddit.com/r/AnimeThemes/wiki/anime_index")

	wg.Wait()
	return animeList
}

func getTitle(a string) string {
	d := []rune(a)
	return string(d[strings.Index(a, `"`)+1 : len(d)-1])
}

// SelectByJapName will find all the anime which contains the string in the japanese name
func (list List) SelectByJapName(name string) List {
	var newList List
	for i := range list {
		if strings.Contains(strings.ToLower(list[i].NameJap), strings.ToLower(name)) {
			newList = append(newList, list[i])
		}
	}
	return newList
}

// SelectByEngName will find all the anime which contains the string in the english name
func (list List) SelectByEngName(name string) List {
	var newList List
	for i := range list {
		if strings.Contains(strings.ToLower(list[i].NameEng), strings.ToLower(name)) {
			newList = append(newList, list[i])
		}
	}
	return newList
}

// SelectByBothNames will find all the anime which contains the string in the japanese or english name
func (list List) SelectByBothNames(name string) List {
	var newList List
	for i := range list {
		if strings.Contains(strings.ToLower(list[i].NameEng), strings.ToLower(name)) {
			newList = append(newList, list[i])
		} else if strings.Contains(strings.ToLower(list[i].NameJap), strings.ToLower(name)) {
			newList = append(newList, list[i])
		}
	}
	return newList
}
