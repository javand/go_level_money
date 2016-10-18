# go_level_money

## Features

 * prints out a summary of what this level user has been purchasing

## Install

Run `go get github.com/javand/go_level_money`
mkdir log

## Usage

Discover possible parameters by using the -help flag

Took out the credentials for security reasons, however you can update the defaults in the config.go

If you set your GOBIN to your go_level_money repo directory you can use 'go install' to build the executable

Example
./go_level_money -api_token GOTMEAPITOKEN -auth_token HAVEMEAUTHTOKEN -uid NEEDSAUSERID