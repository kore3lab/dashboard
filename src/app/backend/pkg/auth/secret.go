package auth

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

/**
<SecretProvider>
	StaticUserSecretProvider			: static username, password
	UserFileSecretProvider				: kubernetes basic-auth secret (username file + password file)
	StaticTokenSecretProvider			: static token
	ServiceAccountTokenSecretProvider	: kubernetes service-account-token secret
*/
type UserSecret struct {
	Reload   func() error
	Username string
	Password string
	mu       sync.RWMutex
}

func StaticUserSecretProvider(username string, password string) SecretProvider {
	h := &UserSecret{Username: username, Password: password}
	h.Reload = func() error { return nil }
	return func(username, realm string) string {
		h.mu.RLock()
		defer h.mu.RUnlock()
		exists := (h.Username == username)
		password := h.Password
		if !exists {
			return ""
		}
		return password
	}
}

type UserSecretFile struct {
	UserSecret
	PathDir      string
	PathUsername string
	PathPassword string
	InfoUsername os.FileInfo
	InfoPassword os.FileInfo
}

func (f *UserSecretFile) ReloadIfNeeded() error {

	if _, err := os.Stat(f.PathDir); os.IsNotExist(err) {
		return err
	}
	f.PathUsername = filepath.Join(f.PathDir, "username")
	f.PathPassword = filepath.Join(f.PathDir, "password")

	info, err := os.Stat(f.PathUsername)
	if err != nil {
		return err
	}

	f.mu.Lock()
	reload := false
	if f.InfoUsername == nil || f.InfoUsername.ModTime() != info.ModTime() {
		f.InfoUsername = info
		reload = true
	}

	if !reload {
		info, err = os.Stat(f.PathPassword)
		if err != nil {
			return err
		}
		if f.InfoPassword == nil || f.InfoPassword.ModTime() != info.ModTime() {
			f.InfoPassword = info
			reload = true
		}
	}

	f.mu.Unlock()

	if reload {
		if err = f.Reload(); err != nil {
			return err
		}
	}
	return nil
}

func reloadUserFileSecret(h *UserSecretFile) error {

	h.mu.Lock()
	h.PathUsername = filepath.Join(h.PathDir, "username")
	if d, err := ioutil.ReadFile(h.PathUsername); err == nil {
		h.Username = string(d)
	}
	h.PathPassword = filepath.Join(h.PathDir, "password")
	if d, err := ioutil.ReadFile(h.PathPassword); err == nil {
		h.Password = string(d)
	}
	h.mu.Unlock()

	if h.Username == "" || h.Password == "" {
		return errors.New("username or password is empty")
	} else {
		return nil
	}

}

// UserFileSecretProvider
func UserFileSecretProvider(dirpath string) SecretProvider {
	h := &UserSecretFile{PathDir: dirpath}
	h.Reload = func() error { return reloadUserFileSecret(h) }
	return func(username, realm string) string {
		h.ReloadIfNeeded()
		h.mu.RLock()
		defer h.mu.RUnlock()
		exists := (h.Username == username)
		password := h.Password
		if !exists {
			return ""
		}
		return password
	}
}

type TokenSecret struct {
	Reload func() error
	Token  string
	mu     sync.RWMutex
}

// StaticTokenSecretProvider
func StaticTokenSecretProvider(token string) SecretProvider {
	h := &TokenSecret{Token: token}
	h.Reload = func() error { return nil }
	return func(token, realm string) string {
		h.mu.RLock()
		defer h.mu.RUnlock()
		return h.Token
	}
}

type ServiceAccountTokenSecret struct {
	TokenSecret
	kubeconfig *rest.Config
}

// ServiceAccountTokenSecretProvider
func ServiceAccountTokenSecretProvider(c *rest.Config) SecretProvider {
	h := &ServiceAccountTokenSecret{kubeconfig: c}
	h.Reload = func() error { return nil }
	return func(token, realm string) string {

		claims, err := GetTokenClaims(token)
		if err != nil {
			log.Errorln(err.Error())
			return ""
		}
		ns := claims["kubernetes.io/serviceaccount/namespace"].(string)
		nm := claims["kubernetes.io/serviceaccount/secret.name"].(string)

		// secret 을 읽어올 cluster 선정
		apiClient, err := kubernetes.NewForConfig(h.kubeconfig)
		if err != nil {
			log.Warnf("cannot create a kubernetes api-client (cause=%s)", err)
			return ""
		}

		se, err := apiClient.CoreV1().Secrets(ns).Get(context.TODO(), nm, v1.GetOptions{})
		if err == nil {
			return string(se.Data["token"])
		} else {
			log.Warnf("cannot load token from service-account (namespace=%s,service-account=%s, cause=%s)", ns, nm, err)
			return ""
		}
	}
}
