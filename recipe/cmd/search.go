// Copyright © 2016 Marcel Hauf <mail@marcelhauf.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"github.com/gophergala2016/recipe/pkg/recipe"
	"github.com/gophergala2016/recipe/pkg/repo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var (
	searchTitle       = true
	searchURL         = false
	searchDescription = false
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for recipes matching the pattern.",
	Long:  `Search the repositories for recipes matching the pattern.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbpath := ""
		cacheMap := viper.Get("cache")
		if cache, ok := cacheMap.(map[interface{}]interface{}); ok {
			dbpath = cache["dbpath"].(string)
		}
		term := strings.Join(args, "")
		options := repo.SearchOptions{
			Title:       searchTitle,
			Description: searchDescription,
			URL:         searchURL,
		}
		recipeLinks, err := recipe.Search(dbpath, term, options)
		if err != nil {
			log.Println(err)
		}

		for _, rl := range recipeLinks {
			fmt.Printf("%s\t%s\n", rl.Title(), rl.URL())
		}
	},
}

func init() {
	searchCmd.Flags().BoolVarP(&searchTitle, "searchtitle", "t", true, "if the title of the recipe should be searched")
	searchCmd.Flags().BoolVarP(&searchURL, "searchurl", "u", true, "if the URL of the recipe should be searched")
	searchCmd.Flags().BoolVarP(&searchTitle, "searchdescription", "d", true, "if the description of the recipe should be searched")
	RootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
