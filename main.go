package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/go-ini/ini"
)

var (
	list        bool
	profileName = "default"
	filePath    string
)

// errors
var (
	ErrMissingProfile         = errors.New("profile not found")
	ErrMissingAccessKeyID     = errors.New("profile does not contain aws_access_key_id")
	ErrMissingSecretAccessKey = errors.New("profile does not contain aws_secret_access_key")
)

const (
	usage               = `don <profile-name>`
	accessKeyIDName     = "aws_access_key_id"
	secretAccessKeyName = "aws_secret_access_key"
	outputFormat        = `export AWS_ACCESS_KEY_ID="%s";
export AWS_SECRET_ACCESS_KEY="%s";`
)

func getCredentials(section *ini.Section) (id, secret string, err error) {
	k := section.Key(accessKeyIDName)
	if k == nil {
		err = ErrMissingAccessKeyID
		return
	}
	id = k.Value()

	k = section.Key(secretAccessKeyName)
	if k == nil {
		err = ErrMissingSecretAccessKey
		return
	}
	secret = k.Value()
	return
}

func init() {
	flag.BoolVar(&list, "l", false, "list available aws profiles")
}

func getCredentialsFilePath() (filePath string, err error) {
	currentUser, err := user.Current()
	if err != nil {
		return
	}
	homedir := currentUser.HomeDir

	filePath = filepath.Join(homedir, ".aws/credentials")
	return
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 1 {
		profileName = args[0]
	} else if len(args) > 1 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	var err error
	filePath, err = getCredentialsFilePath()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fh, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	credentialsFile, err := ini.Load(fh)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if list {
		for _, section := range credentialsFile.SectionStrings() {
			// DEFAULT is implicit section without a section header
			if section != "" && section != "DEFAULT" {
				fmt.Fprintln(os.Stdout, section)
			}
		}
		return
	}

	section, err := credentialsFile.GetSection(profileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, ErrMissingProfile.Error())
		os.Exit(1)
	}

	id, secret, err := getCredentials(section)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, outputFormat, id, secret)
}
