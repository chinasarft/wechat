package flags

import (
	"fmt"
	"os"
	"strconv"

	goflags "github.com/jessevdk/go-flags"
)

type commanLineArgs struct {
	Test         bool   `short:"t" long:"test" description:""`
	Verbose      []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Host         string `short:"h" long:"host" description:"host to bind"`
	Port         int    `short:"p" long:"port" description:"port to listen" default:"80"`
	CertFilePath string `short:"c" long:"cert" description:"Specify Cert File Path for HTTPS connections"`
	KeyFilePath  string `short:"k" long:"key" description:"Specify Key File Path for HTTPS connections"`
}

var Flags *commanLineArgs

func init() {
	var err error
	Flags, err = parseFlags()
	if err != nil {
		panic(err)
	}
	err = check()
	if err != nil {
		panic(err)
	}
}

func parseFlags() (*commanLineArgs, error) {
	flags := commanLineArgs{}
	parser := goflags.NewParser(&flags, goflags.Default|goflags.IgnoreUnknown)
	args, err := parser.ParseArgs(os.Args)
	if err != nil {
		return nil, err
	}
	os.Args = args
	return &flags, nil
}

func check() error {

	if err := checkPort(Flags.Port); err != nil {
		return err
	}
	if Flags.CertFilePath != "" {
		if isFile, err := checkFile(Flags.CertFilePath); err != nil {
			return fmt.Errorf("Failed to cert file: %s", err)
		} else if isFile == false {
			return fmt.Errorf("Failed to cert file, %s is not file", Flags.CertFilePath)
		}
	}
	if Flags.KeyFilePath != "" {
		if isFile, err := checkFile(Flags.KeyFilePath); err != nil {
			return fmt.Errorf("Failed to key file: %s", err)
		} else if isFile == false {
			return fmt.Errorf("Failed to key file, %s is not file", Flags.KeyFilePath)
		}
	}
	return nil
}

func checkPort(port int) error {
	if port <= 0 || port >= 65536 {
		return fmt.Errorf("Invalid port: %d", port)
	}
	return nil
}

func checkDir(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("Cannot open %s: %s", path, err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return false, fmt.Errorf("Cannot read stats of %s: %s", path, err)
	}

	return fileinfo.Mode().IsDir(), nil
}

func checkFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("Cannot open %s: %s", path, err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return false, fmt.Errorf("Cannot read stats of %s: %s", path, err)
	}

	return fileinfo.Mode().IsRegular(), nil
}

func HostAndPort() string {
	return Flags.Host + ":" + strconv.Itoa(Flags.Port)
}

func DoHTTPs() (isHTTPs bool, certFilePath string, keyFilePath string) {
	isHTTPs = Flags.CertFilePath != "" && Flags.KeyFilePath != ""
	certFilePath = Flags.CertFilePath
	keyFilePath = Flags.KeyFilePath
	return
}

func IsTest() bool {
	return Flags.Test
}
