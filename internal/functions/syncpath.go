package functions

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/johnnylin-a/cert-sync/internal/apis"
	"github.com/johnnylin-a/cert-sync/internal/configs"
	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var syncQueue = make(chan string)

func SyncPath(path string) {
	syncQueue <- path
}

func init() {
	go func() {
		f, err := os.Open(configs.GetAppConfig().DotSSHPath + "/config")
		if err != nil {
			log.Println("failed to read " + configs.GetAppConfig().DotSSHPath + "/config")
			panic(err)
		}
		sshCfg, err := ssh_config.Decode(f)
		if err != nil {
			log.Println("failed to decode " + configs.GetAppConfig().DotSSHPath + "/config")
			panic(err)
		}
		sshKnownHosts, err := knownhosts.New(configs.GetAppConfig().DotSSHPath + "/known_hosts")
		if err != nil {
			apis.LogAndSendNotification("WARNING: failed to read " + configs.GetAppConfig().DotSSHPath + "/known_hosts and will allow all server fingerprints!")
			sshKnownHosts = ssh.InsecureIgnoreHostKey()
		}
		for {
			updatedPath := <-syncQueue
			apis.LogAndSendNotification("Syncing " + updatedPath)
			appConfig := configs.GetAppConfig()
			configHosts := appConfig.SyncFilePaths[updatedPath]
			for _, configHost := range configHosts {
				alias := configHost.ConfigHost
				dst := configHost.Dst
				port, _ := sshCfg.Get(alias, "Port")
				hostname, _ := sshCfg.Get(alias, "HostName")
				user, _ := sshCfg.Get(alias, "User")
				privateKey, _ := sshCfg.Get(alias, "IdentityFile")
				if strings.HasPrefix(privateKey, "~/.ssh/") {
					privateKey = strings.ReplaceAll(privateKey, "~/.ssh", appConfig.DotSSHPath)
				}
				if privateKey == "" && configHost.PrivateKey != nil {
					privateKey = *configHost.PrivateKey
				}

				if port == "" {
					port = "22"
				}

				if user == "" {
					user = "root"
				}
				log.Println("Syncing "+updatedPath+" to "+configHost.ConfigHost+":"+dst, alias, port, hostname, user, privateKey)
				scpClientConfig, _ := auth.PrivateKey(user, privateKey, sshKnownHosts)
				scpClientConfig.HostKeyCallback = sshKnownHosts

				privateKeyValue, err := os.ReadFile(privateKey)
				if err != nil {
					apis.LogAndSendNotification("Couldn't read private key file " + privateKey)
					log.Println(err)
					continue
				}
				signer, err := ssh.ParsePrivateKey(privateKeyValue)
				if err != nil {
					apis.LogAndSendNotification("Couldn't parse private key content " + privateKey)
					log.Println(err)
					continue
				}
				scpClientConfig.Auth = []ssh.AuthMethod{
					ssh.PublicKeys(signer),
				}

				scpClient := scp.NewClient(hostname+":"+port, &scpClientConfig)
				err = scpClient.Connect()
				if err != nil {
					apis.LogAndSendNotification("Couldn't establish a connection to the remote server " + configHost.ConfigHost)
					log.Println(err)
					continue
				}

				// copy file over
				f, err := os.Open(updatedPath)
				if err != nil {
					apis.LogAndSendNotification("Couldn't open local file " + updatedPath)
					scpClient.Close()
					continue
				}
				err = scpClient.CopyFile(context.Background(), f, dst, "0644")
				if err != nil {
					apis.LogAndSendNotification("Couldn't copy file " + updatedPath + " to " + configHost.ConfigHost + ":" + dst)
					scpClient.Close()
					f.Close()
					continue
				}
				apis.LogAndSendNotification("Copied " + updatedPath + " to " + configHost.ConfigHost + ":" + dst)
				scpClient.Close()
				f.Close()
			}
		}
	}()
}
