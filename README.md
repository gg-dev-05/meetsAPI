# meetsAPI


An API designed for scheduling meetings.

## USAGE
* /meetings : *(POST)* recieves requests along with JSON request as body to create a new meeting and  
returns back the ObjectID of the meeting.
* /meeting/<id> : *(GET)* sends back the meeting information in JSON format corresponding to the meeting id
* /meeting?start=startTIme&end=endTime : *(GET)* sends back the meetings lying in the given time period
* /meetings?participant=<email id> : *(GET)* sends back the meetings in which the participant with the given email id is the part of 
