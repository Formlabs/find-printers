<!-- PROJECT LOGO -->
<br />
<div align="center">
  <img src="https://github.com/user-attachments/assets/6a9948cc-1ab0-4ab7-be09-e2498def5cc9"/>


  <h3 align="center">find-printers</h3>

  <p align="center">
    CLI tool to find Formlabs printer information
    <br />
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li>
        <a href="#usage">Usage</a>
        </li>
      </ul>
    </li>
    <li><a href="#development">Development</a></li>
    <li><a href="#project-structure">Project Structure</a></li>
    <li><a href="#known-issues">Known Issues</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

![image](https://github.com/user-attachments/assets/3a210419-f1d4-4529-a94f-31552d916ee8)


`find-printers` is a small TUI / CLI tool build on top of [CASE](https://case.factory.priv.prod.gcp.formlabs.cloud/) API to query printer information including IP, SerialName, MachineType and Firmware Version. All of this can be quickly copied with shortcuts.


## Getting Started

Since the tool has no knowledge of Borg secrets, it requires to first setup a Borg / CASE authentication token. The token can be aquired from [Borg / CASE UI](https://case.factory.priv.prod.gcp.formlabs.cloud/).

Install `find-printers`:
```
$ go install github.com/Formlabs/find-printers@latest
```

Setup the token:
```
$ find-printers set-token <token>
```

Use the tool:
```
$ find-printers
```

### Usage
`find-printers` has some handy shortcuts to filter, sort and copy data

`Q / CTRL + C`  - Quit

`/` - Start searching, typing characters filters the data, its possible to search in all columns. Pressing `Enter` after search allows to sort, or copy from filtered data

`S / T / I / F` - Sort by Column (Serial, Machine Type, IP address, Firmware)

`s / t / i / f` - Copy column in selected row (Serial, Machine Type, IP address, Firmware)

## Development

To locally run the tool the easiest is to use the provided `make` commands
```
make build      - builds the binary
make install    - build the binary and links it to `find-printers` command
make run        - runs the tool locally
```


### Project Structure
```
.
├── Makefile
├── README.md
├── borg/          - package for interacting with Borg API
├── go.mod
├── go.sum
├── images/        - images for this README
├── main.go        - entrypoint
└── ui/            - code for the TUI to display information
```
```

## Known issues
Since the data depends on Borg data as a source, it automatically inherits all problems that Borg has.
It can sometimes show outdated data, or missing some printers from the list. Overall I find it pretty reliable, but had some cases where the printers were missing from the datasource.
