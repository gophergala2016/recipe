# Overview
Recipe is a tool for downloading and querying recipes.

Recipe provides:
- Commands to search and download recipes
- Indexers for "allrecipes.com"
- Indexes are cached on disk and can be refreshed.
- Support for schema.org/recipe format

# Install
Use go get to install the latest version.

```
go get github.com/gophergala2016/recipe/recipe
```

# Usage
The `recipe refresh` command will index all "repositories". It will take very long on the first run.
Before you can use the search and get command you need to run this command.

```
recipe refresh
```

The `recipe search <term>` command will search the locally cached index for recipes which contain the search term.

```
recipe search <term>
```

The `recipe get <term>` command will download all recipes which are indexed and match the term.

```
recipe get <term>
```

Use the `recipe -h` command to see additional information.

# The idea
The idea behind recipe is to have a tool similar to software package managers but for recipes.

# Disclaimer
Highly unstable, unoptimized and hacked together. Featuring inconsistent naming conventions and letting errors fall through. I added the bleve indexing for the recipe links last minute which also included saving each recipe link in its own file. That resulted in very slow indexing, probably due not using batches, no concurrent downloading and improper index mapping.

# What is missing
- more recipe repositories (indexers/websites) eg. wikibooks.org (mediawiki)
- support complete schema.org/recipe
- maybe hrecipe microformat
- proper storing and indexing of downloaded (recipe get) recipes
- printing/exporting recipes as html/markdown/printer/...
- repository management
- concurrent indexing/downloading of html documents
- test coverage
- documentation
