/*
Copyright © 2022 Rick Rackow <rick+cobra@rackow.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"flag"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"strings"
	"time"
)

func main() {
	// log in json for easier log parsing
	log.SetFormatter(&log.JSONFormatter{})

	// get all flags and parse
	schedulePtr := flag.String("schedule", "", "Schedule to use")
	aliasPtr := flag.String("alias", "", "Slack alias to adjust")
	slackTokenPtr := flag.String("slack-token", "", "Token to access slack")
	pdTokenPtr := flag.String("pd-token", "", "Token to access Pagerduty")
	debugPtr := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// check if we want debug logs
	if *debugPtr == true {
		log.SetLevel(log.DebugLevel)
	}
	// check if mandatory flags are set
	if *schedulePtr == "" {
		log.Fatal("You have to set a schedule.")
	}

	if *aliasPtr == "" {
		log.Fatal("You have to set an alias to use")
	}

	if *slackTokenPtr == "" {
		log.Fatal("You have to set a Slack token to use")
	}

	if *pdTokenPtr == "" {
		log.Fatal("You have to set a PagerDuty token to use")
	}

	oncallEmails := getOncall(*pdTokenPtr, *schedulePtr)

	// Clean the current members from the group
	err := removeFromAlias(*aliasPtr, *slackTokenPtr)
	if err != nil {
		log.Fatal("Failed to clean current alias")
	}

	// Add who's currently oncall to the slack alias

	err = addToAlias(oncallEmails, *aliasPtr, *slackTokenPtr)

	if err != nil {
		log.Fatal("Couldn't add new members to group")
	}

}

// getOncall gets the users currently oncall in PagerDuty for a given schedule using a given token
func getOncall(pdToken string, schedule string) []string {
	// start pagerduty connection
	pdClient := pagerduty.NewClient(pdToken)
	pdOpts := pagerduty.ListOnCallUsersOptions{
		Since: time.Now().String(),
		Until: time.Now().String(),
	}
	// get all users currently oncall for our schedule
	currentOncallUsers, err := pdClient.ListOnCallUsersWithContext(context.Background(), schedule, pdOpts)

	// fail if we have an error or no one is oncall
	if err != nil {
		log.Fatal(err)
	}
	if len(currentOncallUsers) < 1 {
		log.Fatal("No one is oncall right now")
	}

	// get the names from the users objects
	var currentOncall []string
	for _, user := range currentOncallUsers {
		currentOncall = append(currentOncall, user.Email)
	}
	return currentOncall
}

// getSlackID gets the ID for a given email address
func getSlackID(email string, token string) (string, error) {
	api := slack.New(token)
	user, err := api.GetUserByEmail(email)
	if err != nil {
		log.Fatal("Couldn't get user ID for %s", email)
	}
	return user.ID, err

}

// addToAlias adds a given user by their e-mail to a given group using the given token
func addToAlias(mails []string, alias string, token string) error {

	api := slack.New(token)

	var userList []string

	// convert our array of mails to a string of IDs
	for _, mail := range mails {
		userID, err := getSlackID(mail, token)
		if err != nil {
			log.Error(err)
		}
		userList = append(userList, userID)

	}

	_, err := api.UpdateUserGroupMembers(alias, strings.Join(userList, ","))

	return err

}

// removeFromAlias removes all current users from a group using a given token
func removeFromAlias(alias string, token string) error {
	api := slack.New(token)
	currentMembers, err := api.GetUserGroupMembers(alias)
	if err != nil {
		log.Fatal("Couldn't get current members of slack group")
	}
	log.Info("Current members: %s", currentMembers)
	log.Debug("Removing members")

	_, err = api.UpdateUserGroupMembers(alias, "")
	return err

}
