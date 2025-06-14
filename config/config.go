package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SrcURL string
	DstURL string
}

func Load() *Config {
	_ = godotenv.Load()

	srcDef := getenv("SRC_URL", "https://jsonplaceholder.typicode.com/users")
	dstDef := getenv("DST_URL", "")

	src := flag.String("src", srcDef, "source URL")
	dst := flag.String("dst", dstDef, "destination URL (webhook)")
	flag.Parse()

	return &Config{
		SrcURL: *src,
		DstURL: *dst,
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
