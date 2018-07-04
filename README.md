# RSS2Hook

This project is a self-hosted utility which will make HTTP POST
requests to remote web-hooks when new items appear in an RSS feed.


## Rational

I have a couple of webhooks in-place already which will take incoming
HTTP requests and "do stuff" with them, for example:

* Posting to my alerting system.
   * Which is called [purppura](https://github.com/skx/purppura/) and is pretty neat.
* Posting to IRC.
   * IRC was mattermost before slack before born.

I _also_ have a bunch of RSS feeds that I follow, typically these include
github releases of projects.  For example my git-host is [gitbucket](https://github.com/gitbucket/gitbucket) so I subscribe to the release feed via:

* https://github.com/gitbucket/gitbucket/releases.atom


## Deployment

If you have a working golang setup you should be able to install this
application via:

    go get -u  github.com/skx/rss2hook
    go install github.com/skx/rss2hook


## Setup

There are two parts to the setup:

* Configure the list of feeds and the corresponding hooks to post to.
* Ensure the program is running.

For the first create a configuration-file like so:

    http://example.com/feed.rss = https://webhook.example.com/notify/me
...

For the second you can use your favourite supervision took, but in short
you'll want to run something like this:

     $ rss2hook -config ./sample.cfg


### Sample Webhook

There is a simple webhook example beneath [webhook/](webhook/) which
will listen upon localhost:8080, and dump any POST submission to the
console.

You can launch it like so:

     cd webhook/
     go run webhook.go

Testing it via `curl` would look like this:

      $ curl --header "Content-Type: application/json"  \
      --request POST \
      --data '{"username":"xyz","password":"xyz"}' \
      http://localhost:8080/

Finally you'd use the [sample.cfg](sample.cfg) file to POST to this
server by launching the application:

    $ rss2hook --config=sample.cfg

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
