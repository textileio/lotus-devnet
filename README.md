# Lotus Devnet
Runs a Lotus Devnet using a mocked _sectorbuilder_ which can be accesed as a real Lotus node through the API.
The devnet code is always tried to be updated with Lotus `master` branch.

## Run
You can run the devnet with:
```bash
go run main.go
```

The devnet supports the following configuration:
- `-speed`: Time in milliseconds that blocks are mined. Default is 100ms.
- `-numminers`: Number of miners. Default is 1. (Note: higher values is an experimental feature)
- `-bigsectors`: Miners will use 512Gib sector sizes. Default is _false_ (2Kib sectors)
- `-ipfsaddr`: IPFS multiaddr to allow the client be connected to an IPFS node as a blockstorage.

All flags can be specified by enviroment variables using the `TEXLOTUSDEVNET_` prefix. e.g: `TEXLOTUSDEVNET_SPEED=1000`

## Docker
The Lotus Devnet was originally thought for integration tests in CI pipelines.

### Build locally
You can build the Docker image by running `docker build .`. 

## Public images
A public Docker image is hosted in DockerHub on every `master` commit.
Refer to [textile/lotus-devnet](https://hub.docker.com/repository/docker/textile/lotus-devnet/tags?page=1).

For example, running a Lotus Devnet with 1.5s block generation speed:
```bash
docker run -e TEXLOTUSDEVNET_SPEED=1500 textile/lotus-devnet
```

## Contributing

This project is a work in progress. As such, there's a few things you can do right now to help out:

-   **Ask questions**! We'll try to help. Be sure to drop a note (on the above issue) if there is anything you'd like to work on and we'll update the issue to let others know. Also [get in touch](https://slack.textile.io) on Slack.
-   **Open issues**, [file issues](https://github.com/textileio/go-threads/issues), submit pull requests!
-   **Perform code reviews**. More eyes will help a) speed the project along b) ensure quality and c) reduce possible future bugs.
-   **Take a look at the code**. Contributions here that would be most helpful are **top-level comments** about how it should look based on your understanding. Again, the more eyes the better.
-   **Add tests**. There can never be enough tests.

Before you get started, be sure to read our [contributors guide](./CONTRIBUTING.md) and our [contributor covenant code of conduct](./CODE_OF_CONDUCT.md).

## Changelog

[Changelog is published to Releases.](https://github.com/textileio/go-threads/releases)

## License

[MIT](LICENSE)
