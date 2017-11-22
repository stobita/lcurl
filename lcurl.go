package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli"
)

type paramList struct {
	key   string
	value string
}

func main() {
	app := cli.NewApp()
	app.Name = "Murl"
	app.Usage = "This is Simple cURL"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "parameter,p",
			Value: "",
			Usage: "url parameter",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "get request",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "query,q",
					Value: "",
					Usage: "query parameter",
				},
			},
			Action: func(c *cli.Context) {
				// prepare parameter
				query := c.String("q")

				println(query)
				if query != "" {
					var err error
					query, err = storeQuery(query)
					if err != nil {
						cli.ShowAppHelp(c)
					}
				}
				reqUrl := c.Args().First() + query

				req, err := http.NewRequest("GET", reqUrl, nil)
				if err != nil {
					fmt.Println(err)
					return
				}
				client := new(http.Client)
				res, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer res.Body.Close()

				resBody, err := ioutil.ReadAll(res.Body)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(string(resBody))
			},
		},
		{
			Name:    "post",
			Aliases: []string{"p"},
			Usage:   "post request",
			Action: func(c *cli.Context) {
				fmt.Println("this is post request")
			},
		},
	}
	app.Run(os.Args)
}

func storeQuery(str string) (string, error) {
	list := strings.Split(str, ",")
	result := "&"
	for i := range list {
		if strings.Index(list[i], ":") > 0 {
			list[i] = strings.Replace(list[i], ":", "=", -1)
		} else {
			if strings.Index(list[i], "=") == -1 {
				return "", errors.New("use : or =")
			}
		}
		result += list[i]
	}
	return result, nil
}
