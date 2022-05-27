
## Slack Jobs Failure Notification

Cron is being setup on VM name slacker to send slack notifications for E2E jobs failure every day at 5AM UTC i.e. 10:30 IST.

`0 1 * * * ~/slacker >> /tmp/slacker.log`

Please check slacker.log for logs for any sort of failures

Notification will be sent to #tce-notifier on Vmware slack 
