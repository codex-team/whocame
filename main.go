package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var globalState Config

type Member struct {
	Name       string    `yaml:"name"`
	LastOnline time.Time `yaml:"lastOnline"`
	IsOnline   bool      `yaml:"isOnline"`
}

type Config struct {
	Webhook   string              `yaml:"webhook"`
	LogLevel  string              `yaml:"logLevel"`
	GoneAfter time.Duration       `yaml:"goneAfter"`
	Members   map[string][]string `yaml:"members"`
	State     map[string]Member   `yaml:"-"`
	Interface string              `yaml:"interface"`
}

func (c *Config) load(filename string) *Config {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("YAML error #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
	c.State = make(map[string]Member)
	for name, _ := range c.Members {
		c.State[name] = Member{Name: name}
	}
	return c
}

func (c *Config) memberDeviceFound(member Member) {
	currentTime := time.Now()

	if !member.IsOnline {
		log.Debugf("%s last seen %s, now online", member.Name, member.LastOnline)
		c.notify(fmt.Sprintf("👨‍💻 %s at the CodeX Lab", member.Name))
	}
	member.LastOnline = currentTime
	member.IsOnline = true
	c.State[member.Name] = member
}

func (c *Config) checkIfMemberGone() {
	currentTime := time.Now()

	for _, member := range c.State {
		log.Debugf("%s %s Test: %v %v", member.Name, member.LastOnline, currentTime.Sub(member.LastOnline), c.GoneAfter)
		if member.IsOnline && currentTime.Sub(member.LastOnline) > c.GoneAfter {
			log.Debugf("%s last seen %s, now offline", member.Name, member.LastOnline)

			member.IsOnline = false
			c.State[member.Name] = member

			c.notify(fmt.Sprintf("🏃 %s gone home", member.Name))
		}
	}
}

func (c *Config) processUpDevice(macAddr string) {
	for name, members := range c.Members {
		for _, device := range members {
			if device == macAddr {
				log.Debugf("Found %s's mac %s", name, macAddr)
				c.memberDeviceFound(c.State[name])
				return
			}
		}
	}
}

func (c *Config) notify(message string) {
	data := url.Values{}
	data.Set("message", message)
	MakeHTTPRequest("POST", c.Webhook, []byte(data.Encode()), map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
}

func main() {
	log.SetOutput(os.Stdout)

	globalState.load("storage.yml")
	log.Debugf("%v", globalState)

	level, err := log.ParseLevel(globalState.LogLevel)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)

	go func() {
		for {
			globalState.checkIfMemberGone()
			time.Sleep(15 * time.Second)
		}
	}()

	// start forever loop of ARP searching
	iface, err := net.InterfaceByName(globalState.Interface)
	if err = scan(iface); err != nil {
		log.Printf("interface %v: %v", iface.Name, err)
	}
}
