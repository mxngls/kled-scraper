# KLED-Parser

## Motivation

The abbreviation KLED stands for "Korean Learner's English Dictionary" which is based on the [Korean Leaner's Dictionary](https://krdict.korean.go.kr/mainAction). They even provide an extense API which you can use if you register an user account on their website. I found two particular problems with their API that bothered me a lot:
1. For some reason the multimedia part (videos, pronounciation files etc.) does not work properly.
2. The words are not highlighted in anyway in the example sentences. 

For one of my other [projects](https://github.com/Mxngls/kled-server) I wrote an API endpoint that fetches and parses the html for a given word. For future use I thought it might be greate to just scrape the whole (english-korean) dictionary. 

## Description

The given script does exactly that. As the dictionary itself is licensed under the Creative Common License there are no copyright issues at all to worry about. 

For the example output see the ```dict``` directory. I obtained the necessary ids by sending of an empty prompt with the advanced search. The result contains all entries there are in 'alphabetical' order.[^1]

## Installation & Usage

Given that Go is installed just clone the repository:
```zsh
git clone git@github.com:Mxngls/kled-server.git
```

Version:
```zsh
% go version go1.17.6 darwin/amd64
```

Then run:
```zsh
% go run .
```

The result will be a JSON file that's roughly 70 MB.

[^1]:It should be noted that it is questionable if that behavior is intentioned though. 