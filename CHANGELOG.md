
<a name="EdgeX UART Device Service (found in device-uart) Changelog"></a>
## EdgeX UART Device Service
[Github repository](https://github.com/edgexfoundry/device-uart)

### Change Logs for EdgeX Dependencies
- [device-sdk-go](https://github.com/edgexfoundry/device-sdk-go/blob/main/CHANGELOG.md)
- [go-mod-core-contracts](https://github.com/edgexfoundry/go-mod-core-contracts/blob/main/CHANGELOG.md)
- [go-mod-bootstrap](https://github.com/edgexfoundry/go-mod-bootstrap/blob/main/CHANGELOG.md)  (indirect dependency)
- [go-mod-messaging](https://github.com/edgexfoundry/go-mod-messaging/blob/main/CHANGELOG.md) (indirect dependency)
- [go-mod-registry](https://github.com/edgexfoundry/go-mod-registry/blob/main/CHANGELOG.md)  (indirect dependency)
- [go-mod-secrets](https://github.com/edgexfoundry/go-mod-secrets/blob/main/CHANGELOG.md) (indirect dependency)
- [go-mod-configuration](https://github.com/edgexfoundry/go-mod-configuration/blob/main/CHANGELOG.md) (indirect dependency)

## [4.0.0] Odessa - 2025-03-12 (Only compatible with the 4.x releases)

### ✨  Features

- Enable PIE support for ASLR and full RELRO ([3ed6416…](https://github.com/edgexfoundry/device-uart/commit/3ed64162fe5da762afcd5902d8c2c408f66a7072))

### ♻ Code Refactoring

- Update module to v4 ([53189c1…](https://github.com/edgexfoundry/device-uart/commit/53189c10332bacaee41a57c005c834aa638a9f96))
```text

BREAKING CHANGE: update go module to v4

```

### 🐛 Bug Fixes

- Only one ldflags flag is allowed ([d3f616b…](https://github.com/edgexfoundry/device-uart/commit/d3f616bb6c298463b0160a7450de9f955e758ebd))

### 👷 Build

- Upgrade to go-1.23, Linter1.61.0 and Alpine 3.20 ([44e0fc2…](https://github.com/edgexfoundry/device-uart/commit/44e0fc2ef96c30d275bce862755fe28b7851df0b))


## [v3.1.0] Napa - 2023-11-15 (Only compatible with the 3.x releases)


### ♻ Code Refactoring

- Remove github.com/pkg/errors from Attribution.txt ([447facf…](https://github.com/edgexfoundry/device-uart/commit/447facfb5636b95224070a3f0c345f144de32824))


### 📖 Documentation

- Add badges to README ([c4a221b…](https://github.com/edgexfoundry/device-uart/commit/c4a221b4730497f9ae59da6e1e9b701d59317945))


### 👷 Build

- Add missing .github files ([591fe36…](https://github.com/edgexfoundry/device-uart/commit/591fe365cb0d0d49483bd116b5c6838462170107))
- Upgrade to go-1.21, Linter1.54.2 and Alpine 3.18 ([a4a961d…](https://github.com/edgexfoundry/device-uart/commit/a4a961d24303d7dc17ca7a8bb4f9bd22d5aa1191))


### 🤖 Continuous Integration

- Add automated release workflow on tag creation ([5089f3c…](https://github.com/edgexfoundry/device-uart/commit/5089f3c862426ed73dedf318d808310baa83536c))

