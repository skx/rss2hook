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



## Deployment

If you have a working golang setup you should be able to install this
application via:

    go get -u  github.com/skx/rss2hook
    go install github.com/skx/rss2hook



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



## Feedback?

Turning this into a SaaS project would be interesting.  A simple setup
would be very straight-forward to implement, however at a larger scale
it would get more interesting:

* Assume two people have subscribed to the same feed.
   * But they did so a few days apart.
* That means what is "new" to each of them differs.
   * So you need to keep track of "seen" vs. "new" on a per-user __and__ per-feed basis.

Anyway it would be fun to implement, but I'm not sure there is a decent
revenue model out there for it.  Especially when you can wire up [IFTTT](https://ifttt.com/) or [similar](https://zapier.com/apps/rss/integrations/webhook/1746/send-a-webhook-when-an-rss-feed-is-updated) system to do the same thing.


Steve
--
