package cli

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/herryg91/localhost-proxy/entity"
	"github.com/herryg91/localhost-proxy/pkg/helpers"
	"github.com/spf13/cobra"
)

type CmdRun struct {
	*cobra.Command
}

func NewCmdStatus() *CmdRun {
	c := &CmdRun{}
	c.Command = &cobra.Command{
		Use:   "run",
		Short: "Run Localhost Proxy",
		Long:  "Run Localhost Proxy",
	}
	c.RunE = c.runCommand
	return c
}

func (c *CmdRun) runCommand(cmd *cobra.Command, args []string) error {
	cfg := entity.Config{}.FromFile()

	log.Println("========== Configuration ==========")
	log.Println("Port: ", cfg.Port)
	log.Println("Routes: ")
	if len(cfg.Routes) == 0 {
		log.Println("    No routes")
	}
	for _, route_cfg := range cfg.Routes {
		log.Println(fmt.Sprintf("  %s: %s => %s", route_cfg.Name, route_cfg.Pathname, route_cfg.Destination))
	}
	log.Println("===================================")

	director := func(req *http.Request) {
		for _, route_cfg := range cfg.Routes {
			if strings.HasPrefix(req.URL.Path, route_cfg.Pathname) && strings.HasSuffix(req.URL.Path, route_cfg.Pathname) {
				req.URL.Path = "/"
				req.URL.Host = route_cfg.Destination
				break
			} else if strings.HasPrefix(req.URL.Path, route_cfg.Pathname+"/") {
				req.URL.Path = "/" + helpers.StripUrlPrefix(req.URL.Path, route_cfg.Pathname+"/")
				req.URL.Host = route_cfg.Destination
				break
			}
		}
		// if strings.HasPrefix(req.URL.Path, "/usertoken") {
		// 	req.URL.Path =  "/"s+trip_url_prefix(req.URL.Path, "/usertoken")
		// 	req.URL.Host = "localhost:28001"
		// }

		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", req.URL.Host)
		req.URL.Scheme = "http"

	}

	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
	return nil
}
