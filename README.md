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