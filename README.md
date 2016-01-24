# Overview

Recipe is a tool for downloading and querying recipes.

Recipe provides:
* Commands to search and download recipes
* Indexers for "allrecipes.com"
* Indexes are cached on disk and can be refreshed.

# Install
Use go get to install the latest version.
```
go get github.com/gophergala2016/recipe
```

# Usage
The `recipe refresh` command will index all "repositories".
It will take a while on the first run.
```
recipe refresh
```

The `recipe search <term>` command will search the locally cached index for recipes
which contain the search term.
```
recipe search <term>
```

Use the `recipe -h` command to see additional information.
