# Thesaurum

Thesaurum is a multilayer cache that communicates over HTTP.

## Endpoints:

### `/topic`

- `GET`: Returns a list of lists of string topics in each cache layer

### `/topic/:topic`

- `GET`: Returns the data for the topic, it will continue to pull down from
  deeper layers if it isn't found in earlier layers.
- `DELETE`: Deletes the topic from all layers
- `POST`: Sets the topic in all layers

## Environment Variables

`CACHE_LAYERS`: Comma separated list of the cache layers to use

Valid Values:

- `InMemory`: LRU Cache that keeps a max amount of data cached in memory
- `File`: File backed cache

Default: `InMemory,File`

`MAX_MEMORY_CACHE`: The maximum memory to be cached `InMemory` cache

Default: 1MB (1024 * 1024 bytes)

`DATA_FILEPATH`: File path to store `File` cache data

Default: `./data`

## Local Testing

### Build

`go build -o thesaurum`

### Running

`./thesaurum`

## Docker Testing

### Build

`docker build -t thesaurum .`

### Run with binded mount

`docker run --env-file ./.env --rm --publish 5000:5000 -v "$(pwd)"/data:/data --name thesaurum thesaurum`

You have to create a `.env` file with the necessary environment variables set.
