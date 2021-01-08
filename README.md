# withings-sync
An application to practice Go Lang: Sync's a user's weight from Withings API to a user's profile weight attribute in Strava.

## Strava API
* https://developers.strava.com/docs/reference/

## Withings API
* https://developer.withings.com/

### GetMeas formula
Power of ten to multiply the value field to get the real value.
Formula: value * 10^unit = real value.
Eg: value = 20 and unit = -1 => real value = 2.

## Design

![](static/withings-sync.png)

