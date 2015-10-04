# RedFactor <img src="https://raw.githubusercontent.com/nickrobinson/redfactor/master/public/img/canary.png" alt="Drawing" width="150" height="110"/>
Alarming system built using InfluxDB as a backend



[More Details](http://nickrobinson.github.io/redfactor/) 

#Usage
In order to kick off redfactor you will first need to run `go build .` in the top level directory of the package. After
building run `./redfactor` in order to kick off the server. At this point the Redfactor UI should be up and running on
port 3000. You can then start adding queries and thresholds to the database from the main page.