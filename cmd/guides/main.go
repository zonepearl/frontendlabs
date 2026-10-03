// guides: the Frontend Labs field-guide site generator.
//
//	guides sync                 copy the latest guides from the wiki into docs/
//	guides build                render docs/ into dist/ (a static site)
//	guides serve [-addr :8080]  build, serve, and regenerate on every change
//	guides watch                build and regenerate on change, without serving
package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"frontendlabs.xyz/guides/internal/site"
	"frontendlabs.xyz/guides/web"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	cmd := os.Args[1]
	fl := flag.NewFlagSet(cmd, flag.ExitOnError)
	manifest := fl.String("manifest", "guides.json", "site manifest")
	docs := fl.String("docs", "docs", "documents directory (the source of the site)")
	out := fl.String("out", "dist", "output directory (serve defaults to .preview)")
	addr := fl.String("addr", ":8080", "listen address (serve)")
	syncWiki := fl.Bool("sync", true, "serve/watch: also watch the original guides in the wiki and copy changes into docs/")
	dev := fl.Bool("dev", false, "serve/watch: read templates and assets from ./web on disk and rebuild when they change")
	interval := fl.Duration("interval", 700*time.Millisecond, "how often to check for changes")
	fl.Parse(os.Args[2:])
	outSet := false
	fl.Visit(func(f *flag.Flag) { outSet = outSet || f.Name == "out" })
	if cmd == "serve" && !outSet {
		*out = ".preview" // the live-reload preview never lands in the deployable dist/
	}

	m, err := site.LoadManifest(*manifest)
	check(err)

	var assets fs.FS = web.FS
	webDir := ""
	if *dev {
		webDir = "web"
		assets = os.DirFS(webDir)
	}
	b := &site.Builder{M: m, DocsDir: *docs, OutDir: *out, Assets: assets}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	switch cmd {
	case "sync":
		changed, err := site.Sync(m, *docs)
		check(err)
		for _, c := range changed {
			fmt.Println("synced", c)
		}
		fmt.Printf("%d of %d files updated in %s\n", len(changed), len(m.Files()), *docs)
	case "build":
		check(b.BuildAll())
	case "serve", "watch":
		if *syncWiki {
			changed, err := site.Sync(m, *docs)
			check(err)
			if len(changed) > 0 {
				fmt.Printf("synced %d changed files from the wiki\n", len(changed))
			}
		}
		b.Live = cmd == "serve"
		check(b.BuildAll())
		w := &site.Watcher{B: b, SyncWiki: *syncWiki, WebDir: webDir, Interval: *interval}
		if cmd == "watch" {
			fmt.Printf("watching %s for changes (Ctrl-C to stop)\n", filepath.Clean(*docs))
			w.Run(ctx)
			return
		}
		check(site.Serve(ctx, w, *addr))
	default:
		usage()
	}
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: guides sync|build|serve|watch [flags]   (guides <cmd> -h for flags)")
	os.Exit(2)
}
