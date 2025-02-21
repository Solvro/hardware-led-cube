# Frame generator

Useful for creating animations

## usage

cd to root folder of this repo and

```sh
python -m frame_generator
```

## Data format

### Prototype

json in the format:

```text
[frame][x][y][z]uint32
```

and uint32 represents an rgb color created like so: 0xRRGGBB
