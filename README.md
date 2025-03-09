# hardware-led-cube

![3DLEDCUBE_animation](https://github.com/Solvro/hardware-led-cube/blob/additional-assets/Aura-Cube.gif)

## Description

This repository contains code for configuring hardware, developing drivers, and managing services for a 3D LED Cube. It includes a custom service to control LED patterns, animations, and functionality, providing an interface to determine and program how the LED matrix operates.

## Project with workflow management

[Tasks here](https://github.com/orgs/Solvro/projects/28)

## Usage

### Dependencies

- To build the mock binary,
ensure that you fulfill the requirements for the
[go-gl](https://github.com/go-gl/gl) package.

- To build the binary for the raspberry pi,
ensure to fulfill the requirements described [here](https://github.com/rpi-ws281x/rpi-ws281x-go?tab=readme-ov-file#installing).

### Building

```bash
go build
```

- This produces an executable file named `hardware-led-cube`.
If the project was not compiled on arm64 linux (RPi), executing the binary will open up a window with a mock cube.
 
- Compiling on arm64 linux produces a binary that will use the RPi-ws281x library to control physical LEDs.

### Running

```bash
./hardware-led-cube
```
