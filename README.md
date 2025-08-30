[![Sponsor](https://img.shields.io/badge/Sponsor-Loqa-ff69b4?logo=githubsponsors&style=for-the-badge)](https://github.com/sponsors/annabarnes1138)
[![Ko-fi](https://img.shields.io/badge/Buy%20me%20a%20coffee-Ko--fi-FF5E5B?logo=ko-fi&logoColor=white&style=for-the-badge)](https://ko-fi.com/annabarnes)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL--3.0-blue?style=for-the-badge)](LICENSE)
[![Made with ❤️ by LoqaLabs](https://img.shields.io/badge/Made%20with%20%E2%9D%A4%EF%B8%8F-by%20LoqaLabs-ffb6c1?style=for-the-badge)](https://loqalabs.com)

# 🏠 Loqa Device Service

[![CI/CD Pipeline](https://github.com/loqalabs/loqa-device-service/actions/workflows/ci.yml/badge.svg)](https://github.com/loqalabs/loqa-device-service/actions/workflows/ci.yml)

Device control service that listens on NATS for device commands and executes actions.

## Overview

Loqa Device Service is responsible for:
- Listening to NATS for device control commands (lights, audio, etc.)
- Executing actions on real or simulated devices
- Publishing device status and responses back to the message bus

## Features

- 📡 **NATS Integration**: Subscribes to device command subjects
- 💡 **Device Control**: Handles lights, music, and other smart home devices
- 🎯 **Command Processing**: Parses and executes structured device commands
- 📊 **Status Reporting**: Reports device state changes back to the system
- 🏠 **Extensible**: Easy to add new device types and integrations

## Supported Devices

- Lights (on/off, brightness, color)
- Music/Audio playback
- Temperature control
- Custom device handlers

## Getting Started

See the main [Loqa documentation](https://github.com/loqalabs/loqa) for setup and usage instructions.

## License

Licensed under the GNU Affero General Public License v3.0. See [LICENSE](LICENSE) for details.