[![Go Report Card](https://goreportcard.com/badge/github.com/skx/rss2hook)](https://goreportcard.com/report/github.com/skx/rss2hook)
[![license](https://img.shields.io/github/license/skx/rss2hook.svg)](https://github.com/skx/rss2hook/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/skx/rss2hook.svg)](https://github.com/skx/rss2hook/releases/latest)

* [RSS2Hook](#rss2hook)
* [Rational](#rational)
* [Installation](#installation)
  * [Build without Go Modules (Go before 1.11)](#build-without-go-modules-go-before-111)
  * [Build with Go Modules (Go 1.11 or higher)](#build-with-go-modules-go-111-or-higher)
* [Setup](#setup)
  * [Sample Webhook Receiver](#sample-webhook-receiver)
* [Implementation Notes](#implementation-notes)
* [Github Setup](#github-setup)


# RSS2Hook

This project is a self-hosted utility which will make HTTP POST
requests to remote web-hooks when new items appear in an RSS feed.



## Rational

I have a couple of webhooks in-place already which will take incoming
HTTP submissions and "do stuff" with them, for example:

* Posting to my alerting system.
   * Which is called [purppura](https://github.com/skx/purppura/) and is pretty neat.
* Posting to IRC.
   * IRC was mattermost before slack before born.

I _also_ have a bunch of RSS feeds that I follow, typically these include
github releases of projects.  For example my git-host runs [gitbucket](https://github.com/gitbucket/gitbucket) so I subscribe to the release feed of that, to ensure I'm always up to date:

* https://github.com/gitbucket/gitbucket/releases.atom


## Installation

There are two ways to install this project from source, which depend on the version of the [go](https://golang.org/) version you're using.

If you prefer you can fetch a binary from [our release page](https://github.com/skx/rss2hook/releases).  Currently there is only a binary for Linux (amd64) due to the use of `cgo` in our dependencies.

## Build without Go Modules (Go before 1.11)

    go get -u github.com/skx/rss2hook

## Build with Go Modules (Go 1.11 or higher)

    git clone https://github.com/skx/rss2hook ;# make sure to clone outside of GOPATH
    cd rss2hook
    go install



## Setup

There are two parts to the setup:

* Configure the list of feeds and the corresponding hooks to POST to.
* Ensure the program is running.

For the first create a configuration-file like so:

    http://example.com/feed.rss = https://webhook.example.com/notify/me

(There is a sample configuration file [sample.cfg](sample.cfg) which
will demonstrate this more verbosely.)

You can use your favourite supervision tool to launch the deamon, but you
can test interactively like so:

     $ rss2hook -config ./sample.cfg



### Sample Webhook Receiver

There is a simple webserver located beneath [webhook/](webhook/) which
will listen upon http://localhost:8080, and dump any POST submission to the
console.

You can launch it like so:

     cd webhook/
     go run webhook.go

Testing it via `curl` would look like this:

      $ curl --header "Content-Type: application/json"  \
      --request POST \
      --data '{"username":"blah","password":"blah"}' \
      http://localhost:8080/

The [sample.cfg](sample.cfg) file will POST to this end-point so you can
see how things work:

    $ rss2hook --config=sample.cfg



## Implementation Notes

* By default the server will poll all configured feeds immediately
upon startup.
   * It will look for changes every five minutes.
* To ensure items are only announced once state is kept on the filesystem.
   * Beneath the directory `~/.rss2hook/seen/`.
* Feed items are submitted to the webhook as JSON.



## Github Setup

This repository is configured to run tests upon every commit, and when
pull-requests are created/updated.  The testing is carried out via
[.github/run-tests.sh](.github/run-tests.sh) which is used by the
[github-action-tester](https://github.com/skx/github-action-tester) action.

Releases are automated in a similar fashion via [.github/build](.github/build),
and the [github-action-publish-binaries](https://github.com/skx/github-action-publish-binaries) action.

Steve
--
