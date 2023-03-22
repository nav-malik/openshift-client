package main

import (
	"bufio"
	"flag"
	"fmt"
	"openshift-client/encryption"
	"openshift-client/fileManagement"
	"openshift-client/namespaces"
	"openshift-client/pods"
	"openshift-client/token"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

var accessToken string
var serverUrl string = ""
var errToken error
var username *string
var key = []byte("xyz@#-1230defgh9")
var fileNameWithExtension = "userConfig.txt"

func main() {

	podsCmd := flag.NewFlagSet("pods", flag.ExitOnError)
	namespace := podsCmd.String("ns", "", "the project/namespace name")

	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	username = loginCmd.String("u", "", "username")
	serverUrl2 := loginCmd.String("s", serverUrl, "Server URL")
	serverUrl = *serverUrl2

	// Parse command line arguments
	flag.Parse()

	// Check which subcommand was called
	if len(flag.Args()) < 1 {
		fmt.Println("Please specify a subcommand.")
		return
	}
	switch os.Args[1] {
	case "login":
		// Parse subcommand flags
		loginCmd.Parse(os.Args[2:])
		getAccessToken()
		fmt.Println()
		getNamespaces()
	case "pods":
		podsCmd.Parse(os.Args[2:])
		getPods(*namespace)
	case "namespaces":
		getNamespaces()
	default:
		fmt.Printf("Unknown subcommand: %s\n", os.Args[1])
		return
	}
}

func getAccessToken() {
	if len(*username) == 0 || *username == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("username: ")
		*username, _ = reader.ReadString('\n')
	}

	fmt.Print("password: ")
	password, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	if len(serverUrl) == 0 || serverUrl == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("server: ")
		serverUrl, _ = reader.ReadString('\n')
		serverUrl = strings.TrimRight(serverUrl, "\n\r")
	}

	accessToken, errToken = token.GetToken(*username, password, serverUrl)
	if errToken != nil {
		panic(errToken)
	}

	ciphertext, err := encryption.Encrypt_Base64(key, accessToken)
	if err != nil {
		panic(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err1 := fileManagement.WriteLine(ciphertext, dir+"//"+fileNameWithExtension)
	if err1 != nil {
		panic(err1)
	}
}

func getAccessTokenFromFile() {

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	encryptedToken, err := fileManagement.ReadLine(dir + "//" + fileNameWithExtension)
	if err != nil {
		panic(err)
	}
	decryptedToken, err := encryption.Decrypt_Base64(key, encryptedToken)
	if err != nil {
		panic(err)
	}
	accessToken = decryptedToken
}

func getNamespaces() {
	if len(accessToken) == 0 || accessToken == "" {
		getAccessTokenFromFile()
	}
	namespaces, err := namespaces.GetNamespaces(accessToken, serverUrl)

	if err != nil {
		panic(err)
	}
	for _, namespace := range namespaces.Items {
		fmt.Println(namespace.Name)
	}
}

func getPods(namespace string) {
	if len(accessToken) == 0 || accessToken == "" {
		getAccessTokenFromFile()
	}
	pods, err := pods.GetPods(accessToken, serverUrl, namespace)
	if err != nil {
		panic(err)
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
}
