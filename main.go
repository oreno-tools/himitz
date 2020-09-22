package main

import (
	"bufio"
	kms "cloud.google.com/go/kms/apiv1"
	_ "context"
	"encoding/base64"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	_ "golang.org/x/oauth2/google"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"os"
)

const (
	AppVersion = "0.0.1"
)

var (
	argProject = flag.String("project", "", "Specify a Project ID.")
	argZone    = flag.String("zone", "us-west1-b", "Specify a Zone Name.")
	argKeyRing = flag.String("ring", "", "Specify a Key Ring Name.")
	argKey     = flag.String("key", "", "Specify a Key Name.")
	argEnc     = flag.Bool("encrypt", false, "Encrypt data.")
	argDec     = flag.Bool("decrypt", false, "Decrypt data. ")
	argBatch   = flag.Bool("batch", false, "Enable batch mode.")
	argVersion = flag.Bool("version", false, "Print Version.")
)

// Thank you
// referenced https://qiita.com/ktoshi/items/1fd4f808c955d33c3d28
func decryption(ciphertext []byte, keyName string) (string, error) {
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return "", err
	}
	request := &kmspb.DecryptRequest{
		Name:       keyName,
		Ciphertext: ciphertext,
	}
	response, err := client.Decrypt(ctx, request)
	return string(response.GetPlaintext()), err
}

// Thank you
// referenced https://qiita.com/ktoshi/items/1fd4f808c955d33c3d28
func encryption(str string, keyName string) ([]byte, error) {
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	request := &kmspb.EncryptRequest{
		Name:      keyName,
		Plaintext: []byte(str),
	}

	response, err := client.Encrypt(ctx, request)
	return response.GetCiphertext(), err
}

func genKeyName(project string, ring string, key string) string {
	keyName := fmt.Sprintf("projects/%s/locations/global/keyRings/%s/cryptoKeys/%s", project, ring, key)
	return keyName
}

func main() {
	flag.Parse()
	if *argVersion {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	// read value from stdin
	stat, _ := os.Stdin.Stat()
	var stdInValue string
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		stdin := bufio.NewScanner(os.Stdin)
		stdin.Scan()
		stdInValue = stdin.Text()
		if err := stdin.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Unable to read value from stdin: ", err)
			os.Exit(1)
		}
	}

	var project, ring, key string
	if *argProject == "" && os.Getenv("GCP_PROJECT") == "" {
		fmt.Println("Please set `project` option or environment variable `GCP_PROJECT`")
		os.Exit(1)
	} else {
		if os.Getenv("GCP_PROJECT") != "" {
			project = os.Getenv("GCP_PROJECT")
		} else {
			project = *argProject
		}
	}

	if *argKeyRing == "" {
		fmt.Println("Please set `ring` option")
		os.Exit(1)
	} else {
		ring = *argKeyRing
	}

	if *argKey == "" {
		fmt.Println("Please set `key` option")
		os.Exit(1)
	} else {
		key = *argKey
	}

	keyName := genKeyName(project, ring, key)

	var str string
	if *argEnc {
		encryptedData, err := encryption(stdInValue, keyName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data := base64.StdEncoding.EncodeToString(encryptedData)
		if *argBatch {
			str = data
		} else {
			str = fmt.Sprintf("Encrypted data with base64 encoded: %s", data)
		}
		fmt.Println(str)
	} else if *argDec {
		decryptingData, _ := base64.StdEncoding.DecodeString(stdInValue)
		data, _ := decryption(decryptingData, keyName)
		if *argBatch {
			str = data
		} else {
			str = fmt.Sprintf("Decrypted data: %s", data)
		}
		fmt.Println(str)
	}
}
