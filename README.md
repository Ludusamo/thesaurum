# Thesaurum

## Docker Testing

### Build

`docker build -t thesaurum .`

### Run with binded mount

`docker run --env-file ./.env --rm --publish 5000:5000 -v "$(pwd)"/data:/data --name thesaurum thesaurum`

You have to create a `.env` file with the necessary environment variables set.
