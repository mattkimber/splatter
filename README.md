# Splatter

Splatter is a very small, simple utility for assembling a single spritesheet from multiple PNG files. It does not
attempt to do anything clever with layouts or producing the smallest possible file given the inputs - you will get
a row per input file, with the dimensions of the output being the width of the widest input file, and the combined
height of all inputs.

All inputs must be PNG and all outputs are PNG.

## Usage

`splatter -i input_directory -o output_directory -d def.json -m 8`

`splatter -input_dir input_directory -output_dir output_directory -definition def.json -margin 8`

(Both are equivalent)

`input_directory` is the directory containing your source input files. `output_directory` is where you want to put
the output files. `def.json` is a definition file, of which more details below. `margin` is the amount of vertical space to leave between rows, for input files which do not include this.

## Definition file

A definition file is an array of rules for creating spritesheets. An example file looks like this:

```json
[
  {
    "prefix": "sheet",
	"suffixes": [
		"32bpp",
		"indexed"
		]
  }
]
```

Each item has the following elements:

* `prefix` the start of an image filename.
* `suffix` the end of an image filename. This is useful for dealing with e.g. the output of GoRender which will produce
  files for `_32bpp`, `_8bpp` and `_mask`.
  
All files which start `prefix` and end with the chosen `suffix` will be included (in alphabetical order) in the
resulting spritesheet.

## Indexed vs. RGBA

Indexed inputs will result in an indexed output. RGBA inputs will result in an RGBA output. It's assumed that all
files with a given suffix share the same colour depth and, if indexed, have the same palette. Results may be
inconsistent, or the program may crash if this is not the case.
