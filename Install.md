# Installation instructions


## Pre-requisites

Make sure you have the [App Engine Go SDK][go_sdk] installed and a [Git][git] client.
You'll also need a Google Account to deploy ISumm on App Engine.

[go_sdk]: https://cloud.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go
[git]: https://git-scm.com/downloads

This guide will assume the usage of the command-line interface (CLI) of the above tools in a UNIX operating system.


## Running locally

The first thing we'll do is clone the this git repository and run a local instance of ISumm.

	% git clone git://github.com/danielfireman/isumm.git
	% cd isumm
	% goapp serve

All set!
Access <http://localhost:8080>, use `test@example.com` (default) as your login credentials and play around with the application.

If you create a lot of junk data and decide you want to start from scratch with real data, you can wipe the local datastore by running the application with the following command:

	% goapp serve --clear_datastore


## Deploying to App Engine

Edit `config.go` to limit access to just yourself.
Replace `danielfireman@gmail.com` with your Google Account email address.
Feel free to change any other configurations in that file that that might be relevant to you, such as the currency symbol.

Edit `app.yaml` to change the application name. We suggest you use the pattern `isumm-YOUR_NAME`.

Once you're set, let's create your deployment.
Go to the [Google Developers Console][console] and create a new project with whatever name you want. We suggest you use the pattern `isumm-YOUR_NAME`.


[console]: https://console.developers.google.com

Then, perform these steps to deploy:

	% goapp deploy -application isumm-YOUR_NAME

Point your browser to <http://isumm-YOUR_NAME.appspot.com> and voilà! :)

## Continuous deployment with codeship
codeship.io is free and supports continuous deployment to appengine.