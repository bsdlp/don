package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/go-ini/ini"
)

var (
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
	outputFormat        = `# aws credentials environment variables for profile '%s'
# retrieved from file '%s'
export AWS_ACCESS_KEY_ID="%s"
export AWS_SECRET_ACCESS_KEY="%s"`
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

func main() {
	if len(os.Args) == 2 {
		profileName = os.Args[1]
	} else if len(os.Args) > 2 {
		io.WriteString(os.Stderr, usage)
		os.Exit(1)
	}

	currentUser, err := user.Current()
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}
	homedir := currentUser.HomeDir

	filePath = filepath.Join(homedir, ".aws/credentials")

	fh, err := os.Open(filePath)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}

	credentialsFile, err := ini.Load(fh)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}

	section, err := credentialsFile.GetSection(profileName)
	if err != nil {
		io.WriteString(os.Stderr, ErrMissingProfile.Error())
		os.Exit(1)
	}

	id, secret, err := getCredentials(section)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, outputFormat, profileName, filePath, id, secret)
}
