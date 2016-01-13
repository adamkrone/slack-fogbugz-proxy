Slack Fogbugz Integration proxy
===============================

This is a quick and dirty proxy that we use to relay information from our FogBugz environment
to Slack and help coordinate our support activities.

Enjoy, pull requests are very welcome :)

We (Thirdwave) have expanded it to support a config file, as well as posting to multiple
channels.

A configuration file should include:

- webhook
- fogbugz_url
- default_channel
- channel_mappings

The channel mappings allow you to control what channel should be used, depending on the case_number
sent in the URL Trigger.

An example of this config file is located in the repo.

On the FogBugz side, use the URL Trigger plugin to post to the proxy which
will in turn relay the information to Slack.

In the URL Trigger Plugin configuration add a new trigger and configure the
URL like you normally would but point it to the proxy.

E.g. if your plugin is running at http://some.host:9090 and you want to notify
Slack users when a case is opened/resolved/closed/reactivated you could enter:

http://some.host:9090/?case_number={CaseNumber}&project_name={ProjectName}&title={Title}

as the URL. You should use the GET verb for the trigger.

FogBugz is a bug tracking system from FogCreek Software:
https://www.fogcreek.com/fogbugz/

Slack is a team communication and sharing environment:
https://slack.com/
