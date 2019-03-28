# ICD-10-CM Full-Text Search
A full-text search engine for the ICD-10-CM classification system. A live instance can be found at: https://icd.bemasher.net/

## Features
At the moment it supports queries by term, code, and prefix:
 * Each term in the query is normalized (e.g. "Café" -> "Cafe") and the porter2 word stemming algorithm reduces most terms to their roots (e.g. "gestation" -> "gestat").
 * Any term that ends with an asterisk will be treated as a prefix, so "M1A*" will return all entries containing codes in the M1A sub-category and any term that contains a note with M1A. Searching for M1A alone will return only entries that contain M1A.
 * Results that contain notes such as excludes1 and seventh character notes will have a badge indicating the term can be expanded by clicking on the title. This works well in web browsers on desktops/laptops, but not quite so well on mobile devices where touch is involved.
 * Seventh character notes have been propagated down to each sub-category or code's descendants. (Eventually generic notes will do the same, this is useful on codes like S32 that provide default coding instructions when certain aspects are unspecified.)
 * Hovering over an entry in the tabular index will provide a clipboard button to copy the code and it's description.
 * Entries are structured so that they are as close to what you would find in either index as possible so it can be useful for learning where and how things are represented in the indices.

## License
See: https://github.com/bemasher/icd10/blob/master/LICENSE

## Disclaimer
I have been testing this and it works well for my purposes, but I cannot make any guarantees about correctness or about the absence of errors and bugs.

## Bug Reporting
If you find any bugs or strange behaviors, please submit an issue. Same goes for any suggestions or improvements.