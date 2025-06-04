package link

import (
	"fmt"
	"os"
	"strings"

	"github.com/52funny/pikpakcli/conf"
	"github.com/52funny/pikpakcli/internal/pikpak"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var LinkCommand = &cobra.Command{
	Use:     "link",
	Aliases: []string{"l"},
	Short:   `Get download links for files on the pikpak server`,
	Run: func(cmd *cobra.Command, args []string) {
		p := pikpak.NewPikPak(conf.Config.Username, conf.Config.Password)
		err := p.Login()
		if err != nil {
			logrus.Errorln("Login Failed:", err)
		}
		// Output file handle
		var f = os.Stdout
		if strings.TrimSpace(output) != "" {
			file, err := os.Create(output)
			if err != nil {
				logrus.Errorln("Create file failed:", err)
				return
			}
			defer file.Close()
			f = file
		}

		if len(args) > 0 {
			getFileLinks(&p, args, f)
		} else {
			getFolderLinks(&p, f)
		}
	},
}

// Specifies the folder of the pikpak server
// default is the root folder
var folder string

// Specifies the file to write
// default is the stdout
var output string

var parentId string

func init() {
	LinkCommand.Flags().StringVarP(&folder, "path", "p", "/", "specific the folder of the pikpak server")
	LinkCommand.Flags().StringVarP(&output, "output", "o", "", "specific the file to write")
	LinkCommand.Flags().StringVarP(&parentId, "parent-id", "P", "", "parent folder id")
}

// Get links for all files in a folder
func getFolderLinks(p *pikpak.PikPak, f *os.File) {
	var err error
	if parentId == "" {
		parentId, err = p.GetDeepFolderId("", folder)
		if err != nil {
			logrus.Errorln("Get parent id failed:", err)
			return
		}
	}
	fileStat, err := p.GetFolderFileStatList(parentId)
	if err != nil {
		logrus.Errorln("Get folder file stat list failed:", err)
		return
	}
	for _, stat := range fileStat {
		if stat.Kind == "drive#file" {
			file, err := p.GetFile(stat.ID)
			if err != nil {
				logrus.Errorln("Get file failed:", err)
				continue
			}
			printFileLinks(f, &file)
		}
	}
}

// Get links for specific files
func getFileLinks(p *pikpak.PikPak, args []string, f *os.File) {
	var err error
	if parentId == "" {
		parentId, err = p.GetPathFolderId(folder)
		if err != nil {
			logrus.Errorln("get parent id failed:", err)
			return
		}
	}
	for _, path := range args {
		stat, err := p.GetFileStat(parentId, path)
		if err != nil {
			logrus.Errorln(path, "get file stat error:", err)
			continue
		}
		file, err := p.GetFile(stat.ID)
		if err != nil {
			logrus.Errorln("Get file failed:", err)
			continue
		}
		printFileLinks(f, &file)
	}
}

// Print all available links for a file
func printFileLinks(f *os.File, file *pikpak.File) {
	fmt.Fprintf(f, "File: %s\n", file.Name)

	// Print direct download link
	if file.Links.ApplicationOctetStream.URL != "" {
		fmt.Fprintf(f, "Direct Download Link: %s\n", file.Links.ApplicationOctetStream.URL)
		fmt.Fprintf(f, "Expires at: %s\n", file.Links.ApplicationOctetStream.Expire)
	}

	// Print web content link
	if file.WebContentLink != "" {
		fmt.Fprintf(f, "Web Content Link: %s\n", file.WebContentLink)
	}

	// Print media links
	if len(file.Medias) > 0 {
		fmt.Fprintf(f, "Media Links:\n")
		for i, media := range file.Medias {
			fmt.Fprintf(f, "  %d. %s\n", i+1, media.Link.URL)
			fmt.Fprintf(f, "     Expires at: %s\n", media.Link.Expire)
		}
	}
	fmt.Fprintf(f, "\n")
}
