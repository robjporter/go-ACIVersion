# go-ACIVersion
How often do you check your ACI Domain for the software it is running? When was the last time you checked if you are running a deferred release? Do you know what the suggested or latest software releases are?
These are some of the common questions I ask myself, then I ask, how do I go about checking the answers to these?

This is where this application will help.  It is a simple command line application which, once configured will connect to all of your ACI Domains, gain the current running APIC version, then query http://www.cisco.com for the latest, suggested, deferred and all releases available.

It will then do some simple comparisons between these different versions and produce some output to help you in your upgrade decisions.

## Setting up your GO environment
Depending on your particular environment, there are a number of ways to setup and install GO.  This repo was developed on a MAC and was installed using Brew.  For instructions on installing HomeBrew, please check [here](https://brew.sh/); and then entering;
```fish
> brew install go
```

If you do not want to use HomeBrew or you are running on a different platform, you can install the GO language using a binary from here;
https://golang.org/dl/

Once this has completed, open a cmd or terminal window and check GO has been installed and configured correctly;

Enter <b>echo $GOPATH</b>, hopefully you will be presented with a path and should be ready to go.

```fish
> echo $GOPATH
/path/to/go/bin/src/pkg folders
```

## Testing your GO environment
Once you have completed the above, its time to create a very simple test script to ensure everything is ready.

Go to a path where you are happy to store the source code for your application, this could be anywhere, including your desktop, documents, root folder, etc.

Create a folder and enter the directory.  Create a new file called "main.go" and enter the following code into it;

```go
package main

import "fmt"

func main() {
    fmt.Println("GO is working!")
}
```

At the command line, change directory using cd to the directory where your main.go file is and execute the following;
```fish
> go run main.go
```

You should see as output, something similar to;

"GO is working!"

If you reached this point, everything is working and you are ready to run the included code!

## Getting the code
There are a couple of ways you can get the code, depending on how comfortable you are with the command line and development environments;

You could download the zip file, [here](https://github.com/robjporter/go-UCSVersion/archive/master.zip).

You could use the command line git command to clone the repository to your local machine;
1. At the command line, change directory using cd to the directory where the repository will be stored.
2. Enter, git clone https://github.com/robjporter/go-UCSVersion.git
3. You will see output similar to the following while it is copied.
```fish
Cloning into `go-UCSVersion`...
remote: Counting objects: 10, done.
remote: Compressing objects 100% (8/8), done.
remove: Total 10 (delta 1), reused 10 (delta 1)
unpacking objects: 100% (10/10), done.
```
4. Change into the new directory, cd go-UCSVersion.
5. Move onto setting up the application.

## Application dependencies
For the application to work correctly, we need to get one dependency and we can achieve that with the following, via the cmd line.
```fish
> go get -u github.com/robjporter/go-functions
```

## Setting up the application
You need to add the ACI APICs to the application.  Your password will be encrypted before it is stored, however usernames will remain in plain text.  This should be a read only account on both systems, so should not cause too much of a security risk.

### Add ACI Domain
Repeat this process as many times as needed.
```go
> go run main.go add aci --ip=<IP> --username=<USERNAME> --password=<PASSWORD>
```

### Update ACI Domain
The update process will only succeed if the IP of the ACI APIC is already in the config file.
```go
> go run main.go update aci --ip=<IP> --username=<USERNAME> --password=<PASSWORD>
```

### Delete ACI Domain
The delete process will only succeed if the IP of the ACI APIC is already in the config file.
```go
> go run main.go delete aci --ip=<IP>
```
### Show ACI Domains
To show the current configuration details for a ACI APIC;
```go
> go run main.go show aci --ip=<IP>
```

### Show All discoverable systems
To show all the currently entered system information;
```go
> go run main.go show all
```

## Running the application
Once all the ACI APICs have been added, the application is now ready to run.
```go
> go run main.go run
```

## Building to a Binary
One of the great advantages of GO is the ability to compile the code and all dependencies into a single binary file.  This is enhanced by building for multiple platforms.  I have included a short script to compile to most of the common formats and place them in the ./bin folder.  You may need to add the execute ability onto the script, as this maybe removed during the download process, on a Mac, you can complete this by doing;
```fish
> chmod +x buildall.sh
```
To build the application run this;
```fish
> ./buildall.sh
```

## Interpreting the output
The output from running the application will look similar to the following;
```fish
|---------------------------------------------------------------------------------------------------------------------------|
| ACI APIC     | Current Version | Is Deferred | Suggested Version           | Latest Train Version        | Latest Version |
|---------------------------------------------------------------------------------------------------------------------------|
| 10.52.208.60 | 2.1(2e)         | false       | No Compatible version found | This is the latest release. | 2.2(1o)        |
|---------------------------------------------------------------------------------------------------------------------------|
```

**It is important to understand, this is not performing an inventory of your ACI environment or determining if the version is supported on all the hardware within your ACI environment.  Therefore, some manual validation needs to be completed, before running any update.**

#### The output is shown in 6 columns;
* **Column 1 - ACI APIC** - These are the URL's or IP's that have been entered for each ACI APIC.
* **Column 2 - Current Version** - These are the current ACI APIC versions that are running.
* **Column 3 - Is Deferred** - This column shows whether the executing version has been listed as a deferred release by Cisco.  If true, this would be a strong recommendation for an upgrade.
* **Column 4 - Suggested Version** - This column looks at the currently installed version and recommends the suggested version, it does not automatically suggest the latest or latest suggested release, as there maybe some very valid reasons that the current version is being run.  It tries to match the current version as much as possible, which should minimise upgrade impact.
* **Column 5 - Latest Train Version** - This checks the current version train, matching as many components of the version as possible, matching at least the first two.  Then looking for a newer version train.
* **Column 6 - Latest Version** - This is the very latest version that is available for ACI.