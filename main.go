package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {
	config, err := VMessToSingBox(`vmess://ewogICJ2IjogIjIiLAogICJwcyI6ICJTZXJ2ZXIgQiIsCiAgImFkZCI6ICIxMDQuMjEuMzAuMjI0IiwKICAicG9ydCI6ICI0NDMiLAogICJpZCI6ICJhNmY4YzNhMS02OWE0LTRjN2UtOGFkNi0xYjdhMmQ3ZjliNGMiLAogICJhaWQiOiAiMCIsCiAgInNjeSI6ICJhdXRvIiwKICAibmV0IjogIndzIiwKICAidHlwZSI6ICJub25lIiwKICAiaG9zdCI6ICJzZXJ2ZXItYi50YWJhdGVsZWNvbS5kZXYiLAogICJwYXRoIjogIi8iLAogICJ0bHMiOiAidGxzIiwKICAic25pIjogIiIsCiAgImFscG4iOiAiIiwKICAiZnAiOiAiIgp9`)

	if err != nil {
		log.Println(err)
		return
	}

	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := os.WriteFile("config.json", data, 0655); err != nil {
		log.Fatal(err)
	}
}
