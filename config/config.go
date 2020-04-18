package config

import (
	"math/rand"

	"github.com/slack-go/slack"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	SlackToken    string                   `yaml:"slack_token"`
	Greetings     []string                 `yaml:",flow"`
	Admins        []string                 `yaml:",flow"`

	Channels      map[string]slack.Channel //id to channel details
	AdminChannels map[string]string //admin id to IM channel id
}

/**
 * Parse the contents of the YAML file into the IrcConfig object.
 */
func (c *Config) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}

	if len(c.Greetings) == 0 {
		c.Greetings = []string{"hi"} //put something in there so it can still talk
	}

	c.Channels = make(map[string]slack.Channel)
	c.AdminChannels = make(map[string]string)

	return nil
}

/**
 * Returns a random greeting string from the list of valid greetings.
 */
func (c Config) GetRandomGreeting() string {
	return c.Greetings[rand.Intn(len(c.Greetings))]
}
