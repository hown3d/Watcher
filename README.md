# Watcher
A Terminal UI to watch AWS cloudformation stacks


## Installation
Download the binary according to your OS:

```
$ export OS=$(uname -s) # or "Windows" if you're using windoof
$ curl -L -o watcher https://github.com/hown3d/Watcher/releases/download/v0.0.1-alpha/watcher-$OS
$ chmod +x watcher
$ mv watcher /usr/local/bin
```

## Configuration
Watcher uses your current AWS profile, set by the environment variable **AWS_PROFILE**
```
$ watcher --help                 
Usage of watcher:
  -endpoint string
    	specify a custom endpoint for aws, for example localstack
```


## TODO: Add Demo Gif
