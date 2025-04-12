/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	_ "github.com/rhysmah/CLI-Note-App/cmd/delete"
	_ "github.com/rhysmah/CLI-Note-App/cmd/edit"
	_ "github.com/rhysmah/CLI-Note-App/cmd/list"
	_ "github.com/rhysmah/CLI-Note-App/cmd/new"
	"github.com/rhysmah/CLI-Note-App/cmd/root"
	_ "github.com/rhysmah/CLI-Note-App/cmd/version"
)

func main() {
	root.Execute()
}
