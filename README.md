# CLI-Note-App

## Overview

[24 Feb 2025]

In my previous project, I started writing a simple CLI-based note-taking app using Go.

Link: [CLI-based note-taking app](https://github.com/rhysmah/note-app)

There was one aspect of the project that got rather complicated: the way I was attempting to handle the creation dates. Basically, I wanted to list and sort notes in several ways, including the creation date. Doing so for modified date is simple, because operating systems track this information, but they don't track file creation dates. To work around this -- in an admittedly crude manner -- I appended the creation date to the filename.

This caused enough problems that I decided to look for alternative, easier approaches. The one that seems to make sense -- and has apparently been used by other CLI-based apps -- is to use a locally stored database to save and more easily interact with the notes.

I'm thus starting a new project to implement this approach. I will be using code from the previous project, because I think there's some good stuff there.

In addition, I'm going to try using a more test-driven development approach. In the previous project, I created a custom logging function -- almost a mini side-project -- which was interesting and, for debugging purposes, very helpful. But I want something a little more robust than that; plus, I want to learn something new, and this is an important part of software development.

## Features

[24 Feb 2025]

This note-taking app will be faily simple.

### General Features
- Create a note
- List notes (by modified date, creation date)
- Edit notes (opens a note in the default editor of the OS)
- Delete notes
- Tag notes 
    - Add tag
    - Remove tag
    - List tags

## General Organization

[24 Feb 2025]

This app will have three layers:

1. **CLI Layer**: Go + Cobra to create a CLI commands.
2. **Service Layer**: Go functions for operations on notes.
3. **Database Layer**: uses BBolt, a pure Go key/value store that's fast, lightweight, and easy to use.

## Learnings

[25 Feb 2025]

My initial project organization plan is as follows:

```
cli-note-app/
    cmd/
        add/
            add.go
            validation.go
    db/
        db.go
    models/
        note.go
    go.mod
    go.sum
```

`cmd/` will contain the CLI commands, via Cobra, that the user can run. This will also contain the validation functions,
using the Validator Pattern, use to ensure that commands are valid. For example, when listing notes, we'll want to ensure that the user inputs valid commands, such as `created` or `modified`, and valid flags, such as `-asc` or `--desc`.

`db/` will contain BBolt database code.

`models/` will contain the `Note` struct, which will be used to represent a note object in the app. It will also contain an
interface, yet to be named, that will define the methods that the `Note` struct will implement.


[21 Mar 2025]

