# Lotus Devnet
Runs a Lotus Devnet using a mocked _sectorbuilder_ which can be accesed as a real Lotus node through the API.
This code is experimental.

The devnet code is always up to date with the latest `interopnet` Lotus branch.

## Run
You can run the devnet with:
```bash
go run main.go
```

The devnet supports the following configuration:
- _Speed_ (`-speed`): Time in milliseconds that blocks are mined
- _# Miners_ (`-numminers`): Amount of miners in testnet, default is 1. This feature is experimental, so don't expect to work pefectly for values greater than 1.

It also supports configuration via env-vars with `TEXLOTUSDEVNET_` prefix with capitalized flag names.

## Docker
You can run the image locally or leverage DockerHub images.

### Build locally
The docker-image can be built locally with `docker build .`. 

## Public images
Public images are pushed to DockerHub by the CI. Refer to [textile/lotus-devnet](https://hub.docker.com/repository/docker/textile/lotus-devnet/tags?page=1).
Run example with block generation in 1.5s intervals:
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
