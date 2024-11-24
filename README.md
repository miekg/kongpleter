Package kongpleter generates a yaml description of the command line as described by Kong.
This yaml can then be used by <https://github.com/miekg/gompletely> to generate specific completions.

The package supports an extra struct tag to influence the generated completion:

* A tag named `completion` that tells how to generate values for the completion.
* If a `help` tag is put on a type the help text must end in a dot.
