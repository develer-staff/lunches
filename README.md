# Tinabot9000

## Intro

TODO

## Contributing

If you want to contribute to the repo (thanks!) you can take a look to [open issues](https://github.com/develersrl/lunches/issues).
Add a comment to inform the others yoy're going to work on it and start coding.  
Once you think the code is working and new tests pass, please [open a new PR](https://github.com/develersrl/lunches/compare).  
Someone will review it ASAP and will merge it.

### Setup the environment

TODO

* Install buffalo (what version?)
* Run `buffalo dev` (see below for buffalo setup)
* How to setup psql db locally? Is it required? (I think no)
* How to setup redis? Why do we use it?
* ?

### First good issues

If you're new to the project or Go you can take a look at the [first good issues](https://github.com/develersrl/lunches/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22).

### How to test new changes

TODO

* How to test a `t.bot.RespondTo` action?

## Tests coverage

TODO

## Deploy

When new code is pushed to `master` (directly or when a PR is merged), a new production deployment is triggered.  
The app runs on Heroku. The deploy is triggered automatically using a GitHub WebHook and the `heroku.yml` file.

## UI

Tinabot9000 provides a tiny UI (with debugging purposes atm) made with Buffalo.
You can reach it with the `heroku open` command.

### Database Setup

It looks like you chose to set up your application using a postgres database! Fantastic!

The first thing you need to do is open up the "database.yml" file and edit it to use the correct usernames, passwords, hosts, etc... that are appropriate for your environment.

You will also need to make sure that **you** start/install the database of your choice. Buffalo **won't** install and start postgres for you.

#### Create Your Databases

Ok, so you've edited the "database.yml" file and started postgres, now Buffalo can create the databases in that file for you:

```bash
$ buffalo db create -a
```

### Starting the Application

Buffalo ships with a command that will watch your application and automatically rebuild the Go binary and any assets for you. To do that run the "buffalo dev" command:

```bash
$ buffalo dev
```

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) you should see a "Welcome to Buffalo!" page.

## Roadmap

TBD
