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

A mask (see below) can optionally be supplied with the `-mask` or `-k` parameter.

## Definition file

A definition file is an array of rules for creating spritesheets. An example file looks like this:

```json
[
  {
    "prefix": "sheet",
	"suffixes": [
		"32bpp",
		"indexed"
		],
    "mask": "mask.png"
  }
]
```

Or, for multiple sheets:

```json
[
  {
    "prefixes": [
        "sheet1",
        "sheet2"
    ],
	"suffixes": [
		"32bpp",
		"indexed"
		],
    "mask": "mask.png"
  }
]
```

Each item has the following elements:


* `prefix` the start of an image filename. `prefixes` can alternatively be used as an array to process multiple files
   with the same settings.
* `suffix` the end of an image filename. This is useful for dealing with e.g. the output of GoRender which will produce
  files for `_32bpp`, `_8bpp` and `_mask`.
* `mask` (optional) a .PNG file to use as a mask. When compositing, only areas of the mask that are transparent will
   have pixel data from the input sprites written to them. This can also be globally supplied as a command line
   parameter, `-k`. The definition's mask takes precedence over the command line one.
     
All files which start `prefix` and end with the chosen `suffix` will be included (in alphabetical order) in the
resulting spritesheet. This is useful to note when working with mask files, which add a constraint to the order
in which sprites can be laid out.

## Indexed vs. RGBA

Indexed inputs will result in an indexed output. RGBA inputs will result in an RGBA output. It's assumed that all
files with a given suffix share the same colour depth and, if indexed, have the same palette. Results may be
inconsistent, or the program may crash if this is not the case.
